package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
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
	expiration time.Duration
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

	Log.Infof("current time %v", currentTime.Format(time.RFC3339))
	Log.Infof("session time %v", sessionTime.Format(time.RFC3339))
	Log.Infof("session has expired %v", hasExpired)

	return hasExpired
}

// ID Of the session as an HTTP cookie with secure and http-only (cannot be read by JavaScript) enabled.
// The domain parameter is optional, and only set when it is not an emptry string.
func (m *Manager) ID() string {
	return m.data.Id
}

func (m *Manager) IDCookie(cookiePath, domain string) *http.Cookie {
	c := &http.Cookie{
		Expires:  m.data.Expiration,
		Name:     IDKey,
		Path:     cookiePath,
		Secure:   true,
		HttpOnly: true,
		Value:    m.data.Id,
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
		return NoSessionError{}
	}

	if e := m.Restore(idCookie.Value); e != nil {
		return RestoreError{e.Error()}
	}

	// Verify the session ID in the cookie matches the actual  valid by looking it up in storage,
	// if so, then also compare the browser data for a match, if not,
	// then expire the cookie immediately (tampering).
	if idCookie.Value != m.ID() {
		Log.Errf("%v", stderr.SessionStrange)
		m.RemoveAll()
		return InvalidIDError{idCookie.Value}
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
	m.data = newData(m.expiration)
}

// Restore Restores the session by ID as a string.
func (m *Manager) Restore(id string) error {
	// Validate ID by a regex (see https://stackoverflow.com/questions/136505/searching-for-uuids-in-text-with-regex).
	//re := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	//if !re.MatchString(id) {
	if strings.Trim(id, " \n\t\r") == "" {
		return fmt.Errorf("%v", stderr.EmptySessionID)
	}

	if m.storage == nil {
		return StorageError
	}

	// Load from storage.
	dataBytes, e1 := m.storage.Load(m.storagePath(id))
	if e1 != nil {
		return e1
	}

	var data *Data
	if e := json.Unmarshal(dataBytes, &data); e != nil {
		return fmt.Errorf(stderr.DecodeJSON, e.Error())
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
		return m.storage.Save(m.storagePath(m.ID()), dataBytes)
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
	if idCookie, _ := r.Cookie(IDKey); idCookie != nil && idCookie.Value == m.ID() {
		return
	}

	idCookie := m.IDCookie(IDCookiePath, IDCookieDomain)
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
