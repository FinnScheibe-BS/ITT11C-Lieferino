package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 🛡️ Brute-Force-Schutz: nach so vielen Fehlversuchen wird das Konto kurz gesperrt.
const maxFehlversuche = 5
const sperrDauer = 15 * time.Minute

// Eingabe für die Registrierung.
type registerInput struct {
	Email        string `json:"email" binding:"required,email"`
	Passwort     string `json:"passwort" binding:"required"`
	Username     string `json:"username"`
	Vorname      string `json:"vorname"`
	Nachname     string `json:"nachname"`
	Geburtsdatum string `json:"geburtsdatum"`
}

// 📝 Registrierung: legt einen neuen Nutzer an (Passwort wird gehasht).
func Register(c *gin.Context) {
	var in registerInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte eine gültige E-Mail und ein Passwort angeben."})
		return
	}

	// 🔐 Passwort gegen die Anforderungen prüfen – mit klarer Rückmeldung,
	// was genau noch fehlt (statt nur "ungültig").
	if fehlt := passwortAnforderungen(in.Passwort); len(fehlt) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"fehler":        "Das Passwort erfüllt noch nicht alle Anforderungen.",
			"anforderungen": fehlt,
		})
		return
	}

	// Gibt es die E-Mail schon?
	var vorhanden models.User
	if err := database.DB.Where("email = ?", in.Email).First(&vorhanden).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"fehler": "Diese E-Mail ist bereits registriert."})
		return
	}

	// Passwort sicher hashen (bcrypt).
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Passwort), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "Passwort konnte nicht verarbeitet werden"})
		return
	}

	user := models.User{
		Email:        in.Email,
		PasswortHash: string(hash),
		Username:     in.Username,
		Vorname:      in.Vorname,
		Nachname:     in.Nachname,
		Geburtsdatum: in.Geburtsdatum,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "Nutzer konnte nicht angelegt werden"})
		return
	}

	token, _ := erstelleToken(user.ID)
	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

// Eingabe für den Login.
type loginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Passwort string `json:"passwort" binding:"required"`
}

// 🔑 Login: prüft E-Mail + Passwort und gibt ein JWT zurück.
func Login(c *gin.Context) {
	var in loginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte E-Mail und Passwort angeben."})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", in.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"fehler": "E-Mail oder Passwort falsch"})
		return
	}

	// Vom Admin dauerhaft gesperrt?
	if user.Gesperrt {
		c.JSON(http.StatusForbidden, gin.H{"fehler": "Dieses Konto wurde gesperrt."})
		return
	}

	// 🛡️ Wegen zu vieler Fehlversuche temporär gesperrt?
	if user.GesperrtBis != nil && time.Now().Before(*user.GesperrtBis) {
		rest := int(time.Until(*user.GesperrtBis).Minutes()) + 1
		c.JSON(http.StatusTooManyRequests, gin.H{
			"fehler": fmt.Sprintf("Zu viele Fehlversuche. Konto ist für noch ca. %d Minute(n) gesperrt.", rest),
		})
		return
	}

	// Passwort prüfen.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswortHash), []byte(in.Passwort)); err != nil {
		user.Fehlversuche++
		if user.Fehlversuche >= maxFehlversuche {
			bis := time.Now().Add(sperrDauer)
			user.GesperrtBis = &bis
			user.Fehlversuche = 0
			database.DB.Save(&user)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"fehler": fmt.Sprintf("Zu viele Fehlversuche. Konto ist jetzt %d Minuten gesperrt.", int(sperrDauer.Minutes())),
			})
			return
		}
		database.DB.Save(&user)
		c.JSON(http.StatusUnauthorized, gin.H{
			"fehler": fmt.Sprintf("E-Mail oder Passwort falsch. Noch %d Versuch(e) bis zur Sperre.", maxFehlversuche-user.Fehlversuche),
		})
		return
	}

	// Erfolg: Zähler/Sperre zurücksetzen.
	if user.Fehlversuche != 0 || user.GesperrtBis != nil {
		user.Fehlversuche = 0
		user.GesperrtBis = nil
		database.DB.Save(&user)
	}

	token, _ := erstelleToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// Erstellt ein JWT-Token (24 Stunden gültig) für einen Nutzer.
func erstelleToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "dev-secret-bitte-aendern"
	}
	return []byte(s)
}
