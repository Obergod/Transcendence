package game

import (
    "image/color"
    "strconv"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

    "transcendance/internal/utils"
)

func (g *Game) DrawBullets(screen *ebiten.Image) {
    for _, b := range g.world.Bullets {
        x := utils.FixedToFloat32(b.X)
        y := utils.FixedToFloat32(b.Y)
        var col color.Color
        if (b.OwnerIsPlayer) {
            col = color.RGBA{255, 255, 0, 255}
        } else {
            col = color.RGBA{255, 127, 0, 255}
        }
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

func (g *Game) DrawHP(screen *ebiten.Image) {
    for _, p := range g.world.Players {
        x := utils.FixedToFloat64(p.Hitbox.X)
        y := utils.FixedToFloat64(p.Hitbox.Y)
        txtOp := &text.DrawOptions{}
        txtOp.GeoM.Translate(x-10,y-10)
        col := color.RGBA{0, 0, 0, 255}
        txtOp.ColorScale.ScaleWithColor(col)
        text.Draw(screen, strconv.Itoa(p.HP), g.fontFace, txtOp)

    }
}


func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.Black)
    g.world.RLock()
    defer g.world.RUnlock()
    g.DrawPlayers(screen)
    g.DrawEnemies(screen)
    g.DrawBullets(screen)
    g.DrawHP(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 800, 600
}
