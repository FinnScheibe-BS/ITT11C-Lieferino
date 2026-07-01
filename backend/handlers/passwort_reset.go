package handlers

import (
	"fmt"
	"net/http"
	"time"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 🔑 PASSWORT-RESET (per E-Mail-Einmalcode)
// Ablauf:
//   1. /auth/reset-anfordern {email}          -> Code per E-Mail
//   2. /auth/reset-code       {email, code}   -> Code prüfen -> Reset-Token
//   3. /auth/reset-neu        {passwort}      -> neues Passwort setzen (Reset-Token)

type resetAnfordernInput struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetAnfordern schickt einen Einmal-Code an die E-Mail (falls es das Konto gibt).
// Aus Datenschutz-Gründen antworten wir IMMER gleich (kein Verraten, ob es die E-Mail gibt).
func ResetAnfordern(c *gin.Context) {
	var in resetAnfordernInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte eine gültige E-Mail angeben."})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", in.Email).First(&user).Error; err == nil {
		code := neuerCode()
		ablauf := time.Now().Add(codeGueltigkeit)
		user.ResetCode = code
		user.ResetCodeAblauf = &ablauf
		database.DB.Save(&user)

		text := fmt.Sprintf(
			"Hallo%s,\r\n\r\ndu möchtest dein Lieferino-Passwort zurücksetzen.\r\nDein Einmal-Code lautet:\r\n\r\n    %s\r\n\r\nDer Code ist 15 Minuten gültig. Wenn du das nicht warst, kannst du diese Mail ignorieren.\r\n\r\nDein Lieferino-Team 🍕",
			nameOderLeer(&user), code,
		)
		_, _ = sendeMail(user.Email, "Dein Lieferino-Reset-Code", text)
	}

	c.JSON(http.StatusOK, gin.H{"gesendet": true})
}

type resetCodeInput struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// ResetCodePruefen prüft den Einmal-Code und gibt bei Erfolg ein Reset-Token zurück.
func ResetCodePruefen(c *gin.Context) {
	var in resetCodeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte E-Mail und Code angeben."})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", in.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Der Code ist nicht korrekt."})
		return
	}
	if user.ResetCode == "" || user.ResetCodeAblauf == nil || time.Now().After(*user.ResetCodeAblauf) {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Der Code ist abgelaufen. Bitte fordere einen neuen an."})
		return
	}
	if in.Code != user.ResetCode {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Der Code ist nicht korrekt."})
		return
	}

	rt, _ := resetToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"resetToken": rt})
}

type resetNeuInput struct {
	Passwort string `json:"passwort" binding:"required"`
}

// ResetNeuesPasswort setzt das neue Passwort (braucht ein gültiges Reset-Token).
func ResetNeuesPasswort(c *gin.Context) {
	userID, ok := userAusToken(c, "reset")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"fehler": "Der Reset-Vorgang ist abgelaufen. Bitte starte ihn neu."})
		return
	}

	var in resetNeuInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte ein neues Passwort angeben."})
		return
	}
	if fehlt := passwortAnforderungen(in.Passwort); len(fehlt) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"fehler":        "Das Passwort erfüllt noch nicht alle Anforderungen.",
			"anforderungen": fehlt,
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"fehler": "Nutzer nicht gefunden"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Passwort), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "Passwort konnte nicht verarbeitet werden"})
		return
	}
	user.PasswortHash = string(hash)
	user.ResetCode = ""
	user.ResetCodeAblauf = nil
	// Sicherheit: Fehlversuch-Sperre zurücksetzen.
	user.Fehlversuche = 0
	user.GesperrtBis = nil
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
