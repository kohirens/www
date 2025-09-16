package session

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// Load Will begin a new session, or restore an unexpired session, store the
// session ID in an HTTP cookie to use on the next request.
func (m *Manager) Load(w http.ResponseWriter, r *http.Request) {
	idCookie, _ := r.Cookie(IDKey)

	// ONLY set a new cookie when there is no session, or it has expired.
	if idCookie == nil {
		idCookie = m.IDCookie(IDCookiePath, IDCookieDomain)
		Log.Infof(stdout.IDSet)
		http.SetCookie(w, idCookie)
		return
	}

	if e := m.Restore(idCookie.Value); e != nil {
		Log.Errf(e.Error())
	} else { // Rolling session technique.
		// When we successfully restore a session, we extend it a bit.
		// Have the cookie also reflect this extended time.
		Log.Infof(stdout.Restored)

		// NOTE: Seems setting only the expiration time for a cookie in Go can a new cookie under
		// a path that the cookie was set on. This is not the intended action. We need complete replace it and
		// set the desired time.
		idCookie = m.IDCookie(IDCookiePath, IDCookieDomain)

		// When we restore we also extend the life of the session.
		// Update the session expiration time to match the session.
		// when we do this does it send the cookie with the update or de we need to also set it in the response again?
		idCookie.Expires = m.Expiration()
		// set the cookie so that the update takes effect.
		// locally this may work, but I'm not sure about in lambda/cloudfront world.
		// Update the cookie with the new time.
		// For clarity on updating HTTP Cookies, please see
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies#creating_removing_and_updating_cookies
		// or review https://datatracker.ietf.org/doc/html/rfc6265
		http.SetCookie(w, idCookie)
	}

	// TODO: Validate the session ID is valid by looking it up in storage,
	// if so, then also compare the browser data for a match, if not,
	// then expire the cookie immediately (tampering).
	if idCookie.Value != m.ID() {
		Log.Errf(stderr.SessionStrange)
		idCookie.Expires = time.Now().UTC()
		m.RemoveAll()
	}
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

// Restore Restores the session by ID as a string.
func (m *Manager) Restore(id string) error {
	// TODO validate ID maybe a regex
	if id == "" {
		return fmt.Errorf(stderr.EmptySessionID)
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

	// Verify the session has not expired.
	if time.Now().UTC().After(data.Expiration.UTC()) {
		return ExpiredError{data.Expiration}
	}

	m.data = data

	// Rolling session technique to extend the session a bit more since data was recently accessed.
	Log.Infof(stdout.CurrentTime, data.Expiration.Format("15:04:05"))
	timeLeft := m.Expiration().Sub(time.Now().UTC())
	if timeLeft < time.Minute*5 {
		m.data.Expiration = m.data.Expiration.Add(ExtendTime)
	}

	Log.Infof(stdout.ExtendTime, data.Expiration.Format("15:04:05"))

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

func (m *Manager) storagePath(id string) string {
	if m.location != "" {
		return m.location + "/" + id
	}
	return id
}
