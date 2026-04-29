package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Manager This is the container/interface for your session. Needed to make a
// new session or restore an existing one.
type Manager struct {
	data    *Data
	storage Storage
	// To save on network traffic default this to false, and set to true when
	// Set is called. Save is no-up if this is false.
	hasUpdates bool
	location   string
	mutex      sync.Mutex
	timeout    time.Duration
}

// Get Retrieve data from the session.
func (m *Manager) Get(key string) []byte {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	value, ok := m.data.Items[key]
	if ok {
		return value
	}

	return nil
}

// Expiration Retrieve expiration time.
func (m *Manager) Expiration() time.Time {
	return m.data.Expiration
}

// HasExpired this has been added to replace the 5-minute rolling extension,
// which made working with sessions hard because you could not know how long
// the clients session would last. As it would extend 5 five minutes every
// time a call to the backend was made from the client. This made session length
// indefinite. Not the intended purpose.
func (m *Manager) HasExpired() bool {
	currentTime := time.Now().UTC()
	sessionTime := m.data.Expiration.UTC()
	hasExpired := time.Now().UTC().After(m.data.Expiration)

	Log.Infof(stdout.CurrentTime, currentTime.Format(time.RFC3339))
	Log.Infof(stdout.SessionTime, sessionTime.Format(time.RFC3339))
	Log.Infof(stdout.SessionExpired, hasExpired)

	return hasExpired
}

// ID Of the session as an HTTP cookie with secure and http-only (cannot be read by JavaScript) enabled.
// The domain parameter is optional, and only set when it is not an emptry string.
func (m *Manager) ID() *uuid.UUID {
	return m.data.Id
}

func (m *Manager) IDCookie(cookiePath, domain string) *http.Cookie {
	c := &http.Cookie{
		Expires:  m.data.Expiration,
		Name:     IDKey,
		Path:     cookiePath,
		Secure:   true,
		HttpOnly: true,
		Value:    m.data.Id.String(),
		SameSite: http.SameSiteStrictMode,
	}

	if domain != "" {
		c.Domain = domain
	}

	return c
}

// LoadFromCookie will load a session from an HTTP cookie.
func (m *Manager) LoadFromCookie(r *http.Request) error {
	idCookie, e1 := r.Cookie(IDKey)

	if errors.Is(e1, http.ErrNoCookie) || idCookie == nil {
		return NoSessionCookieError{}
	}

	if e := m.Restore(idCookie.Value); e != nil {
		return RestoreError{e.Error()}
	}

	// Indicate the session has expired.
	if time.Now().UTC().After(m.data.Expiration.UTC()) {
		return ExpiredError{m.data.Expiration}
	}

	return nil
}

// Remove data from a session
func (m *Manager) Remove(key string) error {
	// verify the key exists
	_, ok := m.data.Items[key]
	if !ok {
		return fmt.Errorf(stderr.NoSuchKey, key)
	}

	// Indicate the session data needs to be saved.
	m.hasUpdates = true

	// Remove the key
	delete(m.data.Items, key)

	return nil
}

// RemoveAll When you need to scrub the data from the session and fast.
func (m *Manager) RemoveAll() {
	m.data.Items = Store{}
}

// Reset When you need to scrub the data from the session and fast.
func (m *Manager) Reset() {
	m.data = newData(m.timeout)
}

// Restart an expired session without removing any data.
func (m *Manager) Restart() {
	m.data.Expiration = time.Now().UTC().Add(m.timeout)
}

// Restore Restores the session by ID as a string.
func (m *Manager) Restore(id string) error {
	// Skip empty string.
	if strings.Trim(id, " \n\t\r") == "" {
		return fmt.Errorf("%v", stderr.EmptySessionID)
	}

	if m.storage == nil {
		return StorageError{}
	}

	// Load from storage.
	dataBytes, e1 := m.storage.Load(m.storagePath(id))
	if e1 != nil {
		return e1
	}

	// Validation will occur when the ID is restored with json.Ummarshal back
	// into an uuid.UUID.
	var data *Data
	if e := json.Unmarshal(dataBytes, &data); e != nil {
		return fmt.Errorf(stderr.DecodeJSON, e.Error())
	}

	// Verify the session ID does not morph.
	// TODO: Remove if there unmarshal restored the ID properly.
	// This code is actually testing the uuid library, it acts as a sanity check
	// that you can easily JSON serialize and then unmarshal a UUID.
	// IT can be removed at anytime after we are satisfied.
	if data.Id.String() != id {
		Log.Errf("%v", stderr.SessionStrange)
		return InvalidIDError{data.Id.String()}
	}

	m.data = data

	Log.Infof("%v", stdout.Restored)

	return nil
}

// Save Writes session data to its storage. This is no-op if Set was not previously called.
func (m *Manager) Save() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	dataBytes, e1 := json.Marshal(m.data)
	if e1 != nil {
		return fmt.Errorf(stderr.EncodeJSON, e1)
	}

	if m.hasUpdates {
		return m.storage.Save(m.storagePath(m.ID().String()), dataBytes)
	}

	return nil
}

// Set Store data in the session.
func (m *Manager) Set(key string, value []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.hasUpdates = true
	m.data.Items[key] = value
}

// SetCookie stored the session ID in a secure HTTP cookie.
//
//	This is no-op if the cookie has been previously set and the session has not
//	expired or the cookie deleted.
func (m *Manager) SetCookie(w http.ResponseWriter, r *http.Request) {
	idCookie, e1 := r.Cookie(IDKey)
	// Verify there is no cookie before we set a new one.
	// This is to prevent making orphans of sessions by overwriting them by
	// simply calling this method.
	if errors.Is(e1, http.ErrNoCookie) {
		Log.Dbugf("%v", stderr.NoIDCookieFound)
	}

	if idCookie != nil {
		Log.Dbugf("%v", stdout.IDCookieFound)
		Log.Dbugf(stdout.IDCookieValue, idCookie.Value)
		Log.Dbugf(stdout.IDSessionValue, m.ID().String())
		// This could happen when a session expires, as the transition from
		// session expiration to logging is a work in progress.
		if idCookie.Value != m.ID().String() {
			Log.Warnf(stderr.PhenomenonMismatchCookie, idCookie.Value, m.ID().String())
		}
		return
	}

	idCookie = m.IDCookie(IDCookiePath, IDCookieDomain)
	Log.Infof("%v", stdout.IDSet)
	http.SetCookie(w, idCookie)
}

// storagePath Returns a path to load/save a session to/from.
func (m *Manager) storagePath(id string) string {
	if m.location != "" {
		return m.location + "/" + id + Suffix
	}
	return id + Suffix
}
