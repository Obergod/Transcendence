package enemy

import (
    "golang.org/x/image/math/fixed"
    "transcendance/internal/hitbox"
    "transcendance/internal/weapon"
)

type Enemy struct {
    Hitbox *hitbox.Hitbox
    HP     int
    Speed  fixed.Int26_6
    ID     string
    Weapon *weapon.Weapon
}

// newEnemy (privé)
func newEnemy(x, y fixed.Int26_6, hp int, speed fixed.Int26_6, id string) *Enemy {
    return &Enemy{
        Hitbox: hitbox.NewHitbox(x, y, fixed.I(10)),
        HP:     hp,
        Speed:  speed,
        ID:     id,
        Weapon: nil,
    }
}

// NewRanged crée un ennemi à distance avec 100 PV, vitesse fixe 1 pixels/frame,
// arme : dégâts 5, cadence 1 tir/s, vitesse balle 2 pixels/frame, portée 400 pixels.
func NewRanged(x, y fixed.Int26_6, id string) *Enemy {
    e := newEnemy(x, y, 100, fixed.I(4), id)
    e.Weapon = weapon.NewWeapon(5, 1.0, 2.0, 400.0)
    return e
}