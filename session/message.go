package session

var stderr = struct {
	DecodeJSON,
	EmptySessionID,
	EncodeJSON,
	ExpiredCookie,
	IDCookieFound,
	InvalidSessionID,
	NoIDCookieFound,
	NoStorage,
	NoSuchKey,
	ReadFile,
	RestoreSession,
	SessionStrange,
	UUID,
	WriteFile string
}{
	DecodeJSON: "could not decode JSON from file %v: %w",
	//DecodeJSON:     "could not decode json data: %v",
	EmptySessionID:   "session ID is empty",
	EncodeJSON:       "could not encode JSON: %w",
	ExpiredCookie:    "session has expired at %v",
	IDCookieFound:    "found session ID cookie with a value of %v",
	InvalidSessionID: "invalid session id ",
	NoIDCookieFound:  "no session ID cookie to found",
	NoStorage:        "storage has not been set",
	NoSuchKey:        "the key %v was not found in the session",
	ReadFile:         "could not read file %v: %w",
	RestoreSession:   "cannot to restore session ",
	SessionStrange:   "strangeness detected, the session is out of sync. expiring the current session cookie, the user will have to start a new session",
	UUID:             "cannot generate UUID: %v",
	WriteFile:        "could not write content to file %v: %w",
}

var stdout = struct {
	CurrentTime,
	IDSet,
	Restored string
}{
	CurrentTime: "session current time %v",
	IDSet:       "setting a session ID cookie now",
	Restored:    "session restored",
}
