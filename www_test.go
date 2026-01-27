package www

import (
	"fmt"
	"os"
	"testing"

	"github.com/kohirens/stdlib/fsio"
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
func loadFile(s string) []byte {
	if !fsio.Exist(s) {
		panic(fmt.Sprintf("file %s not found", s))
	}

	content, e1 := os.ReadFile(s)
	if e1 != nil {
		panic(fmt.Sprintf("could not read file %s", s))
	}

	return content
}
