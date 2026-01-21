package awslambda

import "net/http"

type Output struct {
	// StatusCode Is required, when not set, then Lambda will return Internal Error.
	StatusCode int `json:"statusCode"`
	// Headers are key=value, if you want to set it multiple times, you can use
	// a comma to separate each value, newlines may also work, but has not been
	// tested.
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	// IsBase64Encoded set this to true for binary data for sure.
	IsBase64Encoded bool `json:"isBase64Encoded"`

	// Cookies this is only for AWS Lambda, as it does not use cookies set as a header.
	Cookies []string `json:"cookies"`
}

func (res *Output) Header() http.Header {
	//TODO implement me
	return ConvertToHttpHeaders(res.Headers, res.Cookies)
}

// WriteHeader Part of the http.ResponseWriter interface.
func (res *Output) Write(b []byte) (int, error) {
	res.Body = string(b)
	return len(res.Body), nil
}

// WriteHeader Part of the http.ResponseWriter interface.
func (res *Output) WriteHeader(statusCode int) {
	res.StatusCode = statusCode
}
