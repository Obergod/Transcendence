package game

import (
    "math/rand"
    "golang.org/x/image/math/fixed"
    "transcendance/internal/enemy"
)

func (g *Game) SpawnEnemies() {
    // Vérifier s'il reste des ennemis vivants
    aliveEnemies := 0
    for _, e := range g.world.Enemies {
        if e.IsAlive {
            aliveEnemies++
        }
    }

    // S'il reste des ennemis, ne pas en spawn de nouveaux
    if aliveEnemies > 0 {
        return
    }

    screenWidth := 800
    screenHeight := 600

    // Le nombre d'ennemis augmente avec la vague
    numEnemies := 1 + g.waveNumber*2

    // Limiter le nombre maximum
    if numEnemies > 40 {
        numEnemies = 40
    }

    g.world.Lock()
    defer g.world.Unlock()

    for i := 0; i < numEnemies; i++ {
        var x, y int

        // Choisir aléatoirement un bord (0: haut, 1: droite, 2: bas, 3: gauche)
        edge := rand.Intn(4)

        switch edge {
        case 0: // Bord haut
            x = rand.Intn(screenWidth)
            y = 0
        case 1: // Bord droit
            x = screenWidth
            y = rand.Intn(screenHeight)
        case 2: // Bord bas
            x = rand.Intn(screenWidth)
            y = screenHeight
        case 3: // Bord gauche
            x = 0
            y = rand.Intn(screenHeight)
        }

        enemyName := "enemy_wave" + itoa(g.waveNumber) + "_" + itoa(i)

        // Créer un ennemi à distance avec NewRanged
        newEnemy := enemy.NewRanged(fixed.I(x), fixed.I(y), enemyName)

        // Optionnel: Augmenter les stats selon le numéro de vague
        // Plus la vague est élevée, plus les ennemis sont forts
        if g.waveNumber > 3 {
            // Augmenter les PV
            newEnemy.HP = 100 + (g.waveNumber-3)*10
            if newEnemy.HP > 300 {
                newEnemy.HP = 300
            }

            // Augmenter la vitesse
            newEnemy.Speed = fixed.I(4 + (g.waveNumber-3)/2)
            if newEnemy.Speed > fixed.I(8) {
                newEnemy.Speed = fixed.I(8)
            }

            // Augmenter les dégâts de l'arme
            if newEnemy.Weapon != nil {
                newEnemy.Weapon.Damage = 5 + (g.waveNumber-3)/2
                if newEnemy.Weapon.Damage > 15 {
                    newEnemy.Weapon.Damage = 15
                }
            }
        }

        g.world.Enemies[enemyName] = newEnemy
    }

    // Incrémenter le numéro de vague après le spawn
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