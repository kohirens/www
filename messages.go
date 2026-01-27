package www

var Stderr = struct {
	AuthCodeInvalid,
	AuthCodeNotSet,
	AuthHeaderMissing,
	CannotEncodeToJson,
	DecodeBase64,
	FieldNotFound,
	WriteResponseBody string
}{
	AuthCodeInvalid:    "incorrect authorization code was sent",
	AuthCodeNotSet:     "authorization code was not set in the environment",
	AuthHeaderMissing:  "authorization header is missing",
	CannotEncodeToJson: "could not JSON encode content: %v",
	DecodeBase64:       "cannot decode base64 value %v",
	FieldNotFound:      "could not find field %v",
	WriteResponseBody:  "cannot write response body %v",
}
