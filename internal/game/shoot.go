package game

import (
    "transcendance/internal/hitbox"
    "transcendance/internal/player"
	"transcendance/internal/enemy"
)

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

func (g *Game) HandlePlayersShooting() {
    g.world.Lock()
    defer g.world.Unlock()

    for _, p := range g.world.Players{
        if p.Weapon == nil || !p.IsAlive {
            continue
        }
        var closestEnemy *enemy.Enemy
        var closestDistSq int64 = 1 << 62
        for _, e := range g.world.Enemies {
            if !e.IsAlive {
                continue
            }
            dx := int64(e.Hitbox.X - p.Hitbox.X)
            dy := int64(e.Hitbox.Y - p.Hitbox.Y)
            distSq := dx*dx + dy*dy
            if distSq < closestDistSq {
                closestDistSq = distSq
                closestEnemy = e
            }
        }
        if closestEnemy == nil {
            continue
        }

        dx := closestEnemy.Hitbox.X - p.Hitbox.X
        dy := closestEnemy.Hitbox.Y - p.Hitbox.Y
        if dx == 0 && dy == 0 {
            continue
        }

        // Passer l'ownerID et ownerIsPlayer
        bullet, ok := p.Weapon.Shoot(p.Hitbox.X, p.Hitbox.Y, dx, dy, p.ID, true)
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


