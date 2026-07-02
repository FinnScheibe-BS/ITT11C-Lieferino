package handlers

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"net/http"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

// 🔐 MFA (Zwei-Faktor-Authentifizierung) per TOTP / Authenticator-App.
// Ablauf:
//   1. /mfa/setup  (Setup-Token) -> Secret + QR-Code zum Scannen
//   2. /mfa/enable (Setup-Token) -> Code aus der App prüfen, MFA aktivieren
//   3. /mfa/verify (MFA-Token)   -> beim Login den Code prüfen, Zugang geben

// 🔧 MFASetup erzeugt ein neues TOTP-Geheimnis + QR-Code (Data-URL).
// Braucht ein gültiges Setup-Token (E-Mail muss bestätigt sein).
func MFASetup(c *gin.Context) {
	userID, ok := userAusToken(c, "setup")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"fehler": "Setup-Sitzung ungültig oder abgelaufen. Bitte neu anmelden."})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"fehler": "Nutzer nicht gefunden"})
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Lieferino",
		AccountName: user.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "MFA konnte nicht vorbereitet werden"})
		return
	}

	// Geheimnis in der DB speichern (MFA ist erst nach /mfa/enable aktiv).
	user.MFASecret = key.Secret()
	database.DB.Save(&user)

	// QR-Code als PNG erzeugen und als Data-URL einbetten (zeigt das Frontend an).
	img, err := key.Image(220, 220)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "QR-Code konnte nicht erzeugt werden"})
		return
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "QR-Code konnte nicht erzeugt werden"})
		return
	}
	qrDataUrl := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	c.JSON(http.StatusOK, gin.H{
		"secret":     key.Secret(), // zum manuellen Eintippen in die App
		"otpauthUrl": key.URL(),
		"qr":         qrDataUrl,
	})
}

type mfaCodeInput struct {
	Code string `json:"code" binding:"required"`
}

// ✅ MFAEnable prüft den ersten Code aus der App und schaltet MFA scharf.
// Danach ist das Konto vollständig nutzbar -> es kommt ein volles Token zurück.
func MFAEnable(c *gin.Context) {
	userID, ok := userAusToken(c, "setup")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"fehler": "Setup-Sitzung ungültig oder abgelaufen. Bitte neu anmelden."})
		return
	}

	var in mfaCodeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte den 6-stelligen Code aus deiner App eingeben."})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"fehler": "Nutzer nicht gefunden"})
		return
	}
	if user.MFASecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte zuerst die MFA-Einrichtung starten."})
		return
	}
	if !totp.Validate(in.Code, user.MFASecret) {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Der Code stimmt nicht. Bitte den aktuellen Code aus der App eingeben."})
		return
	}

	user.MFAAktiv = true
	database.DB.Save(&user)

	token, _ := vollToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// 🔑 MFAVerify ist Login-Schritt 2: prüft den Code beim Anmelden.
// Braucht ein MFA-Token (Passwort war schon richtig).
func MFAVerify(c *gin.Context) {
	userID, ok := userAusToken(c, "mfa")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"fehler": "Anmeldung abgelaufen. Bitte erneut einloggen."})
		return
	}

	var in mfaCodeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte den 6-stelligen Code aus deiner App eingeben."})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"fehler": "Nutzer nicht gefunden"})
		return
	}
	if !user.MFAAktiv || user.MFASecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Für dieses Konto ist keine MFA aktiv."})
		return
	}
	if !totp.Validate(in.Code, user.MFASecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"fehler": "Der Code stimmt nicht. Bitte den aktuellen Code aus der App eingeben."})
		return
	}

	token, _ := vollToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
