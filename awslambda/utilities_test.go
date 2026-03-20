package awslambda

import (
	"reflect"
	"testing"
)

func TestNotImplemented(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		supported []string
		want      bool
	}{
		{"head", "HEAD", []string{""}, true},
		{"get", "GET", []string{"GET"}, false},
		{"POST", "POST", []string{"POST"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotImplemented(tt.method, tt.supported); got != tt.want {
				t.Errorf("NotImplemented() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestConvertToLambdaHttpHeaders(t *testing.T) {
	cases := []struct {
		name    string
		headers map[string][]string
		want    map[string]string
	}{
		{
			"success",
			map[string][]string{"Content-Type": []string{"application/json"}},
			map[string]string{"Content-Type": "application/json"},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToLambdaHttpHeaders(tt.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToLambdaHttpHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
