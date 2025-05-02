package handlers

import (
	"net/http"
	"time"
	"whatsapp-mougni-api-go/models"
	"whatsapp-mougni-api-go/whatsmeow"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	waClients *whatsmeow.WhatsAppClientManager
	scheduler *whatsmeow.MessageScheduler
}

func NewMessageHandler(clients *whatsmeow.WhatsAppClientManager, scheduler *whatsmeow.MessageScheduler) *MessageHandler {
	return &MessageHandler{
		waClients: clients,
		scheduler: scheduler,
	}
}

type SendMessageRequest struct {
	To        string    `json:"to" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	MediaType string    `json:"media_type"` // "text", "image", "audio", "document"
	MediaData []byte    `json:"media_data"`
	SendAt    time.Time `json:"send_at"` // Optionnel, pour programmation
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	client := h.waClients.GetClient(user.ID)
	if client == nil || !client.Client.IsConnected() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not connected to WhatsApp"})
		return
	}

	if !req.SendAt.IsZero() {
		// Programmer le message
		h.scheduler.ScheduleMessage(user.ID, req.To, req.Content, req.MediaType, req.MediaData, req.SendAt)
		c.JSON(http.StatusOK, gin.H{"message": "Message scheduled successfully"})
		return
	}

	// Envoi immédiat
	var err error
	if req.MediaType == "text" || req.MediaType == "" {
		err = client.SendMessage(req.To, req.Content)
	} else {
		err = client.SendMedia(req.To, req.Content, req.MediaType, req.MediaData)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}
