package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
    "transcendance/internal/player"
    "transcendance/internal/protocol"
    "transcendance/internal/world"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var gameWorld = world.NewWorld()
var idCounter = 0

func generateID() string {
    idCounter++
    return fmt.Sprintf("player-%d", idCounter)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Erreur WebSocket: %v", err)
        return
    }
    defer conn.Close()

    playerID := generateID()
    gameWorld.Lock()
    gameWorld.Players[playerID] = player.NewPlayer(400, 300, 100)
    gameWorld.Unlock()
    log.Printf("Joueur %s connecté", playerID)

    for {
        var msg protocol.Message
        if err := conn.ReadJSON(&msg); err != nil {
            break
        }
        switch msg.Type {
        case protocol.MsgMove:
            var data protocol.MoveData
            b, _ := json.Marshal(msg.Data)
            if err := json.Unmarshal(b, &data); err != nil {
                continue
            }
            gameWorld.Lock()
            if p, ok := gameWorld.Players[playerID]; ok {
                p.X += data.X
                p.Y += data.Y
                // Clamping simple
                if p.X < 0 { p.X = 0 }
                if p.X > 800 { p.X = 800 }
                if p.Y < 0 { p.Y = 0 }
                if p.Y > 600 { p.Y = 600 }
            }
            gameWorld.Unlock()
        case protocol.MsgShoot:
            // À implémenter plus tard
            log.Printf("Tir de %s", playerID)
        }
    }

    gameWorld.Lock()
    delete(gameWorld.Players, playerID)
    gameWorld.Unlock()
    log.Printf("Joueur %s déconnecté", playerID)
}
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Notre route API qui ne peut pas être bloquée
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		// Ce message va s'afficher dans ton terminal de gauche !
		fmt.Println("➡️ BINGO ! Requête reçue sur /api/hello depuis React !")

		w.Header().Set("Content-Type", "application/json")
		// On ajoute les autorisations de sécurité (CORS) au cas où
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]string{"message": "Backend fonctionne"})
	})

	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", fs)

	log.Println("Serveur sur http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}