package awslambda

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

// ConvertBody string to io.Reader
func ConvertBody(body string, isBase64 bool) (io.ReadCloser, int64) {
	if body == "" {
		return nil, 0
	}

	if isBase64 {
		b, e1 := base64.StdEncoding.DecodeString(body)
		if e1 != nil {
			panic(fmt.Errorf(stderr.DecodeBase64, e1))
		}

		return io.NopCloser(bytes.NewReader(b)), int64(len(string(b)))
	}

	return io.NopCloser(strings.NewReader(body)), int64(len(body))
}

// ConvertHttpCookiesForLambda Convert http.Request.Cookies() to []string
// cookies that work with Lambda functions.
// Returns an empty non-nil slice if there are no cookies in the request.
func ConvertHttpCookiesForLambda(httpCookies []*http.Cookie) []string {
	cookies := make([]string, len(httpCookies))

	if len(httpCookies) == 0 {
		return cookies
	}

	for i, cookie := range httpCookies {
		cookies[i] = cookie.String()
	}

	return cookies
}

// ConvertHttpHeaders the http.Response style headers map[string][]string to map[string]string.
func ConvertHttpHeaders(headers map[string][]string) map[string]string {
	flatHeaders := make(map[string]string, len(headers))
	for k, v := range headers {
		// Use a comma to separate multiple field values for a single field name
		// see https://www.rfc-editor.org/rfc/rfc9110.html#name-field-lines-and-combined-fi
		// However, Set-Cookie contains "," we need to make a special case and store them separately.
		if k == "Set-Cookie" {
			continue // Don't put Cookies in the header of the Lambda response object, instead store them in the Cookies slice.
		}
		// In the spec multiline headers can also just be separated with a comma, this is a problem if the header contains commas.
		// The developer should be aware of this an encode the data in base64 or something else.
		// The REAL problem is that this LambdaFunctionUrlResponse does not handle headers properly with its data structure.
		// It should use the same http.Request header structure so that you can set the same header multiple times as needed.
		flatHeaders[k] = strings.Join(v, ", ")
	}
	return flatHeaders
}

// ConvertToHttpHeaders Convert a map of strings to http.Header's.
func ConvertToHttpHeaders(headers map[string]string, cookies []string) http.Header {
	converted := http.Header{}
	// Just initialize if there are no headers
	if len(headers) == 0 {
		return converted
	}

	// Clone headers over to the http.Header
	for k, v := range headers {
		converted[k] = []string{v}
	}

	// Remember that an HTTP request uses Cookie and response uses Set-Cookie.
	if len(cookies) > 0 {
		Log.Dbugf(stdout.LambdaCookies, cookies)
		converted["Cookie"] = cookies
	}

	return converted
}

// GetHeader Retrieve a header from a request.
func GetHeader(headers map[string]string, name string) string {
	value := ""
	lcn := strings.ToLower(name)

	for h, v := range headers {
		lch := strings.ToLower(h)
		if lch == lcn {
			value = v
			break
		}
	}

	return value
}

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

// ParseCookie takes an HTTP cookie string and converts it to a *http.Cookie.
func ParseCookie(cookie string) (*http.Cookie, error) {
	// Split the HTTP cookie up by attributes.
	parts := strings.Split(cookie, ";")

	// The first part of the cookie must be its name and value.
	pair := strings.Split(parts[0], "=")
	if len(pair) != 2 {
		return nil, fmt.Errorf(stderr.BadCookie, pair)
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

// PrepareResponse Convert to a Lambda function URL response.
func PrepareResponse(res *Output) {
	cookies, ok := res.headers["Set-Cookie"]
	if ok {
		res.Cookies = cookies
	}

	for k, h := range res.headers {
		tmp, ok2 := res.Headers[k]
		sep := ","
		if k == "Set-Cookie" {
			sep = ";"
		}
		if ok2 {
			res.Headers[k] = tmp + sep + strings.Join(h, sep)
			continue
		}
		res.Headers[k] = strings.Join(h, sep)
	}

	if res.IsBase64Encoded {
		tmp := base64.StdEncoding.EncodeToString([]byte(res.Body))
		res.Body = tmp
	}
}
