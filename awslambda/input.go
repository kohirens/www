package awslambda

import (
	"fmt"
	"net/http"
	"net/textproto"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// ParseCookies found in the event object.
// Remember that HTTP request use Cookie and response uses Set-Cookie.
func (r *Input) ParseCookies() error {
	Log.Dbugf("%v", stdout.ParseCookies)

	for _, cookie := range r.Cookies {
		c, e := ParseCookie(cookie)
		if e != nil {
			return e
		}

		r.cookies[c.Name] = c
	}

	return nil
}

func ParseCookie(cookie string) (*http.Cookie, error) {
	re := regexp.MustCompile(`;\s`)
	d := re.FindAllStringSubmatch(cookie, 1)
	fmt.Printf("parts regex: %v\n", d)

	parts := strings.Split(cookie, ";")

	pair := strings.Split(parts[0], "=")
	if len(pair) != 2 {
		return nil, fmt.Errorf("not a valid cookie %v", pair)
	}

	c := &http.Cookie{Name: pair[0], Value: pair[1]}
	for _, part := range parts[1:] {
		p := strings.Split(part, "=")
		switch strings.ToLower(textproto.TrimString(p[0])) {
		case "expires":
			t, e := time.Parse("Mon, 02-Jan-2006 15:04:05 MST", p[1])
			if e != nil {
				return nil, e
			}
			c.Expires = t
		case "secure":
			c.Secure = true
		case "domain":
			c.Domain = p[1]
		case "path":
			c.Path = p[1]
		case "samesite":
			switch p[1] {
			case "lax":
				c.SameSite = http.SameSiteLaxMode
			case "strict":
				c.SameSite = http.SameSiteStrictMode
			case "none":
				c.SameSite = http.SameSiteNoneMode
			default:
				c.SameSite = http.SameSiteDefaultMode
			}
		case "httponly":
			c.HttpOnly = true
		case "max-age":
			i, e := strconv.Atoi(p[1])
			if e != nil {
				return nil, e
			}
			c.MaxAge = i
		}
	}
	return c, nil
}
