package handler

import (
	b64 "encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		slog.Error("error", err)
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}

	c.IndentedJSON(http.StatusOK, items)

}

func (h RetrieveHandler) GetByID(c *gin.Context) {
	slog.Info("enter get by id")

	id := c.Param("id")

	if !isValidUUID(id) {
		slog.Error("error", slog.String("validation error", "invalid id"))
		c.JSON(http.StatusBadRequest, errorMessage)
		return
	}

	item, err := h.Client.GetItem(c, id)
	if err != nil {
		slog.Error("error", err)
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}

	decodedPassword, err := b64.StdEncoding.DecodeString(item)
	if err != nil {
		slog.Error("error", err)
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}

	password, err := decryption.Decrypt(string(decodedPassword), h.Key)
	if err != nil {
		slog.Error("error", err)
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}

	c.String(http.StatusOK, password)
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
