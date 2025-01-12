package handler

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"personal-vault/internal/db"
	"testing"
)

func TestSaveHandler_AddItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		requestBody        Request
		putItem            func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
		expectedStatus     int
		expectedResponseID string
	}{
		{
			name: "success case",
			requestBody: Request{
				Name:        "testName",
				Description: "",
				Password:    "testPassword",
			},
			putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				return &dynamodb.PutItemOutput{}, nil
			},
			expectedStatus:     http.StatusCreated,
			expectedResponseID: "001",
		},
		{
			name: "validation error case - missing name",
			requestBody: Request{
				Name:        "",
				Description: "",
				Password:    "testPassword",
			},
			putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				return &dynamodb.PutItemOutput{}, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error case - missing password",
			requestBody: Request{
				Name:        "testName",
				Description: "",
				Password:    "",
			},
			putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
				return &dynamodb.PutItemOutput{}, nil
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error case",
			requestBody: Request{
				Name:        "testName",
				Description: "",
				Password:    "testPassword",
			},
			putItem: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
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
					putItem: tt.putItem,
				}}

			// secret is for testing only
			secret, err := hex.DecodeString("0f6f8edf954592d7523b475bb56fd0486b7a049d67c1e5aa522bbc8bfe961971")
			assert.NoError(t, err)

			key := string(secret)

			saveHandler := SaveHandler{Client: dynamdbMockClient, Validate: validator.New(), Key: key}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{
				Header: make(http.Header),
			}

			//content := map[string]interface{}{"name": tt.requestName, "password": tt.requestPassword}

			jsonbytes, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

			saveHandler.AddItem(ctx)
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				body, err := io.ReadAll(w.Body)
				assert.NoError(t, err)
				assert.NotEmpty(t, body)
			}

		})
	}
}
