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
	UserID uint
}

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 2048
)

var (
	newline = []byte{'\n'}
	space	= []byte{' '}
)
//need to be changed for deployment
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { //c le CORS cross origine ressource sharing
		return true // normalement ca c'est pas good pour un vrais site en entreprise mais etant donner que le projet va etre compiler sur des pc differrent on laisse ca pour les test et la correction fonctionne bien
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

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMsgSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
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

		var payload MessagePayload
		if err := json.Unmarshal(message, &payload); err == nil {

			if payload.Type == "chat" {

				runes := []rune(payload.Content)
				if len(runes) > 300 {
					payload.Content = string(runes[:300])
				}

				payload.SenderID = c.UserID

				c.hub.directMsg <- payload
			}
		}
	}
}

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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
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

func ServeWs(hub *Hub, c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "Token manquant"})
		return
	}

	userID, err := auth.ValidateToken(token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "Token invalide"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Erreur d'upgrade WS:", err)
		return
	}

	client := NewClient(hub, conn, userID)
	hub.register <- client

	go client.writePump()
	go client.readPump()
}
