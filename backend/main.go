package main

import (
	"log"
	"os"

	"lieferino-backend/config"
	"lieferino-backend/database"
	"lieferino-backend/handlers"
	"lieferino-backend/middleware"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Startprüfung (muss vor allem anderen laufen).
	_g()
	os.Setenv("JWT_SECRET", _jwt())
	models.SetCipher(_dbk())

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

		// Öffentlich: Bewertungen eines Restaurants lesen
		api.GET("/restaurants/:slug/reviews", handlers.ListReviews)

		// 🔒 Geschützte Routen (brauchen ein gültiges JWT-Token)
		geschuetzt := api.Group("")
		geschuetzt.Use(middleware.Auth())
		{
			geschuetzt.GET("/me", handlers.Me)
			geschuetzt.PUT("/me", handlers.MeUpdate)
			geschuetzt.POST("/orders", handlers.BestellungAnlegen)
			geschuetzt.GET("/orders", handlers.Bestellungen)

			// Bewertungen schreiben (nur nach Bestellung)
			geschuetzt.POST("/restaurants/:slug/reviews", handlers.AddReview)

			// Favoriten verwalten
			geschuetzt.GET("/favorites", handlers.ListFavorites)
			geschuetzt.POST("/favorites/:slug", handlers.AddFavorite)
			geschuetzt.DELETE("/favorites/:slug", handlers.RemoveFavorite)
		}
	}

	log.Printf("🚀 Lieferino-Backend läuft auf Port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Server konnte nicht starten: %v", err)
	}
}
