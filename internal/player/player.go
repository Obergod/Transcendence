package player

import (
    "github.com/hajimehoshi/ebiten/v2"
    "golang.org/x/image/math/fixed"

    "transcendance/internal/entity"
    "transcendance/internal/weapon"
    "transcendance/internal/logger"
)

type Player struct {
    *entity.Entity
    Image           *ebiten.Image
    Weapon          *weapon.Weapon
}

func NewPlayer(startX, startY fixed.Int26_6, baseHP int, id string) *Player {
    logger.Debugf("Création du joueur %s à (%.2f, %.2f), PV=%d", id, float64(startX)/64, float64(startY)/64, baseHP)
    return &Player{
        Entity:          entity.NewEntity(startX, startY, fixed.I(20), baseHP, fixed.I(1), id),
        Weapon:          weapon.NewWeapon(5, 2.5, 2.0, 2000.0, fixed.I(3)), // Portée augmentée à 2000
    }
}