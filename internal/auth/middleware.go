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
		// 1. Récupérer l'en-tête "Authorization"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Autorisation requise"})
			return
		}

		// 2. Vérifier le format (doit être "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format du token invalide"})
			return
		}
		tokenString := parts[1]

		// 3. Analyser et vérifier la signature du token
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

		// 4. Extraire l'ID de l'utilisateur et le stocker dans le contexte de la requête
		if claims, ok := token.Claims.(*JWTClaims); ok {
			c.Set("userID", claims.UserID)
			c.Next() // Tout est bon, on laisse passer à la route demandée !
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Impossible de lire le token"})
			return
		}
	}
}