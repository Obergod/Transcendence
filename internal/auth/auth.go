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
)

type UpdateProfileRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl"`
}

// Clé secrète pour signer les JWT (À mettre dans un .env plus tard pour la sécurité)
var jwtSecretKey = []byte("super_secret_key_transcendence_42")

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

func SigninHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SigninRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := models.GetUserByLogin(db, req.Login)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := CompareHashAndPassword(user.PasswordHash, req.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Génération du Token JWT
		token, err := GenerateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
			return
		}

		// On renvoie le token ET les infos de l'utilisateur
		c.JSON(http.StatusOK, gin.H{
			"message": "User logged in",
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

func UpdateProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. On récupère l'ID validé par notre "videur" JWT
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne du serveur"})
			return
		}
		userID := userIDInterface.(uint)

		// 2. On lit les nouvelles données envoyées par React
		var req UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 3. On récupère l'utilisateur en base de données
		user, err := models.GetUserByID(db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
			return
		}

		// 4. On met à jour les champs
		user.Username = req.Username
		user.Email = req.Email
		user.AvatarURL = req.AvatarURL

		// 5. On sauvegarde en BDD
		if err := models.UpdateUser(db, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de mettre à jour le profil (pseudo ou email peut-être déjà pris ?)"})
			return
		}

		// 6. On renvoie le profil mis à jour au Frontend
		c.JSON(http.StatusOK, gin.H{
			"message": "Profil mis à jour avec succès",
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

func GetMeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Le middleware a déjà vérifié le token et extrait l'ID
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Non autorisé"})
			return
		}
		userID := userIDInterface.(uint)

		// On cherche l'utilisateur dans la base
		user, err := models.GetUserByID(db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
			return
		}

		// On renvoie son profil
		c.JSON(http.StatusOK, gin.H{
			"id":        user.ID,
			"pseudo":    user.Username,
			"email":     user.Email,
			"avatarUrl": user.AvatarURL,
			"status":    "online",
		})
	}
}

func ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue")
		}
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("token invalide")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims.UserID, nil
	}
	return 0, fmt.Errorf("impossible de lire les claims")
}
