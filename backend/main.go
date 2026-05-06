package main

import (
    "log"
    //"golang.org/x/image/math/fixed"

    "github.com/hajimehoshi/ebiten/v2"

    "transcendance/backend/player"
    "transcendance/backend/game"
)

func main() {
    // Initialize the player at starting coordinates (e.g., the center of the screen)
    // For a 800x600 window, (400, 300) is the center
    startX, startY, baseHP := 400, 300, 100
    initialPlayer := player.NewPlayer(startX, startY, baseHP)

    game := game.NewGame(initialPlayer);

    // Set the window title and size
    ebiten.SetWindowTitle("Player Movement with Arrow Keys")
    ebiten.SetWindowSize(800, 600)

    // Run the game
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}