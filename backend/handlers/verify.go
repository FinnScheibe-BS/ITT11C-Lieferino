package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
)

// 📧 E-MAIL-VERIFIZIERUNG
// Beim Registrieren/Login (unbestätigt) bekommt der Nutzer einen 6-stelligen
// Code per E-Mail. Erst nach Eingabe gilt die E-Mail als bestätigt.

const codeGueltigkeit = 15 * time.Minute

// neuerCode erzeugt einen zufälligen 6-stelligen Code (z.B. "048213").
func neuerCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "000000"
	}
	return fmt.Sprintf("%06d", n.Int64())
}

// sendeVerifizierungsCode erzeugt einen Code, speichert ihn (mit Ablaufzeit) in
// der DB und verschickt ihn per E-Mail.
func sendeVerifizierungsCode(user *models.User) {
	code := neuerCode()
	ablauf := time.Now().Add(codeGueltigkeit)
	user.EmailCode = code
	user.EmailCodeAblauf = &ablauf
	database.DB.Save(user)

	text := fmt.Sprintf(
		"Hallo%s,\r\n\r\ndein Lieferino-Bestätigungscode lautet:\r\n\r\n    %s\r\n\r\nDer Code ist 15 Minuten gültig.\r\n\r\nViele Grüße\r\nDein Lieferino-Team 🍕",
		nameOderLeer(user), code,
	)
	// Versand-Fehler bewusst ignorieren (Konto ist trotzdem angelegt; der Nutzer
	// kann den Code erneut anfordern).
	_, _ = sendeMail(user.Email, "Dein Lieferino-Bestätigungscode", text)
}

func nameOderLeer(user *models.User) string {
	if user.Vorname != "" {
		return " " + user.Vorname
	}
	return ""
}

type verifyInput struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// ✅ VerifyEmail prüft den eingegebenen Code. Stimmt er (und ist nicht
// abgelaufen), gilt die E-Mail als bestätigt. Danach folgt die MFA-Einrichtung.
func VerifyEmail(c *gin.Context) {
	var in verifyInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte E-Mail und Code angeben."})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", in.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Code ungültig."})
		return
	}

	if user.EmailVerifiziert {
		c.JSON(http.StatusOK, gin.H{"verifiziert": true, "hinweis": "E-Mail war bereits bestätigt."})
		return
	}

	if user.EmailCode == "" || user.EmailCodeAblauf == nil || time.Now().After(*user.EmailCodeAblauf) {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Der Code ist abgelaufen. Bitte fordere einen neuen an."})
		return
	}
	if in.Code != user.EmailCode {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Der Code ist nicht korrekt."})
		return
	}

	user.EmailVerifiziert = true
	user.EmailCode = ""
	user.EmailCodeAblauf = nil
	database.DB.Save(&user)

	// Nächster Pflicht-Schritt: MFA einrichten -> dafür ein Setup-Token geben.
	st, _ := setupToken(user.ID)
	c.JSON(http.StatusOK, gin.H{
		"verifiziert":   true,
		"needsMfaSetup": true,
		"setupToken":    st,
	})
}

type resendInput struct {
	Email string `json:"email" binding:"required,email"`
}

// 🔁 ResendCode schickt einen neuen Bestätigungscode (falls noch nicht bestätigt).
func ResendCode(c *gin.Context) {
	var in resendInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte E-Mail angeben."})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", in.Email).First(&user).Error; err != nil {
		// Aus Datenschutz-Gründen nicht verraten, ob es die E-Mail gibt.
		c.JSON(http.StatusOK, gin.H{"gesendet": true})
		return
	}
	if user.EmailVerifiziert {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Diese E-Mail ist bereits bestätigt."})
		return
	}

	sendeVerifizierungsCode(&user)
	c.JSON(http.StatusOK, gin.H{"gesendet": true})
}
