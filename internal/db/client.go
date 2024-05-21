package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	tableName = "personal-vault-dynamodb"
)

type DynamoDBAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type DynamoDBClient struct {
	API       DynamoDBAPI
	TableName string
}

func NewClient(svc DynamoDBAPI) *DynamoDBClient {

	return &DynamoDBClient{
		API:       svc,
		TableName: tableName,
	}
}
