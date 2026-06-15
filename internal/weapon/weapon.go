package weapon

import (
    "math"
    "time"

    "golang.org/x/image/math/fixed"
    "transcendance/internal/logger"
)

type Weapon struct {
    Damage      int
    Cooldown    time.Duration
    LastShot    time.Time
    BulletSpeed float64
    BulletRange float64
    BulletSize  fixed.Int26_6
}

func NewWeapon(damage int, fireRate float64, bulletSpeed, bulletRange float64, bulletSize fixed.Int26_6) *Weapon {
    interval := time.Duration(float64(time.Second) / fireRate)
    logger.Debugf("Nouvelle arme créée: dégâts=%d, cadence=%v, vitesse=%.1f, portée=%.1f",
        damage, interval, bulletSpeed, bulletRange)
    return &Weapon{
        Damage:      damage,
        Cooldown:    interval,
        LastShot:    time.Now(),
        BulletSpeed: bulletSpeed,
        BulletRange: bulletRange,
        BulletSize:  bulletSize,
    }
}

func (w *Weapon) CanShoot() bool {
    return time.Since(w.LastShot) >= w.Cooldown
}

func (w *Weapon) Shoot(x, y, dirX, dirY fixed.Int26_6, ownerID string, ownerIsPlayer bool) (*Bullet, bool) {
    if !w.CanShoot() {
        logger.Debugf("Arme de %s en cooldown, tir ignoré", ownerID)
        return nil, false
    }
    w.LastShot = time.Now()

    dx := float64(dirX) / 64.0
    dy := float64(dirY) / 64.0
    length := math.Hypot(dx, dy)
    if length == 0 {
        logger.Debugf("Direction nulle pour %s, tir ignoré", ownerID)
        return nil, false
    }
    ux := dx / length
    uy := dy / length
    moveX := ux * w.BulletSpeed
    moveY := uy * w.BulletSpeed
    moveXFixed := fixed.Int26_6(int64(moveX * 64))
    moveYFixed := fixed.Int26_6(int64(moveY * 64))
    stepLengthFixed := fixed.Int26_6(int64(w.BulletSpeed * 64))

    maxRangeFixed := fixed.Int26_6(int64(w.BulletRange * 64))

    bullet := NewBullet(x, y, w.BulletSize, moveXFixed, moveYFixed, stepLengthFixed, w.Damage, maxRangeFixed, ownerID, ownerIsPlayer)
    logger.Debugf("Tir effectué par %s: direction (%.2f, %.2f), vitesse %.1f", ownerID, ux, uy, w.BulletSpeed)
    return bullet, true
}