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
	DirNoExist:   "could not decode JSON: %v",
	DecodeJSON:   "%v directory does not exist",
	EncodeJSON:   "could not encode JSON: %v",
	ReadFile:     "could not read file %v",
	S3Key:        "key %v not found",
	S3PutObject:  "could not put object: %v",
	S3ReadObject: "could not read object %v: %v",
	WriteFile:    "could not write the file %v",
}
var stdout = struct {
	S3Key string
}{
	S3Key: "S3 key %v",
}
