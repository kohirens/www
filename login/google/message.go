package google

var stderr = struct {
	DecodeJSON,
	ECCookie,
	EncodeJSON,
	EndpointNotFound,
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
	LoginRegistration: "cannot complete login registration: %v",
	ParseSignInData:   "could not parse login form data: %v",
	SignOut:           "Sign Out: %v",
	ValidEmail:        "email failed validation: %v",
	WriteResponseBody: "could not write response body: %v",
}

var stdout = struct {
	GoogleCallback string
}{
	GoogleCallback: "Google is calling back",
}
