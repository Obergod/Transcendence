package auth

import (
	"net/http"
	"time"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"transcendance/internal/models"

    "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Claims personnalisés pour le JWT
type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Fonction pour générer le Token JWT
func GenerateJWT(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Le token expire dans 24h
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func SignupHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := models.GetUserByUsername(db, req.Username)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"message": "Username already registered"})
			return
		}

		_, err = models.GetUserByEmail(db, req.Email)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"message": "Email already registered"})
			return
		}

		encryptPw, err := PasswordEncrypt(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Pswd encryption error"})
			return
		}

		user := models.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: encryptPw,
			AvatarURL:    "https://upload.wikimedia.org/wikipedia/commons/8/89/Portrait_Placeholder.png",
		}

		if err := models.CreateUser(db, &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user"})
			return
		}

		// Auto-login après l'inscription
		token, err := GenerateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User signed up and logged in",
			"token":   token,
			"user": gin.H{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"avatarUrl": user.AvatarURL,
				"status":    "online",
			},
		})
	}
}
