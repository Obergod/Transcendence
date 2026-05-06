package player

import (
	//"golang.org/x/image/math/fixed"
)

type Player struct {
    X, Y	int
	HP		int
}

func NewPlayer(startX, startY int, baseHP int) *Player {
    return &Player{X: startX, Y: startY, HP: baseHP}
}

func (p *Player) MoveUp()    { p.Y++ }
func (p *Player) MoveDown()  { p.Y-- }
func (p *Player) MoveLeft()  { p.X-- }
func (p *Player) MoveRight() { p.X++ }