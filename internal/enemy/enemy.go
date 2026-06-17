package enemy

import (
    "golang.org/x/image/math/fixed"

    "transcendance/internal/entity"
    "transcendance/internal/weapon"
)

type Enemy struct {
    *entity.Entity
    Weapon *weapon.Weapon
}

func newEnemy(x, y fixed.Int26_6, hp int, speed fixed.Int26_6, id string) *Enemy {
    return &Enemy{
        Entity: entity.NewEntity(x, y, fixed.I(10), hp, speed, id),
        Weapon: nil,
    }
}

// NewRanged crée un ennemi à distance avec 100 PV, vitesse 4 pixels/frame,
// arme : dégâts 5, cadence 1 tir/s, vitesse balle 2 pixels/frame, portée 400 pixels, taille balle 3.
func NewRanged(x, y fixed.Int26_6, id string) *Enemy {
    e := newEnemy(x, y, 100, fixed.I(12), id)
    e.Weapon = weapon.NewWeapon(5, 1.0, 2.0, 400.0, fixed.I(3))
    return e
}