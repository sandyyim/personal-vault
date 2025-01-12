package handler

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"personal-vault/internal/db"
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

func TestRetrieveHandler_GetAll(t *testing.T) {
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
		name                 string
		scan                 func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
		expectedStatus       int
		expectedResponseSize int
		expectedResponseID   string
		expectedResponseName string
	}{
		{
			name: "success case",
			scan: func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
				return &dynamodb.ScanOutput{
					Items: items,
				}, nil
			},
			expectedStatus:       http.StatusOK,
			expectedResponseSize: 1,
			expectedResponseID:   "001",
			expectedResponseName: "TestName",
		},
		{
			name: "error case",
			scan: func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
				return nil, errors.New("this is mock error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dynamdbMockClient := db.DynamoDBClient{
				API: &dynamoDBMockAPI{
					scan: tt.scan,
				}}

			retrieveHandler := RetrieveHandler{Client: dynamdbMockClient}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			retrieveHandler.GetAll(ctx)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				body, err := io.ReadAll(w.Body)
				assert.NoError(t, err)

				var responses []db.VaultMetadata
				err = json.Unmarshal(body, &responses)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseSize, len(responses))

				for _, response := range responses {
					assert.Equal(t, response.ID, tt.expectedResponseID)
					assert.Equal(t, response.Name, tt.expectedResponseName)
				}
			}

		})
	}
}

func TestRetrieveHandler_GetByID(t *testing.T) {
	t.Parallel()

	item := map[string]types.AttributeValue{
		"ID":          &types.AttributeValueMemberS{Value: "001"},
		"Name":        &types.AttributeValueMemberS{Value: "TestName"},
		"Description": &types.AttributeValueMemberS{Value: "TestDescr."},
		"Password":    &types.AttributeValueMemberS{Value: "gA8vgNGMxa3W0M0t7059MhLqYruaVgFRaVzuGcTAIXzIhY2mKAVqbw=="},
	}

	tests := []struct {
		name             string
		testId           string
		getItem          func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:   "success case",
			testId: "6b2bfbc0-8c23-414b-9c39-cf9b76520b39",
			getItem: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{
					Item: item,
				}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "testPassword",
		},
		{
			name:   "invalid id case",
			testId: "6b2bfbc0-414b-9c39-cf9b76520b39",
			getItem: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return &dynamodb.GetItemOutput{
					Item: item,
				}, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "db error case",
			testId: "6b2bfbc0-8c23-414b-9c39-cf9b76520b39",
			getItem: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("this is mock error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dynamdbMockClient := db.DynamoDBClient{
				API: &dynamoDBMockAPI{
					getItem: tt.getItem,
				}}

			// secret is for testing only
			secret, err := hex.DecodeString("0f6f8edf954592d7523b475bb56fd0486b7a049d67c1e5aa522bbc8bfe961971")
			assert.NoError(t, err)

			key := string(secret)

			retrieveHandler := RetrieveHandler{Client: dynamdbMockClient, Key: key}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			params := []gin.Param{
				{
					Key:   "id",
					Value: tt.testId,
				},
			}

			ctx.Params = params

			retrieveHandler.GetByID(ctx)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				body, err := io.ReadAll(w.Body)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedResponse, string(body))
			}

		})
	}
}
