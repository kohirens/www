package storage

var stderr = struct {
	AwsConfig,
	DecodeJSON,
	DirNoExist,
	EncodeJSON,
	ReadFile,
	S3Key,
	S3ReadObject,
	S3PutObject,
	WriteFile string
}{
	AwsConfig:    "failed to load AWS config: %v",
	DecodeJSON:   "directory not decode JSON: %v",
	DirNoExist:   "%v directory does not exist",
	EncodeJSON:   "cannot encode JSON: %v",
	ReadFile:     "cannot read file %v",
	S3Key:        "key %v not found",
	S3PutObject:  "cannot put object: %v",
	S3ReadObject: "cannot read object %v: %v",
	WriteFile:    "cannot write the file %v",
}
var stdout = struct {
	S3Key string
}{
	S3Key: "S3 key %v",
}
