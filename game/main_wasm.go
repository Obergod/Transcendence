package main

import (
    "log"

    "github.com/gorilla/websocket"
    "github.com/hajimehoshi/ebiten/v2"
    "transcendance/internal/game"
    "transcendance/internal/player"
    "transcendance/internal/protocol"
)

func main() {
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    var initMsg protocol.Message
    if err := conn.ReadJSON(&initMsg); err != nil {
        log.Fatal(err)
    }
    if initMsg.Type != protocol.MsgInit {
        log.Fatal("expected init message")
    }
    playerID, ok := initMsg.Data.(string)
    if !ok {
        log.Fatal("invalid init data")
    }
    log.Printf("Connected as %s", playerID)

    // Conversion explicite : les coordonnées du joueur sont des entiers (pixels)
    startX, startY, baseHP := 400, 300, 100
    initialPlayer := player.NewPlayer(startX, startY, baseHP)
    gameInstance := game.NewGame(initialPlayer)
    gameInstance.SetWebSocket(conn, playerID)

    go func() {
        for {
            var msg protocol.Message
            if err := conn.ReadJSON(&msg); err != nil {
                log.Println("read error:", err)
                return
            }
            if msg.Type == protocol.MsgState {
                state, ok := msg.Data.(protocol.StateData)
                if !ok {
                    continue
                }
                gameInstance.UpdateWorld(state)
            }
        }
    }()

    ebiten.SetWindowTitle("Bullet Heaven Multiplayer")
    ebiten.SetWindowSize(800, 600)
    if err := ebiten.RunGame(gameInstance); err != nil {
        log.Fatal(err)
    }
}