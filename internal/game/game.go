package game

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
<<<<<<< HEAD
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	//"golang.org/x/image/math/fixed"
	"github.com/gorilla/websocket"
	"transcendance/internal/player"
	"transcendance/internal/protocol"
=======
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"

    "transcendance/internal/player"
	"transcendance/internal/utils"
    "github.com/hajimehoshi/ebiten/v2/vector"
>>>>>>> ee7fd06e13e46e9e34ef9291a815f4dccf0c8c82
)

type Game struct {
	player *player.Player

	wsConn       *websocket.Conn
	playerID     string
	mu           sync.RWMutex
	otherPlayers map[string]struct{ X, Y float64 }
	enemies      []struct{ X, Y float64 }
}

func NewGame(p *player.Player) *Game {
	return &Game{
		player:       p,
		otherPlayers: make(map[string]struct{ X, Y float64 }),
		enemies:      []struct{ X, Y float64 }{},
	}
}

func (g *Game) SetWebSocket(conn *websocket.Conn, id string) {
	g.wsConn = conn
	g.playerID = id
}

func (g *Game) UpdateWorld(state protocol.StateData) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Notre propre position n'est pas modifiée ici (le serveur ne la renvoie pas forcément)
	// On suppose que le serveur envoie les positions des autres joueurs.
	// Adaptez selon votre structure StateData réelle.
	// Par exemple, si StateData a un champ Players (map[string]struct{X,Y float64})...
	// Ici je suppose que state.OtherPlayers existe, mais vous devez corriger selon votre code.
	// Pour l'instant, on commente cette partie pour compiler.
	/*
	g.otherPlayers = make(map[string]struct{ X, Y float64 })
	for id, pos := range state.OtherPlayers {
		if id != g.playerID {
			g.otherPlayers[id] = pos
		}
	}
	g.enemies = make([]struct{ X, Y float64 }, len(state.Enemies))
	for i, e := range state.Enemies {
		g.enemies[i] = struct{ X, Y float64 }{X: e.X, Y: e.Y}
	}
	*/
}

func (g *Game) Update() error {
	// Sauvegarde ancienne position (en float64 pour comparaison)
	oldX := float64(g.player.X) / 64.0  // car fixed.Int26_6 = pixels * 64
	oldY := float64(g.player.Y) / 64.0

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.player.MoveUp()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.player.MoveDown()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.MoveLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.MoveRight()
	}

	// Convertir nouvelle position en float64
	newX := float64(g.player.X) / 64.0
	newY := float64(g.player.Y) / 64.0

	if (newX != oldX || newY != oldY) && g.wsConn != nil {
		msg := protocol.Message{
			Type: protocol.MsgMove,
			Data: protocol.MoveData{
				X: newX,
				Y: newY,
			},
		}
		_ = g.wsConn.WriteJSON(msg)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Clear the screen with a dark background
    screen.Fill(color.Black)

	pX := utils.FixedToFloat(g.player.X)
	pY := utils.FixedToFloat(g.player.Y)

    // Prepare a debug string with the player's coordinates
    debugStr := fmt.Sprintf("Player Position: (%d, %d), HP: %d", pX, pY, g.player.HP)

    // Draw the debug string on the screen
    ebitenutil.DebugPrint(screen, debugStr)
	vector.FillRect(screen, float32(pX), float32(pY),
		float32(32), float32(32), color.RGBA{255, 0, 0 , 255}, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}