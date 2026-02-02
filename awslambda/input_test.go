package awslambda

import (
	"encoding/json"
	"os"
	"testing"
)

func TestInput_Cookie(runner *testing.T) {
	cases := []struct {
		name       string
		fixture    string
		cookieName string
		want       string
		wantErr    bool
	}{
		{
			"not-found",
			"event-in-01.json",
			"test",
			"",
			true,
		},
		{
			"found",
			"event-in-01.json",
			"nevergonna",
			"bringyoudown",
			false,
		},
	}
	for _, tc := range cases {
		runner.Run(tc.name, func(t *testing.T) {
			var input *Input
			jsonFixture, _ := os.ReadFile(fixtureDir + "/" + tc.fixture)
			_ = json.Unmarshal(jsonFixture, &input)

			if e := input.ParseCookies(); e != nil {
				t.Fatal(e)
				return
			}

			got, err := input.Cookie(tc.cookieName)
			if (err != nil) != tc.wantErr {
				t.Errorf("Cookie() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != nil && (got.Value != tc.want) {
				t.Errorf("Cookie() got = %v, want %v", got, tc.want)
				return
			}
		})
	}
}
