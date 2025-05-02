package whatsmeow

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	Client   *whatsmeow.Client
	QRCode   string
	QRCodeMu sync.Mutex
	QRChan   chan string
	UserID   int
}

type WhatsAppClientManager struct {
	clients map[int]*WhatsAppClient
	mu      sync.Mutex
	store   *sqlstore.Container
}

func NewWhatsAppClientManager(dbPath string) (*WhatsAppClientManager, error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite", fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", dbPath), dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	return &WhatsAppClientManager{
		clients: make(map[int]*WhatsAppClient),
		store:   container,
	}, nil
}

func (m *WhatsAppClientManager) NewClientForUser(userID int) (*WhatsAppClient, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, exists := m.clients[userID]; exists {
		return client, nil
	}

	logger := waLog.Stdout(fmt.Sprintf("WhatsApp-User-%d", userID), "DEBUG", true)

	// Utilise NewDevice() sans arguments avec la version actuelle de whatsmeow
	device := m.store.NewDevice()

	client := whatsmeow.NewClient(device, logger)

	wc := &WhatsAppClient{
		Client: client,
		QRChan: make(chan string, 1),
		UserID: userID,
	}
	client.AddEventHandler(wc.eventHandler)

	m.clients[userID] = wc
	return wc, nil
}

func (m *WhatsAppClientManager) GetClient(userID int) *WhatsAppClient {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clients[userID]
}

func (wc *WhatsAppClient) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.QR:
		wc.QRCodeMu.Lock()
		wc.QRCode = v.Codes[0] // Stocke le premier code QR
		wc.QRCodeMu.Unlock()
		fmt.Printf("QR Code for user %d: %s\n", wc.UserID, v.Codes[0])
		select {
		case wc.QRChan <- v.Codes[0]: // Envoie le code QR au canal
		default:
			// Évite le blocage si le canal est déjà rempli
		}
	case *events.Connected:
		fmt.Printf("User %d connected to WhatsApp!\n", wc.UserID)
	case *events.Disconnected:
		fmt.Printf("User %d disconnected from WhatsApp.\n", wc.UserID)
	}
}

func (wc *WhatsAppClient) Connect() error {
	if wc.Client.IsConnected() {
		return nil
	}
	return wc.Client.Connect()
}

func (wc *WhatsAppClient) GetQRCode() string {
	wc.QRCodeMu.Lock()
	defer wc.QRCodeMu.Unlock()
	return wc.QRCode
}

func (wc *WhatsAppClient) WaitForQRCode(timeout time.Duration) (string, error) {
	select {
	case qrCode := <-wc.QRChan:
		return qrCode, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timeout waiting for QR code")
	}
}

func (wc *WhatsAppClient) SendMessage(to, message string) error {
	recipient, err := types.ParseJID(to)
	if err != nil {
		return fmt.Errorf("invalid JID: %v", err)
	}

	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	resp, err := wc.Client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	fmt.Printf("Message sent for user %d, ID: %s\n", wc.UserID, resp.ID)
	return nil
}

func (wc *WhatsAppClient) SendMedia(to, caption, mediaType string, mediaData []byte) error {
	recipient, err := types.ParseJID(to)
	if err != nil {
		return fmt.Errorf("invalid JID: %v", err)
	}

	var msg *waProto.Message
	switch mediaType {
	case "image":
		msg = &waProto.Message{
			ImageMessage: &waProto.ImageMessage{
				Caption:       proto.String(caption),
				JPEGThumbnail: mediaData, // Champ corrigé
			},
		}
	case "audio":
		msg = &waProto.Message{
			AudioMessage: &waProto.AudioMessage{
				Mimetype: proto.String("audio/ogg; codecs=opus"), // Champ corrigé
			},
		}
	case "document":
		msg = &waProto.Message{
			DocumentMessage: &waProto.DocumentMessage{
				Title:    proto.String(caption),
				FileName: proto.String(caption),
			},
		}
	default:
		return fmt.Errorf("unsupported media type")
	}

	_, err = wc.Client.SendMessage(context.Background(), recipient, msg)
	if err != nil {
		return fmt.Errorf("failed to send media: %v", err)
	}
	return nil
}
