package sss

var stdout = struct {
}{}

var stderr = struct {
	DecodeJSON  string
	DownLoadKey string
	EncodeJSON  string
	ReadObject  string
	PutObject   string
}{
	DecodeJSON:  "could not decode JSON: %v",
	DownLoadKey: "cannot download key %v from bucket %v: %v",
	EncodeJSON:  "could not encode JSON: %v",
	ReadObject:  "cannot read object key %v: %v",
	PutObject:   "could not upload object to s3: %v",
}
