package game

import (
    "math/rand"
    "golang.org/x/image/math/fixed"
    "transcendance/internal/enemy"
    "transcendance/internal/logger"
)

func (g *Game) SpawnEnemies() {
    aliveEnemies := 0
    for _, e := range g.world.Enemies {
        if e.IsAlive {
            aliveEnemies++
        }
    }

    if aliveEnemies > 0 {
        return
    }

    logger.Infof("Vague %d - spawn de nouveaux ennemis", g.waveNumber+1)

    screenWidth := 800
    screenHeight := 600

    numEnemies := 1 + g.waveNumber*2
    if numEnemies > 40 {
        numEnemies = 40
    }

    g.world.Lock()
    defer g.world.Unlock()

    for i := 0; i < numEnemies; i++ {
        var x, y int
        edge := rand.Intn(4)
        switch edge {
        case 0:
            x = rand.Intn(screenWidth)
            y = 0
        case 1:
            x = screenWidth
            y = rand.Intn(screenHeight)
        case 2:
            x = rand.Intn(screenWidth)
            y = screenHeight
        case 3:
            x = 0
            y = rand.Intn(screenHeight)
        }

        enemyName := "enemy_wave" + itoa(g.waveNumber + 1) + "_" + itoa(i + 1)
        newEnemy := enemy.NewRanged(fixed.I(x), fixed.I(y), enemyName)

        

        g.world.Enemies[enemyName] = newEnemy
        logger.Debugf("Ennemi %s spawné à (%d, %d)", enemyName, x, y)
    }

    g.waveNumber++
}

func itoa(n int) string {
    if n == 0 {
        return "0"
    }
    var digits [10]byte
    i := len(digits)
    for n > 0 {
        i--
        digits[i] = byte('0' + n%10)
        n /= 10
    }
    return string(digits[i:])
}