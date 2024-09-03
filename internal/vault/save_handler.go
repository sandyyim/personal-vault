package vault

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"net/http"
	"personal-vault/internal/db"

	"github.com/gin-gonic/gin"
)

const errorMessage = "there is error"

type SaveHandler struct {
	Client   db.DynamoDBClient
	Validate *validator.Validate
}

type Request struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Password    string `json:"password" validate:"required"`
}

func (h SaveHandler) ServeHTTP(c *gin.Context) {
	slog.Info("enter save")

	var request Request

	// call BindJSON to bind the received JSON to request
	if err := c.BindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, errorMessage)
		return
	}

	err := h.Validate.Struct(request)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, errorMessage)
		return
	}

	id := uuid.NewString()

	// save into DB
	vaultEntity := db.VaultEntity{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		Password:    request.Password,
	}

	err = h.Client.InsertItem(c, vaultEntity)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}

	response := fmt.Sprintf("path: %s", id)

	// return success/error
	c.IndentedJSON(http.StatusCreated, response)
}
