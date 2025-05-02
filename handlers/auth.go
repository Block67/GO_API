package handlers

import (
	"encoding/base64"
	"net/http"
	"time"
	"whatsapp-mougni-api-go/models"
	"whatsapp-mougni-api-go/whatsmeow"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type AuthHandler struct {
	waClients     map[int]*whatsmeow.WhatsAppClient
	clientManager *whatsmeow.WhatsAppClientManager
}

func NewAuthHandler(clientManager *whatsmeow.WhatsAppClientManager) *AuthHandler {
	return &AuthHandler{
		waClients:     make(map[int]*whatsmeow.WhatsAppClient),
		clientManager: clientManager,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	client, exists := h.waClients[user.ID]
	if !exists {
		var err error
		client, err = h.clientManager.NewClientForUser(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.waClients[user.ID] = client
	}

	if client.Client.IsConnected() {
		c.JSON(http.StatusOK, gin.H{"message": "Already connected"})
		return
	}

	err := client.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	qrCode, err := client.WaitForQRCode(5 * time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	qrImage, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code image"})
		return
	}
	qrBase64 := base64.StdEncoding.EncodeToString(qrImage)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Scan the QR code",
		"qr_code":      qrCode,
		"qr_image_b64": qrBase64,
	})
}

func (h *AuthHandler) ShowQRPage(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	client, exists := h.waClients[user.ID]
	if !exists || client.GetQRCode() == "" {
		c.String(http.StatusNotFound, "No QR code available")
		return
	}

	qrImage, _ := qrcode.Encode(client.GetQRCode(), qrcode.Medium, 256)
	qrBase64 := base64.StdEncoding.EncodeToString(qrImage)

	html := `
    <html>
    <head><title>WhatsApp QR Code</title></head>
    <body style="text-align: center; font-family: sans-serif;">
        <h2>Scan the QR code with WhatsApp</h2>
        <img src="data:image/png;base64,` + qrBase64 + `" />
    </body>
    </html>
    `

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
