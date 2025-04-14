package handlers

import (
	"net/http"

	"whatsapp-mougni-api-go/whatsmeow"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	waClient *whatsmeow.WhatsAppClient
}

func NewMessageHandler(waClient *whatsmeow.WhatsAppClient) *MessageHandler {
	return &MessageHandler{waClient: waClient}
}

type SendMessageRequest struct {
	To      string `json:"to" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if !h.waClient.Client.IsConnected() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not connected to WhatsApp"})
		return
	}

	err := h.waClient.SendMessage(req.To, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}
