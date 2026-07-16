package models

// ==========================================
// AUTHENTIFIZIERUNG & USER STRUCTS
// ==========================================

type RegisterRequest struct {
	Vorname       string `json:"vorname"`
	Nachname      string `json:"nachname"`
	Strasse       string `json:"strasse"`
	Hausnummer    string `json:"hausnummer"`
	Ort	      string `json:"ort"`
	PLZ    	      string `json:"plz"`
	Telefonnummer string `json:"telefonnummer"`
	EmailAdresse  string `json:"email_adresse"`
	Passwort      string `json:"passwort"`
}

type LoginRequest struct {
	Email    string `json:"email_adresse"`
	Passwort string `json:"passwort"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	KundeID int    `json:"kunde_id"`
	Vorname string `json:"vorname"`
}

type UserInfo struct {
	KundeID int
	Email   string
}

// ==========================================
// WARENKORB, SPEISEN & CHECKOUT STRUCTS
// ==========================================

type CartItemRequest struct {
	GerichtID int `json:"gericht_id"`
	Anzahl    int `json:"anzahl"`
}

type CheckoutRequest struct {
	KundeID      int    `json:"kunde_id"`
	RestaurantID int    `json:"restaurant_id"`
	Zahlungsart  string `json:"zahlungsart"`
}


type Gericht struct {
    PK_ID_Gericht    int     `json:"id"`
    FK_ID_Restaurant int     `json:"restaurant_id"`
    Name             string  `json:"name"`
    Vegetarisch      bool    `json:"vegetarisch"`
    Vegan            bool    `json:"vegan"`
    Preis            float64 `json:"preis"`
    FK_ID_Kategorie  int     `json:"kategorie_id"`
}

// Restaurant angepasst mit ID und Nationalitaet laut dishes.go:89
type Restaurant struct {
	ID               int    `json:"id"`
	PK_ID_Restaurant int    `json:"restaurant_id"`
	Name             string `json:"name"`
	Adresse          string `json:"adresse"`
	Nationalitaet    string `json:"nationalitaet"`
}

// ==========================================
// GUTSCHEIN STRUCTS
// ==========================================

type GutscheinRequest struct {
	Code string `json:"code"`
}