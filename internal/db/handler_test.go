package db

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type dynamoDBMockAPI struct {
	getItem func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	putItem func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	scan    func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

func (m *dynamoDBMockAPI) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.getItem(ctx, params, optFns...)
}

func (m *dynamoDBMockAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.putItem(ctx, params, optFns...)
}

func (m *dynamoDBMockAPI) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return m.scan(ctx, params, optFns...)
}

func TestDynamoDBClient_PutItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		vaultEntity VaultEntity
		putItem     func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
		expectedErr error
	}{
		{
			name: "success case",
			vaultEntity: VaultEntity{
				ID:          "001",
				Name:        "TestName",
				Description: "TestDescr.",
				Password:    "TestPassword",
			},
			putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				return &dynamodb.PutItemOutput{}, nil
			},
			expectedErr: nil,
		},
		{
			name: "error case",
			vaultEntity: VaultEntity{
				ID:          "001",
				Name:        "TestName",
				Description: "TestDescr.",
				Password:    "TestPassword",
			},
			putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				return nil, errors.New("this is mock error")
			},
			expectedErr: errors.New("this is mock error"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dynamdbMockClient := DynamoDBClient{
				API: &dynamoDBMockAPI{
					putItem: tt.putItem,
				}}
			err := dynamdbMockClient.PutItem(context.Background(), tt.vaultEntity)
			assert.Equal(t, err, tt.expectedErr)
		})
	}
}

func TestDynamoDBClient_ScanItems(t *testing.T) {
	t.Parallel()

	var (
		items []map[string]types.AttributeValue
		item  map[string]types.AttributeValue
	)

	item = map[string]types.AttributeValue{
		"ID":          &types.AttributeValueMemberS{Value: "001"},
		"Name":        &types.AttributeValueMemberS{Value: "TestName"},
		"Description": &types.AttributeValueMemberS{Value: "TestDescr."},
		"Password":    &types.AttributeValueMemberS{Value: "TestPassword"},
	}

	items = append(items, item)

	tests := []struct {
		name        string
		scan        func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
		expectedErr error
	}{
		{
			name: "success case",
			scan: func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
				return &dynamodb.ScanOutput{
					Items: items,
				}, nil
			},
			expectedErr: nil,
		},
		{
			name: "error case",
			scan: func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
				return nil, errors.New("this is mock error")
			},
			expectedErr: errors.New("this is mock error"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dynamdbMockClient := DynamoDBClient{
				API: &dynamoDBMockAPI{
					scan: tt.scan,
				}}
			metadatas, err := dynamdbMockClient.ScanItems(context.Background())
			if tt.expectedErr != nil {
				assert.Equal(t, err, tt.expectedErr)
				assert.Nil(t, metadatas)
			} else {
				assert.NotEmpty(t, metadatas)
			}

		})
	}
}

func TestDynamoDBClient_GetItem(t *testing.T) {
	t.Parallel()

	item := map[string]types.AttributeValue{
		"ID":          &types.AttributeValueMemberS{Value: "001"},
		"Name":        &types.AttributeValueMemberS{Value: "TestName"},
		"Description": &types.AttributeValueMemberS{Value: "TestDescr."},
		"Password":    &types.AttributeValueMemberS{Value: "TestPassword"},
	}

	tests := []struct {
		name        string
		getItem     func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
		expectedErr error
	}{
		{
			name: "success case",
			getItem: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{
					Item: item,
				}, nil
			},
			expectedErr: nil,
		},
		{
			name: "error case",
			getItem: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("this is mock error")
			},
			expectedErr: errors.New("this is mock error"),
		},
		{
			name: "item not found",
			getItem: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{
					Item: nil,
				}, nil
			},
			expectedErr: errors.New("unable to find the record"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dynamdbMockClient := DynamoDBClient{
				API: &dynamoDBMockAPI{
					getItem: tt.getItem,
				}}
			password, err := dynamdbMockClient.GetItem(context.Background(), "001")
			if tt.expectedErr != nil {
				assert.Equal(t, err, tt.expectedErr)
				assert.Empty(t, password)
			} else {
				assert.Equal(t, password, "TestPassword")
			}

		})
	}
}
