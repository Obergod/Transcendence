package ws

import (
	"encoding/json"
	"log"
	"transcendance/internal/models"

	"gorm.io/gorm"
)

// Le format du message qu'on va s'échanger avec React
type MessagePayload struct {
	Type       string `json:"type"`        // ex: "chat"
	TargetID   uint   `json:"target_id"`   // L'ID de l'ami à qui on parle
	Content    string `json:"content"`     // Le texte du message
	SenderID   uint   `json:"sender_id"`   // Rempli par le Hub (sécurité)
	SenderName string `json:"sender_name"` // Rempli par le Hub (pour l'affichage)
}

type Hub struct {
	clients    map[uint]*Client
	directMsg  chan MessagePayload // <-- Remplace l'ancien "broadcast"
	register   chan *Client
	unregister chan *Client
	db         *gorm.DB            // <-- Le Hub a accès à la BDD
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
			// 1. Trouver le pseudo de l'expéditeur pour l'affichage React
			var sender models.User
			if err := h.db.First(&sender, payload.SenderID).Error; err == nil {
				payload.SenderName = sender.Username
			} else {
				payload.SenderName = "Inconnu"
			}

			// 2. Sauvegarder le message en Base de Données
			dm := models.DirectMessage{
				SenderID:   payload.SenderID,
				ReceiverID: payload.TargetID,
				Content:    payload.Content,
			}
			h.db.Create(&dm)

			// 3. Livrer le message (Conversion en JSON)
			jsonMsg, _ := json.Marshal(payload)

			// 3a. On l'envoie au destinataire s'il est en ligne !
			if targetClient, ok := h.clients[payload.TargetID]; ok {
				targetClient.send <- jsonMsg
			}

			// 3b. On le renvoie aussi à l'expéditeur pour confirmer l'envoi
			if senderClient, ok := h.clients[payload.SenderID]; ok {
				senderClient.send <- jsonMsg
			}
		}
	}
}