package handler

import (
	b64 "encoding/base64"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
	"personal-vault/internal/db"
	"personal-vault/internal/decryption"
)

type RetrieveHandler struct {
	Client db.DynamoDBClient
	Key    string
}

func (h RetrieveHandler) GetAll(c *gin.Context) {
	slog.Info("enter get all")

	items, err := h.Client.ScanItems(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}

	c.IndentedJSON(http.StatusOK, items)

}

func (h RetrieveHandler) GetByID(c *gin.Context) {
	slog.Info("enter get by id")

	id := c.Param("id")

	item, err := h.Client.GetItem(c, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}

	decodedPassword, err := b64.StdEncoding.DecodeString(item)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}

	password := decryption.Decrypt(string(decodedPassword), h.Key)

	c.IndentedJSON(http.StatusOK, password)
}
