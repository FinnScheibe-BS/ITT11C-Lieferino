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
