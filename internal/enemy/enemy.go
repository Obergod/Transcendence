package enemy

import (
	"golang.org/x/image/math/fixed"
	
	"transcendance/internal/hitbox"
)

type Enemy struct {
    Hitbox	*hitbox.Hitbox
	HP		int
	Speed	fixed.Int26_6
	Type    int
    ID      string
}

func NewEnemy(startX, startY fixed.Int26_6, baseHP int, id string) *Enemy {
    return &Enemy{
		Hitbox: hitbox.NewHitbox(startX, startY, fixed.I(10)),
		HP: baseHP,
		Speed: fixed.Int26_6(5 * 64),
		ID: id,
	}
}