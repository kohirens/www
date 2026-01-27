package www

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime"
	"mime/multipart"
	"net/url"
	"strings"
)

var (
	MaxMemory = int64(1e+6 * 8 * 5)
)

type Form interface {
	Field(key string) (string, error)
}

type FormData struct {
	data *multipart.Form
}

func (form *FormData) Field(key string) (string, error) {
	fv, ok := form.data.Value[key]
	if !ok {
		return "", fmt.Errorf(Stderr.FieldNotFound, key)
	}
	return fv[0], nil
}

func (form *FormData) File(key string) (*multipart.FileHeader, error) {
	fh, ok := form.data.File[key]
	if !ok {
		return nil, fmt.Errorf(Stderr.FieldNotFound, key)
	}

	if len(fh) > 0 {
		return fh[0], nil
	}
	return nil, fmt.Errorf(Stderr.FieldNotFound, key)
}

type FormUrlEncoded struct {
	data *url.Values
}

// Field Get a field from the form.
func (fd *FormUrlEncoded) Field(key string) (string, error) {
	if fd.data.Has(key) {
		return fd.data.Get(key), nil
	}

	return "", fmt.Errorf(Stderr.FieldNotFound, key)
}

func ParseForm(encodedData []byte) (*FormUrlEncoded, error) {
	decodedData := make([]byte, base64.StdEncoding.DecodedLen(len(encodedData)))
	_, e1 := base64.StdEncoding.Decode(decodedData, encodedData)
	if e1 != nil {
		return nil, fmt.Errorf("could not decode 64 bit string: %v", e1.Error())
	}

	formData, e3 := url.ParseQuery(string(decodedData))
	if e3 != nil {
		return nil, fmt.Errorf(": %v", e3.Error())
	}

	return &FormUrlEncoded{data: &formData}, nil
}

func ParseFormWithFiles(encodedData []byte, contentType string) (*FormData, error) {
	decodedData := make([]byte, base64.StdEncoding.DecodedLen(len(encodedData)))
	_, e1 := base64.StdEncoding.Decode(decodedData, encodedData)
	if e1 != nil {
		return nil, fmt.Errorf(Stderr.DecodeBase64, e1.Error())
	}

	reader := bytes.NewReader(decodedData)

	mediaType, params, e1 := mime.ParseMediaType(contentType)
	if e1 != nil {
		return nil, e1
	}

	var formData *multipart.Form

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(reader, params["boundary"])

		f, e2 := mr.ReadForm(MaxMemory)
		if e2 != nil {
			return nil, e2
		}

		formData = f
	}

	return &FormData{data: formData}, nil
}
