package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/kohirens/www"
	"io"
	"time"
)

type BucketStorage struct {
	Duration time.Duration
	Name     string
	S3       *s3.Client
	Prefix   string
}

// NewBucketStorage Initializes an S3 client to use as storage.
// Credentials are expected to be configured in the environment to be picked up
// by the AWS SDK. Panics on failure.
func NewBucketStorage(bucket string, ctx context.Context) (*BucketStorage, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf(stderr.AwsConfig, err)
	}

	return &BucketStorage{
		Duration: time.Second * 10,
		Name:     bucket,
		S3:       s3.NewFromConfig(cfg),
	}, nil
}

// Load data from S3. It may be best to use a prefix, like the site domain,
// to prevent key name collision in the bucket. See an example at
// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/s3-checksums.html#use-service-S3-checksum-download
func (c *BucketStorage) Load(key string) ([]byte, error) {
	fullKey := c.Prefix + key

	Log.Infof(stdout.S3Key, fullKey)

	obj, e1 := c.S3.GetObject(
		www.GetContextWithTimeout(c.Duration),
		&s3.GetObjectInput{
			Bucket:       &c.Name,
			Key:          &fullKey,
			ChecksumMode: types.ChecksumModeEnabled,
		},
	)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.S3Key, key, c.Name, e1.Error())
	}

	content, e2 := io.ReadAll(obj.Body)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.S3ReadObject, key)
	}

	return content, nil
}

// Save Uploads an object to S3, validating the checksum on success.
// For an example, see
// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/s3-checksums.html#use-service-S3-checksum-upload
func (c *BucketStorage) Save(key string, content []byte) error {
	fullKey := c.Prefix + key

	Log.Infof(stdout.S3Key, fullKey)

	_, e1 := c.S3.PutObject(
		www.GetContextWithTimeout(c.Duration),
		&s3.PutObjectInput{
			Bucket:               &c.Name,
			Key:                  &fullKey,
			Body:                 bytes.NewReader(content),
			ChecksumAlgorithm:    types.ChecksumAlgorithmCrc32,
			ServerSideEncryption: "AES256",
		},
	)
	if e1 != nil {
		return fmt.Errorf(stderr.S3PutObject, e1.Error())
	}

	return nil
}
