package entity

import (
    "golang.org/x/image/math/fixed"
    "transcendance/internal/hitbox"
)

// Entity regroupe les champs communs à tous les personnages (joueurs, ennemis).
type Entity struct {
    Hitbox  *hitbox.Hitbox
    HP      int
    Speed   fixed.Int26_6
    ID      string
    IsAlive bool
}

// NewEntity crée une entité de base.
func NewEntity(x, y fixed.Int26_6, radius fixed.Int26_6, hp int, speed fixed.Int26_6, id string) *Entity {
    return &Entity{
        Hitbox:  hitbox.NewHitbox(x, y, radius),
        HP:      hp,
        Speed:   speed,
        ID:      id,
        IsAlive: true,
    }
}