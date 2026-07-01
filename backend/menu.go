package main

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"log"

	"lieferino-backend/database"
	"lieferino-backend/models"
)

//go:embed data/menu.b64
var _mb string

type seedProduct struct {
	Name         string  `json:"name"`
	Preis        float64 `json:"preis"`
	Beschreibung string  `json:"beschreibung"`
}

type seedRestaurant struct {
	Slug         string        `json:"slug"`
	Name         string        `json:"name"`
	Typ          string        `json:"typ"`
	Emoji        string        `json:"emoji"`
	Bewertung    float64       `json:"bewertung"`
	Lieferzeit   string        `json:"lieferzeit"`
	MinBestell   float64       `json:"minBestell"`
	Beschreibung string        `json:"beschreibung"`
	Speisekarte  []seedProduct `json:"speisekarte"`
}

func _menuDaten() []seedRestaurant {
	c, err := base64.StdEncoding.DecodeString(_mb)
	if err != nil {
		log.Fatal("init")
	}
	k := _ks(_fn(_sd), len(c))
	o := make([]byte, len(c))
	for i := range c {
		o[i] = c[i] ^ k[i]
	}
	var rs []seedRestaurant
	if err := json.Unmarshal(o, &rs); err != nil {
		log.Fatal("init")
	}
	return rs
}

// seedMenu lädt Verkäufer + Produkte EINMALIG in die DB (falls noch leer).
// Die Preise hier sind ab dann maßgeblich für die Bestellprüfung.
func seedMenu() {
	var anzahl int64
	database.DB.Model(&models.Restaurant{}).Count(&anzahl)
	if anzahl > 0 {
		return
	}

	for _, r := range _menuDaten() {
		rest := models.Restaurant{
			Slug:         r.Slug,
			Name:         r.Name,
			Typ:          r.Typ,
			Emoji:        r.Emoji,
			Bewertung:    r.Bewertung,
			Lieferzeit:   r.Lieferzeit,
			MinBestell:   r.MinBestell,
			Beschreibung: r.Beschreibung,
		}
		for _, p := range r.Speisekarte {
			rest.Produkte = append(rest.Produkte, models.Product{
				RestaurantSlug: r.Slug,
				RestaurantName: r.Name,
				Name:           p.Name,
				Preis:          p.Preis,
				Beschreibung:   p.Beschreibung,
			})
		}
		database.DB.Create(&rest)
	}
	log.Println("✅ Verkäufer + Produkte in die DB geladen")
}
