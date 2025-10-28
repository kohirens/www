package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kohirens/sso"
	"github.com/kohirens/sso/pkg/google"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/backend"
	"github.com/kohirens/www/storage"
	"github.com/kohirens/www/validation"
	"net/http"
)

// Legend:
// * f - field
// * max - maximum
const (
	// Used as a hint when the user attempts to login with the provider.
	fEmail = "email"
	fCode  = "code"
	fState = "state"
	name   = "google"
)

var (
	// CallbackRedirect A location the client will be sent after a successful callback.
	CallbackRedirect = "/"
	// Log Set a logger, must be compatible with Kohirens stdlib/logger.
	Log = &logger.Standard{}
	// SignOutRedirect A location to send the client after they sign out.
	SignOutRedirect = "/"
)

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
	Log.Dbugf("%v", stdout.GoogleCallback)

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

	// Get client account info
	ams, e3 := a.ServiceManager().Get(backend.KeyAccountManager)
	if e3 != nil {
		return e3
	}
	am := ams.(backend.AccountManager)

	// Retrieve the storage manager.
	sd, e5 := a.Service(backend.KeyStorage)
	if e5 != nil {
		return e5
	}
	store := sd.(storage.Storage)

	account, e6 := GetAccount(am, gp, store)
	if e6 != nil {
		return e6
	}

	userAgent := r.Header.Get("User-Agent")
	Log.Infof("user-agent: %v", userAgent)

	// Store that token away for safe keeping
	if e7 := gp.SaveLoginInfo(); e7 != nil {
		return e7
	}

	encodeMessage, e8 := a.Encrypt(account.ID)
	if e8 != nil {
		return e8
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "_aid_",
		Value:  string(encodeMessage),
		Path:   "/",
		Secure: true,
	})

	// send user to a predetermined link or the dashboard.
	w.Header().Set("Location", CallbackRedirect)
	w.WriteHeader(http.StatusSeeOther)
	return nil
}

// GetAccount Lookup an existing account or make a new account only when
// a client has a successful login and an existing account cannot be found.
func GetAccount(
	am backend.AccountManager,
	gp *google.Provider,
	store storage.Storage,
) (*backend.Account, error) {
	account, e1 := am.Lookup(gp.ClientID())

	var notFound *backend.AccountNotFoundError

	// Make a new account only when a client has a successful login and
	// an existing account cannot be found.
	if errors.As(e1, &notFound) {
		Log.Errf("%v", e1.Error())
		acct, e3 := am.AddWithProvider(gp.ClientID(), gp.Name())
		if e3 != nil {
			return nil, e3
		}

		acct.GoogleId = gp.ClientID()
		acct.Email = gp.ClientEmail()

		aData, e4 := json.Marshal(account)
		if e4 != nil {
			return nil, fmt.Errorf(stderr.EncodeJSON, e4.Error())
		}

		// Save the account in <storage>/accounts/<account-id>
		if e := store.Save("accounts/"+account.ID, aData); e != nil {
			return nil, e
		}
	}

	return account, nil
}
