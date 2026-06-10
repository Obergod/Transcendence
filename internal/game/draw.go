package game

import (
    "image/color"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"

    "transcendance/internal/utils"
)

func (g *Game) DrawBullets(screen *ebiten.Image) {
    for _, b := range g.world.Bullets {
        x := utils.FixedToFloat32(b.X)
        y := utils.FixedToFloat32(b.Y)
        col := color.RGBA{255, 255, 0, 255}
        vector.FillCircle(screen, x, y, 3, col, true)
    }
}

func (g *Game) DrawEnemies(screen *ebiten.Image) {
    for _, e := range g.world.Enemies {
        x := utils.FixedToFloat32(e.Hitbox.X)
        y := utils.FixedToFloat32(e.Hitbox.Y)
        col := color.RGBA{255, 0, 0, 255}
        vector.FillCircle(screen, x, y, 10, col, true)
    }
}

func (g *Game) DrawPlayers(screen *ebiten.Image) {
    for id, p := range g.world.Players {
        x := utils.FixedToFloat32(p.Hitbox.X)
        y := utils.FixedToFloat32(p.Hitbox.Y)
        var col color.Color
        if !p.IsAlive {
            col = color.RGBA{255, 255, 255, 255}
        } else if id == "p1" {
            col = color.RGBA{0, 255, 0, 255}
        } else {
            col = color.RGBA{0, 0, 255, 255}
        }
        vector.FillCircle(screen, x, y, 20, col, true)
    }
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.Black)
    g.world.RLock()
    defer g.world.RUnlock()
    g.DrawPlayers(screen)
    g.DrawEnemies(screen)
    g.DrawBullets(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 800, 600
}
