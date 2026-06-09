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

// NOUVEAU : Fonction qui envoie la liste de tous les IDs connectés à tout le monde
func (h *Hub) broadcastOnlineUsers() {
	var onlineIDs []uint
	for id := range h.clients {
		onlineIDs = append(onlineIDs, id)
	}

	msg := map[string]interface{}{
		"type":  "online_users",
		"users": onlineIDs,
	}
	jsonMsg, _ := json.Marshal(msg)

	for _, client := range h.clients {
		client.send <- jsonMsg
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

			// NOUVEAU : On met à jour tout le monde !
			h.broadcastOnlineUsers()

		case client := <-h.unregister:
			if currentClient, ok := h.clients[client.UserID]; ok && currentClient == client {
				delete(h.clients, client.UserID)
				close(client.send)
				log.Printf("🔴 Joueur %d déconnecté du Hub", client.UserID)

				// NOUVEAU : On met à jour tout le monde !
				h.broadcastOnlineUsers()
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