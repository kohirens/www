package dynamo

import (
	"context"
	"fmt"
	"github.com/kohirens/stdlib/logger"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type StorageTable struct {
	Context context.Context
	name    string
	svc     *dynamodb.Client
}

var Log = logger.Standard{}

// Load Session data from a DynamoDB table, the ID serves as the ID key for the
// table.
func (c *StorageTable) Load(id string) ([]byte, error) {
	// Load data from the DynamoDB table
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.name),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := c.svc.GetItem(context.Background(), input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("no item found with ID %s", id)
	}

	itemData := result.Item["Data"].(*types.AttributeValueMemberN).Value

	return []byte(itemData), nil
}

// Save Session data to a DynamoDB table.
func (c *StorageTable) Save(name string, data []byte) error {
	// Save data to the DynamoDB table
	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.name),
		Item: map[string]types.AttributeValue{
			"ID":   &types.AttributeValueMemberS{Value: name},
			"Data": &types.AttributeValueMemberN{Value: string(data)},
		},
	}

	// Verify the results.
	_, e2 := c.svc.PutItem(context.Background(), input)
	if e2 != nil {
		return fmt.Errorf(stderr.PutItem, c.name, e2.Error())
	}

	return nil
}

// NewStorageClient Initializes an S3 client to use as session storage.
// Credentials are expected to be configured in the environment to be picked up
// by the AWS SDK. Panics on failure.
func NewStorageClient(table string) *StorageTable {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}

	svc := dynamodb.NewFromConfig(cfg)
	// Create a DynamoDB client
	return &StorageTable{
		name: table,
		svc:  svc,
	}
}
