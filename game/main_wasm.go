package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/hajimehoshi/ebiten/v2"
    "golang.org/x/image/math/fixed"

    "transcendance/internal/enemy"
    "transcendance/internal/game"
    "transcendance/internal/player"
    "transcendance/internal/world"
)

func main() {
    // Appel HTTP exemple (backend)
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

    // --- Initialisation du monde ---
    w := world.NewWorld()

    // Joueur local
    initialPlayer := player.NewPlayer(fixed.I(400), fixed.I(300), 100, "local")
    w.AddPlayer(initialPlayer.ID, initialPlayer)

    // Second joueur (autre)
    extraPlayer := player.NewPlayer(fixed.I(600), fixed.I(600), 100, "extra")
    w.AddPlayer(extraPlayer.ID, extraPlayer)

    // Trois ennemis "ranged"
    firstEnemy := enemy.NewRanged(fixed.I(100), fixed.I(100), "ranger1")
    w.AddEnemy(firstEnemy.ID, firstEnemy)

    secondEnemy := enemy.NewRanged(fixed.I(150), fixed.I(100), "ranger2")
    w.AddEnemy(secondEnemy.ID, secondEnemy)

    thirdEnemy := enemy.NewRanged(fixed.I(200), fixed.I(100), "ranger3")
    w.AddEnemy(thirdEnemy.ID, thirdEnemy)

    gameInstance := game.NewGame(w, initialPlayer.ID)

    ebiten.SetTPS(60) // force 60 ticks par seconde
    ebiten.SetWindowSize(800, 600)
    ebiten.SetWindowTitle("Multiplayer with Ranged Enemies")
    if err := ebiten.RunGame(gameInstance); err != nil {
        log.Fatal(err)
    }
}