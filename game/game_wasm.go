package main

import (
    "encoding/json"
    "github.com/gorilla/websocket"
    "github.com/hajimehoshi/ebiten/v2"
    "transcendance/internal/protocol"
    // plus d'imports
)

var wsConn *websocket.Conn
var currentState protocol.StateData

func main() {
    // Connexion WebSocket au serveur
    conn, _, _ := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
    wsConn = conn
    go readMessages()

    // Initialisation Ebiten (pas de logique de jeu propre)
    ebiten.SetWindowTitle("Bullet Heaven Multi")
    ebiten.SetWindowSize(800, 600)
    if err := ebiten.RunGame(&Game{}); err != nil { panic(err) }
}

func readMessages() {
    for {
        var msg protocol.Message
        if err := wsConn.ReadJSON(&msg); err != nil { break }
        if msg.Type == protocol.MsgState {
            // met à jour currentState
        }
    }
}