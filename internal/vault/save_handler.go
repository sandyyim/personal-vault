package vault

import (
	"log"
	"personal-vault/internal/db"

	"github.com/gin-gonic/gin"
)

type SaveHandler struct {
	Client db.DynamoDBClient
}

type VaultRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h SaveHandler) ServeHTTP(c *gin.Context) {
	log.Println("enter save")

	// validation

	// save into DB

	// return success/error
}
