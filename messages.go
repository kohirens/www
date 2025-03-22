package www

var Stdout = struct {
	BytesRead            string
	ConnectTo            string
	DomainOnRedirectList string
	EnvVarEmpty          string
	LoadPage             string
	RunCli               string
}{
	BytesRead:            "number of bytes read from %v is %d",
	ConnectTo:            "connecting to %v",
	DomainOnRedirectList: "domain %v in in the list of domains to redirect to %v",
	EnvVarEmpty:          "environment variable %v is empty",
	LoadPage:             "loading the %v page",
	RunCli:               "Running CLI",
}

var Stderr = struct {
	AuthCodeInvalid     string
	AuthCodeNotSet      string
	AuthHeaderMissing   string
	CannotCloseFile     string
	CannotEncodeToJson  string
	CannotGetExt        string
	CannotLoadPage      string
	CannotOpenFile      string
	CannotParseFile     string
	CannotReadFile      string
	DecodeBase64        string
	EnvVarUnset         string
	FatalHeader         string
	FieldNotFound       string
	FileHasNoContent    string
	FileNotFound        string
	HostNotSet          string
	DoNotRedirectToSelf string
	NewRequest          string
	RedirectToEmpty     string
}{
	AuthCodeInvalid:     "incorrect authorization code was sent",
	AuthCodeNotSet:      "authorization code was not set in the environment",
	AuthHeaderMissing:   "authorization header is missing",
	CannotCloseFile:     "could not close file: %v",
	CannotEncodeToJson:  "could not JSON encode content: %v",
	CannotGetExt:        "could not close file: %v",
	CannotLoadPage:      "could not load the page %v: %v",
	CannotOpenFile:      "could not open file %v: %v",
	CannotParseFile:     "could not parse XSD: %v",
	CannotReadFile:      "could not read file %v: %v",
	DecodeBase64:        "could not decode base64: %v",
	EnvVarUnset:         "environment variable %v has not been set",
	FatalHeader:         "fatal error detected: %v",
	FieldNotFound:       "could not find field %v",
	FileHasNoContent:    "the field %v points to a file %v that contains no content (it is empty)",
	FileNotFound:        "could not find file %v",
	HostNotSet:          "could not retrieve the host from the request",
	DoNotRedirectToSelf: "will not redirect %v to host %v",
	NewRequest:          "could not initialize a new request: %v",
	RedirectToEmpty:     "the REDIRECT_TO environment variables was empty",
}
