package awslambda

var stderr = struct {
	BadCookie,
	CookieNotFound,
	DistroRequest,
	DecodeBase64,
	EnvVarUnset,
	HostNotSet,
	NewRequest,
	NoSuchKey,
	MissingEnv,
	RedirectToEmpty string
}{
	BadCookie:       "not a valid cookie %v",
	CookieNotFound:  "Cookie %v not found",
	DistroRequest:   "a request was made using the CloudFront distribution domain name, which is not authorized: %v",
	DecodeBase64:    "could not decode base64: %v",
	EnvVarUnset:     "environment variable %v has not been set",
	HostNotSet:      "could not retrieve the host from the request",
	MissingEnv:      "environment variable %v is not set",
	NewRequest:      "cannot init a new http.Request",
	NoSuchKey:       "no such key %v",
	RedirectToEmpty: "the REDIRECT_TO environment variables was empty",
}

var stdout = struct {
	DistDomain,
	LambdaCookies,
	ParseCookies,
	PreChecks string
}{
	DistDomain:    "distribution domain = %v",
	LambdaCookies: "AWS Lambda Cookies: %v",
	ParseCookies:  "parsing cookies...",
	PreChecks:     "preliminary checks have completed",
}
