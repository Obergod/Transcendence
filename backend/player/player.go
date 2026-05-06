package player

import (
	"math"
	"golang.org/x/image/math/fixed"
)

type Player struct {
    X, Y	fixed.Int26_6
	HP		int
	Speed	fixed.Int26_6
}

func NewPlayer(startX, startY int, baseHP int) *Player {
    return &Player{
		X: fixed.I(startX),
		Y: fixed.I(startY),
		HP: baseHP,
		Speed: fixed.Int26_6(math.Round(1 * 64)),
	}
}

func (p *Player) MoveUp()    { p.Y+= p.Speed }
func (p *Player) MoveDown()  { p.Y-= p.Speed }
func (p *Player) MoveLeft()  { p.X-= p.Speed }
func (p *Player) MoveRight() { p.X+= p.Speed }