package db

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type VaultEntity struct {
	ID     string `dynamodbav:"id"`
	Name   string `dynamodbav:"name"`
	Secret string `dynamodbav:"secret"`
}

func (dbClient DynamoDBClient) InsertItem(ctx context.Context, vaultEntity VaultEntity) error {
	item, err := attributevalue.MarshalMap(vaultEntity)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(dbClient.TableName),
		Item:      item,
	}

	_, err = dbClient.API.PutItem(ctx, input)
	if err != nil {
		return err
	}

	return nil

}
