package whatsmeow

import (
	"context"
	"fmt"

	_ "modernc.org/sqlite"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type WhatsAppClient struct {
	Client *whatsmeow.Client
}

// Initialiser le client WhatsApp
func NewWhatsAppClient() (*WhatsAppClient, error) {
	// Logger
	logger := waLog.Stdout("WhatsApp", "DEBUG", true)

	// Base de données SQLite pour les sessions
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite", "file:whatsapp_sessions.db?_pragma=foreign_keys(1)", dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Récupération ou création d’un appareil
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %v", err)
	}

	// Initialisation du client
	client := whatsmeow.NewClient(deviceStore, logger)
	client.AddEventHandler(eventHandler)

	return &WhatsAppClient{Client: client}, nil
}

// Gestion des événements
func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.QR:
		fmt.Printf("QR Code: %s\n", v.Codes[0])
	case *events.Connected:
		fmt.Println("Connected to WhatsApp!")
	case *events.Disconnected:
		fmt.Println("Disconnected from WhatsApp.")
	}
}

// Connexion
func (wc *WhatsAppClient) Connect() error {
	if wc.Client.IsConnected() {
		return nil
	}
	return wc.Client.Connect()
}

// Envoi de message texte
func (wc *WhatsAppClient) SendMessage(to, message string) error {
	// Parse le JID
	recipient, err := types.ParseJID(to)
	if err != nil {
		return fmt.Errorf("invalid JID: %v", err)
	}

	// Crée le message texte
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	// Envoie du message (note : seulement 3 paramètres requis)
	resp, err := wc.Client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	fmt.Printf("Message sent, ID: %s\n", resp.ID)
	return nil
}
