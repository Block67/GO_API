package models

// Message représente un message WhatsApp
type Message struct {
	To      string `json:"to" binding:"required"`
	Content string `json:"message" binding:"required"`
}
