package awslambda

import (
	"fmt"
	"net/http"
)

type Authorizer struct {
	Iam struct {
		AccessKey       string      `json:"accessKey"`
		AccountId       string      `json:"accountId"`
		CallerId        string      `json:"callerId"`
		CognitoIdentity interface{} `json:"cognitoIdentity"`
		PrincipalOrgId  interface{} `json:"principalOrgId"`
		UserArn         string      `json:"userArn"`
		UserId          string      `json:"userId"`
	} `json:"iam"`
}

type Context struct {
	AccountId    string      `json:"accountId"`
	ApiId        string      `json:"apiId"`
	Authorizer   *Authorizer `json:"authorizer"`
	DomainName   string      `json:"domainName"`
	DomainPrefix string      `json:"domainPrefix"`
	HTTP         *Http       `json:"http"`
	RequestId    string      `json:"requestId"`
	RouteKey     string      `json:"routeKey"`
	Stage        string      `json:"stage"`
	Time         string      `json:"time"`
	TimeEpoch    int64       `json:"timeEpoch"`
}

type Http struct {
	Method    string `json:"method"`
	Path      string `json:"path"`
	Protocol  string `json:"protocol"`
	SourceIp  string `json:"sourceIp"`
	UserAgent string `json:"userAgent"`
}

type Input struct {
	Version               string                  `json:"version"` // Version is expected to be `"2.0"`
	RawPath               string                  `json:"rawPath"`
	RawQueryString        string                  `json:"rawQueryString"`
	RouteKey              string                  `json:"routeKey"`
	Cookies               []string                `json:"cookies,omitempty"`
	cookies               map[string]*http.Cookie `json:"-"`
	Headers               map[string]string       `json:"headers"`
	QueryStringParameters map[string]string       `json:"queryStringParameters,omitempty"`
	Body                  string                  `json:"body,omitempty"`
	IsBase64Encoded       bool                    `json:"isBase64Encoded"`
	RequestContext        *Context                `json:"requestContext"`
}

// Cookie returns an HTTP cookie if found.
// Remember that HTTP request use Cookie and response uses Set-Cookie.
func (r *Input) Cookie(name string) (*http.Cookie, error) {
	c, ok := r.cookies[name]
	if !ok {
		return nil, fmt.Errorf(stderr.CookieNotFound, name)
	}
	return c, nil
}

// ParseCookies takes the AWS Lambda event object cookies and turns them
// into []*httpCookies. This makes them easier to work with in your code.
// This simply uses the awslambda.ParseCookie, which actually parses each
// cookie.
func (r *Input) ParseCookies() error {
	Log.Dbugf("%v", stdout.ParseCookies)
	if r.cookies == nil {
		r.cookies = make(map[string]*http.Cookie)
	}

	for _, cookie := range r.Cookies {
		c, e := ParseCookie(cookie)
		if e != nil {
			return e
		}

		r.cookies[c.Name] = c
	}

	return nil
}
