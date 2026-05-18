package ws

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

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

var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true }, // Maati: on autorise React a se connecter
}

type Client struct {
	hub *Hub
	conn *websocket.Conn

	//buffered channel of outbound msg
	send chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client {
		hub: hub,
		conn: conn,
		send: make(chan []byte, 256),
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
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
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
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) { // Maati: jai mis le s de server en maj sinon c'est priv
	conn, err := upgrader.Upgrade(w, r, nil) // upgrades to websocket
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(hub, conn)
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
