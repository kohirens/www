package session

var stderr = struct {
	DecodeJSON,
	EmptySessionID,
	EncodeJSON,
	ExpiredCookie,
	InvalidSessionID,
	NoIDCookieFound,
	NoStorage,
	NoSuchKey,
	PhenomenonMismatchCookie,
	ReadFile,
	RestoreSession,
	SessionStrange,
	UUID,
	WriteFile string
}{
	DecodeJSON: "could not decode JSON from file %v: %w",
	//DecodeJSON:     "could not decode json data: %v",
	EmptySessionID:           "session ID is empty",
	EncodeJSON:               "could not encode JSON: %w",
	ExpiredCookie:            "session has expired at %v",
	InvalidSessionID:         "invalid session id ",
	NoIDCookieFound:          "no session ID cookie found",
	NoStorage:                "storage has not been set",
	NoSuchKey:                "the key %v was not found in the session",
	PhenomenonMismatchCookie: "session ID found in cookie value %v, does not match current session ID value %v",
	ReadFile:                 "could not read file %v: %w",
	RestoreSession:           "cannot to restore session ",
	SessionStrange:           "strangeness detected, the session is out of sync. expiring the current session cookie, the user will have to start a new session",
	UUID:                     "cannot generate UUID: %v",
	WriteFile:                "could not write content to file %v: %w",
}

var stdout = struct {
	CurrentTime,
	IDCookieFound,
	IDCookieValue,
	IDSessionValue,
	IDSet,
	Restored string
}{
	CurrentTime:    "session current time %v",
	IDCookieFound:  "attempting to set a session ID cookie, but one has been found",
	IDCookieValue:  "cookie session ID: %v",
	IDSessionValue: "current session ID: %v",
	IDSet:          "setting a session ID cookie now",
	Restored:       "session restored",
}
