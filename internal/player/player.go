package player

import (
    "github.com/hajimehoshi/ebiten/v2"
    "golang.org/x/image/math/fixed"

    "transcendance/internal/entity"
    "transcendance/internal/weapon"
)

type Player struct {
    *entity.Entity
<<<<<<< HEAD
    Weapon *weapon.Weapon
    Image *ebiten.Image  
=======
    Image *ebiten.Image
    Weapon *weapon.Weapon
>>>>>>> mafioron
}

func NewPlayer(startX, startY fixed.Int26_6, baseHP int, id string) *Player {
    return &Player{
        Entity: entity.NewEntity(startX, startY, fixed.I(20), baseHP, fixed.I(1), id),
<<<<<<< HEAD
        Weapon: weapon.NewWeapon(5, 5.0, 2.0, 400.0, fixed.I(3)),
=======
        Weapon:  weapon.NewWeapon(5, 1.0, 2.0, 400.0, fixed.I(3)),
>>>>>>> mafioron
    }
}
