package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"gorm.io/gorm"
)

var _ck []byte

const _mk = "enc:"

// SetCipher hinterlegt den Schlüssel, mit dem Adressdaten in der DB abgelegt werden.
func SetCipher(k []byte) { _ck = k }

func _enc(s string) string {
	if len(_ck) == 0 || s == "" || strings.HasPrefix(s, _mk) {
		return s
	}
	block, err := aes.NewCipher(_ck)
	if err != nil {
		return s
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return s
	}
	out := make([]byte, len(s))
	cipher.NewCTR(block, iv).XORKeyStream(out, []byte(s))
	return _mk + base64.StdEncoding.EncodeToString(append(iv, out...))
}

func _dec(s string) string {
	if len(_ck) == 0 || !strings.HasPrefix(s, _mk) {
		return s
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(s, _mk))
	if err != nil || len(raw) < aes.BlockSize {
		return s
	}
	block, err := aes.NewCipher(_ck)
	if err != nil {
		return s
	}
	ct := raw[aes.BlockSize:]
	out := make([]byte, len(ct))
	cipher.NewCTR(block, raw[:aes.BlockSize]).XORKeyStream(out, ct)
	return string(out)
}

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

	// 🛡️ Brute-Force-Schutz: zählt Fehlversuche und sperrt das Konto kurz.
	Fehlversuche int        `json:"-"`
	GesperrtBis  *time.Time `json:"-"`

	// 📧 E-Mail-Verifizierung: Konto erst nutzbar, wenn die E-Mail bestätigt ist.
	EmailVerifiziert bool       `gorm:"default:false" json:"emailVerifiziert"`
	EmailCode        string     `json:"-"` // aktueller 6-stelliger Bestätigungscode
	EmailCodeAblauf  *time.Time `json:"-"` // wann der Code abläuft

	// 🔐 MFA (Zwei-Faktor per Authenticator-App / TOTP): Pflicht für Zugang.
	MFAAktiv  bool   `gorm:"default:false" json:"mfaAktiv"`
	MFASecret string `json:"-"` // TOTP-Geheimnis (Base32)
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

// ⭐ Eine Bewertung eines Restaurants durch einen Nutzer.
type Review struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `json:"-"`
	RestaurantSlug string    `gorm:"index" json:"slug"`
	RestaurantName string    `json:"restaurantName"`
	Name           string    `json:"name"`
	Sterne         int       `json:"sterne"`
	Text           string    `json:"text"`
	CreatedAt      time.Time `json:"datum"`
}

// ❤️ Ein Lieblings-Restaurant eines Nutzers (per Slug).
type Favorite struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	UserID         uint   `gorm:"index" json:"-"`
	RestaurantSlug string `json:"slug"`
}

// Adressdaten werden in der Datenbank nicht im Klartext abgelegt.
func (a *Address) BeforeSave(tx *gorm.DB) error {
	a.Strasse = _enc(a.Strasse)
	a.Hausnummer = _enc(a.Hausnummer)
	a.PLZ = _enc(a.PLZ)
	a.Ort = _enc(a.Ort)
	return nil
}

func (a *Address) AfterSave(tx *gorm.DB) error {
	a.Strasse = _dec(a.Strasse)
	a.Hausnummer = _dec(a.Hausnummer)
	a.PLZ = _dec(a.PLZ)
	a.Ort = _dec(a.Ort)
	return nil
}

func (a *Address) AfterFind(tx *gorm.DB) error {
	a.Strasse = _dec(a.Strasse)
	a.Hausnummer = _dec(a.Hausnummer)
	a.PLZ = _dec(a.PLZ)
	a.Ort = _dec(a.Ort)
	return nil
}
