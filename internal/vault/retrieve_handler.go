package vault

import (
	"log"
	"log/slog"
	"net/http"
	"personal-vault/internal/db"

	"github.com/gin-gonic/gin"
)

type RetrieveHandler struct {
	Client db.DynamoDBClient
}

func (h RetrieveHandler) ServeHTTP(c *gin.Context) {
	slog.Info("enter retrieve")

	items, err := h.Client.ScanItems(c)
	if err != nil {
		log.Println(err)
		return
	}

	// return success/error
	c.IndentedJSON(http.StatusOK, items)
}
