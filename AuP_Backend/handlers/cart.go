package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"aup-backend/database"
	"aup-backend/models"
)

func ValidateGutschein(c *gin.Context) {
	var input models.GutscheinRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Ungültiger Code"})
		return
	}

	switch input.Code {
	case "PROMO10":
		c.JSON(200, gin.H{"valid": true, "discount_percent": 10, "message": "10% Rabatt abgezogen!"})
	case "PROMO20":
		c.JSON(200, gin.H{"valid": true, "discount_percent": 20, "message": "20% Rabatt abgezogen!"})
	default:
		c.JSON(400, gin.H{"valid": false, "discount_percent": 0, "error": "Gutscheincode ist ungültig oder abgelaufen."})
	}
}

func GetCart(c *gin.Context) {
	kundeID, exists := c.Get("kunde_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nicht autorisiert"})
		return
	}

	gutscheinCode := c.Query("gutschein")
	var rabattProzent float64 = 0

	if gutscheinCode != "" {
		queryGutschein := "SELECT Rabatt_Prozent FROM Gutschein WHERE Code = ? AND Gueltig = 1"
		err := database.DB.QueryRow(queryGutschein, gutscheinCode).Scan(&rabattProzent)
		if err != nil {
			rabattProzent = 0
		}
	}

	query := `
		SELECT b.FK_ID_Gericht, g.Name, g.Preis, COUNT(b.FK_ID_Gericht) as Anzahl 
		FROM Bestellung b 
		JOIN Gericht g ON b.FK_ID_Gericht = g.PK_ID_Gericht 
		WHERE b.FK_ID_Kunde = ?
		GROUP BY b.FK_ID_Gericht, g.Name, g.Preis`
	
	rows, err := database.DB.Query(query, kundeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bestellungen konnten nicht geladen werden", "details": err.Error()})
		return
	}
	defer rows.Close()

	var artikel []gin.H
	var gesamtpreisNormal float64 = 0

	for rows.Next() {
		var id int
		var name string
		var preis float64
		var anzahl int

		if err := rows.Scan(&id, &name, &preis, &anzahl); err == nil {
			summeArtikel := preis * float64(anzahl)
			gesamtpreisNormal += summeArtikel

			artikel = append(artikel, gin.H{
				"gericht_id": id,
				"name":       name,
				"preis":      preis,
				"anzahl":     anzahl,
				"total":      summeArtikel,
			})
		}
	}

	rabattBetrag := gesamtpreisNormal * (rabattProzent / 100.0)
	endpreis := gesamtpreisNormal - rabattBetrag

	c.JSON(http.StatusOK, gin.H{
		"artikel":             artikel,
		"gesamtpreis_normal":  gesamtpreisNormal,
		"gutschein_angewandt": gutscheinCode,
		"rabatt_prozent":      rabattProzent,
		"gespart":             rabattBetrag,
		"endpreis":            endpreis,
	})
}

func AddToCart(c *gin.Context) {
	kundeID, exists := c.Get("kunde_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nicht autorisiert"})
		return
	}

	var input models.CartItemRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ungültige Gericht-ID"})
		return
	}

	location, _ := time.LoadLocation("Europe/Berlin")
	jetzt := time.Now().In(location)
	datum := jetzt.Format("2006-01-02")
	uhrzeit := jetzt.Format("15:04")
	restaurantID := 1 // Dummy-ID für die Bestellungstabelle

	query := "INSERT INTO Bestellung (FK_ID_Restaurant, FK_ID_Kunde, FK_ID_Gericht, Datum, Uhrzeit) VALUES (?, ?, ?, ?, ?)"
	
	for i := 0; i < input.Anzahl; i++ {
		_, err := database.DB.Exec(query, restaurantID, kundeID, input.GerichtID, datum, uhrzeit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Konnte nicht in Bestellung gespeichert werden", "details": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gericht erfolgreich hinzugefügt!"})
}

func CheckoutCart(c *gin.Context) {
	// Da du die Artikel nun direkt über AddToCart in der Tabelle "Bestellung" speicherst,
	// ist der Checkout hier im Grunde schon beim Hinzufügen passiert.
	c.JSON(http.StatusOK, gin.H{"message": "🎉 Checkout erfolgreich (wird direkt gebucht)!"})
}