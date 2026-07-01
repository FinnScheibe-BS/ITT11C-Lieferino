package database

import (
	"fmt"
	"log"

	"lieferino-backend/config"
	"lieferino-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB ist die globale Datenbank-Verbindung, die alle Handler nutzen.
var DB *gorm.DB

// Verbinden baut die Verbindung zur PostgreSQL-Datenbank auf und legt die
// Tabellen automatisch an (AutoMigrate).
func Verbinden(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPasswort, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Datenbank-Verbindung fehlgeschlagen: %v", err)
	}

	// Tabellen automatisch anlegen / aktualisieren.
	if err := db.AutoMigrate(&models.User{}, &models.Address{}, &models.Order{}, &models.OrderItem{}, &models.Review{}, &models.Favorite{}, &models.Restaurant{}, &models.Product{}); err != nil {
		log.Fatalf("❌ Migration fehlgeschlagen: %v", err)
	}

	DB = db
	log.Println("✅ Mit der Datenbank verbunden + migriert")
}
