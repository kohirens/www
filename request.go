package www

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
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
