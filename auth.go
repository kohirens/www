package www

import (
	"fmt"
	"os"
)

const (
	authHeader = "Authorization"
)

// Authenticate Challenge user access.
//
//	Expected format = "Basic " + base64("username:password")
//	In Go  for example example:
//	  auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))
func Authenticate(headers map[string]string) error {
	headerVal := GetHeader(headers, authHeader)
	if headerVal == "" {
		return fmt.Errorf("%v", Stderr.AuthHeaderMissing)
	}

	envVal, ok := os.LookupEnv(authHeader)
	if !ok {
		return fmt.Errorf("%v", Stderr.AuthCodeNotSet)
	}

	if headerVal != envVal {
		return fmt.Errorf("%v", Stderr.AuthCodeInvalid)
	}

	return nil
}
