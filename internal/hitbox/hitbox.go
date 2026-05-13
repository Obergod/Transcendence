package hitbox

import (
    "golang.org/x/image/math/fixed"
)

type Hitbox struct {
	X, Y, R fixed.Int26_6
}

func NewHitbox(x, y, r fixed.Int26_6) *Hitbox {
	return &Hitbox{X: x, Y: y, R: r}
}
