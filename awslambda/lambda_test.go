package awslambda

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

const (
	fixtureDir = "testdata"
)

func TestPreliminaryChecks(runner *testing.T) {
	runner.Setenv("REDIRECT_TO", "www.example.com")
	runner.Setenv("REDIRECT_HOSTS", "example.com")
	runner.Setenv("HTTP_METHODS_ALLOWED", "GET,HEAD,POST")

	tests := []struct {
		name    string
		event   *Input
		want    int
		wantNil bool
	}{
		{
			"not-implemented",
			&Input{
				RequestContext: &Context{
					HTTP: &Http{
						Method: "PUT",
					},
				},
			},
			501,
			false,
		},
		{
			"redirect-301",
			&Input{
				Headers: map[string]string{headerAltHost: "example.com"},
				RequestContext: &Context{
					HTTP: &Http{
						Method: "GET",
					},
				},
			},
			301,
			false,
		},
		{
			"redirect-308",
			&Input{
				Headers: map[string]string{headerAltHost: "example.com"},
				RequestContext: &Context{
					HTTP: &Http{
						Method: "POST",
					},
				},
			},
			308,
			false,
		},
		{
			"request-using-cloudfront-domain-not-allowed",
			&Input{
				Headers: map[string]string{
					headerAltHost:  "cfd.cloudfront.aws",
					headerCfDomain: "cfd.cloudfront.aws",
				},
				RequestContext: &Context{
					HTTP: &Http{
						Method: "GET",
					},
				},
			},
			401,
			false,
		},
		{
			"ok",
			&Input{
				Headers: map[string]string{
					headerAltHost:  "www.example.com",
					headerCfDomain: "cfd.cloudfront.aws",
				},
				RequestContext: &Context{
					HTTP: &Http{
						Method: "GET",
					},
				},
				RawPath: fixtureDir + string(os.PathSeparator) + "index.html",
			},
			0,
			true,
		},
	}

	for _, c := range tests {
		runner.Run(c.name, func(t *testing.T) {
			got := PreliminaryChecks(c.event)

			if (got == nil) != c.wantNil {
				t.Errorf("PreliminaryChecks() error, want %v", c.wantNil)
				return
			}

			if (got != nil) && got.StatusCode != c.want {
				t.Errorf("PreliminaryChecks() got %v, want %v", got.StatusCode, c.want)
				return
			}
		})
	}
}

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name    string
		lambda  string
		wantErr bool
	}{
		{
			"canParsePost",
			"lambda-meal-plan-upload-2024-01-18T12_31_43-6e03d9cc.json",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := loadLambdaFixture(tt.lambda)
			got, err := NewRequest(l)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequestFromLambdaFunctionURLRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Do you get POST data in the http.Request
			_ = got.ParseForm()
			gotData := got.Form.Get("name")
			fmt.Println(got.Form)
			wantData := "Menu 1"
			if gotData != wantData {
				t.Errorf("NewRequestFromLambdaFunctionURLRequest().Request.Form.Get(\"\") = %v,  want %v", gotData, wantData)
				return
			}
		})
	}
}

func loadLambdaFixture(fn string) *Input {
	j, err := os.ReadFile("testdata/" + fn)
	if err != nil {
		panic(err)
	}

	l := &Input{}
	if e := json.Unmarshal(j, l); e != nil {
		panic(e)
	}

	return l
}

func TestDoRedirect(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		rt      string
		rh      string
		want    bool
		wantErr bool
	}{
		{
			"env var REDIRECT_TO not set",
			"www.example.com",
			"",
			"",
			false,
			true,
		},
		{
			"does not redirect host",
			"www.example.com",
			"www.example.com",
			"example.com",
			false,
			false,
		},
		{
			"redirect host",
			"example.com",
			"www.example.com",
			"example.com",
			true,
			false,
		},
		{
			"redirect host",
			"example.com",
			"www.example.com",
			"example.com",
			true,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.rt != "" {
				_ = os.Setenv("REDIRECT_TO", tt.rt)
				defer func() { _ = os.Unsetenv("REDIRECT_TO") }()
			}
			if tt.rh != "" {
				_ = os.Setenv("REDIRECT_HOSTS", tt.rh)
				defer func() { _ = os.Unsetenv("REDIRECT_HOSTS") }()
			}

			got, err := ShouldRedirect(tt.host)

			if (err != nil) != tt.wantErr {
				t.Errorf("ShouldRedirect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != got {
				t.Errorf("ShouldRedirect() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestDoRedirect2(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		to      string
		hosts   string
		want    bool
		wantErr bool
	}{
		{
			"cannot get host from request",
			"",
			"www.example.com",
			"example.com",
			false,
			true,
		},
		{
			"REDIRECT_TO is set to empty string",
			"www.example.com",
			"",
			"example.com",
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("REDIRECT_TO", tt.to)
			defer func() { _ = os.Unsetenv("REDIRECT_TO") }()
			_ = os.Setenv("REDIRECT_HOSTS", tt.hosts)
			defer func() { _ = os.Unsetenv("REDIRECT_HOSTS") }()

			got, err := ShouldRedirect(tt.host)

			if (err != nil) != tt.wantErr {
				t.Errorf("ShouldRedirect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != got {
				t.Errorf("ShouldRedirect() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
