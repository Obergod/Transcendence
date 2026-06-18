package auth

import (
	"os"
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

// NOUVEAU : Fonction sécurisée pour récupérer la clé secrète depuis le .env
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Valeur de secours au cas où (ne devrait pas arriver si le .env est bien configuré)
		return []byte("fallback_secret_key")
	}
	return []byte(secret)
}

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
	// CORRECTION ICI : on utilise getJWTSecret()
	return token.SignedString(getJWTSecret())
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
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne du serveur"})
			return
		}
		userID := userIDInterface.(uint)

		// 1. On ne lit plus du JSON, on lit un FormData
		username := c.PostForm("username")
		email := c.PostForm("email")

		user, err := models.GetUserByID(db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
			return
		}

		user.Username = username
		user.Email = email

		// 2. Gestion de l'upload de l'image (si fournie)
		file, err := c.FormFile("avatar")
		if err == nil {
			// Verification du poids (2 MB max)
			const maxSize = 2 * 1024 * 1024
			if file.Size > maxSize {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Image trop lourde (max 2MB)"})
				return
			}

			// Verification du MIME type reel
			src, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de lire le fichier"})
				return
			}
			defer src.Close()

			// lire les 512 premiers octets
			buffer := make([]byte, 512)
			_, err = src.Read(buffer)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible d'analyser le fichier"})
				return
			}

			contentType := http.DetectContentType(buffer)
			allowedTypes := map[string]bool{
				"image/jpeg": true,
				"image/png": true,
				"image/gif": true,
				"image/webp": true,
			}
			if !allowedTypes[contentType] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Format non autorise (JPG, PNG, GIF, WEBP uniquement)"})
				return
			}

			// sauvegarde
			filename := fmt.Sprintf("avatar_%d_%s", user.ID, file.Filename)
			savePath := "uploads/avatars/" + filename
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de sauvegarder l'image"})
				return
			}
			user.AvatarURL = "/" + savePath
		}

		// 3. On sauvegarde en BDD
		if err := models.UpdateUser(db, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de mettre à jour le profil (pseudo ou email déjà pris ?)"})
			return
		}

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
		// CORRECTION ICI : on utilise getJWTSecret()
		return getJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("token invalide")
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims.UserID, nil
	}
	return 0, fmt.Errorf("impossible de lire les claims")
}