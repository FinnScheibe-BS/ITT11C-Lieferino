package middleware

import "github.com/gin-gonic/gin"

// 🌍 CORS erlaubt dem Frontend (anderer Port), auf das Backend zuzugreifen.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		// Vorab-Anfragen (OPTIONS) direkt beantworten.
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
