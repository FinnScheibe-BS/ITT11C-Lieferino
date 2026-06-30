package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"lieferino-backend/database"
	"lieferino-backend/models"

	"github.com/gin-gonic/gin"
)

type orderItemInput struct {
	Name       string  `json:"name"`
	Preis      float64 `json:"preis"`
	Menge      int     `json:"menge"`
	Restaurant string  `json:"restaurant"`
}

type orderInput struct {
	Zwischensumme float64          `json:"zwischensumme"`
	Trinkgeld     float64          `json:"trinkgeld"`
	Summe         float64          `json:"summe"`
	Gutschein     string           `json:"gutschein"`
	Zahlungsart   string           `json:"zahlungsart"`
	Liefertermin  string           `json:"liefertermin"`
	Positionen    []orderItemInput `json:"positionen"`
}

// 📦 BestellungAnlegen speichert eine neue Bestellung für den eingeloggten Nutzer.
func BestellungAnlegen(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var in orderInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": err.Error()})
		return
	}

	// Bestellnummer erzeugen, z.B. "LF-7F3A9C".
	nummer := fmt.Sprintf("LF-%06X", rand.Intn(0xFFFFFF))

	order := models.Order{
		UserID:        userID,
		Nummer:        nummer,
		Summe:         in.Summe,
		Zwischensumme: in.Zwischensumme,
		Trinkgeld:     in.Trinkgeld,
		Gutschein:     in.Gutschein,
		Zahlungsart:   in.Zahlungsart,
		Liefertermin:  in.Liefertermin,
		Status:        "erhalten",
	}
	for _, p := range in.Positionen {
		order.Positionen = append(order.Positionen, models.OrderItem{
			Name:       p.Name,
			Preis:      p.Preis,
			Menge:      p.Menge,
			Restaurant: p.Restaurant,
		})
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"fehler": "Bestellung konnte nicht gespeichert werden"})
		return
	}

	// 📧 Bestellbestätigung per E-Mail (im Hintergrund, blockiert die Antwort nicht).
	var user models.User
	if err := database.DB.First(&user, userID).Error; err == nil && user.Email != "" {
		betreff := "Deine Lieferino-Bestellung " + order.Nummer + " 🍕"
		text := bestellBestaetigungText(&user, &order)
		go func() { _, _ = sendeMail(user.Email, betreff, text) }()
	}

	c.JSON(http.StatusCreated, order)
}

// Baut den Text der Bestellbestätigungs-E-Mail zusammen.
func bestellBestaetigungText(user *models.User, order *models.Order) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Hallo%s,\r\n\r\nvielen Dank für deine Bestellung bei Lieferino!\r\n\r\n", nameOderLeer(user))
	fmt.Fprintf(&b, "Bestellnummer: %s\r\n\r\nDeine Artikel:\r\n", order.Nummer)
	for _, p := range order.Positionen {
		fmt.Fprintf(&b, "  - %dx %s (%.2f €)\r\n", p.Menge, p.Name, p.Preis*float64(p.Menge))
	}
	if order.Trinkgeld > 0 {
		fmt.Fprintf(&b, "  - Trinkgeld: %.2f €\r\n", order.Trinkgeld)
	}
	fmt.Fprintf(&b, "\r\nGesamt: %.2f €\r\n", order.Summe)
	if order.Liefertermin != "" {
		fmt.Fprintf(&b, "Liefertermin: %s Uhr\r\n", order.Liefertermin)
	}
	fmt.Fprintf(&b, "Bezahlt mit: %s\r\n\r\nGuten Appetit! 🍕\r\nDein Lieferino-Team", order.Zahlungsart)
	return b.String()
}

// 🧾 Bestellungen liefert alle Bestellungen des eingeloggten Nutzers (neueste zuerst).
func Bestellungen(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var orders []models.Order
	database.DB.Preload("Positionen").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders)

	c.JSON(http.StatusOK, orders)
}
