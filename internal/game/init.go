package game

import (
    "fmt"
    "math"

    "golang.org/x/image/math/fixed"

    "transcendance/internal/enemy"
    "transcendance/internal/player"
    "transcendance/internal/world"
)

func InitGame(mode int) (*world.World, *player.Player) {
    w := world.NewWorld()

    // joueur 1
    initialPlayer := player.NewPlayer(fixed.I(400), fixed.I(300), 100, "local")
    w.AddPlayer(initialPlayer.ID, initialPlayer)

    // jia add un if tout con si on appuie bien sur 2 joueur
    if mode == 2 {
        extraPlayer := player.NewPlayer(fixed.I(450), fixed.I(300), 100, "extra")
        w.AddPlayer(extraPlayer.ID, extraPlayer)
    }

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

func Reset(g *Game, mode int) {
	w, p := InitGame(mode)

	g.world = w
	g.localID = p.ID
	g.isGameover = false
	g.ticks = 0
}