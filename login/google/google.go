package google

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/sso"
	"github.com/kohirens/sso/pkg/google"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/backend"
	"github.com/kohirens/www/gpg"
	"github.com/kohirens/www/session"
	"github.com/kohirens/www/storage"
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
	name             = "google"
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
	if e3 := gp.SaveLoginInfo(backend.PrefixLogin); e3 != nil {
		return e3
	}

	// Get client account info
	ams, e4 := a.ServiceManager().Get(backend.KeyAccountManager)
	if e4 != nil {
		return e4
	}
	am := ams.(backend.AccountManager)

	sms, e5 := a.Service(backend.KeySessionManager)
	if e5 != nil {
		return e5
	}
	sm := sms.(*session.Manager)

	// Retrieve the storage manager.
	sd, e6 := a.Service(backend.KeyStorage)
	if e6 != nil {
		return e6
	}
	store := sd.(storage.Storage)

	account, e7 := GetAccount(am, gp, r, sm.ID(), store)
	if e7 != nil {
		return e7
	}

	// Pull the GPG key from <storage>/secret/<app-name>
	gpgData, e8 := store.Load(backend.PrefixGPGKey + "/" + a.Name())
	if e8 != nil {
		return e8
	}

	var gpgKey = &appKey{}
	if e := json.Unmarshal(gpgData, &gpgKey); e != nil {
		return fmt.Errorf(stderr.DecodeJSON, e.Error())
	}

	// Encrypt the data and store in a secure cookie.
	capsule, e9 := gpg.NewCapsule(gpgKey.PublicKey, gpgKey.PrivateKey, gpgKey.PassPhrase)
	if e9 != nil {
		return e9
	}
	encodeMessage, e10 := capsule.Encrypt(account.ID)
	if e10 != nil {
		return e10
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

type appKey struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	PassPhrase string `json:"pass_phrase"`
}

// GetAccount Lookup an existing account or make a new account only when
// a client has a successful login and an existing account cannot be found.
func GetAccount(
	am backend.AccountManager,
	gp *google.Provider,
	r *http.Request,
	sessionID string,
	store storage.Storage,
) (*backend.Account, error) {
	account, e1 := am.Lookup(gp.ClientID())

	switch e1.(type) {
	// Make a new account only when a client has a successful login and
	// an existing account cannot be found and we, do we make a new account.
	case *backend.AccountNotFoundError:
		Log.Errf(e1.Error())

		uaMeta := r.Header.Get("User-Agent")
		Log.Infof("user-agent: %v", uaMeta)

		device, e2 := backend.NewDevice([]byte(uaMeta), sessionID, backend.KeyGoogleProvider)
		if e2 != nil {
			return nil, e2
		}
		acct, e3 := am.Add(gp.ClientID(), name, device)
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
