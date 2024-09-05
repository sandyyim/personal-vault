package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
	"personal-vault/internal/db"
	"personal-vault/internal/vault"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func healthcheckHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

func notFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
}

func notMethodHandler(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"code": "METHOD_NOT_ALLOWED", "message": "405 method not allowed"})
}

func main() {

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return
	}

	os.Setenv("AWS_ENDPOINT_URL_DYNAMODB", "http://localhost:8000")

	svc := dynamodb.NewFromConfig(awsConfig)
	dbClient := db.NewClient(svc)

	validate := validator.New()

	saveHandler := vault.SaveHandler{Client: *dbClient, Validate: validate}
	retrieveHandler := vault.RetrieveHandler{Client: *dbClient}

	router := gin.Default()

	router.GET("/healthcheck", healthcheckHandler)

	router.POST("/save", saveHandler.ServeHTTP)

	retrieve := router.Group("/retrieve")
	{
		retrieve.GET("/all", retrieveHandler.GetAll)
		retrieve.GET("/:id", retrieveHandler.GetByID)
	}

	router.NoRoute(notFoundHandler)
	router.NoMethod(notMethodHandler)

	err = router.Run("localhost:8080")
	if err != nil {
		return
	}
}
