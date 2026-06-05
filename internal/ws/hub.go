package ws

import (
	"encoding/json"
	"log"
	"transcendance/internal/models"

	"gorm.io/gorm"
)

// Le format du message qu'on va s'échanger avec React
type MessagePayload struct {
	Type       string `json:"type"`
	TargetID   uint   `json:"target_id"`
	Content    string `json:"content"`
	SenderID   uint   `json:"sender_id"`
	SenderName string `json:"sender_name"`
}

type Hub struct {
	clients    map[uint]*Client
	directMsg  chan MessagePayload
	register   chan *Client
	unregister chan *Client
	db         *gorm.DB
}

func NewHub(db *gorm.DB) *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		directMsg:  make(chan MessagePayload),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		db:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if oldClient, ok := h.clients[client.UserID]; ok {
				close(oldClient.send)
			}
			h.clients[client.UserID] = client
			log.Printf("🟢 Joueur %d connecté au Hub", client.UserID)

		case client := <-h.unregister:
			if currentClient, ok := h.clients[client.UserID]; ok && currentClient == client {
				delete(h.clients, client.UserID)
				close(client.send)
				log.Printf("🔴 Joueur %d déconnecté du Hub", client.UserID)
			}

		case payload := <-h.directMsg:
			var sender models.User
			if err := h.db.First(&sender, payload.SenderID).Error; err == nil {
				payload.SenderName = sender.Username
			} else {
				payload.SenderName = "Inconnu"
			}

			dm := models.DirectMessage{
				SenderID:   payload.SenderID,
				ReceiverID: payload.TargetID,
				Content:    payload.Content,
			}
			h.db.Create(&dm)

			jsonMsg, _ := json.Marshal(payload)

			if targetClient, ok := h.clients[payload.TargetID]; ok {
				targetClient.send <- jsonMsg
			}

			if senderClient, ok := h.clients[payload.SenderID]; ok {
				senderClient.send <- jsonMsg
			}
		}
	}
}
