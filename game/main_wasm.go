package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/hajimehoshi/ebiten/v2"
    "transcendance/internal/game"
    "transcendance/internal/player"
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
    startX, startY, baseHP := 400, 300, 100
    initialPlayer := player.NewPlayer(startX, startY, baseHP)
    gameInstance := game.NewGame(initialPlayer)

    ebiten.SetWindowTitle("Player Movement with Arrow Keys")
    ebiten.SetWindowSize(800, 600)

    if err := ebiten.RunGame(gameInstance); err != nil {
        log.Fatal(err)
    }
}