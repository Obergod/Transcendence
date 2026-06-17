package game

import (
    "golang.org/x/image/math/fixed"


    "transcendance/internal/player"
    "transcendance/internal/world"
    "transcendance/internal/logger"
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
	return w, localIDs
}

func Reset(g *Game) {
	logger.Infof("Reset du jeu - vague=%d, ticks=%d", g.waveNumber, g.ticks)
	w, IDs := InitGame(g.nbPlayer)
	g.world = w
	g.localIDs = IDs
	g.isGameover = false
	g.ticks = 0
	g.waveNumber = 0
	g.lastShotTicks = make(map[string]int)
	logger.Debugf("Reset terminé, nouvelles localIDs=%v, lastShotTicks vide", g.localIDs)
}
