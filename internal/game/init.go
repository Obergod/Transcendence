package game

import (
    "fmt"
    "math"

    "golang.org/x/image/math/fixed"

    "transcendance/internal/enemy"
    "transcendance/internal/player"
    "transcendance/internal/world"
)

func InitGame(nbPlayer int) (*world.World, []string) {
    w := world.NewWorld()
	var localIDs []string

    // joueur 1
    p1 := player.NewPlayer(fixed.I(600), fixed.I(360), 100, "p1")
    w.AddPlayer(p1.ID, p1)
	localIDs = append(localIDs, p1.ID)

    // jia add un if tout con si on appuie bien sur 2 joueur
    if nbPlayer == 2 {
		p2 := player.NewPlayer(fixed.I(680), fixed.I(360), 100, "p2")
        w.AddPlayer(p2.ID, p2)
		localIDs = append(localIDs, p2.ID)
	}

    centerX, centerY := 640, 360
    radius := 350.0
    numEnemies := 30

    for i := 0; i < numEnemies; i++ {
        angle := 2 * math.Pi * float64(i) / float64(numEnemies)
        x := centerX + int(math.Round(radius*math.Cos(angle)))
        y := centerY + int(math.Round(radius*math.Sin(angle)))
        s := fmt.Sprintf("enemy_%d", i)
        w.AddEnemy(s, enemy.NewRanged(fixed.I(x), fixed.I(y), s))
    }
	return w, localIDs
}

func Reset(g *Game) {
	w, IDs := InitGame(g.nbPlayer)

	g.world = w
	g.localIDs = IDs 
	g.isGameover = false
	g.ticks = 0

}
