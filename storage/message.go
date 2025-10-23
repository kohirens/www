package storage

var stderr = struct {
	AwsConfig,
	DecodeJSON,
	DirNoExist,
	EncodeJSON,
	ReadFile,
	LoadKey,
	ReadObject,
	PutObject,
	WriteFile string
}{
	AwsConfig:  "failed to load AWS config: %v",
	DecodeJSON: "cannot decode JSON: %v",
	DirNoExist: "%v directory does not exist",
	EncodeJSON: "cannot encode JSON: %v",
	ReadFile:   "cannot read file %v",
	LoadKey:    "cannot load object key %v in bucket %v: %v",
	PutObject:  "cannot put object: %v",
	ReadObject: "cannot read object %v: %v",
	WriteFile:  "cannot write the file %v",
}
var stdout = struct {
	LoadKey,
	SaveKey string
}{
	LoadKey: "loading object from key %v",
	SaveKey: "saving object to key %v",
}
