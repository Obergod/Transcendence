package main

import (
    "log"
	"os"
	"syscall/js"

    "github.com/hajimehoshi/ebiten/v2"

    "transcendance/internal/game"
    "transcendance/internal/logger"
)

func main() {
	// Activer les logs si la variable d'environnement est présente
	if os.Getenv("DEBUG") == "1" {
		logger.EnableDebug()
	}

	// initialaze all variables for a game instance
	nbPlayer := js.Global().Get("gameMode").Int()
	if nbPlayer == 0 {
		nbPlayer  = 1
	}
	logger.Infof("Démarrage du jeu avec %d joueur(s)", nbPlayer)

	w, IDs := game.InitGame(nbPlayer)

	gameInstance := game.NewGame(w, IDs, nbPlayer)

	// mode vite fait le reset pour que quand on retry ca nous garde dans le bon mode
	js.Global().Set("restartGame", js.FuncOf(func(this js.Value, args []js.Value) any {
		logger.Infof("Reset demandé depuis JavaScript")
		if gameInstance != nil {
			game.Reset(gameInstance)
		}
		return nil
	}))

	js.Global().Set("quitGame", js.FuncOf(func(this js.Value, args []js.Value) any {
		logger.Infof("Arrêt du moteur Ebitengine demandé par React")
		if gameInstance != nil {
			gameInstance.ShouldQuit = true
		}
		return nil
	}))

	ebiten.SetTPS(60) // force 60 ticks par seconde
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Multiplayer with Ranged Enemies")
	if err := ebiten.RunGame(gameInstance); err != nil {
		logger.Errorf("Erreur lors de l'exécution du jeu: %v", err)
		log.Fatal(err)
	}
}