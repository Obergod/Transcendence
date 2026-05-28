package weapon

import "golang.org/x/image/math/fixed"

type Bullet struct {
    X, Y          fixed.Int26_6
    Size          fixed.Int26_6
    MoveX, MoveY  fixed.Int26_6
    StepLength    fixed.Int26_6
    Damage        int
    Traveled      fixed.Int26_6
    MaxRange      fixed.Int26_6
    OwnerID       string // ID de l'entité qui a tiré
    OwnerIsPlayer bool   // true = joueur, false = ennemi
}

func NewBullet(x, y, size, moveX, moveY, stepLength fixed.Int26_6, damage int, maxRange fixed.Int26_6, ownerID string, ownerIsPlayer bool) *Bullet {
    return &Bullet{
        X:             x,
        Y:             y,
        Size:          size,
        MoveX:         moveX,
        MoveY:         moveY,
        StepLength:    stepLength,
        Damage:        damage,
        Traveled:      0,
        MaxRange:      maxRange,
        OwnerID:       ownerID,
        OwnerIsPlayer: ownerIsPlayer,
    }
}

func (b *Bullet) Update() bool {
    b.X += b.MoveX
    b.Y += b.MoveY
    b.Traveled += b.StepLength
    return b.Traveled < b.MaxRange
}