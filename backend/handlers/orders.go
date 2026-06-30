package handlers

import (
	"fmt"
	"math/rand"
	"net/http"

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
	c.JSON(http.StatusCreated, order)
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
