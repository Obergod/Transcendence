package ws

import (
	//"bytes"
	"log"
	"net/http"
	"time"

	"encoding/json"
	"github.com/gorilla/websocket"
    "github.com/gin-gonic/gin"
	"transcendance/internal/auth"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	UserID uint // <-- NOUVEAU : Le client connait son identité
}

const (
	// Time allowed to write a message to peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// send pings to peer with this persiod must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Max msg size allowed from peer (change for the json gameInfo size ?)
	maxMsgSize = 512
)

//	usefull for chat msg, readapt for game communication
var (
	newline = []byte{'\n'}
	space	= []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CORRECTION : On n'autorise que ton frontend React !
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:5173"
	},
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uint) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		UserID: userID,
	}
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait)) //launch "timer" before client considered afk
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait)) // reset timer when pong sent
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// NOUVEAU : On décode le JSON envoyé par React
		var payload MessagePayload
		if err := json.Unmarshal(message, &payload); err == nil {

			// Si c'est un message de chat :
			if payload.Type == "chat" {
				// SÉCURITÉ ABSOLUE : Même si le frontend a essayé de tricher sur l'expéditeur,
				// on écrase le SenderID avec le VRAI ID du client connecté !
				payload.SenderID = c.UserID

				// On l'envoie au cerveau (le Hub)
				c.hub.directMsg <- payload
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.[27;5;106~
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				//hub closed channell
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// maybe not necessary if no batching (just writeMessage)
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
					return
				}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

	//serveWs handles websocket request from peer
func ServeWs(hub *Hub, c *gin.Context) {
	// 1. On récupère le token JWT passé dans l'URL (ex: ws://localhost:8081/ws?token=XYZ)
	token := c.Query("token")
	if token == "" {
		c.JSON(401, gin.H{"error": "Token manquant"})
		return
	}

	// 2. On vérifie qui est ce joueur
	userID, err := auth.ValidateToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "Token invalide"})
		return
	}

	// 3. On accepte la connexion WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Erreur d'upgrade WS:", err)
		return
	}

	// 4. On enregistre le client avec son ID !
	client := NewClient(hub, conn, userID)
	hub.register <- client

	go client.writePump()
	go client.readPump()
}