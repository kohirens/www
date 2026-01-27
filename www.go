// Package www
package www

import (
	"strings"
)

type StringMap map[string]string

const (
	Http200Debug      = `<!DOCTYPE html><html><head><title>Debugging Implemented (%[1]d)</title></head><body style="text-align:center"><h1>Debugging is currently Implemented</h1><p>%[2]v</p><div>%[3]v</div></body></html>`
	HttpStatusContent = `<!DOCTYPE html><html><head><title>%[1]v %[2]v</title></head><body><center><h1>%[1]v %[2]v</h1><hr /><div>%[3]v</div></center></body></html>`

	// See [Media Types](https://www.iana.org/assignments/media-types/media-types.xhtml)
	// Also see [IETF Media Types](https://www.rfc-editor.org/rfc/rfc9110.html#media.type)

	ContentTypeCSS  = "text/css;charset=utf-8"
	ContentTypeGif  = "image/gif;charset=utf-8"
	ContentTypeHtml = "text/html;charset=utf-8"
	ContentTypeJpg  = "image/jpeg;charset=utf-8"
	ContentTypeJS   = "text/javascript;charset=utf-8"
	ContentTypeJson = "application/json;charset=utf-8"
	ContentTypePng  = "image/png;charset=utf-8"
	ContentTypeSvg  = "image/svg+xml;charset=utf-8"
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
func NewResponse() *Response {
	return &Response{
		Headers: make(map[string][]string),
	}
}
