package game

import (
	"transcendance/internal/hitbox"
	"transcendance/internal/player"
	"transcendance/internal/enemy"
	"transcendance/internal/weapon"
)

func (g *Game) HandleEnemyShooting() {
	type shootInfo struct {
		bullet *weapon.Bullet
	}
	var shoots []shootInfo

	g.world.RLock()
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
		bullet, ok := e.Weapon.Shoot(e.Hitbox.X, e.Hitbox.Y, dx, dy, e.ID, false)
		if ok {
			bullet.SpawnTick = g.ticks
			shoots = append(shoots, shootInfo{bullet: bullet})
		}
	}
	g.world.RUnlock()

	for _, s := range shoots {
		g.world.AddBullet(s.bullet)
	}
}

func (g *Game) HandlePlayersShooting() {
	const fireCooldown = 10

	for _, id := range g.localIDs {
		g.world.RLock()
		player, exists := g.world.Players[id]
		g.world.RUnlock()
		if !exists || !player.IsAlive || player.Weapon == nil {
			continue
		}
		lastTick := g.lastShotTicks[id]
		if g.ticks-lastTick < fireCooldown {
			continue
		}
		g.world.RLock()
		var closestEnemy *enemy.Enemy
		var closestDistSq int64 = 1 << 62
		for _, e := range g.world.Enemies {
			if !e.IsAlive {
				continue
			}
			dx := int64(e.Hitbox.X - player.Hitbox.X)
			dy := int64(e.Hitbox.Y - player.Hitbox.Y)
			distSq := dx*dx + dy*dy
			if distSq < closestDistSq {
				closestDistSq = distSq
				closestEnemy = e
			}
		}
		g.world.RUnlock()

		if closestEnemy == nil {
			continue
		}
		dx := closestEnemy.Hitbox.X - player.Hitbox.X
		dy := closestEnemy.Hitbox.Y - player.Hitbox.Y
		if dx == 0 && dy == 0 {
			continue
		}
		bullet, ok := player.Weapon.Shoot(player.Hitbox.X, player.Hitbox.Y, dx, dy, id, true)
		if ok {
			bullet.SpawnTick = g.ticks
			g.world.AddBullet(bullet)
			g.lastShotTicks[id] = g.ticks
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

		if b.SpawnTick == g.ticks {
			continue
		}

		hit := false
		bulletHitbox := hitbox.NewHitbox(b.X, b.Y, b.Size)

		if b.OwnerIsPlayer {
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
			for _, p := range g.world.Players {
				if !p.IsAlive {
					continue
				}
				if hitbox.CheckCollision(bulletHitbox, p.Hitbox) {
					if g.ticks >= p.InvincibleUntil {
						p.HP -= b.Damage
						if p.HP <= 0 {
							p.IsAlive = false
						} else {
							p.InvincibleUntil = g.ticks + 10
						}
					}
					hit = true
					break
				}
			}
		}

		if hit {
			*bullets = append((*bullets)[:i], (*bullets)[i+1:]...)
			i--
		}
	}
}