package middleware

import (
	"net/http"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
)

// 🛠️ AdminOnly lässt nur Admin-Nutzer durch (läuft NACH Auth()).
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("userID")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"fehler": "Nicht angemeldet"})
			return
		}
		var user models.User
		if err := database.DB.First(&user, uid).Error; err != nil || !user.IstAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"fehler": "Kein Admin-Zugang"})
			return
		}
		c.Next()
	}
}
