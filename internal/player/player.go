package player

import (
    "math"

    "github.com/hajimehoshi/ebiten/v2"
    "golang.org/x/image/math/fixed"
)

type Player struct {
    X, Y	fixed.Int26_6
	HP		int
	Speed	fixed.Int26_6
	ID       string
	Image	*ebiten.Image
}

func NewPlayer(startX, startY int, baseHP int) *Player {
    return &Player{
		X: fixed.I(startX),
		Y: fixed.I(startY),
		HP: baseHP,
		Speed: fixed.Int26_6(math.Round(1 * 64)),
	}
}