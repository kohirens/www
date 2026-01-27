package www

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestGetHeader(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		header  string
		want    string
	}{
		{"get-lowercase-host-header", map[string]string{"Host": "example.com"}, "host", "example.com"},
		{"get-uppercase-host-header", map[string]string{"HOST": "example.com"}, "host", "example.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHeader(tt.headers, tt.header)

			if got != tt.want {
				t.Errorf("GetHeader() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestGetPageType(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    string
	}{
		{"html-type", map[string]string{"content-type": "text/html"}, "text/html"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPageType(tt.headers); got != tt.want {
				t.Errorf("GetPageType() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestGetPageTypeByExt(t *testing.T) {
	tests := []struct {
		name     string
		pagePath string
		want     string
	}{
		{"html", "page.html", "text/html;charset=utf-8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPageTypeByExt(tt.pagePath)

			if got != tt.want {
				t.Errorf("GetPageTypeByExt() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestRespond301Or308(t *testing.T) {
	tests := []struct {
		name     string
		call     func(http.ResponseWriter, string)
		location string
		w        *Response
		wantCode int
		status   string
	}{
		{"201", Respond201, "https://www.example.com", NewResponse(), 201, "Created"},
		{"301", Respond301, "https://www.example.com", NewResponse(), 301, "Moved Permanently"},
		{"302", Respond302, "https://www.example.com", NewResponse(), 302, "Found"},
		{"308", Respond308, "https://www.example.com", NewResponse(), 308, "Permanent Redirect"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.call(tt.w, tt.location)

			if tt.w.StatusCode != tt.wantCode {
				t.Errorf("Respond%v() = %v, want %v", tt.name, tt.w.StatusCode, tt.wantCode)
				return
			}

			if !strings.Contains(tt.w.Body, tt.status) {
				t.Errorf("Respond%v() = does not contain %v", tt.name, tt.status)
				return
			}
		})
	}
}

// This suite of test ensure any refactoring of these methods leave the
// required HTTP status code and recommended status message are left intact.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/500
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/501
func TestRespondXXX(t *testing.T) {
	tests := []struct {
		call       func(w http.ResponseWriter, body []byte, contentType string)
		w          *Response
		wantCode   int
		wantStatus string
	}{
		{Respond401, NewResponse(), 401, http.StatusText(401)},
		{Respond404, NewResponse(), 404, http.StatusText(404)},
		{Respond500, NewResponse(), 500, http.StatusText(500)},
		{Respond501, NewResponse(), 501, http.StatusText(501)},
	}
	for _, tt := range tests {
		t.Run(tt.wantStatus, func(t *testing.T) {
			tt.call(tt.w, []byte{}, "")

			gotCode := tt.w.StatusCode
			if gotCode != tt.wantCode {
				t.Errorf("Respond%v() = %v, want %v", tt.wantCode, gotCode, tt.wantCode)
				return
			}

			gotBody := tt.w.Body
			if !strings.Contains(gotBody, tt.wantStatus) {
				t.Errorf("Respond%v() = does not contain %v", tt.wantCode, tt.wantStatus)
				return
			}
		})
	}
}

// This suite of test ensure any refactoring of these methods leave the
// required HTTP status code and recommended status message are left intact.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/405
func TestRespond405(t *testing.T) {
	tests := []struct {
		name       string
		methods    string
		w          *Response
		wantCode   int
		wantStatus string
	}{
		{"405", "GET, HEAD, POST,", NewResponse(), 405, "405"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Respond405(tt.w, tt.methods)

			if tt.w.StatusCode != tt.wantCode {
				t.Errorf("Respond%v() = %v, want %v", tt.name, tt.w.StatusCode, tt.wantCode)
				return
			}
		})
	}
}

func TestRespondJSON(t *testing.T) {
	type jsonMsg struct {
		Msg string `json:"msg"`
	}

	fixedBody := &jsonMsg{"Salam"}

	tests := []struct {
		name     string
		content  *jsonMsg
		w        *Response
		wantBody string
	}{
		{"can-encode", fixedBody, NewResponse(), `{"msg":"Salam"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RespondWithJSON(tt.w, tt.content)

			if tt.w.Body != tt.wantBody {
				t.Errorf("RespondWithJSON() = %v, want %v", tt.w.Body, tt.wantBody)
				return
			}
		})
	}
}

func TestRespondDebug(t *testing.T) {
	tests := []struct {
		name    string
		message string
		footer  string
		code    int

		w *Response
	}{
		{"Debug200", "status ok", "Acme", 401, NewResponse()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RespondDebug(tt.w, tt.code, tt.message, tt.footer)

			if !strings.Contains(tt.w.Body, fmt.Sprintf("%v", tt.code)) {
				t.Errorf("RespondDebug() = does not contain %v", tt.code)
				return
			}

			if !strings.Contains(tt.w.Body, fmt.Sprintf("%v", tt.message)) {
				t.Errorf("RespondDebug() = does not contain %v", tt.footer)
				return
			}

			if !strings.Contains(tt.w.Body, fmt.Sprintf("%v", tt.footer)) {
				t.Errorf("RespondDebug() = does not contain %v", tt.footer)
				return
			}
		})
	}
}
