package google

import (
	"github.com/kohirens/stdlib/test"
	"github.com/kohirens/www/backend"
	"net/http"
	"net/url"
	"testing"
)

func TestAuthLink(t *testing.T) {
	goodAuth := backend.NewAuthManager()
	goodAuth.Add(backend.KeyGoogleProvider, &MockProvider{
		ExpectedAuthLink: "good-link",
	})

	cases := []struct {
		name    string
		w       http.ResponseWriter
		r       *http.Request
		a       backend.App
		wantErr bool
	}{
		{
			"provider_not_found",
			&test.MockResponseWriter{
				ExpectedBody:       nil,
				ExpectedHeaders:    nil,
				Headers:            nil,
				ExpectedStatusCode: 200,
			},
			&http.Request{
				URL: &url.URL{
					Scheme:   "https",
					Host:     "google.com",
					Path:     "/auth/google/callback",
					RawQuery: "email=test@example.com",
				},
			},
			&MockApp{
				Authorizer: backend.NewAuthManager(),
			},
			true,
		},
		{
			"return_link",
			&test.MockResponseWriter{
				Headers:            nil,
				ExpectedStatusCode: 200,
			},
			&http.Request{
				URL: &url.URL{
					RawQuery: "email=test@example.com",
				},
				Header: http.Header{},
			},
			&MockApp{
				Authorizer: goodAuth,
			},
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := AuthLink(c.w, c.r, c.a); (err != nil) != c.wantErr {
				t.Errorf("AuthLink() error = %v, wantErr %v", err, c.wantErr)
				return
			}
		})
	}
}

func TestSignIn(t *testing.T) {
	goodAuth := backend.NewAuthManager()
	goodAuth.Add(backend.KeyGoogleProvider, &MockProvider{
		ExpectedAuthLink: "good-link",
	})

	cases := []struct {
		name    string
		w       http.ResponseWriter
		r       *http.Request
		a       backend.App
		wantErr bool
	}{
		{
			"provider_not_found",
			&test.MockResponseWriter{
				ExpectedBody:       nil,
				ExpectedHeaders:    nil,
				Headers:            nil,
				ExpectedStatusCode: 200,
			},
			&http.Request{
				URL: &url.URL{
					Scheme:   "https",
					Host:     "google.com",
					Path:     "/auth/google/callback",
					RawQuery: "email=test@example.com",
				},
			},
			&MockApp{
				Authorizer: backend.NewAuthManager(),
			},
			true,
		},
		{
			"redirect_to_google_auth_server",
			&test.MockResponseWriter{
				Headers:            nil,
				ExpectedStatusCode: 307,
			},
			&http.Request{
				URL: &url.URL{
					RawQuery: "email=test@example.com",
				},
				Header: http.Header{},
			},
			&MockApp{
				Authorizer: goodAuth,
			},
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := SignIn(c.w, c.r, c.a); (err != nil) != c.wantErr {
				t.Errorf("AuthLink() error = %v, wantErr %v", err, c.wantErr)
				return
			}
		})
	}
}

func TestSignOut(t *testing.T) {
	goodAuth := backend.NewAuthManager()
	goodAuth.Add(backend.KeyGoogleProvider, &MockProvider{
		ExpectedAuthLink: "good-link",
	})

	cases := []struct {
		name    string
		w       http.ResponseWriter
		r       *http.Request
		a       backend.App
		wantErr bool
	}{
		{
			"provider_not_found",
			nil,
			nil,
			&MockApp{
				Authorizer: backend.NewAuthManager(),
			},
			true,
		},
		{
			"redirect_to_google_auth_server",
			&test.MockResponseWriter{
				Headers:            nil,
				ExpectedStatusCode: 307,
			},
			nil,
			&MockApp{
				Authorizer: goodAuth,
			},
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := SignOut(c.w, c.r, c.a); (err != nil) != c.wantErr {
				t.Errorf("AuthLink() error = %v, wantErr %v", err, c.wantErr)
				return
			}
		})
	}
}
