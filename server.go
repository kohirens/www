package www

import (
	"fmt"
	"os"
	"strings"
)

const (
	RedirectToEnvVar   = "REDIRECT_TO"
	RedirectHostEnvVar = "REDIRECT_HOSTS"
)

// NotImplemented Return true if the HTTP method is supported by this server
// and false otherwise.
func NotImplemented(method string, supported []string) bool {
	missing := true
	for _, sm := range supported {
		if strings.EqualFold(sm, method) {
			missing = false
		}
	}
	return missing
}

// ResponseFromFile Load a page from storage.
func ResponseFromFile(pagePath, contentType string) (*Response, error) {
	content, e1 := os.ReadFile(pagePath)
	if e1 != nil {
		return nil, fmt.Errorf(Stderr.CannotLoadPage, e1.Error())
	}

	if content == nil {
		return Respond404(), nil
	}

	res := Respond200(content, contentType)

	return res, nil
}

// ShouldRedirect Perform a redirect if the host matches any of the domains in
// the REDIRECT_HOST environment variable.
func ShouldRedirect(host string) (bool, error) {
	if host == "" {
		return false, fmt.Errorf(Stderr.HostNotSet)
	}

	rt, ok1 := os.LookupEnv(RedirectToEnvVar)
	if !ok1 {
		return false, fmt.Errorf(Stderr.EnvVarUnset, RedirectToEnvVar)
	}

	if rt == "" {
		return false, fmt.Errorf(Stderr.RedirectToEmpty, RedirectToEnvVar)
	}

	if strings.EqualFold(host, rt) {
		return false, nil
	}

	rh, ok2 := os.LookupEnv(RedirectHostEnvVar)
	if !ok2 {
		return false, fmt.Errorf(Stderr.EnvVarUnset, RedirectHostEnvVar)
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
