package handlers

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

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
	Zwischensumme   float64          `json:"zwischensumme"`
	Trinkgeld       float64          `json:"trinkgeld"`
	Summe           float64          `json:"summe"`
	Gutschein       string           `json:"gutschein"`
	PunkteEinloesen bool             `json:"punkteEinloesen"`
	Zahlungsart     string           `json:"zahlungsart"`
	Liefertermin    string           `json:"liefertermin"`
	Positionen      []orderItemInput `json:"positionen"`
}

// 🎟️ Server-seitige Regeln (müssen mit dem Checkout übereinstimmen).
const lieferGebuehr = 2.49
const punkteSchritt = 100  // so viele Punkte kann man einlösen ...
const punkteWert = 5.0     // ... = so viel € Rabatt
const punkteProEuro = 1

type gutscheinDef struct {
	Typ  string // "prozent" | "betrag" | "lieferung"
	Wert float64
}

var gutscheine = map[string]gutscheinDef{
	"LIEFERINO10": {"prozent", 10},
	"WILLKOMMEN5": {"betrag", 5},
	"GRATIS":      {"lieferung", 0},
	"PIZZAPARTY":  {"prozent", 25},
	"GEBURTSTAG":  {"prozent", 20},
}

// 📦 BestellungAnlegen speichert eine neue Bestellung für den eingeloggten Nutzer.
func BestellungAnlegen(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var in orderInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": err.Error()})
		return
	}

	// 🛡️ PREISPRÜFUNG: Jedes Produkt + jeder Preis wird gegen die DATENBANK
	// geprüft. So kann niemand im Frontend (z.B. per Inspektor) den Preis
	// manipulieren und teure Sachen billig kaufen.
	if len(in.Positionen) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"fehler": "Die Bestellung enthält keine Artikel."})
		return
	}
	var echteZwischensumme float64
	for i := range in.Positionen {
		pos := &in.Positionen[i]
		if pos.Menge < 1 {
			pos.Menge = 1
		}
		var prod models.Product
		if err := database.DB.
			Where("(restaurant_name = ? OR restaurant_slug = ?) AND name = ?", pos.Restaurant, pos.Restaurant, pos.Name).
			First(&prod).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"fehler": fmt.Sprintf("Produkt \"%s\" wurde nicht gefunden und kann nicht bestellt werden.", pos.Name),
			})
			return
		}
		// Der Preis aus der DB ist maßgeblich (überschreibt den vom Client).
		pos.Preis = prod.Preis
		echteZwischensumme += prod.Preis * float64(pos.Menge)
	}
	// Hat der Client eine andere Zwischensumme behauptet? -> Manipulation -> ablehnen.
	if math.Abs(echteZwischensumme-in.Zwischensumme) > 0.01 {
		c.JSON(http.StatusBadRequest, gin.H{
			"fehler": "Die Preise stimmen nicht mit unseren Daten überein. Bitte lade die Seite neu.",
		})
		return
	}
	in.Zwischensumme = echteZwischensumme
	if in.Trinkgeld < 0 {
		in.Trinkgeld = 0
	}

	// 🛡️ Mindestbestellwert des Restaurants serverseitig prüfen.
	var rest models.Restaurant
	if err := database.DB.
		Where("name = ? OR slug = ?", in.Positionen[0].Restaurant, in.Positionen[0].Restaurant).
		First(&rest).Error; err == nil && echteZwischensumme < rest.MinBestell {
		c.JSON(http.StatusBadRequest, gin.H{
			"fehler": fmt.Sprintf("Mindestbestellwert von %.2f € nicht erreicht.", rest.MinBestell),
		})
		return
	}

	// 🎟️ Gutschein serverseitig auswerten (nur bekannte Codes wirken).
	rabatt := 0.0
	lieferkosten := lieferGebuehr
	code := strings.ToUpper(strings.TrimSpace(in.Gutschein))
	if g, ok := gutscheine[code]; ok {
		switch g.Typ {
		case "prozent":
			rabatt = echteZwischensumme * g.Wert / 100
		case "betrag":
			rabatt = math.Min(g.Wert, echteZwischensumme)
		case "lieferung":
			lieferkosten = 0
		}
		in.Gutschein = code
	} else {
		in.Gutschein = "" // ungültiger/leerer Code -> kein Rabatt
	}

	// ⭐ Treuepunkte einlösen (nur wenn genug vorhanden) – serverseitig geprüft.
	var user models.User
	database.DB.First(&user, userID)
	punkteRabatt := 0.0
	punkteAbzug := 0
	if in.PunkteEinloesen && user.Treuepunkte >= punkteSchritt {
		punkteRabatt = punkteWert
		punkteAbzug = punkteSchritt
	}

	// 🧮 Autoritative Endsumme (NICHT der Client-Wert).
	in.Summe = echteZwischensumme - rabatt - punkteRabatt + lieferkosten + in.Trinkgeld
	if in.Summe < 0 {
		in.Summe = 0
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

	// ⭐ Treuepunkte aktualisieren: eingelöste abziehen + neue sammeln (1 pro € Artikelwert).
	neuePunkte := user.Treuepunkte - punkteAbzug + int(echteZwischensumme)*punkteProEuro
	if neuePunkte < 0 {
		neuePunkte = 0
	}
	database.DB.Model(&models.User{}).Where("id = ?", userID).Update("treuepunkte", neuePunkte)

	// 📧 Bestellbestätigung per E-Mail (im Hintergrund, blockiert die Antwort nicht).
	if user.Email != "" {
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

// 🚚 Liefer-Phasen + serverseitige Status-Berechnung aus der vergangenen Zeit.
// So ist der Status manipulationssicher (kein Frontend-Wert) und braucht keinen
// Hintergrund-Job. Gesamte Lieferzeit: 30 Minuten, in 4 Phasen.
var lieferPhasen = []string{"Bestellung erhalten", "Wird zubereitet", "Unterwegs", "Geliefert"}

func statusBerechnen(erstellt time.Time) (int, string) {
	const gesamt = 30 * time.Minute
	anteil := float64(time.Since(erstellt)) / float64(gesamt)
	if anteil < 0 {
		anteil = 0
	}
	idx := int(anteil * float64(len(lieferPhasen)))
	if idx > len(lieferPhasen)-1 {
		idx = len(lieferPhasen) - 1
	}
	return idx, lieferPhasen[idx]
}

// 🧾 Bestellungen liefert alle Bestellungen des eingeloggten Nutzers (neueste zuerst).
func Bestellungen(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var orders []models.Order
	database.DB.Preload("Positionen").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&orders)

	// Aktuellen Status je Bestellung serverseitig setzen.
	for i := range orders {
		_, orders[i].Status = statusBerechnen(orders[i].CreatedAt)
	}

	c.JSON(http.StatusOK, orders)
}

// 🚚 BestellStatus liefert den aktuellen Liefer-Status EINER Bestellung (für Live-Tracking).
func BestellStatus(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	nummer := c.Param("nummer")

	var order models.Order
	if err := database.DB.Where("nummer = ? AND user_id = ?", nummer, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"fehler": "Bestellung nicht gefunden"})
		return
	}

	phase, text := statusBerechnen(order.CreatedAt)
	c.JSON(http.StatusOK, gin.H{
		"nummer": order.Nummer,
		"phase":  phase,
		"status": text,
		"phasen": lieferPhasen,
	})
}
