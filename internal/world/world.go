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

func (w *World) AddPlayer(id string, p *player.Player) {
    w.Lock()
    defer w.Unlock()
    w.Players[id] = p
}

func (w *World) AddEnemy(id string, e *enemy.Enemy) {
    w.Lock()
    defer w.Unlock()
    w.Enemies[id] = e
}