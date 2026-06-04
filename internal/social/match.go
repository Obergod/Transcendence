package social

import (
	"net/http"
	"transcendance/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SaveMatchRequest struct {
	Duration int `json:"duration" binding:"required"`
	Score    int `json:"score" binding:"required"` // NOUVEAU
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
			Score:    req.Score, // Sauvegarde du score pur
		}
		db.Create(&match)
		c.JSON(http.StatusOK, gin.H{"message": "Score sauvegardé !"})
	}
}

// NOUVEAU : 4. Récupérer les 3 statistiques globales pour le profil
type UserStatsResponse struct {
	BestScore     int `json:"best_score"`
	BestDuration  int `json:"best_duration"`
	TotalDuration int `json:"total_duration"`
}

func UserStatsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, _ := c.Get("userID")
		userID := userIDInterface.(uint)

		var stats UserStatsResponse

		// Récupérer le meilleur score (COALESCE évite d'avoir un bug si la table est vide au début)
		db.Model(&models.Match{}).Where("user_id = ?", userID).Select("COALESCE(MAX(score), 0)").Row().Scan(&stats.BestScore)

		// Récupérer le meilleur temps (en secondes)
		db.Model(&models.Match{}).Where("user_id = ?", userID).Select("COALESCE(MAX(duration), 0)").Row().Scan(&stats.BestDuration)

		// Récupérer la somme totale du temps passé ingame
		db.Model(&models.Match{}).Where("user_id = ?", userID).Select("COALESCE(SUM(duration), 0)").Row().Scan(&stats.TotalDuration)

		c.JSON(http.StatusOK, stats)
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

		var friendships []models.Friendship
		db.Where("(user_id = ? OR friend_id = ?) AND status = 'accepted'", userID, userID).Find(&friendships)

		friendIDs := []uint{userID}
		for _, f := range friendships {
			if f.UserID == userID {
				friendIDs = append(friendIDs, f.FriendID)
			} else {
				friendIDs = append(friendIDs, f.UserID)
			}
		}

		var bestMatches []models.Match
		db.Preload("User").
			Where("user_id IN ?", friendIDs).
			Order("duration desc").
			Limit(10).
			Find(&bestMatches)

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