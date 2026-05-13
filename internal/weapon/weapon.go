package weapon

import (
    "math"
    "time"
    "golang.org/x/image/math/fixed"
)

type Weapon struct {
    Damage      int
    Cooldown    time.Duration
    LastShot    time.Time
    BulletSpeed float64   // pixels par frame (en flottant pour la précision)
    BulletRange float64   // pixels
}

func NewWeapon(damage int, fireRate float64, bulletSpeed, bulletRange float64) *Weapon {
    interval := time.Duration(float64(time.Second) / fireRate)
    return &Weapon{
        Damage:      damage,
        Cooldown:    interval,
        LastShot:    time.Now(),
        BulletSpeed: bulletSpeed,
        BulletRange: bulletRange,
    }
}

func (w *Weapon) CanShoot() bool {
    return time.Since(w.LastShot) >= w.Cooldown
}

func (w *Weapon) Shoot(x, y, dirX, dirY fixed.Int26_6) (*Bullet, bool) {
    if !w.CanShoot() {
        return nil, false
    }
    w.LastShot = time.Now()

    dx := float64(dirX) / 64.0
    dy := float64(dirY) / 64.0
    length := math.Hypot(dx, dy)
    if length == 0 {
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

    bullet := NewBullet(x, y, moveXFixed, moveYFixed, stepLengthFixed, w.Damage, maxRangeFixed)
    return bullet, true
}