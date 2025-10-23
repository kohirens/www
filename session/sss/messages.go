package sss

var stderr = struct {
	AWSConfig,
	DecodeJSON,
	DownLoadKey,
	EncodeJSON,
	PutObject,
	ReadObject string
}{
	AWSConfig:   "failed to load AWS config: %v",
	DecodeJSON:  "cannot decode JSON: %v",
	DownLoadKey: "cannot download key %v from bucket %v: %v",
	EncodeJSON:  "cannot encode JSON: %v",
	PutObject:   "cannot upload object: %v",
	ReadObject:  "cannot read object key %v: %v",
}

var stdout = struct {
	Saving string
}{
	Saving: "saving data for key %v",
}
