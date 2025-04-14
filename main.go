package main

import (
	"log"

	"whatsapp-mougni-api-go/handlers"
	"whatsapp-mougni-api-go/whatsmeow"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialiser le client WhatsApp
	waClient, err := whatsmeow.NewWhatsAppClient()
	if err != nil {
		log.Fatalf("Failed to initialize WhatsApp client: %v", err)
	}

	// Configurer Gin
	r := gin.Default()

	// Initialiser les gestionnaires
	authHandler := handlers.NewAuthHandler(waClient)
	messageHandler := handlers.NewMessageHandler(waClient)

	// Définir les routes
	r.POST("/auth/login", authHandler.Login)
	r.POST("/messages/send", messageHandler.SendMessage)

	// Lancer le serveur
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
