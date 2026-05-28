package auth

import (
	"os"
	"time"

	"transcendance/internal/models"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"	
	"gorm.io/gorm"
)

const identityKey = "user_id"

func InitParams(db *gorm.DB) *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:			"Transcendance",
		Key:			[]byte(os.Getenv("JWT_SECRET")),
		Timeout:		24 * time.Hour,
		MaxRefresh: 	24 * time.Hour,
		IdentityKey:	identityKey,
		PayloadFunc:	payloadFunc(),

		Authenticator: authenticator(db),
		IdentityHandler: identityHandler(),
		Unauthorized:	unauthorized(),
	}
}

func authenticator(db *gorm.DB) func(c *gin.Context) (any, error) {
	return func(c *gin.Context) (any, error) {
		var req	SigninRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, jwt.ErrMissingLoginValues
		}

		user, err := models.GetUserByLogin(db, req.Login)
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}

		if err := CompareHashAndPassword(user.PasswordHash, req.Password); err != nil {
			return nil, jwt.ErrFailedAuthentication
		}
		return user, nil
	}
}

func payloadFunc() func(data any) gojwt.MapClaims {
	return func(data any) gojwt.MapClaims {
		if u, ok := data.(*models.User); ok {
			return gojwt.MapClaims{
				identityKey: u.ID,
				"username": u.Username,
			}
		}
		return gojwt.MapClaims{}
	}
}

func identityHandler() func(c *gin.Context) any {
	return func(c *gin.Context) any {
		claims := jwt.ExtractClaims(c)
		return &models.User{
			ID: uint(claims[identityKey].(float64)),
			Username: claims["username"].(string),
		}
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{"code": code, "message": message})
	}
}
