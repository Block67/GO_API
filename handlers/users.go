package handlers

import (
	"net/http"
	"whatsapp-mougni-api-go/models"
	"whatsapp-mougni-api-go/storage"
	"whatsapp-mougni-api-go/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userStore *storage.UserStore
}

func NewUserHandler(userStore *storage.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	apiKey, err := utils.GenerateAPIKey(13)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}
	user.APIKey = apiKey

	if err := h.userStore.CreateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"api_key": user.APIKey,
	})
}
