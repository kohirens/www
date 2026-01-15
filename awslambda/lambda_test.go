package awslambda

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
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
		event   *events.LambdaFunctionURLRequest
		want    int
		wantNil bool
	}{
		{
			"not-implemented",
			&events.LambdaFunctionURLRequest{
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "PUT",
					},
				},
			},
			501,
			false,
		},
		{
			"redirect-301",
			&events.LambdaFunctionURLRequest{
				Headers: map[string]string{headerAltHost: "example.com"},
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "GET",
					},
				},
			},
			301,
			false,
		},
		{
			"redirect-308",
			&events.LambdaFunctionURLRequest{
				Headers: map[string]string{headerAltHost: "example.com"},
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "POST",
					},
				},
			},
			308,
			false,
		},
		{
			"request-using-cloudfront-domain-not-allowed",
			&events.LambdaFunctionURLRequest{
				Headers: map[string]string{
					headerAltHost:  "cfd.cloudfront.aws",
					headerCfDomain: "cfd.cloudfront.aws",
				},
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
						Method: "GET",
					},
				},
			},
			401,
			false,
		},
		{
			"ok",
			&events.LambdaFunctionURLRequest{
				Headers: map[string]string{
					headerAltHost:  "www.example.com",
					headerCfDomain: "cfd.cloudfront.aws",
				},
				RequestContext: events.LambdaFunctionURLRequestContext{
					HTTP: events.LambdaFunctionURLRequestContextHTTPDescription{
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
