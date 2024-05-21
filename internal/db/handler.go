package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
	"log/slog"
)

type VaultEntity struct {
	ID     string `dynamodbav:"id"`
	Name   string `dynamodbav:"name"`
	Secret string `dynamodbav:"secret"`
}

type VaultMetadata struct {
	ID   string `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}

func (dbClient DynamoDBClient) InsertItem(ctx context.Context, vaultEntity VaultEntity) error {
	item, err := attributevalue.MarshalMap(vaultEntity)
	if err != nil {
		slog.Error("error", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	_, err = dbClient.API.PutItem(ctx, input)
	if err != nil {
		slog.Error("error", err)
		return err
	}

	return nil

}

func (dbClient DynamoDBClient) ScanItems(ctx context.Context) ([]VaultMetadata, error) {
	var metadatas []VaultMetadata

	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	output, err := dbClient.API.Scan(ctx, input)
	if err != nil {
		slog.Error("error", err)
		return nil, err
	}

	for _, i := range output.Items {
		var metadata VaultMetadata

		err = attributevalue.UnmarshalMap(i, &metadata)

		if err != nil {
			log.Fatalf("Got error unmarshalling: %s", err)
		}

		metadatas = append(metadatas, metadata)
	}

	return metadatas, err
}
