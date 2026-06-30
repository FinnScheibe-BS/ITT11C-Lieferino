package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
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

// ───────────────────────────────────────────────────────────────────────────
// 🎫 TOKEN-TYPEN
// "full"  = voller Zugang (24 h). Nur damit kommt man an geschützte Routen.
// "setup" = darf NUR die MFA-Einrichtung machen (kurzlebig).
// "mfa"   = Passwort war richtig, es fehlt nur noch der MFA-Code (kurzlebig).
// ───────────────────────────────────────────────────────────────────────────

func tokenMitTyp(userID uint, typ string, dauer time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"typ": typ,
		"exp": time.Now().Add(dauer).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret())
}

func vollToken(userID uint) (string, error)  { return tokenMitTyp(userID, "full", 24*time.Hour) }
func setupToken(userID uint) (string, error) { return tokenMitTyp(userID, "setup", 15*time.Minute) }
func mfaToken(userID uint) (string, error)   { return tokenMitTyp(userID, "mfa", 5*time.Minute) }

// userAusToken liest die User-ID aus dem Bearer-Token UND prüft, dass es vom
// erwarteten Typ ist (z.B. "setup"). So kann man mit einem Setup-Token nicht
// an die echten Daten kommen.
func userAusToken(c *gin.Context, erwarteterTyp string) (uint, bool) {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return 0, false
	}
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret(), nil
	})
	if err != nil || !token.Valid {
		return 0, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["typ"] != erwarteterTyp {
		return 0, false
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, false
	}
	return uint(sub), true
}

// Eingabe für die Registrierung.
type registerInput struct {
	Email        string `json:"email" binding:"required,email"`
	Passwort     string `json:"passwort" binding:"required"`
	Username     string `json:"username"`
	Vorname      string `json:"vorname"`
	Nachname     string `json:"nachname"`
	Geburtsdatum string `json:"geburtsdatum"`
}

// 📝 Registrierung: legt einen neuen Nutzer an (Passwort gehasht) und schickt
// einen E-Mail-Bestätigungscode. Es gibt NOCH KEIN Zugangstoken – das Konto ist
// erst nutzbar, wenn E-Mail bestätigt UND MFA eingerichtet wurde.
func Register(c *gin.Context) {
	var in registerInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Bitte eine gültige E-Mail und ein Passwort angeben."})
		return
	}

	// 🔐 Passwort gegen die Anforderungen prüfen – mit klarer Rückmeldung.
	if fehlt := passwortAnforderungen(in.Passwort); len(fehlt) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"fehler":        "Das Passwort erfüllt noch nicht alle Anforderungen.",
			"anforderungen": fehlt,
		})
		return
	}

	var vorhanden models.User
	if err := database.DB.Where("email = ?", in.Email).First(&vorhanden).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"fehler": "Diese E-Mail ist bereits registriert."})
		return
	}

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

	// 📧 Bestätigungscode erzeugen + per E-Mail schicken.
	sendeVerifizierungsCode(&user)

	c.JSON(http.StatusCreated, gin.H{
		"needsVerification": true,
		"email":             user.Email,
		"hinweis":           "Wir haben dir einen Bestätigungscode per E-Mail geschickt.",
	})
}

// Eingabe für den Login.
type loginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Passwort string `json:"passwort" binding:"required"`
}

// 🔑 Login Schritt 1: prüft E-Mail + Passwort. Bei Erfolg wird NICHT sofort
// eingeloggt – es kommt erst die MFA-Abfrage (oder ein Hinweis, dass noch
// E-Mail-Bestätigung / MFA-Einrichtung fehlt).
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

	if user.Gesperrt {
		c.JSON(http.StatusForbidden, gin.H{"fehler": "Dieses Konto wurde gesperrt."})
		return
	}

	if user.GesperrtBis != nil && time.Now().Before(*user.GesperrtBis) {
		rest := int(time.Until(*user.GesperrtBis).Minutes()) + 1
		c.JSON(http.StatusTooManyRequests, gin.H{
			"fehler": fmt.Sprintf("Zu viele Fehlversuche. Konto ist für noch ca. %d Minute(n) gesperrt.", rest),
		})
		return
	}

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

	// Passwort stimmt -> Fehlversuche zurücksetzen.
	if user.Fehlversuche != 0 || user.GesperrtBis != nil {
		user.Fehlversuche = 0
		user.GesperrtBis = nil
		database.DB.Save(&user)
	}

	// 📧 E-Mail noch nicht bestätigt? -> neuen Code schicken, Zugang verweigern.
	if !user.EmailVerifiziert {
		sendeVerifizierungsCode(&user)
		c.JSON(http.StatusForbidden, gin.H{
			"fehler":            "Deine E-Mail ist noch nicht bestätigt. Wir haben dir einen neuen Code geschickt.",
			"needsVerification": true,
			"email":             user.Email,
		})
		return
	}

	// 🔐 MFA noch nicht eingerichtet? -> Setup-Token geben, Zugang verweigern.
	if !user.MFAAktiv {
		st, _ := setupToken(user.ID)
		c.JSON(http.StatusForbidden, gin.H{
			"fehler":        "Bitte richte zuerst die Zwei-Faktor-Authentifizierung (MFA) ein.",
			"needsMfaSetup": true,
			"setupToken":    st,
		})
		return
	}

	// ✅ Alles da -> jetzt nur noch der MFA-Code (Schritt 2).
	mt, _ := mfaToken(user.ID)
	c.JSON(http.StatusOK, gin.H{
		"mfaRequired": true,
		"mfaToken":    mt,
	})
}

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "dev-secret-bitte-aendern"
	}
	return []byte(s)
}
