package social

import (
	"net/http"
	"strconv"
	"transcendance/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetHistoryHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, _ := c.Get("userID")
		userID := userIDInterface.(uint)

		friendID64, err := strconv.ParseUint(c.Param("friendId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID ami invalide"})
			return
		}
		friendID := uint(friendID64)

		messages, err := models.GetConversation(db, userID, friendID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur base de donnees"})
			return
		}

		// formatage comme msg WebSocket -> React affiche msgs avec le meme code
		type ChatItem struct {
			Type		string	`json:"type"`
			SenderID	uint	`json:"sender_id"`
			TargetID	uint	`json:"target_id"`
			Content		string	`json:"content"`
			SenderName	string	`json:"sender_name"`
		}

		me, _ := models.GetUserByID(db, userID)
		friend, _ := models.GetUserByID(db, friendID)
		nameFor := func(id uint) string {
			if me != nil && id == me.ID {
				return me.Username
			}
			if friend != nil && id == friend.ID {
				return friend.Username
			}
			return "Inconnu"
		}

		result := make([]ChatItem, 0, len(messages))
		for _, m := range messages {
			result = append(result, ChatItem{
				Type:		"chat",
				SenderID:	m.SenderID,
				TargetID:	m.ReceiverID,
				Content:	m.Content,
				SenderName:	nameFor(m.SenderID),
			})
		}

		friendInfo := gin.H{"avatarUrl": "", "username": ""}
		if friend != nil {
			friendInfo = gin.H{
				"avatarUrl":	friend.AvatarURL,
				"username":		friend.Username,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"data": result,
			"friend": friendInfo,
		})
	}
}
