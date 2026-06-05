package game

import (
    "golang.org/x/image/math/fixed"

    "transcendance/internal/player"
    "transcendance/internal/world"
)

func InitGame() (*world.World, *player.Player) {
    // --- Initialisation du monde ---
    w := world.NewWorld()

    // Joueur local
    initialPlayer := player.NewPlayer(fixed.I(400), fixed.I(300), 100, "local")
    w.AddPlayer(initialPlayer.ID, initialPlayer)

    // Second joueur (autre)
    extraPlayer := player.NewPlayer(fixed.I(600), fixed.I(600), 100, "extra")
    w.AddPlayer(extraPlayer.ID, extraPlayer)

	return w, initialPlayer
}

func Reset(g *Game) {
	w, p := InitGame()

	g.world = w
	g.localID = p.ID
	g.isGameover = false
}
