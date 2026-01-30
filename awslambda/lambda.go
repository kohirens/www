package awslambda

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/kohirens/stdlib/logger"
)

type Handler struct {
	PageSource PageSource
}

type PageSource interface {
	Load(pagePath string) ([]byte, error)
}

const (
	envHttpMethods     = "HTTP_METHODS_ALLOWED"
	envRedirectTo      = "REDIRECT_TO"
	headerAltHost      = "viewer-host"
	headerCfDomain     = "distribution-domain"
	redirectToEnvVar   = "REDIRECT_TO"
	redirectHostEnvVar = "REDIRECT_HOSTS"
)

var (
	Log = logger.Standard{}
)

// GetCookie from a list of cookies.
func GetCookie(cookies []string, name string) string {
	value := ""

	re := regexp.MustCompile(`[^=]+=([^;]+);?.*$`)
	for _, cookie := range cookies {
		if strings.Contains(cookie, name+"=") { // we got a hit
			//
			d := re.FindAllStringSubmatch(cookie, 1)
			if d != nil {
				value = d[0][1]
				break
			}
		}
	}

	return value
}

// NewRequest Work with this type of request as though it were of type http.Request.
func NewRequest(l *Input) (*http.Request, error) {
	origin := GetHeader(l.Headers, "Origin")
	uri := origin + l.RequestContext.HTTP.Path

	if l.RawQueryString != "" {
		uri += "?" + l.RawQueryString
	}

	headers := ConvertToHttpHeaders(l.Headers, l.Cookies)
	method := l.RequestContext.HTTP.Method
	body, _ := ConvertBody(l.Body, l.IsBase64Encoded)
	//body, bodyLength := ConvertBody(l.Body, l.IsBase64Encoded)

	// TODO: Find out why the parseForm does not work with this method.
	r, e2 := http.NewRequest(l.RequestContext.HTTP.Method, uri, body)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.NewRequest, e2)
	}
	// TODO: I assume headers are not properly set, so parse form does not know that it should parse.
	r.Header = headers

	//u, e1 := url.Parse(uri)
	//if e1 != nil {
	//	return nil, e1
	//}

	//r := &http.Request{
	//	Method:        method,
	//	Proto:         l.RequestContext.HTTP.Protocol,
	//	Body:          body,
	//	ContentLength: bodyLength,
	//	Host:          GetHeader(l.Headers, "Host"),
	//	Header:        headers,
	//	URL:           u,
	//}

	if method == "POST" || method == "PUT" {
		b := l.Body
		if l.IsBase64Encoded {
			tmp, _ := base64.StdEncoding.DecodeString(l.Body)
			b = string(tmp)
		}
		formData, e0 := url.ParseQuery(b)
		if e0 != nil {
			return nil, e0
		}

		r.Form = formData
		r.PostForm = formData
	}

	return r, nil
}

func NewResponse() *Output {
	return &Output{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "text/html; charset=utf-8",
		},
		Body:            "",
		IsBase64Encoded: false,
		Cookies:         make([]string, 0),
	}
}

func PreliminaryChecks(event *Input) *Output {
	method := event.RequestContext.HTTP.Method
	httpAllowedMethods, ok := os.LookupEnv(envHttpMethods)

	lResponse := NewResponse()
	if !ok {
		Log.Errf(stderr.MissingEnv, envHttpMethods)
		lResponse.StatusCode = http.StatusInternalServerError
		return lResponse
	}

	if strings.ToUpper(method) == "OPTIONS" {
		lResponse.StatusCode = 204
		lResponse.Headers["Allow"] = httpAllowedMethods
		return lResponse
	}

	supportedMethods := strings.Split(httpAllowedMethods, ",")
	if NotImplemented(method, supportedMethods) {
		lResponse.StatusCode = http.StatusNotImplemented
		return lResponse
	}

	host := GetHeader(event.Headers, headerAltHost)

	doIt, e1 := ShouldRedirect(host)
	if e1 != nil {
		Log.Errf("%v", e1.Error())
		lResponse.StatusCode = http.StatusInternalServerError
		return lResponse
	}

	if doIt {
		serverHost, _ := os.LookupEnv(envRedirectTo)
		if !strings.Contains(serverHost, "https://") {
			serverHost = "https://" + serverHost
		}
		switch method {
		case "POST":
			lResponse.StatusCode = http.StatusPermanentRedirect
			return lResponse
		}
		lResponse.StatusCode = http.StatusMovedPermanently
		return lResponse
	}

	distributionDomain := GetHeader(event.Headers, headerCfDomain)

	Log.Infof(stdout.DistDomain, distributionDomain)

	if host == distributionDomain {
		Log.Errf(stderr.DistroRequest, distributionDomain)
		lResponse.StatusCode = http.StatusUnauthorized
		return lResponse
	}

	Log.Infof("%v", stdout.PreChecks)

	return nil
}

// ShouldRedirect Perform a redirect if the host matches any of the domains in
// the REDIRECT_HOST environment variable.
func ShouldRedirect(host string) (bool, error) {
	if host == "" {
		return false, fmt.Errorf("%v", stderr.HostNotSet)
	}

	rt, ok1 := os.LookupEnv(redirectToEnvVar)
	if !ok1 {
		return false, fmt.Errorf(stderr.EnvVarUnset, redirectToEnvVar)
	}

	if rt == "" {
		return false, fmt.Errorf(stderr.RedirectToEmpty, redirectToEnvVar)
	}

	if strings.EqualFold(host, rt) {
		return false, nil
	}

	rh, ok2 := os.LookupEnv(redirectHostEnvVar)
	if !ok2 {
		return false, fmt.Errorf(stderr.EnvVarUnset, redirectHostEnvVar)
	}

	if rh == "" {
		return false, nil
	}

	retVal := false
	rhs := strings.Split(rh, ",")
	for _, h := range rhs {
		if h == host {
			retVal = true
		}
	}

	return retVal, nil
}
