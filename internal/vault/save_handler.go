package vault

import (
	"github.com/google/uuid"
	"log"
	"net/http"
	"personal-vault/internal/db"

	"github.com/gin-gonic/gin"
)

type SaveHandler struct {
	Client db.DynamoDBClient
}

type Request struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

func (h SaveHandler) ServeHTTP(c *gin.Context) {
	log.Println("enter save")
	var request Request

	// call BindJSON to bind the received JSON to request
	if err := c.BindJSON(&request); err != nil {
		log.Println(err)
		return
	}

	// validation
	// validation

	id := uuid.NewString()

	// save into DB
	vaultEntity := db.VaultEntity{
		ID:     id,
		Name:   request.Name,
		Secret: request.Secret,
	}

	err := h.Client.InsertItem(c, vaultEntity)
	if err != nil {
		log.Println(err)
		return
	}

	// return success/error
	c.IndentedJSON(http.StatusCreated, id)
}
