package world

import (
    "sync"
    "transcendance/internal/player"
    "transcendance/internal/enemy"
)

type World struct {
    sync.RWMutex
    Players map[string]*player.Player
    Enemies map[string]*enemy.Enemy
    // Ajoutez projectiles, etc.
}

func NewWorld() *World {
    return &World{
        Players: make(map[string]*player.Player),
        Enemies: make(map[string]*enemy.Enemy),
    }
}

// Mise à jour du monde à chaque tick (ex: 60 fois par seconde)
func (w *World) Update() {
    w.Lock()
    defer w.Unlock()
    // Déplacement des ennemis, collisions, spawn, etc.
}