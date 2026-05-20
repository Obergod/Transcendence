package auth

import (
	"net/http"
    "golang.org/x/crypto/bcrypt"

	"transcendance/internal/models"

    "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type SigninRequest struct {
	Login string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func PasswordEncrypt(password string) (string, error) {
	hash, err :=bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CompareHashAndPassword(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func	SignupHandler(db *gorm.DB) gin.HandlerFunc {
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
		// need to check if real error

		encryptPw, err := PasswordEncrypt(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message" : "Pswd encryption error"})
			return
		}

		user := models.User{
			Username: req.Username,
			Email: req.Email,
			PasswordHash: encryptPw,
		}
		models.CreateUser(db, &user)

		c.JSON(http.StatusOK, gin.H{"message" : "User signed up"})
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

		// !!! Need to implement JWT (maybe use gin-jwt)

	//	_, err = auth.GenerateJWT(user.ID)
	//	if err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
	//		return
	//	}

		c.JSON(http.StatusOK, gin.H{"message": "User logged in"})
	}
}
