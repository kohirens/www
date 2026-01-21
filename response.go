package www

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// Response Serves as a middle ground between types such as http.Response and events.LambdaFunctionURLResponse;
// with methods to easily convert to each type. In addition, it also works as a http.ResponseWriter.
type Response struct {
	Body            string              `json:"body"`
	Headers         map[string][]string `json:"headers"`
	IsBase64Encoded bool                `json:"isBase64Encoded"`
	StatusCode      int                 `json:"statusCode"`
}

var _ http.ResponseWriter = &Response{}

// Base64Encode Encodes the response body.
func (res *Response) Base64Encode() {
	res.IsBase64Encoded = true
	res.Body = base64.StdEncoding.EncodeToString([]byte(res.Body))
}

// Cookies Get all the cookies.
func (res *Response) Cookies() []string {
	return res.Headers["Set-Cookie"]
}

// Header Part of the http.ResponseWriter interface.
func (res *Response) Header() http.Header {
	return res.Headers
}

// ToHttpResponse Convert to an HTTP response.
func (res *Response) ToHttpResponse() *http.Response {
	r := &http.Response{
		StatusCode: res.StatusCode,
		Header:     res.Headers,
	}

	b := bytes.NewBufferString(res.Body)
	if e := r.Write(b); e != nil {
		panic(e)
	}

	return r
}

func (res *Response) ToJSON() (string, error) {
	data, e1 := json.Marshal(res)
	if e1 != nil {
		return "", fmt.Errorf(Stderr.CannotEncodeToJson, e1.Error())
	}
	return string(data), nil
}

// WriteHeader Part of the http.ResponseWriter interface.
func (res *Response) Write(b []byte) (int, error) {
	res.Body = string(b)
	return len(res.Body), nil
}

// WriteHeader Part of the http.ResponseWriter interface.
func (res *Response) WriteHeader(statusCode int) {
	res.StatusCode = statusCode
}
