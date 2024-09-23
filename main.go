package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"personal-vault/internal/configuration"
	"personal-vault/internal/db"
	"personal-vault/internal/handler"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

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
	cfg, err := configuration.LoadConfig()
	if err != nil {
		slog.Error("error", err)
		return
	}

	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return
	}

	svc := dynamodb.NewFromConfig(awsConfig)
	dbClient := db.NewClient(svc)

	validate := validator.New()

	saveHandler := handler.SaveHandler{Client: *dbClient, Validate: validate, Key: cfg.Secret}
	retrieveHandler := handler.RetrieveHandler{Client: *dbClient, Key: cfg.Secret}

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
