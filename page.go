package www

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var FooterText = "&copy; " + time.Now().Format("2006")

// GetPageType Get the content type via header.
func GetPageType(headers StringMap) string {
	ct := GetHeader(headers, "content-type")

	fct := strings.Split(ct, ",")
	if fct != nil {
		ct = fct[0]
	}

	return ct
}

// GetPageTypeByExt Get the content type by the extension of the file being
// requested.
func GetPageTypeByExt(pagePath string) string {
	var ct string

	ext := filepath.Ext(pagePath)

	switch ext {
	case ".css":
		ct = ContentTypeCSS
	case ".html":
		ct = ContentTypeHtml
	case ".js":
		ct = ContentTypeJS
	case ".json":
		ct = ContentTypeJson
	case ".jpg":
		ct = ContentTypeJpg
	case ".gif":
		ct = ContentTypeGif
	case ".png":
		ct = ContentTypePng
	case ".svg", ".svgz":
		ct = ContentTypeSvg
	default:
		ct = ""
	}

	return ct
}

// GetHeader Retrieve a header from a request.
func GetHeader(headers StringMap, name string) string {
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

// GetMapItem Retrieve an item from a string map.
func GetMapItem(mapData StringMap, name string) string {
	value := ""
	ln := strings.ToLower(name)

	for k, v := range mapData {
		lk := strings.ToLower(k)
		if lk == ln {
			value = v
			break
		}
	}

	return value
}

// Respond200 Send an OK HTTP response.
func Respond200(content []byte, contentType string) *Response {
	code := http.StatusOK
	res := &Response{
		Headers: MapOfStrings{
			"Content-Type": {contentType},
		},
		StatusCode: code,
	}

	switch contentType {
	case ContentTypeGif, ContentTypeJpg, ContentTypePng:
		res.Body = base64.StdEncoding.EncodeToString(content)
		res.IsBase64Encoded = true
	default:
		res.Body = string(content)
	}

	return res
}

// Respond201 Send a Created HTTP response.
func Respond201(location string) *Response {
	return RespondWithLocation(location, http.StatusCreated)
}

// Respond301 Send a "Moved Permanently" HTTP response.
func Respond301(location string) *Response {
	return RespondWithLocation(location, http.StatusMovedPermanently)
}

// Respond302 Send a Found HTTP response0.
func Respond302(location string) *Response {
	return RespondWithLocation(location, http.StatusFound)
}

// Respond308 Send an 308 HTTP response redirect to another location (full URL).
func Respond308(location string) *Response {
	return RespondWithLocation(location, http.StatusPermanentRedirect)
}

// Respond401 Send a 401 Unauthorized HTTP response.
func Respond401() *Response {
	return RespondWithStatus(http.StatusUnauthorized)
}

// Respond404 Send a 404 Not Found HTTP response.
func Respond404() *Response {
	return RespondWithStatus(http.StatusNotFound)
}

// Respond405 Send a 405 Method Not Allowed HTTP response.
//
//	allowedMethods is a comma-delimited string of HTTP methods that are allowed.
//	Example:
//	  GET, HEAD, PUT
func Respond405(allowedMethods string) *Response {
	return RespondWithStatus(http.StatusMethodNotAllowed)
}

// Respond500 Send a 500 Internal Server Error HTTP response.
func Respond500() *Response {
	return RespondWithStatus(http.StatusInternalServerError)
}

// Respond501 Send a 501 Not Implemented HTTP response.
//
//	501 is the appropriate response when the server does not recognize the
//	request method and is incapable of supporting it for any resource. The only
//	methods that servers are required to support (and therefore that must not
//	return 501) are GET and HEAD.
func Respond501() *Response {
	return RespondWithStatus(http.StatusNotImplemented)
}

// RespondDebug Respond with a debug message and whatever code your like.
//
//	This was handy when testing AWS Lambda function or initial set up of the
//	Lambda URL feature.
func RespondDebug(code int, message, footer string) *Response {
	return &Response{
		Body: fmt.Sprintf(Http200Debug, code, message, footer),
		Headers: MapOfStrings{
			"Content-Type": {ContentTypeHtml},
		},
		StatusCode: code,
	}
}

// RespondWithJSON Send a JSON HTTP response.
func RespondWithJSON(content interface{}) (*Response, error) {
	jsonEncodedContent, e1 := json.Marshal(content)
	if e1 != nil {
		return nil, fmt.Errorf(Stderr.CannotEncodeToJson, e1.Error())
	}

	return Respond200(jsonEncodedContent, ContentTypeJson), nil
}

// ResponseOptions Respond with an HTTP Allow header listing all HTTP methods
// allowed for a request.
func ResponseOptions(options string) *Response {
	return &Response{
		Body: "",
		Headers: MapOfStrings{
			"Allow": {options},
		},
		StatusCode: 204,
	}
}

// RespondWithLocation Send a status to go to another location HTTP response0.
func RespondWithLocation(location string, code int) *Response {
	return &Response{
		Body: fmt.Sprintf(HttpStatusContent, code, http.StatusText(code)+"<br />"+location, FooterText),
		Headers: map[string][]string{
			"Content-Type": {ContentTypeHtml},
			"Location":     {location},
		},
		StatusCode: code,
	}
}

// RespondWithStatus Send a status HTTP response.
//
//	See: https://www.rfc-editor.org/rfc/rfc9110
//	Also see https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
func RespondWithStatus(code int) *Response {
	return &Response{
		Body: fmt.Sprintf(HttpStatusContent, code, http.StatusText(code), FooterText),
		Headers: MapOfStrings{
			"Content-Type": {ContentTypeHtml},
		},
		StatusCode: code,
	}
}
