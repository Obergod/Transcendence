package game

import (
    "image/color"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"
    "golang.org/x/image/math/fixed"

    "transcendance/internal/world"
    "transcendance/internal/utils"
)

type Game struct {
    world   *world.World
    localID string
}

func NewGame(w *world.World, ID string) *Game {
    return &Game{world: w, localID: ID}
}

// Update logic runs every tick (1/60 second by default)
func (g *Game) Update() error {
    dx := fixed.Int26_6(0)
    dy := fixed.Int26_6(0)

    if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
        dy = -1
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
        dy = 1
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        dx = -1
    }
    if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        dx = 1
    }

    g.world.Lock()
    defer g.world.Unlock()

    localPlayer, exists := g.world.Players[g.localID] 
    // seul le joueur local bouge quand on input des touches
    if !exists {
        return nil
    }

    // Application de la vitesse (propre au joueur)

    var moveX, moveY fixed.Int26_6
    
    if (dx != 0 && dy != 0) {
        // vitesse diagonale = speed * 707 / 1000 (707 ≈ 1000/√2)
        moveX = fixed.Int26_6(int64(dx) * int64(localPlayer.Speed) * 707 / 1000)
        moveY = fixed.Int26_6(int64(dy) * int64(localPlayer.Speed) * 707 / 1000)
    } else { // else sur la meme ligne que {} sinon ca compile pas (hihi le Go)
        moveX = fixed.Int26_6(int64(dx) * int64(localPlayer.Speed))
        moveY = fixed.Int26_6(int64(dy) * int64(localPlayer.Speed))
    }

    localPlayer.X += moveX
    localPlayer.Y += moveY

    // Limites de l'écran
    minX := fixed.I(0)
    maxX := fixed.I(800)
    minY := fixed.I(0)
    maxY := fixed.I(600)

    if localPlayer.X < minX {
        localPlayer.X = minX
    }
    if localPlayer.X > maxX {
        localPlayer.X = maxX
    }
    if localPlayer.Y < minY {
        localPlayer.Y = minY
    }
    if localPlayer.Y > maxY {
        localPlayer.Y = maxY
    }

    return nil
}

func (g *Game) DrawPlayers(screen *ebiten.Image) {
    for id, p := range g.world.Players {
        x := utils.FixedToFloat(p.X)
        y := utils.FixedToFloat(p.Y)

        var col color.Color
        if id == g.localID {
            col = color.RGBA{0, 255, 0, 255} // Vert pour le joueur local
        } else {
            col = color.RGBA{255, 0, 0, 255} // Rouge pour les autres
        }

        // Dessiner un carré de 40x40 centré sur (x, y)
        vector.FillRect(screen, x-20, y-20, 40, 40, col, false)
    }
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.Black)

    g.world.RLock()
    defer g.world.RUnlock()

    g.DrawPlayers(screen)
}

// Layout defines the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 800, 600
}