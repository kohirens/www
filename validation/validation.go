package validation

import (
	"regexp"
	"strings"
)

// Email Validate an email address. While this is not fully IETF RFC compliant,
// it should allow most emails without allowing anything too malicious to pass
// through.
func Email(email string) (string, bool) {
	email = strings.TrimSpace(email)
	re := regexp.MustCompile(`[a-z-A-Z-!#$%&'*+-/=?^_{|}~]+@[a-z-A-Z0-9-]+\.[a-z-A-Z0-9-]+`)

	if !re.MatchString(email) {
		return email, false
	}

	parts := strings.Split(email, "@")
	domain := parts[1]
	labels := strings.Split(domain, ".")
	for _, l := range labels {
		if len(l) > 63 {
			return email, false
		}
	}

	return email, true
}

// MaxLen Subject length does not exceed the max length.
func MaxLen(subject string, max int) bool {
	return len(subject) <= max
}
