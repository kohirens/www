package awslambda

var stderr = struct {
	DistroRequest,
	DecodeBase64,
	EnvVarUnset,
	HostNotSet,
	NoSuchKey,
	MissingEnv,
	RedirectToEmpty string
}{
	DistroRequest:   "a request was made using the CloudFront distribution domain name, which is not authorized: %v",
	DecodeBase64:    "could not decode base64: %v",
	EnvVarUnset:     "environment variable %v has not been set",
	HostNotSet:      "could not retrieve the host from the request",
	MissingEnv:      "environment variable %v is not set",
	NoSuchKey:       "no such key %v",
	RedirectToEmpty: "the REDIRECT_TO environment variables was empty",
}

var stdout = struct {
	DistDomain,
	PreChecks string
}{
	DistDomain: "distribution domain = %v",
	PreChecks:  "preliminary checks have completed",
}
