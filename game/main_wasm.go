package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/hajimehoshi/ebiten/v2"
    "golang.org/x/image/math/fixed"

    "transcendance/internal/game"
    "transcendance/internal/player"
    "transcendance/internal/world"
    "transcendance/internal/enemy"
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

    initialPlayer := player.NewPlayer(fixed.I(400), fixed.I(300), 100, "local")
    w.AddPlayer(initialPlayer.ID, initialPlayer)

    extraPlayer := player.NewPlayer(fixed.I(600), fixed.I(600), 100, "extra")
    w.AddPlayer("other", extraPlayer)

    firstEnemy := enemy.NewEnemy(fixed.I(100), fixed.I(100), 100, "fils")
    w.AddEnemy("first", firstEnemy)

    secondEnemy := enemy.NewEnemy(fixed.I(150), fixed.I(100), 100, "de")
    w.AddEnemy("second", secondEnemy)

    thirdEnemy := enemy.NewEnemy(fixed.I(200), fixed.I(100), 100, "con")
    w.AddEnemy("third", thirdEnemy)

    gameInstance := game.NewGame(w, initialPlayer.ID)

    ebiten.SetWindowSize(800, 600)
    ebiten.SetWindowTitle("Multiplayer Test")
    if err := ebiten.RunGame(gameInstance); err != nil {
        log.Fatal(err)
    }
}