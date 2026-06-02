package social

import (
	"net/http"
	"transcendance/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SaveMatchRequest struct {
	Duration int `json:"duration" binding:"required"`
}

// 1. Sauvegarder le score à la fin d'une partie
func SaveMatchHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, _ := c.Get("userID")
		userID := userIDInterface.(uint)

		var req SaveMatchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalide"})
			return
		}

		match := models.Match{
			UserID:   userID,
			Duration: req.Duration,
		}
		db.Create(&match)
		c.JSON(http.StatusOK, gin.H{"message": "Score sauvegardé !"})
	}
}

// 2. Récupérer l'historique personnel
func MyHistoryHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, _ := c.Get("userID")
		userID := userIDInterface.(uint)

		var matches []models.Match
		db.Where("user_id = ?", userID).Order("created_at desc").Find(&matches)

		c.JSON(http.StatusOK, gin.H{"data": matches})
	}
}

// 3. Récupérer le TOP 10 (Moi + Mes amis)
func LeaderboardHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, _ := c.Get("userID")
		userID := userIDInterface.(uint)

		// A. On récupère la liste de mes amis acceptés
		var friendships []models.Friendship
		db.Where("(user_id = ? OR friend_id = ?) AND status = 'accepted'", userID, userID).Find(&friendships)

		// B. On construit une liste des IDs autorisés (Moi inclus)
		friendIDs := []uint{userID}
		for _, f := range friendships {
			if f.UserID == userID {
				friendIDs = append(friendIDs, f.FriendID)
			} else {
				friendIDs = append(friendIDs, f.UserID)
			}
		}

		// C. On cherche les 10 meilleures parties de ce groupe !
		var bestMatches []models.Match
		db.Preload("User").
			Where("user_id IN ?", friendIDs).
			Order("duration desc").
			Limit(10).
			Find(&bestMatches)

		// On formate pour React
		type LeaderboardEntry struct {
			Username string `json:"username"`
			Duration int    `json:"duration"`
			Date     string `json:"date"`
		}

		var result []LeaderboardEntry
		for _, m := range bestMatches {
			result = append(result, LeaderboardEntry{
				Username: m.User.Username,
				Duration: m.Duration,
				Date:     m.CreatedAt.Format("02/01/2006"),
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": result})
	}
}