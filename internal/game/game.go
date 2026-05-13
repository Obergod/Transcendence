package game

import (
    "image/color"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"
    "golang.org/x/image/math/fixed"

    "transcendance/internal/hitbox"
    "transcendance/internal/player"
    "transcendance/internal/utils"
    "transcendance/internal/world"
)

type Game struct {
    world   *world.World
    localID string
}

func NewGame(w *world.World, ID string) *Game {
    return &Game{
        world:   w,
        localID: ID,
    }
}

// pushOutCollisionFixed déplace uniquement e1 hors de e2 (corCircleion totale).
func pushOutCollisionFixed(e1, e2 *hitbox.Hitbox) bool {
    dx := int64(e2.X - e1.X)
    dy := int64(e2.Y - e1.Y)
    rsum := int64(e1.R + e2.R)
    distSq := dx*dx + dy*dy
    if distSq >= rsum*rsum {
        return false
    }
    if distSq == 0 {
        e1.X += fixed.I(1)
        return true
    }
    dist := int64(utils.FixedSqrt(fixed.Int26_6(distSq)))
    if dist == 0 {
        return false
    }
    overlap := rsum - dist
    // Déplacement total (coefficient 1/3) pour sortir complètement
    moveX := (overlap * dx) / dist / 3
    moveY := (overlap * dy) / dist / 3
    e1.X += fixed.Int26_6(moveX)
    e1.Y += fixed.Int26_6(moveY)
    return true
}

func (g *Game) MovePlayer() error {
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
    if !exists {
        return nil
    }

    var moveX, moveY fixed.Int26_6
    if dx != 0 && dy != 0 {
        moveX = fixed.Int26_6(int64(dx) * int64(localPlayer.Speed) * 707 / 1000)
        moveY = fixed.Int26_6(int64(dy) * int64(localPlayer.Speed) * 707 / 1000)
    } else {
        moveX = fixed.Int26_6(int64(dx) * int64(localPlayer.Speed))
        moveY = fixed.Int26_6(int64(dy) * int64(localPlayer.Speed))
    }

    localPlayer.Hitbox.X += moveX
    localPlayer.Hitbox.Y += moveY

    // Limites écran
    minX := fixed.I(0)
    maxX := fixed.I(800)
    minY := fixed.I(0)
    maxY := fixed.I(600)
    if localPlayer.Hitbox.X < minX {
        localPlayer.Hitbox.X = minX
    }
    if localPlayer.Hitbox.X > maxX {
        localPlayer.Hitbox.X = maxX
    }
    if localPlayer.Hitbox.Y < minY {
        localPlayer.Hitbox.Y = minY
    }
    if localPlayer.Hitbox.Y > maxY {
        localPlayer.Hitbox.Y = maxY
    }
    return nil
}

func (g *Game) MoveEnemies() {
    if len(g.world.Players) == 0 {
        return
    }
    g.world.Lock()
    defer g.world.Unlock()

    for _, e := range g.world.Enemies {
        // Trouver le joueur le plus proche
        var closestPlayer *player.Player
        var closestDistSq int64 = 1 << 62

        for _, p := range g.world.Players {
            dx := int64(p.Hitbox.X - e.Hitbox.X)
            dy := int64(p.Hitbox.Y - e.Hitbox.Y)
            distSq := dx*dx + dy*dy
            if distSq < closestDistSq {
                closestDistSq = distSq
                closestPlayer = p
            }
        }

        if closestPlayer == nil {
            continue
        }

        dx := int64(closestPlayer.Hitbox.X - e.Hitbox.X)
        dy := int64(closestPlayer.Hitbox.Y - e.Hitbox.Y)
        if dx == 0 && dy == 0 {
            continue
        }

        dist := int64(utils.FixedSqrt(fixed.Int26_6(dx*dx + dy*dy)))
        if dist == 0 {
            continue
        }
        speed := int64(e.Speed)
        moveX := fixed.Int26_6((dx * speed) / dist)
        moveY := fixed.Int26_6((dy * speed) / dist)

        e.Hitbox.X += moveX
        e.Hitbox.Y += moveY

        for _, p := range g.world.Players {
            pushOutCollisionFixed(e.Hitbox, p.Hitbox)
        }
        for _, other := range g.world.Enemies {
            if e.ID != other.ID {
                pushOutCollisionFixed(e.Hitbox, other.Hitbox)
            }
        }

        // Limites écran
        minX := fixed.I(0)
        maxX := fixed.I(800)
        minY := fixed.I(0)
        maxY := fixed.I(600)
        if e.Hitbox.X < minX {
            e.Hitbox.X = minX
        }
        if e.Hitbox.X > maxX {
            e.Hitbox.X = maxX
        }
        if e.Hitbox.Y < minY {
            e.Hitbox.Y = minY
        }
        if e.Hitbox.Y > maxY {
            e.Hitbox.Y = maxY
        }
    }
}

func (g *Game) HandleEnemyShooting() {
    g.world.Lock()
    defer g.world.Unlock()

    for _, e := range g.world.Enemies {
        if e.Weapon == nil {
            continue
        }
        var closestPlayer *player.Player
        var closestDistSq int64 = 1 << 62
        for _, p := range g.world.Players {
            dx := int64(p.Hitbox.X - e.Hitbox.X)
            dy := int64(p.Hitbox.Y - e.Hitbox.Y)
            distSq := dx*dx + dy*dy
            if distSq < closestDistSq {
                closestDistSq = distSq
                closestPlayer = p
            }
        }
        if closestPlayer == nil {
            continue
        }

        dx := closestPlayer.Hitbox.X - e.Hitbox.X
        dy := closestPlayer.Hitbox.Y - e.Hitbox.Y
        if dx == 0 && dy == 0 {
            continue
        }

        bullet, ok := e.Weapon.Shoot(e.Hitbox.X, e.Hitbox.Y, dx, dy)
        if ok {
            g.world.Bullets = append(g.world.Bullets, bullet)
        }
    }
}

func (g *Game) UpdateBullets() {
    g.world.Lock()
    defer g.world.Unlock()
    bullets := &g.world.Bullets
    for i := 0; i < len(*bullets); i++ {
        if !(*bullets)[i].Update() {
            *bullets = append((*bullets)[:i], (*bullets)[i+1:]...)
            i--
        }
    }
}

func (g *Game) Update() error {
    g.MovePlayer()
    g.MoveEnemies()
    g.HandleEnemyShooting()
    g.UpdateBullets()
    return nil
}

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
        if id == g.localID {
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