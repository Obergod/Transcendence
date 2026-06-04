package main

import (
    "encoding/json"
    "log"
    "net/http"
	"syscall/js"

    "github.com/hajimehoshi/ebiten/v2"

    "transcendance/internal/game"
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

	// initialaze all variables for a game instance
	gameMode := js.Global().Get("gameMode").Int()
	if gameMode == 0 {
		gameMode = 1
	}

	w, initialPlayer := game.InitGame(gameMode)

	gameInstance := game.NewGame(w, initialPlayer.ID)

	// mode vite fait le reset pour que quand on retry ca nous garde dans le bon mode
	js.Global().Set("restartGame", js.FuncOf(func(this js.Value, args []js.Value) any {
		if gameInstance != nil {
			game.Reset(gameInstance, gameMode)
		}
		return nil
	}))

	ebiten.SetTPS(60) // force 60 ticks par seconde
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Multiplayer with Ranged Enemies")
	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}
