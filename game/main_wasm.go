package main

import (
    "encoding/json"
    "log"
    "net/http"
    "fmt"
    "math"

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

    centerX, centerY := 400, 300
    radius := 250.0
    numEnemies := 30

    for i := 0; i < numEnemies; i++ {
        angle := 2 * math.Pi * float64(i) / float64(numEnemies)
        x := centerX + int(math.Round(radius*math.Cos(angle)))
        y := centerY + int(math.Round(radius*math.Sin(angle)))
        s := fmt.Sprintf("enemy_%d", i)
        w.AddEnemy(s, enemy.NewRanged(fixed.I(x), fixed.I(y), s))
    }

    gameInstance := game.NewGame(w, initialPlayer.ID)

    ebiten.SetTPS(60) // force 60 ticks par seconde
    ebiten.SetWindowSize(800, 600)
    ebiten.SetWindowTitle("Multiplayer with Ranged Enemies")
    if err := ebiten.RunGame(gameInstance); err != nil {
        log.Fatal(err)
    }
}