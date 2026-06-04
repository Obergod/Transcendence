package game

import (
    "fmt"
    "math"

    "golang.org/x/image/math/fixed"

    "transcendance/internal/enemy"
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
	return w, initialPlayer
}

func Reset(g *Game) {
	w, p := InitGame()

	g.world = w
	g.localID = p.ID
	g.isGameover = false
	g.ticks = 0
}
