package dynamo

var stdout = struct {
}{}

var stderr = struct {
	DecodeJSON,
	EncodeJSON,
	PutItem string
}{
	DecodeJSON: "could not decode JSON: %v",
	EncodeJSON: "could not encode JSON: %v",
	PutItem:    "could not put item %v to dynamodb table %v: %v",
}
