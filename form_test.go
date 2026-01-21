package www

import (
	"net/url"
	"testing"

	"github.com/kohirens/stdlib/fsio"
)

func TestParseUrlEncodedForm(t *testing.T) {
	tests := []struct {
		name        string
		form        string
		contentType string
		want        url.Values
		wantErr     bool
	}{
		{
			"load-form-with-file",
			"request-meal-plan-upload-2024-01-18T12_31_43-6e03d9cc.json",
			"application/x-www-form-urlencoded",
			map[string][]string{"doc": []string{"menu-01.jpg"}, "due-date": []string{"2024-01-19"}, "name": []string{"Menu 1"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := loadEvent("testdata/" + tt.form)

			form, err := ParseForm(event.Body, tt.contentType)
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
			"lambda-meal-plan-upload-2024-01-20T16-25-57-1d58ba82.json",
			"multipart/form-data; boundary=----WebKitFormBoundarydkoqSwxjXfp9UkJb",
			"testdata/menu-01.jpg",
			243015,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := loadEvent("testdata/" + tt.form)
			_, _ = fsio.CopyToDir(tt.upload, "./", "/")

			got, err := ParseFormWithFiles(event.Body, tt.contentType)
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
