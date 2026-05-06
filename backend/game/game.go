package game

import (
	"fmt"
    "image/color"

	"github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"

    "transcendance/backend/player"
)

type Game struct {
    player *player.Player
}

func NewGame(p *player.Player) *Game {
	return &Game{player: p}
}

// Update logic runs every tick (1/60 second by default)
func (g *Game) Update() error {
    // Handle arrow key input
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
    return nil
}

// Draw renders the screen every frame
func (g *Game) Draw(screen *ebiten.Image) {
    // Clear the screen with a dark background
    screen.Fill(color.RGBA{0, 0, 0, 255})

    // Prepare a debug string with the player's coordinates
    debugStr := fmt.Sprintf("Player Position: (%d, %d)", g.player.X, g.player.Y)

    // Draw the debug string on the screen
    ebitenutil.DebugPrint(screen, debugStr)
}

// Layout defines the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    // Set the game window size to 800x600 logical pixels
    return 800, 600
}
