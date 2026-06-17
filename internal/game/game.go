package game

import (
	"syscall/js"
	"fmt"
	"bytes"
	"log"

	"transcendance/internal/world"
	"transcendance/internal/logger"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
)

type Game struct {
	world         *world.World
	localIDs      []string
	waveNumber    int
	isGameover    bool
	ticks         int
	nbPlayer      int
	lastShotTicks map[string]int
	initialized	  bool
	ShouldQuit	  bool

	fontSource *text.GoTextFaceSource
	fontFace *text.GoTextFace
}

func NewGame(w *world.World, IDs []string, nb int) *Game {
	logger.Debugf("Création d'une nouvelle partie, nbPlayer=%d, localIDs=%v", nb, IDs)
	return &Game{
		world:         w,
		localIDs:      IDs,
		waveNumber:    1,
		isGameover:    false,
		ticks:         0,
		nbPlayer:      nb,
		lastShotTicks: make(map[string]int),
		initialized: false,
		ShouldQuit:  false,
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
			logger.Debugf("Joueur %s est mort", id)
		}
	}
	if playersdead == g.nbPlayer && !g.isGameover {
		g.isGameover = true
		durationInSeconds := g.ticks / 60
		score := g.ticks
		logger.Infof("Game Over! Durée=%d sec, Score=%d ticks", durationInSeconds, score)
		if js.Global().Get("onGameover").Type() == js.TypeFunction {
			js.Global().Call("onGameover", durationInSeconds, score)
		}
	}
}

func (g *Game) RemoveDeadEnemies() {
	for id, e := range g.world.Enemies {
		if !e.IsAlive {
			logger.Debugf("Ennemi %s retiré du monde", id)
			delete(g.world.Enemies, id)
		}
	}
}

func (g *Game) UpdateScoreTimer() {
	if !g.isGameover {
		g.ticks++
		if g.ticks%6 == 0 {
			jsDoc := js.Global().Get("document")
			if jsDoc.Type() != js.TypeUndefined {
				timerEl := jsDoc.Call("getElementById", "game-timer")
				if timerEl.Type() != js.TypeNull {
					seconds := g.ticks / 60
					minutes := seconds / 60
					secRemainder := seconds % 60
					timerLabel := js.Global().Get("timerLabel").String()
					timerEl.Set("innerText", fmt.Sprintf("%s: %02d:%02d", timerLabel, minutes, secRemainder))
				}
				scoreEl := jsDoc.Call("getElementById", "game-score")
				if scoreEl.Type() != js.TypeNull {
					scoreLabel := js.Global().Get("scoreLabel").String()
					scoreEl.Set("innerText", fmt.Sprintf("%s: %05d", scoreLabel, g.ticks))
				}
			}
		}
	}
}

func (g *Game) initialize() {
	// Load variable-width font embedded in Ebitengine.
	var err error
	g.fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Panic(err)
	}

	// Create 16px font face from the above source.
	g.fontFace = &text.GoTextFace{
		Source: g.fontSource,
		Size:   16,
	}
	g.initialized = true
}

func (g *Game) Update() error {
	if g.ShouldQuit {
		return ebiten.Termination
	}
	if !g.initialized {
		g.initialize()
	}
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