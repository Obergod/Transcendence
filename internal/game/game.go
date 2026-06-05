package game

import (
	"syscall/js"
	"fmt"


    "transcendance/internal/world"
    "transcendance/internal/enemy"
)

type Game struct {
    world         *world.World
    localIDs      []string
    waveNumber    int
    isGameover    bool
    ticks         int // compteur de tick
    nbPlayer      int
    lastShotTicks map[string]int // cooldown par joueur
}

func NewGame(w *world.World, IDs []string, nb int) *Game {
    return &Game{
        world:         w,
        localIDs:      IDs,
        waveNumber:    1,
        isGameover:    false,
        ticks:         0,
        nbPlayer:      nb,
        lastShotTicks: make(map[string]int),
    }
}



func (g *Game) CheckPlayersAlive() {
    playersdead := 0
    for _, id := range g.localIDs {
        Player, exists := g.world.Players[id]
        if !exists {
            continue
        }
        if !Player.IsAlive {
            playersdead++
        }
    }
    if playersdead == g.nbPlayer && !g.isGameover {
        g.isGameover = true

        durationInSeconds := g.ticks / 60
        score := g.ticks // Le score final est égal au nombre de ticks survécus

		if js.Global().Get("onGameover").Type() == js.TypeFunction {
			js.Global().Call("onGameover", durationInSeconds, score)
		}
	}
}

func (g *Game) RemoveDeadEnemies() {
    for id, e := range g.world.Enemies {
        if !e.IsAlive {
            delete(g.world.Enemies, id)
        }
    }
}

func (g *Game) UpdateScoreTimer() {
    if !g.isGameover {
        g.ticks++

        // INJECTION EN TEMPS RÉEL DANS LE DOM DE REACT
        if g.ticks%6 == 0 {
            jsDoc := js.Global().Get("document")
            if jsDoc.Type() != js.TypeUndefined {
                // 1. Mise à jour du Timer
                timerEl := jsDoc.Call("getElementById", "game-timer")
                if timerEl.Type() != js.TypeNull {
                    seconds := g.ticks / 60
                    minutes := seconds / 60
                    secRemainder := seconds % 60
                    timerEl.Set("innerText", fmt.Sprintf("TEMPS: %02d:%02d", minutes, secRemainder))
                }

                // 2. Mise à jour du Score (accumulateur de ticks)
                scoreEl := jsDoc.Call("getElementById", "game-score")
                if scoreEl.Type() != js.TypeNull {
                    scoreEl.Set("innerText", fmt.Sprintf("SCORE: %05d", g.ticks))
                }
            }
        }
    }
}

func (g *Game) Update() error {
    g.CheckPlayersAlive()
    g.UpdateScoreTimer()
    for _, id := range g.localIDs {
        g.MovePlayer(id)
    }
    g.SpawnEnemies()
    g.MoveEnemies()
    g.HandleEnemyShooting()
    g.HandlePlayersShooting()
    g.UpdateBullets()
    g.HandleBulletCollisions()
    g.RemoveDeadEnemies()
    return nil
}

