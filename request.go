package www

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Request Serves as a medium between the different types of HTTP request that
// you may run into. Initially this works with Go's HTTP request and AWS Lambda
// function URL request. More will be added as encountered; or until a better
// solution is found.
// The main goal is to always use http.Request no matter what type of Request
// you are given, so that you code across projects is more consistent and highly
// reusable.
// This is just a wrapper around http.Request, as that is used under the hood.
// When you need another type use any of the Request.To* methods to convert to
// another type that is supported.
// While the main goal it not so server as a way to convert one type to another
// it does work out that way.
type Request struct {
	Request                  *http.Request
	LambdaFunctionURLRequest *events.LambdaFunctionURLRequest
}

// NewRequest Wrap an http.Request.
func NewRequest(r *http.Request) *Request {
	return &Request{
		Request: r,
	}
}

// NewRequestFromLambdaFunctionURLRequest Work with this type of request as though it were of type http.Request.
func NewRequestFromLambdaFunctionURLRequest(l *events.LambdaFunctionURLRequest) (*Request, error) {
	origin := GetHeader(l.Headers, "Origin")
	uri := origin + l.RawPath

	if l.RawQueryString != "" {
		uri += "?" + l.RawQueryString
	}

	//// TODO: Find out why the parseForm does not work with this method.
	//// WARNING: This is a time sink-hole.
	//body := convertBody(l.Body, l.IsBase64Encoded)
	//r, e2 := http.NewRequest(l.RequestContext.HTTP.Method, uri, body)
	//if e2 != nil {
	//	return nil, fmt.Errorf(Stderr.NewRequest, e2)
	//}

	headers := convertToHttpHeaders(l)
	method := l.RequestContext.HTTP.Method
	body, _ := base64.StdEncoding.DecodeString(l.Body)

	u, e1 := url.ParseRequestURI(uri)
	if e1 != nil {
		return nil, e1
	}

	r := &http.Request{
		Method:        method,
		Proto:         l.RequestContext.HTTP.Protocol,
		Body:          io.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(string(body))),
		Host:          GetHeader(l.Headers, "Host"),
		Header:        headers,
		URL:           u,
	}

	if method == "POST" || method == "PUT" {
		formData, e0 := url.ParseQuery(string(body))
		if e0 != nil {
			return nil, e0
		}

		r.Form = formData
		r.PostForm = formData
	}

	r.Header = headers

	return &Request{
		Request:                  r,
		LambdaFunctionURLRequest: l,
	}, nil
}

// ToLambdaFunctionURLRequest Get the Lambda function URL request you put in or a new one with properties set.
func (r *Request) ToLambdaFunctionURLRequest() *events.LambdaFunctionURLRequest {
	if r.LambdaFunctionURLRequest != nil {
		return r.LambdaFunctionURLRequest
	}

	l := &events.LambdaFunctionURLRequest{
		Cookies: convertCookiesToStringArray(r.Request.Cookies()),
	}

	return l
}

// getBody Convert string body to io.Reader
func convertBody(body string, isBase64 bool) io.Reader {
	if body == "" {
		return nil
	}

	if isBase64 {
		b, e1 := base64.StdEncoding.DecodeString(body)
		if e1 != nil {
			panic(fmt.Errorf(Stderr.DecodeBase64, e1))
		}

		return bytes.NewReader(b)
	}

	return strings.NewReader(body)
}

// convertCookiesToStringArray Convert http.Request.Cookies() to []string
// cookies.
// Returns an empty non-nil slice if there are no cookies in the request.
func convertCookiesToStringArray(rcs []*http.Cookie) []string {
	cookies := make([]string, len(rcs))

	if len(rcs) == 0 {
		return cookies
	}

	for i, cookie := range rcs {
		cookies[i] = cookie.String()
	}

	return cookies
}

// convertToHttpHeaders Convert a map of strings to http.Header's.
func convertToHttpHeaders(l *events.LambdaFunctionURLRequest) http.Header {
	headers := http.Header{}
	// Just initialize if there are no headers
	if len(l.Headers) == 0 {
		return headers
	}

	// Remember that HTTP request use Cookie and Response uses Set-Cookie.
	cookieHeader := "Cookie"

	// Clone headers over to the http.Header
	for k, v := range l.Headers {
		headers[k] = []string{v}
	}

	// Lambda stores cookies in a separate array, so make sure to grab them.
	if len(l.Cookies) > 0 {
		headers[cookieHeader] = l.Cookies
	}

	return headers
}

// Wrappers Methods that simply wrap the http.Request, nothing special below this line.

// AddCookie Wraps http.Request.AddCookie()
func (r *Request) AddCookie(c *http.Cookie) {
	r.Request.AddCookie(c)
}

// Cookie Wraps http.Request.Cookie()
func (r *Request) Cookie(name string) (*http.Cookie, error) {
	return r.Request.Cookie(name)
}

// Cookies Wraps http.Request.Cookies()
func (r *Request) Cookies() []*http.Cookie {
	return r.Request.Cookies()
}

// CookiesNamed Wraps http.Request.CookiesNamed()
func (r *Request) CookiesNamed(name string) []*http.Cookie {
	return r.Request.CookiesNamed(name)
}

// ParseForm Wraps http.Request.ParseForm()
func (r *Request) ParseForm() error {
	return r.Request.ParseForm()
}

//TODO: Add the rest of the wrappers
