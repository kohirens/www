package backend

import "net/http"

func InArray(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// isLoggedIn indicates when a client is logged in or not.
func isLoggedIn(a *Api, w http.ResponseWriter, r *http.Request) bool {
	sm, e1 := a.Session()
	if e1 != nil {
		return false
	}
	// Do not check logged-in state if this page is on the public list.
	uriPath := r.URL.Path
	Log.Dbugf(stdout.UriPath, uriPath)
	if InArray(uriPath, PublicPages) {
		Log.Dbugf(stdout.SkipLogin, uriPath)
		return true
	}

	// Not logged in when the session has expired.
	if sm.HasExpired() {
		return false
	}

	loggedIn := string(sm.Get(skLoggedIn))

	Log.Dbugf(stdout.LoggedIn, loggedIn)

	return loggedIn == "true"
}
