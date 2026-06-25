package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ❤️ Einfacher Health-Check, um zu sehen, ob das Backend läuft.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "lieferino-backend",
	})
}
