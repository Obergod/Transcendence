package player

type Player struct {
    X, Y int
}

func NewPlayer(startX, startY int) *Player {
    return &Player{X: startX, Y: startY}
}

func (p *Player) MoveUp()    { p.Y++ }
func (p *Player) MoveDown()  { p.Y-- }
func (p *Player) MoveLeft()  { p.X-- }
func (p *Player) MoveRight() { p.X++ }