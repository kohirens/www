package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
)

type BucketStorage struct {
	Name   string
	S3     *s3.Client
	Prefix string
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
		Name: bucket,
		S3:   s3.NewFromConfig(cfg),
	}, nil
}

// Load data from S3. It may be best to use a prefix, like the site domain,
// to prevent key name collision in the bucket. See an example at
// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/s3-checksums.html#use-service-S3-checksum-download
func (c *BucketStorage) Load(key string) ([]byte, error) {
	fullKey := c.Prefix + key

	Log.Infof(stdout.LoadKey, fullKey)

	obj, e1 := c.S3.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket:       &c.Name,
			Key:          &fullKey,
			ChecksumMode: types.ChecksumModeEnabled,
		},
	)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.LoadKey, key, c.Name, e1.Error())
	}

	content, e2 := io.ReadAll(obj.Body)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.ReadObject, key)
	}

	return content, nil
}

// Save Uploads an object to S3, validating the checksum on success.
// For an example, see
// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/s3-checksums.html#use-service-S3-checksum-upload
func (c *BucketStorage) Save(key string, content []byte) error {
	fullKey := c.Prefix + key

	Log.Infof(stdout.SaveKey, fullKey)

	_, e1 := c.S3.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket:               &c.Name,
			Key:                  &fullKey,
			Body:                 bytes.NewReader(content),
			ChecksumAlgorithm:    types.ChecksumAlgorithmCrc32,
			ServerSideEncryption: "AES256",
		},
	)
	if e1 != nil {
		return fmt.Errorf(stderr.PutObject, e1.Error())
	}

	return nil
}

func (c *BucketStorage) Location(key string) string {
	return c.Prefix + key
}

// Remove Delete an object from S3.
func (c *BucketStorage) Remove(key string) error {
	fullKey := c.Location(key)

	Log.Infof(stdout.SaveKey, fullKey)

	_, e1 := c.S3.DeleteObject(
		context.Background(),
		&s3.DeleteObjectInput{
			Bucket: &c.Name,
			Key:    &fullKey,
		},
	)
	if e1 != nil {
		return fmt.Errorf(stderr.DeleteObject, e1.Error())
	}

	return nil
}
