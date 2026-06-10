package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/math/fixed"
	"transcendance/internal/hitbox"
	"transcendance/internal/player"
	"transcendance/internal/utils"
)

func KeyPressp1() (fixed.Int26_6, fixed.Int26_6) {
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
	return dx, dy
}

func KeyPressp2() (fixed.Int26_6, fixed.Int26_6) {
	dx := fixed.Int26_6(0)
	dy := fixed.Int26_6(0)
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx = 1
	}
	return dx, dy
}

func (g *Game) MovePlayer(id string) error {
	g.world.Lock()
	defer g.world.Unlock()

	var dx, dy fixed.Int26_6
	switch id {
	case "p1":
		dx, dy = KeyPressp1()
	case "p2":
		dx, dy = KeyPressp2()
	}

	localPlayer, exists := g.world.Players[id]
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