package google

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/www/backend"
	"net/http"
)

type EncryptedCookie struct {
	AID string
	DID string
	UA  string
}

const ecCookieName = "_ec_"

func GetEncryptedCookie(r *http.Request, a backend.App) (*EncryptedCookie, error) {
	cookie, e1 := r.Cookie(ecCookieName)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.ECCookie, e1.Error())
	}

	message, e2 := a.Decrypt([]byte(cookie.Value))
	if e2 != nil {
		return nil, e2
	}

	ec := &EncryptedCookie{}
	if err := json.Unmarshal(message, ec); err != nil {
		return nil, fmt.Errorf(stderr.DecodeJSON, err.Error())
	}

	return ec, nil
}

func PutEncryptedCookie(
	accountID,
	deviceID,
	userAgent string,
	w http.ResponseWriter,
	a backend.App,
) error {
	ec := &EncryptedCookie{
		AID: accountID,
		DID: deviceID,
		UA:  userAgent,
	}

	ecBytes, e15 := json.Marshal(ec)
	if e15 != nil {
		return fmt.Errorf(stderr.EncodeJSON, e15.Error())
	}

	encodeMessage, e10 := a.Encrypt(ecBytes)
	if e10 != nil {
		return e10
	}

	http.SetCookie(w, &http.Cookie{
		Name:   ecCookieName,
		Value:  string(encodeMessage),
		Path:   "/",
		Secure: true,
	})
	return nil
}
