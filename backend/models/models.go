package models

import "time"

// 👤 Ein Nutzerkonto. Das Passwort wird NUR als bcrypt-Hash gespeichert.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswortHash string    `json:"-"` // wird nie nach außen gegeben
	Username     string    `json:"username"`
	Vorname      string    `json:"vorname"`
	Nachname     string    `json:"nachname"`
	Geburtsdatum string    `json:"geburtsdatum"`
	Gesperrt     bool      `gorm:"default:false" json:"gesperrt"`
	Adressen     []Address `json:"adressen"`
	CreatedAt    time.Time `json:"createdAt"`
}

// 🏠 Eine Lieferadresse, die zu einem Nutzer gehört.
type Address struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	UserID     uint   `json:"-"`
	Label      string `json:"label"`
	Strasse    string `json:"strasse"`
	Hausnummer string `json:"hausnummer"`
	PLZ        string `json:"plz"`
	Ort        string `json:"ort"`
}

// 📦 Eine Bestellung eines Nutzers.
type Order struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	UserID        uint        `json:"-"`
	Nummer        string      `json:"nummer"`
	Summe         float64     `json:"summe"`
	Zwischensumme float64     `json:"zwischensumme"`
	Trinkgeld     float64     `json:"trinkgeld"`
	Gutschein     string      `json:"gutschein"`
	Zahlungsart   string      `json:"zahlungsart"`
	Liefertermin  string      `json:"liefertermin"`
	Status        string      `json:"status"`
	Positionen    []OrderItem `json:"positionen"`
	CreatedAt     time.Time   `json:"datum"`
}

// Ein einzelner Artikel innerhalb einer Bestellung.
type OrderItem struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	OrderID    uint    `json:"-"`
	Name       string  `json:"name"`
	Preis      float64 `json:"preis"`
	Menge      int     `json:"menge"`
	Restaurant string  `json:"restaurant"`
}
