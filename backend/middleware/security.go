package middleware

import "github.com/gin-gonic/gin"

// 🛡️ SecurityHeaders setzt ein paar Standard-Schutz-Header bei jeder Antwort.
// Die schützen u.a. gegen Clickjacking, MIME-Sniffing und ungewolltes Leaken
// der Herkunfts-URL.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Next()
	}
}
