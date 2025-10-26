package backend

import (
	"bytes"
	"github.com/kohirens/www/storage"
	"testing"
)

func TestHandler_Render(t *testing.T) {
	fixtureLs, _ := storage.NewLocalStorage(fixtureDir)

	cases := []struct {
		name     string
		store    storage.Storage
		filename string
		vars     Variables
		wantW    string
		wantErr  bool
	}{
		{
			"simple-render",
			fixtureLs,
			"test-render-01",
			Variables{"TestVar": "1234"},
			"1234",
			false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := NewTemplateManager(c.store, "", "tmpl")
			w := &bytes.Buffer{}
			err := h.Render(c.filename, w, c.vars)

			if (err != nil) != c.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if gotW := w.String(); gotW != c.wantW {
				t.Errorf("Render() gotW = %v, want %v", gotW, c.wantW)
				return
			}
		})
	}
}
