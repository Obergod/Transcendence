package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/hajimehoshi/ebiten/v2"

    "transcendance/internal/game"
    "transcendance/internal/player"
    "transcendance/internal/world"
)

func main() {
    // --- Appel HTTP vers le backend (exemple) ---
    resp, err := http.Get("/api/hello")
    if err != nil {
        log.Println("Erreur HTTP:", err)
    } else {
        defer resp.Body.Close()
        var result map[string]string
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            log.Println("Erreur décodage JSON:", err)
        } else {
            log.Printf("Réponse backend: %v", result)
        }
    }

    // --- Initialisation du jeu ---
    w := world.NewWorld()
    localID := "local_player"

    initialPlayer := player.NewPlayer(400, 300, 100)
    w.AddPlayer(localID, initialPlayer)

    extraPlayer := player.NewPlayer(600, 600, 100)
    w.AddPlayer("other", extraPlayer)

    gameInstance := game.NewGame(w, localID)

    ebiten.SetWindowSize(800, 600)
    ebiten.SetWindowTitle("Multiplayer Test")
    if err := ebiten.RunGame(gameInstance); err != nil {
        log.Fatal(err)
    }
}