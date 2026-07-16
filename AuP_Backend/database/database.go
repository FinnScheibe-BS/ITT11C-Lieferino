package database

import (
	"database/sql"
	"fmt"
	"log"
	"time" // <-- Hinzugefügt für das time.Sleep
)

// DB muss großgeschrieben sein, damit die Handler darauf zugreifen können
var DB *sql.DB

// InitDB initialisiert die Datenbankverbindung mit einer Warteschleife
func InitDB() {
	var err error
	
	// Falls dein MySQL-Service in docker-compose.yml anders heißt (z.B. "mysql-db"), passe es hier an!
	
	dsn := "root:rootpass@tcp(mysql-db:3306)/restaurant_db"

	// Bis zu 10 Versuche, um MySQL Zeit zum Booten zu geben
	for i := 1; i <= 10; i++ {
		fmt.Printf("🔄 Verbindungsversuch zur Datenbank (%d/10)...\n", i)
		
		DB, err = sql.Open("mysql", dsn)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				fmt.Println("🚀 Erfolgreich mit der MySQL-Datenbank verbunden!")
				return // Verbindung steht, wir können die Funktion beenden!
			}
		}

		fmt.Printf("⚠️ Datenbank noch nicht bereit, neuer Versuch in 2 Sekunden... (%v)\n", err)
		time.Sleep(2 * time.Second) 
	}

	// Falls nach 10 Versuchen immer noch nichts geht
	log.Fatalf("❌ Kritischer Fehler: Datenbank nach mehreren Versuchen nicht erreichbar!")
}