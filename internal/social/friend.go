package social

import (
	"net/http"
	"transcendance/internal/models"
	"strconv"

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

		type FriendItem struct {
			FriendshipID	uint	`json:"friendship_id"`
			Status			string	`json:"status"`
			UserID			uint	`json:"user_id"`
			FriendID		uint	`json:"friend_id"`
			OtherID			uint	`json:"other_id"`
			OtherUsername	string	`json:"other_username"`
			OtherAvatar		string	`json:"other_avatar"`
			OtherLevel		int		`json:"other_level"`
		}

		result := make([]FriendItem, 0, len(friendships))
		for _, f := range friendships {
			item := FriendItem{
				FriendshipID:	f.ID,
				Status:			f.Status,
				UserID:			f.UserID,
				FriendID:		f.FriendID,
			}
			if f.UserID == userID && f.Friend != nil {
				item.OtherID = f.Friend.ID
				item.OtherUsername = f.Friend.Username
				item.OtherAvatar = f.Friend.AvatarURL
			} else if f.User != nil {
				item.OtherID = f.User.ID
				item.OtherUsername = f.User.Username
				item.OtherAvatar = f.User.AvatarURL
			}
			item.OtherLevel = LevelForUser(db, item.OtherID)
			result = append(result, item)
		}

		c.JSON(http.StatusOK, gin.H{"data": result})
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

func RemoveFriendHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorisé"})
			return
		}
		myID := userIDInterface.(uint)
		
        friendIDStr := c.Param("friendId")
        friendID, err := strconv.Atoi(friendIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ID d'ami invalide"})
            return
        }

        result := db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", myID, friendID, friendID, myID).Delete(&models.Friendship{})
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
            return
        }

        // Si aucune ligne n'a été affectée, c'est qu'ils n'étaient pas amis
        if result.RowsAffected == 0 {
            c.JSON(http.StatusNotFound, gin.H{"error": "Amitié introuvable"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Ami supprimé avec succès"})
    }
}
