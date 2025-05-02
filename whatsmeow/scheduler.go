package whatsmeow

import (
	"sync"
	"time"
)

type ScheduledMessage struct {
	UserID    int
	To        string
	Content   string
	MediaType string
	MediaData []byte
	SendAt    time.Time
}

type MessageScheduler struct {
	messages []ScheduledMessage
	mu       sync.Mutex
	clients  map[int]*WhatsAppClient
}

func NewMessageScheduler() *MessageScheduler {
	scheduler := &MessageScheduler{
		clients: make(map[int]*WhatsAppClient),
	}
	go scheduler.run()
	return scheduler
}

func (s *MessageScheduler) ScheduleMessage(userID int, to, content, mediaType string, mediaData []byte, sendAt time.Time) {
	s.mu.Lock()
	s.messages = append(s.messages, ScheduledMessage{
		UserID:    userID,
		To:        to,
		Content:   content,
		MediaType: mediaType,
		MediaData: mediaData,
		SendAt:    sendAt,
	})
	s.mu.Unlock()
}

func (s *MessageScheduler) run() {
	for {
		s.mu.Lock()
		now := time.Now()
		for i := 0; i < len(s.messages); i++ {
			msg := s.messages[i]
			if now.After(msg.SendAt) {
				client := s.clients[msg.UserID]
				if client != nil && client.Client.IsConnected() {
					switch msg.MediaType {
					case "text":
						client.SendMessage(msg.To, msg.Content)
					case "image", "audio", "document":
						client.SendMedia(msg.To, msg.Content, msg.MediaType, msg.MediaData)
					}
				}
				// Supprimer le message envoyé
				s.messages = append(s.messages[:i], s.messages[i+1:]...)
				i--
			}
		}
		s.mu.Unlock()
		time.Sleep(time.Second)
	}
}
