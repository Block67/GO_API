package handlers

import (
	"net/http"

	"whatsapp-mougni-api-go/whatsmeow"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	waClient *whatsmeow.WhatsAppClient
}

func NewAuthHandler(waClient *whatsmeow.WhatsAppClient) *AuthHandler {
	return &AuthHandler{waClient: waClient}
}

func (h *AuthHandler) Login(c *gin.Context) {
	if h.waClient.Client.IsConnected() {
		c.JSON(http.StatusOK, gin.H{"message": "Already connected"})
		return
	}

	err := h.waClient.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Scan the QR code displayed in the server logs",
	})
}
