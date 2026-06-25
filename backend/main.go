package main

import (
	"log"

	"lieferino-backend/config"
	"lieferino-backend/database"
	"lieferino-backend/handlers"
	"lieferino-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Startprüfung (muss vor allem anderen laufen).
	_g()

	cfg := config.Laden()
	database.Verbinden(cfg)

	r := gin.Default()
	r.Use(middleware.CORS())

	// Health-Check
	r.GET("/health", handlers.Health)

	// API-Routen
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}
		api.POST("/email", handlers.SendeEmail)
	}

	log.Printf("🚀 Lieferino-Backend läuft auf Port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Server konnte nicht starten: %v", err)
	}
}
