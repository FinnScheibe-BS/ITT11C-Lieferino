package handlers

import (
	"net/http"
	"strings" // <-- Wichtig für strings.ToLower

	"aup-backend/database"
	"aup-backend/models" // <-- Wichtig für die Structs
	"github.com/gin-gonic/gin"
)

// GetGerichte liefert alle Gerichte basierend auf den Filtern zurück
func GetGerichte(c *gin.Context) {
	// Filter aus der URL holen und direkt in Kleinschreibung umwandeln
	filterNation := strings.ToLower(c.Query("nationalitaet")) 
	filterVegan := strings.ToLower(c.Query("vegan"))          
	filterVeggie := strings.ToLower(c.Query("vegetarisch"))   

	query := `
		SELECT 
			g.PK_ID_Gericht, g.FK_ID_Restaurant, g.Name, g.Vegetarisch, g.Vegan, g.Preis 
		FROM Gericht g
		INNER JOIN Restaurant r ON g.FK_ID_Restaurant = r.PK_ID_Restaurant
		WHERE 1=1
	`
	var args []interface{}

	if filterNation != "" {
		query += " AND LOWER(r.Nationalitaet) = ?"
		args = append(args, filterNation)
	}
	if filterVegan == "true" {
		query += " AND g.Vegan = 1"
	}
	if filterVeggie == "true" {
		query += " AND g.Vegetarisch = 1"
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Filter-Abfrage fehlgeschlagen: " + err.Error()})
		return
	}
	defer rows.Close()

	var gerichte []models.Gericht // models. nutzen

	for rows.Next() {
		var g models.Gericht // models. nutzen
		var vegetarischBytes []byte
		var veganBytes []byte

		// Ersetze deine Scan-Zeile durch diese:
err := rows.Scan(&g.PK_ID_Gericht, &g.FK_ID_Restaurant, &g.Name, &vegetarischBytes, &veganBytes, &g.Preis)

if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan-Fehler: " + err.Error()})
    return
}

		if len(vegetarischBytes) > 0 {
			g.Vegetarisch = vegetarischBytes[0] == 1
		}
		if len(veganBytes) > 0 {
			g.Vegan = veganBytes[0] == 1
		}

		gerichte = append(gerichte, g)
	}

	if gerichte == nil {
		gerichte = []models.Gericht{}
	}
	c.JSON(http.StatusOK, gerichte)
}

// GetRestaurants liefert alle Restaurants aus der Datenbank zurück
func GetRestaurants(c *gin.Context) { // <-- Name korrigiert und Klammern aufgeräumt
	query := "SELECT PK_ID_Restaurant, Name, Adresse, Nationalitaet FROM Restaurant"
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Datenbankabfrage fehlgeschlagen: " + err.Error()})
		return
	}
	defer rows.Close()

	var restaurants []models.Restaurant // models. nutzen

	for rows.Next() {
		var res models.Restaurant // models. nutzen
		err := rows.Scan(&res.ID, &res.Name, &res.Adresse, &res.Nationalitaet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Einlesen der Daten: " + err.Error()})
			return
		}
		restaurants = append(restaurants, res)
	}

	if restaurants == nil {
		restaurants = []models.Restaurant{}
	}
	c.JSON(http.StatusOK, restaurants)
}
	