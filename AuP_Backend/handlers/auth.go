package handlers

import (
	"net/http"
	"time"
	
	"aup-backend/database"
	"aup-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("dein_super_geheimes_zufaelliges_geheimnis")

// ==========================================
// 1. REGISTRIERUNG
// ==========================================
func RegisterKunde(c *gin.Context) {
	var input models.RegisterRequest 
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Eingabedaten"})
		return
	}


	// Passwort sicher hashen
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Passwort), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Verschlüsseln des Passworts"})
		return
	}

	// 9 Spalten -> 9 Fragezeichen
query := "INSERT INTO Kunde (Vorname, Nachname, Email_Adresse, Passwort, Ort, PLZ, Hausnummer, Strasse, Telefonnummer) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

// Hier müssen alle 9 Felder in der gleichen Reihenfolge wie oben stehen
_, err = database.DB.Exec(query, 
    input.Vorname, 
    input.Nachname, 
    input.EmailAdresse, 
    hashedPassword, 
    input.Ort, 
    input.PLZ, 
    input.Hausnummer, 
    input.Strasse, 
    input.Telefonnummer,
)

if err != nil {
    // Gibt die Fehlermeldung direkt als Antwort an Curl zurück
    c.JSON(http.StatusInternalServerError, gin.H{"status": "Fehler", "details": err.Error()})
    return
}

	c.JSON(http.StatusCreated, gin.H{"message": "Kunde erfolgreich registriert!"})
}

// ==========================================
// 2. LOGIN (mit JWT Token)
// ==========================================
func LoginKunde(c *gin.Context) {
	var input models.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Eingabedaten"})
		return
	}

	var kundeID int
	var kundeVorname string
	var kundePasswort string
	var kundeEmail string

	// WICHTIG: Die Reihenfolge im SELECT muss EXAKT mit der im .Scan() übereinstimmen!
	// 1. PK_ID_Kunde, 2. Vorname, 3. Passwort, 4. Email_Adresse
	query := "SELECT PK_ID_Kunde, Vorname, Passwort, Email_Adresse FROM Kunde WHERE Email_Adresse = ?"
	
	// input.Email matcht dein LoginRequest-Struct aus structs.go
	err := database.DB.QueryRow(query, input.Email).Scan(&kundeID, &kundeVorname, &kundePasswort, &kundeEmail)
	if err != nil {
		// Wenn kein User gefunden wird
		c.JSON(http.StatusUnauthorized, gin.H{"error": "E-Mail oder Passwort falsch"})
		return
	}

	// Passwort-Hash abgleichen
	err = bcrypt.CompareHashAndPassword([]byte(kundePasswort), []byte(input.Passwort))
	if err != nil {
		// Wenn das Passwort nicht zum Hash passt
		c.JSON(http.StatusUnauthorized, gin.H{"error": "E-Mail oder Passwort falsch"})
		return
	}

	// JWT Token generieren
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"kunde_id": kundeID,
		"email":    kundeEmail,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Generieren des Tokens"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Message: "Login erfolgreich!",
		Token:   tokenString,
		KundeID: kundeID,
		Vorname: kundeVorname,
	})
}