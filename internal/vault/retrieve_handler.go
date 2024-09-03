package vault

import (
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
	"personal-vault/internal/db"
)

type RetrieveHandler struct {
	Client db.DynamoDBClient
}

func (h RetrieveHandler) GetAll(c *gin.Context) {
	slog.Info("enter get all")

	items, err := h.Client.ScanItems(c)
	if err != nil {
		log.Println(err)
		return
	}

	// return success/error
	c.IndentedJSON(http.StatusOK, items)

}

func (h RetrieveHandler) GetByID(c *gin.Context) {
	slog.Info("enter get by id")

	id := c.Param("id")

	item, err := h.Client.GetItem(c, id)
	if err != nil {
		log.Println(err)
		return
	}

	// return success/error
	c.IndentedJSON(http.StatusOK, item)
}
