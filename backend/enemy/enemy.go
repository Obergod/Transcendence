package	enemy

import
(
	"fmt"
	"image/color"
	"golang.org/x/image/math/fixed"
)

type Enemy struct{
	X, Y fixed.Int26_6
	Speed fixed.Int26_6
	HP int
	Damage int
}

func NewWorm(x, y int)	*Enemy{
	return &Enemy{
		X: fixed.I(x), Y: fixed.I(y),
		Speed: fixed.I(1),
		HP: 10,
		Damage: 1,
	}
}
