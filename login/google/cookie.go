package login

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EncryptedCookie struct {
	AID string
	DID string
	UA  string
}

func GetEncryptedCookie() *http.Cookie {

}

const CookieNameAccount = "_ed_"

func PutEncryptedCookie() *http.Cookie {

	ec := &EncryptedCookie{
		AID: account.ID,
		DID: gp.DeviceID(),
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
		Name:   CookieNameAccount,
		Value:  string(encodeMessage),
		Path:   "/",
		Secure: true,
	})
}
