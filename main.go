package main

import (
	"log"
	"whatsapp-mougni-api-go/config"
	"whatsapp-mougni-api-go/handlers"
	"whatsapp-mougni-api-go/middleware"
	"whatsapp-mougni-api-go/storage"
	"whatsapp-mougni-api-go/whatsmeow"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Initialiser le stockage utilisateur
	userStore, err := storage.NewUserStore(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize user store: %v", err)
	}

	// Initialiser le gestionnaire de clients WhatsApp
	clientManager, err := whatsmeow.NewWhatsAppClientManager(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize WhatsApp client manager: %v", err)
	}

	// Initialiser le planificateur de messages
	scheduler := whatsmeow.NewMessageScheduler()

	// Configurer Gin
	r := gin.Default()

	// Initialiser les gestionnaires
	userHandler := handlers.NewUserHandler(userStore)
	authHandler := handlers.NewAuthHandler(clientManager) // Passer clientManager ici
	messageHandler := handlers.NewMessageHandler(clientManager, scheduler)

	// Routes publiques
	r.POST("/users", userHandler.CreateUser)

	// Routes protégées
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(userStore))
	{
		protected.POST("/auth/login", authHandler.Login)
		protected.GET("/auth/qr", authHandler.ShowQRPage)
		protected.POST("/messages/send", messageHandler.SendMessage)
	}

	// Lancer le serveur
	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
