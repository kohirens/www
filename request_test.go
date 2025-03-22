package www

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"os"
	"reflect"
	"testing"
)

func TestNewRequestFromLambdaFunctionURLRequest(t *testing.T) {
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
			got, err := NewRequestFromLambdaFunctionURLRequest(l)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequestFromLambdaFunctionURLRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Do you get back what you put in?
			want := got.ToLambdaFunctionURLRequest()
			if !reflect.DeepEqual(l, want) {
				t.Errorf("NewRequestFromLambdaFunctionURLRequest() got = %v, want %v", got, want)
				return
			}

			// Do you get POST data in the http.Request
			_ = got.Request.ParseForm()
			gotData := got.Request.Form.Get("name")
			fmt.Println(got.Request.Form)
			wantData := "Menu 1"
			if gotData != wantData {
				t.Errorf("NewRequestFromLambdaFunctionURLRequest().Request.Form.Get(\"\") = %v,  want %v", gotData, wantData)
				return
			}
		})
	}
}

func loadLambdaFixture(fn string) *events.LambdaFunctionURLRequest {
	j, err := os.ReadFile("testdata/" + fn)
	if err != nil {
		panic(err)
	}

	l := &events.LambdaFunctionURLRequest{}
	if e := json.Unmarshal(j, l); e != nil {
		panic(e)
	}

	return l
}
