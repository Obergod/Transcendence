package weapon

import "golang.org/x/image/math/fixed"

type Bullet struct {
    X, Y       fixed.Int26_6
    MoveX      fixed.Int26_6
    MoveY      fixed.Int26_6
    StepLength fixed.Int26_6
    Damage     int
    Traveled   fixed.Int26_6
    MaxRange   fixed.Int26_6
}

func NewBullet(x, y, moveX, moveY, stepLength fixed.Int26_6, damage int, maxRange fixed.Int26_6) *Bullet {
    return &Bullet{
        X:          x,
        Y:          y,
        MoveX:      moveX,
        MoveY:      moveY,
        StepLength: stepLength,
        Damage:     damage,
        Traveled:   0,
        MaxRange:   maxRange,
    }
}

func (b *Bullet) Update() bool {
    b.X += b.MoveX
    b.Y += b.MoveY
    b.Traveled += b.StepLength
    return b.Traveled < b.MaxRange
}