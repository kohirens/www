package google

import (
	"fmt"
	"github.com/kohirens/sso/pkg/google"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/backend"
	"github.com/kohirens/www/validation"
	"net/http"
)

// Legend:
// * f - field
// * max - maximum
const (
	CallbackRedirect = "/"
	SignOutRedirect  = "/"
	fEmail           = "email"
	fCode            = "code"
	fState           = "state"
)

var Log = &logger.Standard{}

// AuthLink Build link to authenticate with Google.
func AuthLink(w http.ResponseWriter, r *http.Request, a backend.App) error {
	email, emailOK := validation.Email(r.URL.Query().Get(fEmail))
	if !emailOK {
		email = "" // It's not required, so it is O.K. to leave it out.
	}

	p, e1 := a.AuthManager().Get(backend.KeyGoogleProvider)
	if e1 != nil {
		return e1
	}
	gp := p.(*google.Provider)

	authURI, e2 := gp.AuthLink(email)
	if e2 != nil {
		return e2
	}

	s := fmt.Sprintf(`{"status": %q, "link": %q}`, "ok", authURI)

	_, e3 := w.Write([]byte(s))
	if e3 != nil {
		return fmt.Errorf(stderr.EncodeJSON, e3.Error())
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	return nil
}

// SignIn Begin the authentication process for a client.
func SignIn(w http.ResponseWriter, r *http.Request, a backend.App) error {
	if e := r.ParseForm(); e != nil {
		return fmt.Errorf(stderr.ParseSignInData, e.Error())
	}

	email, emailOK := validation.Email(r.PostForm.Get(fEmail))
	if !emailOK {
		w.Header().Set("Location", "/?m=invalid-email")
		return backend.NewReferralError(stderr.ValidEmail, "/?m=invalid-email", http.StatusTemporaryRedirect, true)
	}

	p, e1 := a.AuthManager().Get(backend.KeyGoogleProvider)
	if e1 != nil {
		return e1
	}
	gp := p.(*google.Provider)

	authURI, e2 := gp.AuthLink(email)
	if e2 != nil {
		return e2
	}

	// set a redirect for the browser.
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Location", authURI)
	w.WriteHeader(http.StatusTemporaryRedirect)

	return nil
}

// SignOut Invalidate a authentication token.
func SignOut(w http.ResponseWriter, r *http.Request, a backend.App) error {
	endpoint := "/?signed-out=1"

	p, e2 := a.AuthManager().Get(backend.KeyGoogleProvider)
	if e2 != nil {
		return e2
	}
	gp := p.(*google.Provider)

	if e := gp.SignOut(); e != nil {
		Log.Errf(stderr.SignOut, e)
	}

	body := []byte(fmt.Sprintf(backend.MetaRefresh, endpoint))
	_, e3 := w.Write(body)
	if e3 != nil {
		return fmt.Errorf(stderr.WriteResponseBody, e3.Error())
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Location", SignOutRedirect)
	w.WriteHeader(http.StatusTemporaryRedirect)

	return nil
}

// Callback Handles callback request initiated from a Google
// authentication server when the client chose to sign in with Google.
func Callback(w http.ResponseWriter, r *http.Request, a backend.App) error {
	Log.Dbugf(stdout.GoogleCallback)

	queryParams := r.URL.Query()
	code := queryParams.Get(fCode)
	state := queryParams.Get(fState)

	p, e1 := a.AuthManager().Get(backend.KeyGoogleProvider)
	if e1 != nil {
		return e1
	}
	gp := p.(*google.Provider)

	// Exchange the 1 time code for an ID and refresh tokens.
	if e3 := gp.ExchangeCodeForToken(state, code); e3 != nil {
		return e3
	}

	// send user to a predetermined link or the dashboard.
	w.Header().Set("Location", CallbackRedirect)
	w.WriteHeader(http.StatusSeeOther)
	return nil
}
