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

var (
	FooterText = "&copy; " + time.Now().Format("2006")
)

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
func Respond200(w http.ResponseWriter, content []byte, contentType string) {
	var body []byte
	switch contentType {
	case ContentTypeGif, ContentTypeJpg, ContentTypePng:
		base64.StdEncoding.Encode(content, body)
	default:
		body = content
	}

	_, e1 := w.Write(body)
	if e1 != nil {
		panic(fmt.Errorf(Stderr.WriteResponseBody, e1.Error()))
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
}

// Respond201 Send a Created HTTP response.
func Respond201(w http.ResponseWriter, location string) {
	RespondWithLocation(w, location, http.StatusCreated)
}

// Respond301 Send a "Moved Permanently" HTTP response.
func Respond301(w http.ResponseWriter, location string) {
	RespondWithLocation(w, location, http.StatusMovedPermanently)
}

// Respond302 Send a Found HTTP response0.
func Respond302(w http.ResponseWriter, location string) {
	RespondWithLocation(w, location, http.StatusFound)
}

// Respond308 Send an 308 HTTP response redirect to another location (full URL).
func Respond308(w http.ResponseWriter, location string) {
	RespondWithLocation(w, location, http.StatusPermanentRedirect)
}

// Respond401 Send a 401 Unauthorized HTTP response.
func Respond401(w http.ResponseWriter, body []byte, contentType string) {

	RespondWithStatus(w, http.StatusUnauthorized, body, contentType)
}

// Respond404 Send a 404 Not Found HTTP response.
func Respond404(w http.ResponseWriter, body []byte, contentType string) {
	RespondWithStatus(w, http.StatusNotFound, body, contentType)
}

// Respond405 Send a 405 Method Not Allowed HTTP response.
//
//	allowedMethods is a comma-delimited string of HTTP methods that are allowed.
//	Example:
//	Allow: GET, HEAD, PUT
//
//	Also see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/405
func Respond405(w http.ResponseWriter, allowedMethods string) {
	w.Header().Set("Allow", allowedMethods)
	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// Respond500 Send a 500 Internal Server Error HTTP response.
func Respond500(w http.ResponseWriter, body []byte, contentType string) {
	RespondWithStatus(w, http.StatusInternalServerError, body, contentType)
}

// Respond501 Send a 501 Not Implemented HTTP response.
//
//	501 is the appropriate response when the server does not recognize the
//	request method and is incapable of supporting it for any resource. The only
//	methods that servers are required to support (and therefore that must not
//	return 501) are GET and HEAD.
func Respond501(w http.ResponseWriter, body []byte, contentType string) {
	RespondWithStatus(w, http.StatusNotImplemented, body, contentType)
}

// RespondDebug Respond with a debug message and whatever code your like.
//
//	This was handy when testing AWS Lambda function or initial set up of the
//	Lambda URL feature.
func RespondDebug(w http.ResponseWriter, code int, message, footer string) {
	body := fmt.Sprintf(Http200Debug, code, message, footer)
	_, e1 := w.Write([]byte(body))
	if e1 != nil {
		panic(fmt.Errorf(Stderr.WriteResponseBody, e1.Error()))
	}

	w.Header().Set("Content-Type", ContentTypeHtml)
	w.WriteHeader(code)
}

// RespondWithJSON Send a JSON HTTP response.
func RespondWithJSON(w http.ResponseWriter, content any) {
	jsonEncodedContent, e1 := json.Marshal(content)
	if e1 != nil {
		panic(fmt.Sprintf(Stderr.CannotEncodeToJson, e1.Error()))
	}

	Respond200(w, jsonEncodedContent, ContentTypeJson)
}

// RespondWithLocation Send a status to go to another location HTTP response0.
func RespondWithLocation(w http.ResponseWriter, location string, code int) {
	writeBody(w, code, http.StatusText(code)+"<br />"+location)
	w.Header().Set("Content-Type", ContentTypeHtml)
	w.Header().Set("Location", location)
	w.WriteHeader(code)
}

// RespondWithStatus Send a status HTTP response.
//
//	See: https://www.rfc-editor.org/rfc/rfc9110
//	Also see https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
func RespondWithStatus(w http.ResponseWriter, code int, body []byte, contentType string) {
	writeBody(w, code, http.StatusText(code))
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
}

func writeBody(w http.ResponseWriter, code int, content ...string) {
	body := fmt.Sprintf(HttpStatusContent, code, content, FooterText)
	_, e1 := w.Write([]byte(body))
	if e1 != nil {
		panic(fmt.Errorf(Stderr.WriteResponseBody, e1.Error()))
	}
}
