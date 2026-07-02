package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

// 📧 GMAIL-ZUGANGSDATEN (im Klartext, außerhalb .env – mit Projektleitung abgesprochen).
// ⚠️ Das normale Konto-Passwort funktioniert NICHT für SMTP! Hier muss ein
//    16-stelliges Gmail-APP-PASSWORT stehen, damit echt verschickt wird.
//    Solange das App-Passwort leer ist, läuft der Versand im Test-Modus.
const standardGmailUser = "lieferino5@gmail.com"
const standardGmailAppPasswort = "ahzcpeqhzyyuykls" // 16-stelliges Gmail-App-Passwort

// gmailZugang liefert Nutzer + App-Passwort (erst Umgebungsvariablen, sonst die
// Klartext-Standardwerte oben).
func gmailZugang() (string, string) {
	user := os.Getenv("GMAIL_USER")
	if user == "" {
		user = standardGmailUser
	}
	pass := os.Getenv("GMAIL_APP_PASSWORD")
	if pass == "" {
		pass = standardGmailAppPasswort
	}
	return user, pass
}

// sendeMail verschickt eine E-Mail über Gmail-SMTP.
// Rückgabe: (true, nil) = echt verschickt, (false, nil) = Test-Modus (keine
// Zugangsdaten), (false, err) = Fehler beim Versand.
func sendeMail(an, betreff, text string) (bool, error) {
	user, pass := gmailZugang()
	if user == "" || pass == "" {
		return false, nil // Test-Modus
	}

	auth := smtp.PlainAuth("", user, pass, "smtp.gmail.com")
	nachricht := []byte(fmt.Sprintf(
		"From: Lieferino <%s>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		user, an, betreff, text,
	))
	if err := smtp.SendMail("smtp.gmail.com:587", auth, user, []string{an}, nachricht); err != nil {
		log.Printf("❌ E-Mail-Versand an %s fehlgeschlagen: %v", an, err)
		return false, err
	}
	return true, nil
}

// Eingabe für den E-Mail-Versand (gleiche Form wie im Frontend-Service).
type emailInput struct {
	An      string `json:"an" binding:"required,email"`
	Betreff string `json:"betreff"`
	Text    string `json:"text"`
}

// 📧 Verschickt eine echte E-Mail über das Gmail-Konto (SMTP).
// Ohne gesetzte Zugangsdaten läuft es im "Test-Modus" (sendet nichts).
func SendeEmail(c *gin.Context) {
	var in emailInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": err.Error()})
		return
	}

	gesendet, err := sendeMail(in.An, in.Betreff, in.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"gesendet": false, "fehler": err.Error()})
		return
	}
	if !gesendet {
		c.JSON(http.StatusOK, gin.H{"gesendet": false, "testModus": true})
		return
	}
	c.JSON(http.StatusOK, gin.H{"gesendet": true})
}
