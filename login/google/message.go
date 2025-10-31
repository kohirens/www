package google

var stderr = struct {
	DecodeJSON,
	ECCookie,
	EncodeJSON,
	EndpointNotFound,
	GetAccount,
	LookupLoginInfo,
	LoginRegistration,
	ParseSignInData,
	SignOut,
	ValidEmail,
	WriteResponseBody string
}{
	DecodeJSON:        "cannot decode json: %v",
	ECCookie:          "cannot get the encrypted cookie: %v",
	EncodeJSON:        "cannot encode json: %v",
	EndpointNotFound:  "no api endpoint %v can be found",
	GetAccount:        "failed to get the account",
	LookupLoginInfo:   "failed to get the login information",
	LoginRegistration: "cannot complete login registration: %v",
	ParseSignInData:   "could not parse login form data: %v",
	SignOut:           "Sign Out: %v",
	ValidEmail:        "email failed validation: %v",
	WriteResponseBody: "could not write response body: %v",
}

var stdout = struct {
	AccountID,
	DeviceID,
	EncryptedCookie,
	EncryptedCookieValue,
	GoogleCallback,
	LookupAccount,
	LookupLoginInfo,
	MakeAccount,
	MakeLoginInfo,
	NewAccount,
	RegisterAccount,
	SessionID,
	UpdateLoginInfo,
	UserAgent string
}{
	AccountID:            "account ID: %v",
	DeviceID:             "device ID: %v",
	EncryptedCookie:      "looking for an encrypted cookie...",
	EncryptedCookieValue: "setting encrypted value cookie",
	GoogleCallback:       "Google is calling back",
	LookupAccount:        "lookup account...%v",
	LookupLoginInfo:      "lookup login information...",
	MakeAccount:          "making a new account",
	MakeLoginInfo:        "making %v login info",
	NewAccount:           "registered a new account %v",
	RegisterAccount:      "register new account",
	SessionID:            "session ID: %v",
	UpdateLoginInfo:      "update login information",
	UserAgent:            "user-agent: %v",
}
