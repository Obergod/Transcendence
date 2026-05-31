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
    world   	*world.World
    localID 	string
	isGameover	bool
	//need to add score
}

func NewGame(w *world.World, ID string) *Game {
    return &Game{
        world:   w,
        localID: ID,
		isGameover: false,
    }
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

    if !localPlayer.IsAlive {
		// !!! STP Matti implemet dans le ts du game un popup
		// js.Global().Call("onGameover")
		Reset(g)
      //  return nil
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
        if !e.IsAlive {
            continue
        }
        // Trouver le joueur le plus proche
        var closestPlayer *player.Player
        var closestDistSq int64 = 1 << 62

        for _, p := range g.world.Players {
            if !p.IsAlive {
                continue
            }
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
            if p.IsAlive {
                hitbox.PushOutCollisionFixed(e.Hitbox, p.Hitbox)
            }
        }
        for _, other := range g.world.Enemies {
            if e.ID != other.ID && other.IsAlive {
                hitbox.PushOutCollisionFixed(e.Hitbox, other.Hitbox)
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
        if e.Weapon == nil || !e.IsAlive {
            continue
        }
        var closestPlayer *player.Player
        var closestDistSq int64 = 1 << 62
        for _, p := range g.world.Players {
            if !p.IsAlive {
                continue
            }
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

        // Passer l'ownerID et ownerIsPlayer
        bullet, ok := e.Weapon.Shoot(e.Hitbox.X, e.Hitbox.Y, dx, dy, e.ID, false)
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

func (g *Game) HandleBulletCollisions() {
    g.world.Lock()
    defer g.world.Unlock()

    bullets := &g.world.Bullets
    for i := 0; i < len(*bullets); i++ {
        b := (*bullets)[i]
        hit := false

        // Créer la hitbox de la balle
        bulletHitbox := hitbox.NewHitbox(b.X, b.Y, b.Size)

        if b.OwnerIsPlayer {
            // Balle de joueur : touche les ennemis
            for _, e := range g.world.Enemies {
                if !e.IsAlive {
                    continue
                }
                if hitbox.CheckCollision(bulletHitbox, e.Hitbox) {
                    e.HP -= b.Damage
                    if e.HP <= 0 {
                        e.IsAlive = false
                    }
                    hit = true
                    break
                }
            }
        } else {
            // Balle d'ennemi : touche les joueurs
            for _, p := range g.world.Players {
                if !p.IsAlive {
                    continue
                }
                if hitbox.CheckCollision(bulletHitbox, p.Hitbox) {
                    p.HP -= b.Damage
                    if p.HP <= 0 {
                        p.IsAlive = false
                    }
                    hit = true
                    break
                }
            }
        }

        if hit {
            // Supprimer la balle
            *bullets = append((*bullets)[:i], (*bullets)[i+1:]...)
            i--
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

func (g *Game) Update() error {
    g.MovePlayer()
    g.MoveEnemies()
    g.HandleEnemyShooting()
    g.UpdateBullets()
    g.HandleBulletCollisions()
    g.RemoveDeadEnemies()
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
        if !p.IsAlive {
            col = color.RGBA{255, 255, 255, 255}
        } else if id == g.localID {
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
