package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired est le middleware qui vérifie le JWT
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// recupérer l'en-tête "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Autorisation requise"})
			return
		}

		// verifier le format (doit être "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format du token invalide"})
			return
		}
		tokenString := parts[1]

		// analyser et vérifier la signature du token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("méthode de signature inattendue")
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalide ou expiré"})
			return
		}

		// extraire l'ID de l'utilisateur et le stocker dans le contexte de la requete
		if claims, ok := token.Claims.(*JWTClaims); ok {
			c.Set("userID", claims.UserID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Impossible de lire le token"})
			return
		}
	}
}