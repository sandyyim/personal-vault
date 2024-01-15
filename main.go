package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	router := gin.Default()
	router.GET("/healthcheck", healthcheckHandler)

	// router.POST("/os", func(c *gin.Context) {
	// 	c.String(200, runtime.GOOS)
	// })

	router.NoRoute(notFoundHandler)
	router.NoMethod(notMethodHandler)

	ginLambda = ginadapter.New(router)
	lambda.Start(Handler)
}
