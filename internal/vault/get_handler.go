package vault

import (
	"log"
	"log/slog"
	"net/http"
	"personal-vault/internal/db"

	"github.com/gin-gonic/gin"
)

type GetHandler struct {
	Client db.DynamoDBClient
}

type GetRequest struct {
	Id string `json:"id"`
}

func (h GetHandler) ServeHTTP(c *gin.Context) {
	slog.Info("enter get")

	var request GetRequest

	// call BindJSON to bind the received JSON to request
	if err := c.BindJSON(&request); err != nil {
		log.Println(err)
		return
	}

	item, err := h.Client.GetItem(c, request.Id)
	if err != nil {
		log.Println(err)
		return
	}

	// return success/error
	c.IndentedJSON(http.StatusOK, item)
}
