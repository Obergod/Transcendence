package player

import (
    "github.com/hajimehoshi/ebiten/v2"
    "golang.org/x/image/math/fixed"

    "transcendance/internal/entity"
    "transcendance/internal/weapon"
)

type Player struct {
    *entity.Entity
    Image           *ebiten.Image
    Weapon          *weapon.Weapon
    InvincibleUntil int
}

func NewPlayer(startX, startY fixed.Int26_6, baseHP int, id string) *Player {
    return &Player{
        Entity:          entity.NewEntity(startX, startY, fixed.I(20), baseHP, fixed.I(1), id),
        Weapon:          weapon.NewWeapon(5, 2.5, 2.0, 400.0, fixed.I(3)),
        InvincibleUntil: 0,
    }
}