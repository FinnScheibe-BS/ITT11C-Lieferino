# 📋 Fortschritt – Lieferino

Übersicht über alles, was am Frontend neu dazugekommen ist. Sortiert nach Bereich.
Backend-Themen sind im Code mit großen `🚨`-Hinweisen markiert (dort steht jeweils,
welcher Endpunkt gebraucht wird).

---

## 🎲 Easter-Egg (global)
- Liegt in `src/routes/+layout.svelte` und läuft dadurch auf **jeder** Seite.
- Pro Sekunde ein „Tick", pro Tick eine 1-zu-10000-Chance, dass ein GIF + Sound
  kurz abgespielt werden.
- Einstellbar über die Konstanten `TICK_DAUER_MS` und `CHANCE` oben im Script.

## 👤 Registrierung / Login (`src/routes/+page.svelte`)
- **E-Mail-Verifizierung:** Nach Schritt 1 muss ein 6-stelliger Code eingegeben
  werden (Service: `src/lib/services/email.js`). Der echte Mailversand kommt später
  vom Backend.
- **Passwort-Sicherheit:** Live-Stärke-Anzeige + Checkliste; „Weiter" ist gesperrt,
  bis das Passwort sicher genug ist (Service: `src/lib/services/passwort.js`).
- **Adress-Check:** Beim Abschluss wird über OpenStreetMap geprüft, ob die Adresse
  existiert (Service: `src/lib/services/adresse.js`).

## 🍽️ Restaurants & Daten
- **Eine zentrale Datenquelle** für alle Restaurants + Speisekarten unter
  `src/lib/data` (vorher gab es mehrere getrennte Listen). Die Daten bearbeitet man
  in `src/lib/data/quelle.js`.
- **Restaurant-Liste** (`/restaurants`): Suche, Filter nach Küche, Sortierung,
  Favoriten-Filter.
- **Detailseite** (`/restaurant/[slug]`): echte Speisekarte, „In den Warenkorb"
  mit Bestätigung, Favoriten-Herz und **Bewertungen** (Formular + Liste).

## 🛒 Warenkorb (`/cart`)
- Mengen ändern (+/−), Artikel entfernen, Zwischensumme, Liefergebühr, Gesamt.
- **Mindestbestellwert-Check** (Button zur Kasse erst ab erreichtem Wert).
- Bleibt nach dem Neuladen erhalten (localStorage).

## 🧾 Checkout (`/checkout`)
- Lieferadresse (aus den Registrierungsdaten), Zahlungsart-Auswahl, Bestellübersicht.
- **Gutscheincodes:** `LIEFERINO10` (10 %), `WILLKOMMEN5` (5 €), `GRATIS` (gratis Lieferung).
- **Voraussichtlicher Liefertermin** wird berechnet.
- **Live-Lieferstatus** nach der Bestellung: erhalten → wird zubereitet → unterwegs → geliefert.
- E-Mail-Bestätigung wird angestoßen (echter Versand = Backend).

## ❤️ Favoriten
- Store: `src/lib/stores/favoriten.js`. Herz-Buttons auf den Karten + „Nur Favoriten"-Filter.

## ⭐ Bewertungen
- Store: `src/lib/stores/bewertungen.js`. Reviews je Restaurant (Name, Sterne, Text).

## 📦 Bestellverlauf (`/bestellungen`)
- Liste aller früheren Bestellungen mit „🔁 Nochmal bestellen".

## 🛠️ Admin (`/admin`)
- Dashboard mit Kennzahlen (Restaurants, Bestellungen, verkaufte Artikel, Umsatz)
  und einer Restaurant-Tabelle.

## 🌙 Dark-/Light-Mode
- Runder Schalter unten links auf **jeder** Seite (Store: `src/lib/stores/theme.js`).
- Auswahl bleibt gespeichert; das Farbschema zieht sich über alle Seiten.

---

## ⚠️ Hinweis fürs Backend-Team
Favoriten, Bewertungen, Bestellungen und Nutzerdaten liegen aktuell nur lokal im
Browser (localStorage). Für echte, geräteübergreifende Daten müssen die im Code
markierten Endpunkte gebaut werden (E-Mail-Versand, Bestellungen, Bewertungen,
Adress-Proxy, Admin-Verwaltung).
