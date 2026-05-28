package player

import (
    "github.com/hajimehoshi/ebiten/v2"
    "golang.org/x/image/math/fixed"

    "transcendance/internal/entity"
)

type Player struct {
    *entity.Entity
    Image *ebiten.Image
}

func NewPlayer(startX, startY fixed.Int26_6, baseHP int, id string) *Player {
    return &Player{
        Entity: entity.NewEntity(startX, startY, fixed.I(20), baseHP, fixed.I(1), id),
    }
}