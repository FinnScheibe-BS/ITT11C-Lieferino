package handlers

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

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

	user := os.Getenv("GMAIL_USER")
	pass := os.Getenv("GMAIL_APP_PASSWORD")

	// Keine Zugangsdaten -> Test-Modus (es wird NICHTS verschickt).
	if user == "" || pass == "" {
		c.JSON(http.StatusOK, gin.H{"gesendet": false, "testModus": true})
		return
	}

	// Gmail-SMTP mit App-Passwort.
	auth := smtp.PlainAuth("", user, pass, "smtp.gmail.com")
	nachricht := []byte(fmt.Sprintf(
		"From: Lieferino <%s>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		user, in.An, in.Betreff, in.Text,
	))

	if err := smtp.SendMail("smtp.gmail.com:587", auth, user, []string{in.An}, nachricht); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"gesendet": false, "fehler": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"gesendet": true})
}
