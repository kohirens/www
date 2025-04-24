package sss

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kohirens/stdlib/logger"
	"github.com/kohirens/www/session"
	"io"
	"log"
	"time"
)

type StorageBucket struct {
	Context context.Context
	Name    string
	S3      *s3.Client
	prefix  string
}

var Log = logger.Standard{}

// Load Session data from S3, the ID serves as the object name. We recommend the
// site domain be used as a prefix to prevent collision in the bucket.
// The timeout on the Context will interrupt the request if it expires.
// See also https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#example_S3_GetObject_shared00
func (c *StorageBucket) Load(key string) (*session.Data, error) {
	fullKey := c.prefix + key

	Log.Infof("Loading session for key %v", fullKey)

	obj, e1 := c.S3.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: &c.Name,
			Key:    &fullKey,
		},
	)

	if e1 != nil {
		return nil, fmt.Errorf(stderr.DownLoadKey, key, c.Name, e1.Error())
	}

	b, e2 := io.ReadAll(obj.Body)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.ReadObject, key)
	}

	data := &session.Data{}
	if e := json.Unmarshal(b, data); e != nil {
		return nil, fmt.Errorf(stderr.DecodeJSON, key)
	}
	return data, nil
}

// Save Session data to S3.
func (c *StorageBucket) Save(data *session.Data) error {
	// TODO: Lock the object on now
	// TODO: Check if the object is locked.
	// if it is then wait and try again.
	// If not locked, then lock it. then unlock when done.
	content, e1 := json.Marshal(data)
	if e1 != nil {
		return fmt.Errorf(stderr.EncodeJSON, e1)
	}

	_, e := c.Upload(content, data.Id)
	if e != nil {
		return e
	}

	return nil
}

// Upload Uploads an object to S3, returning the eTag on success. The Context
// will interrupt the request if the timeout expires.
// For more info, see https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#example_S3_PutObject_shared00
func (c *StorageBucket) Upload(b []byte, key string) (string, error) {
	fullKey := c.prefix + key

	Log.Infof("Saving data for key %v", fullKey)

	put, e1 := c.S3.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:               &c.Name,
		Key:                  &fullKey,
		Body:                 bytes.NewReader(b),
		ServerSideEncryption: "AES256",
		Expires:              aws.Time(time.Now().UTC().AddDate(0, 0, 7)),
	})

	if e1 != nil {
		return "", fmt.Errorf(stderr.PutObject, e1.Error())
	}

	return *put.ETag, nil
}

// NewStorageClient Initializes an S3 client to use as session storage.
// Credentials are expected to be configured in the environment to be picked up
// by the AWS SDK. Panics on failure.
func NewStorageClient(bucket string, ctx context.Context) *StorageBucket {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}
	client := s3.NewFromConfig(cfg)
	return &StorageBucket{
		Name:    bucket,
		S3:      client,
		Context: ctx,
	}
}

// Prefix Set a prefix for the bucket to prepend to keys before downloaded or uploading objects.
func (c *StorageBucket) Prefix(prefix string) *StorageBucket {
	c.prefix = prefix
	return c
}
