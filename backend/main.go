package main

import (
	"log"
	"os"
	"time"

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
	seedMenu() // Verkäufer + Produkte in die DB laden (einmalig)

	r := gin.Default()
	// 🛡️ Keinem Proxy-Header blind vertrauen (Backend ist direkt erreichbar).
	r.SetTrustedProxies(nil)
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS())

	// Health-Check
	r.GET("/health", handlers.Health)

	// API-Routen
	api := r.Group("/api")
	{
		// 🛡️ Login/Registrierung gegen Brute-Force begrenzen (pro IP).
		auth := api.Group("/auth")
		auth.Use(middleware.RateLimit(10, time.Minute))
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			// 📧 E-Mail-Verifizierung
			auth.POST("/verify-email", handlers.VerifyEmail)
			auth.POST("/resend-code", handlers.ResendCode)
			// 🔐 MFA (TOTP): Einrichtung (Setup-Token) + Login-Schritt-2 (MFA-Token)
			auth.POST("/mfa/setup", handlers.MFASetup)
			auth.POST("/mfa/enable", handlers.MFAEnable)
			auth.POST("/mfa/verify", handlers.MFAVerify)
		}
		// 🛡️ E-Mail-Versand gegen Spam begrenzen (pro IP).
		api.POST("/email", middleware.RateLimit(5, time.Minute), handlers.SendeEmail)

		// Öffentlich: Verkäufer + Produkte (aus der DB) und Bewertungen lesen
		api.GET("/restaurants", handlers.ListRestaurants)
		api.GET("/restaurants/:slug/products", handlers.ListProducts)
		api.GET("/restaurants/:slug/reviews", handlers.ListReviews)
		api.GET("/reviews/schnitt", handlers.ReviewSchnitt)

		// 🔒 Geschützte Routen (brauchen ein gültiges JWT-Token)
		geschuetzt := api.Group("")
		geschuetzt.Use(middleware.Auth())
		{
			geschuetzt.GET("/me", handlers.Me)
			geschuetzt.PUT("/me", handlers.MeUpdate)
			geschuetzt.POST("/orders", handlers.BestellungAnlegen)
			geschuetzt.GET("/orders", handlers.Bestellungen)
			geschuetzt.GET("/orders/:nummer/status", handlers.BestellStatus)

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
