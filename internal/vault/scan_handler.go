package vault

import (
	"log"
	"log/slog"
	"net/http"
	"personal-vault/internal/db"

	"github.com/gin-gonic/gin"
)

type ScanHandler struct {
	Client db.DynamoDBClient
}

func (h ScanHandler) ServeHTTP(c *gin.Context) {
	slog.Info("enter scan")

	items, err := h.Client.ScanItems(c)
	if err != nil {
		log.Println(err)
		return
	}

	// return success/error
	c.IndentedJSON(http.StatusOK, items)
}
