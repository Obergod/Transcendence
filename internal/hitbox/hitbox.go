package hitbox

import (
    "golang.org/x/image/math/fixed"

	"transcendance/internal/utils"

)

type Hitbox struct {
	X, Y, R fixed.Int26_6
}

func NewHitbox(x, y, r fixed.Int26_6) *Hitbox {
	return &Hitbox{X: x, Y: y, R: r}
}

func PushOutCollisionFixed(e1, e2 *Hitbox) bool {
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