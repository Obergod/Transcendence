package social

import (
	"math"
	"net/http"
	"transcendance/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Calcul XP d'une partie (score x bonus) bonus incremente de 5% tous les 10s
func MatchXP(durationSeconds int) float64 {
	if durationSeconds <= 0 {
		return 0
	}
	score := float64(durationSeconds) * 60.0
	steps := durationSeconds / 10
	multiplier := 1.2 * math.Pow(1.05, float64(steps))
	return score * multiplier
}

// xp actuel et xp pour prochain niveau
func ComputeLevel(totalXP float64) (level int, xpInLevel float64, xpForNext float64) {
	level = 1
	cost := 10000.0
	remaining := totalXP
	for remaining >= cost {
		remaining -= cost
		level ++
		cost *= 1.1
	}
	return level, remaining, cost
}

func GetLevelHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, _ := c.Get("userID")
		userID := userIDInterface.(uint)

		var matches []models.Match
		db.Where("user_id = ?", userID).Find(&matches)

		var totalXP float64
		for _, m := range matches {
			totalXP += MatchXP(m.Duration)
		}

		level, xpInLevel, xpForNext := ComputeLevel(totalXP)
		c.JSON(http.StatusOK, gin.H{
			"level":		level,
			"xp_in_level":	int(xpInLevel),
			"xp_for_next":	int(xpForNext),
		})
	}
}

func LevelForUser(db *gorm.DB, userID uint) int {
	var matches []models.Match
	db.Where("user_id = ?", userID).Find(&matches)
	var total float64
	for _, m := range matches {
		total += MatchXP(m.Duration)
	}
	level, _, _ := ComputeLevel(total)
	return level
}
