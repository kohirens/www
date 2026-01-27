package www

import (
	"net/url"
	"testing"

	"github.com/kohirens/stdlib/fsio"
)

func TestParseUrlEncodedForm(t *testing.T) {
	tests := []struct {
		name    string
		body    []byte
		want    url.Values
		wantErr bool
	}{
		{
			"load-form-with-file",
			[]byte("ZG9jPW1lbnUtMDEuanBnJm5hbWU9TWVudSsxJmR1ZS1kYXRlPTIwMjQtMDEtMTk="),
			map[string][]string{"doc": {"menu-01.jpg"}, "due-date": {"2024-01-19"}, "name": {"Menu 1"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form, err := ParseForm(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, _ := form.Field("doc")
			if got != tt.want["doc"][0] {
				t.Errorf("ParseForm() got = %v, want %v", got, tt.want["doc"][0])
				return
			}
		})
	}
}

func TestParseForm2(t *testing.T) {
	tests := []struct {
		name        string
		form        string
		contentType string
		upload      string
		want        int64
		wantErr     bool
	}{
		{
			"load-form-with-file",
			"lambda-html-form-base64.txt",
			"multipart/form-data; boundary=----WebKitFormBoundarydkoqSwxjXfp9UkJb",
			"testdata/menu-01.jpg",
			243015,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := loadFile("testdata/" + tt.form)
			_, _ = fsio.CopyToDir(tt.upload, "./", "/")

			got, err := ParseFormWithFiles(body, tt.contentType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotFile, _ := got.File("doc")

			if gotFile.Size != tt.want {
				t.Errorf("ParseForm() got = %v, want %v", gotFile.Size, tt.want)
				return
			}
		})
	}
}
