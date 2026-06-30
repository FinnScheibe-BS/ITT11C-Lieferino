package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 🔒 Auth prüft das JWT-Token im "Authorization: Bearer ..."-Header.
// Bei Erfolg legt es die User-ID im Context ab ("userID").
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"fehler": "Kein Token"})
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "dev-secret-bitte-aendern"
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"fehler": "Token ungültig"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"fehler": "Token ungültig"})
			return
		}
		sub, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"fehler": "Token ungültig"})
			return
		}

		c.Set("userID", uint(sub))
		c.Next()
	}
}
