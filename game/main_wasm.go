package main

import (
    "log"
	"syscall/js"

    "github.com/hajimehoshi/ebiten/v2"

    "transcendance/internal/game"
)

func main() {
	// initialaze all variables for a game instance
	nbPlayer := js.Global().Get("gameMode").Int()
	if nbPlayer == 0 {
		nbPlayer  = 1
	}

	w, IDs := game.InitGame(nbPlayer)

	gameInstance := game.NewGame(w, IDs, nbPlayer)

	// mode vite fait le reset pour que quand on retry ca nous garde dans le bon mode
	js.Global().Set("restartGame", js.FuncOf(func(this js.Value, args []js.Value) any {
		if gameInstance != nil {
			game.Reset(gameInstance)
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
