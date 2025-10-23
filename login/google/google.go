package google

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/sso"
	"github.com/kohirens/sso/pkg/google"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/backend"
	"github.com/kohirens/www/gpg"
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
	gp := p.(sso.OIDCProvider)

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

	email := r.PostForm.Get(fEmail)
	_, emailOK := validation.Email(email)
	if email != "" && !emailOK {
		w.Header().Set("Location", "/?m=invalid-email")
		return backend.NewReferralError(
			"",
			stderr.ValidEmail,
			"/?m=invalid-email",
			http.StatusTemporaryRedirect,
			true,
		)
	}

	p, e1 := a.AuthManager().Get(backend.KeyGoogleProvider)
	if e1 != nil {
		return e1
	}
	gp := p.(sso.OIDCProvider)

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
	gp := p.(sso.OIDCProvider)

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
	if e2 := gp.ExchangeCodeForToken(state, code); e2 != nil {
		return e2
	}

	// Store that token away for safe keeping
	if e3 := gp.SaveLoginInfo(backend.KeyLoginPrefix); e3 != nil {
		return e3
	}
	// Get client account info
	x, e4 := a.ServiceManager().Get(backend.KeyAccountManager)
	if e4 != nil {
		return e4
	}
	am := x.(backend.AccountManager)

	account, e5 := am.Lookup(gp.ClientID())
	switch e5.(type) {
	case *backend.AccountNotFoundError:
		Log.Errf(stderr.SignOut, e5)
		// TODO Make a new one and save it in the account store.
		deviceID := NewDeviceId(r)
		acct, e6 := am.Add()
		if e6 != nil {

			return e6
		}
		acct.GoogleId = gp.ClientID()
		acct.Email = gp.ClientEmail()
	}

	aData, e7 := json.Marshal(account)
	if e7 != nil {
		return fmt.Errorf(stderr.EncodeJSON, e7.Error())
	}
	// TODO: Encrypt the account ID and store it in a cookie.
	gpg.NewCapsule(pubKeyFile, privKeyFile, passPhrase)

	// send user to a predetermined link or the dashboard.
	w.Header().Set("Location", CallbackRedirect)
	w.WriteHeader(http.StatusSeeOther)
	return nil
}
