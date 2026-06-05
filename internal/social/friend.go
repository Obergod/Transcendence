package social

import (
	"net/http"
	"transcendance/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FriendRequest struct {
	TargetUsername string `json:"username" binding:"required"`
}

type RespondFriendRequest struct {
	FriendshipID uint   `json:"friendship_id" binding:"required"`
	Action       string `json:"action" binding:"required"` // "accept" ou "reject"
}

// Handler pour accepter ou refuser une demande
func RespondFriendRequestHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, _ := c.Get("userID")
		userID := userIDInterface.(uint)

		var req RespondFriendRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalide"})
			return
		}

		var friendship models.Friendship
		if err := db.First(&friendship, req.FriendshipID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Demande introuvable"})
			return
		}

		if friendship.FriendID != userID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Cette demande ne t'est pas adressée"})
			return
		}

		if req.Action == "accept" {
			friendship.Status = "accepted"
			db.Save(&friendship)
			c.JSON(http.StatusOK, gin.H{"message": "Demande d'ami acceptée"})
		} else if req.Action == "reject" {
			db.Delete(&friendship)
			c.JSON(http.StatusOK, gin.H{"message": "Demande d'ami refusée"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Action inconnue"})
		}
	}
}

// Handler pour récupérer la liste d'amis et les demandes
func ListFriendsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, _ := c.Get("userID")
		userID := userIDInterface.(uint)

		var friendships []models.Friendship
		if err := db.Preload("User").Preload("Friend").
			Where("user_id = ? OR friend_id = ?", userID, userID).
			Find(&friendships).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des amis"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": friendships})
	}
}

// Handler pour envoyer une demande d'ami
func SendFriendRequestHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Récupérer l'ID de l'utilisateur qui fait la requête (grâce au JWT)
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorisé"})
			return
		}
		userID := userIDInterface.(uint)

		var req FriendRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalide"})
			return
		}

		targetUser, err := models.GetUserByUsername(db, req.TargetUsername)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Joueur introuvable"})
			return
		}

		if targetUser.ID == userID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tu ne peux pas t'ajouter toi-même"})
			return
		}

		friendship := models.Friendship{
			UserID:   userID,
			FriendID: targetUser.ID,
			Status:   "pending",
		}

		if err := db.Create(&friendship).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Une demande existe déjà avec ce joueur"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Demande d'ami envoyée à " + targetUser.Username})
	}
}