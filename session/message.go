package session

var stderr = struct {
	DecodeJSON,
	EmptySessionID,
	EncodeJSON,
	ExpiredCookie,
	NoStorage,
	NoSuchKey,
	ReadFile,
	SessionStrange,
	WriteFile string
}{
	DecodeJSON: "could not decode JSON from file %v: %w",
	//DecodeJSON:     "could not decode json data: %v",
	EmptySessionID: "session ID is empty",
	EncodeJSON:     "could not encode JSON: %w",
	ExpiredCookie:  "session has expired at %v",
	NoStorage:      "storage has not been set",
	NoSuchKey:      "the key %v was not found in the session",
	ReadFile:       "could not read file %v: %w",
	SessionStrange: "strangeness detected, the session is out of sync. expiring the current session cookie, the user will have to start a new session",
	WriteFile:      "could not write content to file %v: %w",
}

var stdout = struct {
	IDSet,
	Restored string
}{
	IDSet:    "setting a session ID cookie now",
	Restored: "session restored",
}
