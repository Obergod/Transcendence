package enemy

import (
	"golang.org/x/image/math/fixed"
)

type Enemy struct {
    X, Y	fixed.Int26_6
	HP		int
	Speed	fixed.Int26_6
	Type     int
    ID       string
}

func NewEnemy(startX, startY int, baseHP int) *Enemy {
    return &Enemy{
		X: fixed.I(startX),
		Y: fixed.I(startY),
		HP: baseHP,
		Speed: fixed.Int26_6(1 * 64),
	}
}