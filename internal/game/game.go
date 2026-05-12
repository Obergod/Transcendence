package game

import (
    "image/color"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"
    "golang.org/x/image/math/fixed"

    "transcendance/internal/world"
    "transcendance/internal/utils"
    "transcendance/internal/player"
)

type Game struct {
    world   *world.World
    localID string
}

func NewGame(w *world.World, ID string) *Game {
    return &Game{world: w, localID: ID}
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

func (g *Game) MoveEnemies() {
    if len(g.world.Players) == 0 {
        return
    }

    for _, e := range g.world.Enemies {
        // Trouver le joueur le plus proche
        var closestPlayer *player.Player
        var closestDistSq fixed.Int26_6 = fixed.Int26_6(0x7fffffff) // très grand

        for _, p := range g.world.Players {
            dx := p.X - e.X
            dy := p.Y - e.Y

            // on compare les hypothenuses (pythagore takaptew)
            // (dx*dx + dy*dy) en fixed, attention aux débordements 64-bit
            // On convertit en int64 pour le produit
            dx64 := int64(dx)
            dy64 := int64(dy)
            distSq := fixed.Int26_6(dx64*dx64 + dy64*dy64)

            if distSq < closestDistSq {
                closestDistSq = distSq
                closestPlayer = p
            }
        }

        if closestPlayer == nil {
            continue
        }

        dx := closestPlayer.X - e.X
        dy := closestPlayer.Y - e.Y

        // Normalisation approximée en fixed : on veut (dx, dy) * speed / longueur
        // Longueur = sqrt(dx^2 + dy^2). Éviter division par zéro.
        if dx == 0 && dy == 0 {
            continue
        }

        lenFixed := utils.FixedSqrt(dx*dx + dy*dy)
        if lenFixed == 0 {
            continue
        }

        // On travaille en int64 pour garder la précision
        speed64 := int64(e.Speed)
        len64 := int64(lenFixed)

        moveX := fixed.Int26_6((int64(dx) * speed64) / len64)
        moveY := fixed.Int26_6((int64(dy) * speed64) / len64)

        e.X += moveX
        e.Y += moveY
    }
}

// Update logic runs every tick (1/60 second by default)
func (g *Game) Update() error {
    g.MovePlayer()
    g.MoveEnemies()
    return nil
}

func (g *Game) DrawEnemies(screen *ebiten.Image) {
    for _, e := range g.world.Enemies {
        x := utils.FixedToFloat32(e.X)
        y := utils.FixedToFloat32(e.Y)

        var col color.Color
        col = color.RGBA{255, 0, 0, 255} // Ennemis en rouge

        vector.FillRect(screen, x-10, y-10, 20, 20, col, false)
    }
}

func (g *Game) DrawPlayers(screen *ebiten.Image) {
    for id, p := range g.world.Players {
        x := utils.FixedToFloat32(p.X)
        y := utils.FixedToFloat32(p.Y)

        var col color.Color
        if id == g.localID {
            col = color.RGBA{0, 255, 0, 255} // Vert pour le joueur local
        } else {
            col = color.RGBA{0, 0, 255, 255} // Bleu pour les autres
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
    g.DrawEnemies(screen)
}

// Layout defines the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 800, 600
}