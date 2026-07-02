package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// 🌍 CORS erlaubt NUR dem eigenen Frontend (bekannte Origins), auf das Backend
// zuzugreifen – nicht mehr jeder beliebigen Website ("*").
// Zusätzliche Origins können über die Umgebungsvariable CORS_ORIGINS
// (kommagetrennt) erlaubt werden.
func CORS() gin.HandlerFunc {
	erlaubt := map[string]bool{
		"http://localhost:5173": true, // Dev-Server (Vite)
		"http://localhost:4173": true, // Vorschau-Build (vite preview)
	}
	if extra := os.Getenv("CORS_ORIGINS"); extra != "" {
		for _, o := range strings.Split(extra, ",") {
			erlaubt[strings.TrimSpace(o)] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if erlaubt[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		}

		// Vorab-Anfragen (OPTIONS) direkt beantworten.
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
