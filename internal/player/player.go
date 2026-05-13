package player

import (
    "github.com/hajimehoshi/ebiten/v2"
    "golang.org/x/image/math/fixed"

	"transcendance/internal/hitbox"
)

type Player struct {
    Hitbox	*hitbox.Hitbox
	HP		int
	Speed	fixed.Int26_6
	ID       string
	Image	*ebiten.Image
}

func NewPlayer(startX, startY fixed.Int26_6, baseHP int, id string) *Player {
    return &Player{
		Hitbox: hitbox.NewHitbox(startX, startY, fixed.I(20)),
		HP: baseHP,
		Speed: fixed.I(1),
		ID: id,
	}
}