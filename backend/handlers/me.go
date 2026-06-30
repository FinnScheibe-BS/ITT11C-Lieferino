package handlers

import (
	"net/http"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
)

// 👤 Me liefert die Daten des aktuell eingeloggten Nutzers (inkl. Adressen).
func Me(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var user models.User
	if err := database.DB.Preload("Adressen").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"fehler": "Nutzer nicht gefunden"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Eingabe zum Aktualisieren des Profils.
type meUpdate struct {
	Username     string           `json:"username"`
	Vorname      string           `json:"vorname"`
	Nachname     string           `json:"nachname"`
	Geburtsdatum string           `json:"geburtsdatum"`
	Adressen     []models.Address `json:"adressen"`
}

// ✏️ MeUpdate aktualisiert Profil + Adressen des eingeloggten Nutzers.
func MeUpdate(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var in meUpdate
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"fehler": "Nutzer nicht gefunden"})
		return
	}

	user.Username = in.Username
	user.Vorname = in.Vorname
	user.Nachname = in.Nachname
	user.Geburtsdatum = in.Geburtsdatum
	database.DB.Save(&user)

	// Adressen einfach komplett ersetzen (einfachste, robuste Variante).
	database.DB.Where("user_id = ?", userID).Delete(&models.Address{})
	for _, a := range in.Adressen {
		a.ID = 0
		a.UserID = userID
		database.DB.Create(&a)
	}

	database.DB.Preload("Adressen").First(&user, userID)
	c.JSON(http.StatusOK, user)
}
