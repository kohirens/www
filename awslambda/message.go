package awslambda

var stderr = struct {
	DistroRequest,
	NoSuchKey,
	MissingEnv string
}{
	DistroRequest: "a request was made using the CloudFront distribution domain name, which is not authorized: %v",
	MissingEnv:    "environment variable %v is not set",
	NoSuchKey:     "no such key %v",
}

var stdout = struct {
	DistDomain,
	PreChecks string
}{
	DistDomain: "distribution domain = %v",
	PreChecks:  "preliminary checks have completed",
}
