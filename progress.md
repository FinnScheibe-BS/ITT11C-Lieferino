# 📋 Fortschritt – Lieferino

> ### ⚡ WICHTIG ZUERST
> Falls nach Änderungen an den Restaurant-/Menü-Daten in der App etwas **nicht passt**,
> im Projektordner `Lieferino` einmal diesen Befehl ausführen:
> ```
> node tools/refresh.mjs
> ```
> Danach `npm run dev` neu laden – dann passt alles wieder. ✅

Übersicht über alles, was am Frontend dazugekommen ist. Backend-Themen sind im Code
mit großen `🚨`-Hinweisen markiert (dort steht jeweils, welcher Endpunkt gebraucht wird).

---

## 👤 Registrierung & Konto
- **3-Schritt-Registrierung** mit Validierung:
  - **E-Mail-Verifizierung** per 6-stelligem Code.
  - **Passwort-Sicherheit**: Live-Stärke-Anzeige + Checkliste, „Weiter" erst ab sicherem Passwort.
  - **Namens-Check**: Vor-/Nachname min. 2 Buchstaben, keine Zahlen.
  - **Adress-Check** über OpenStreetMap + **PLZ-Format** (5 Ziffern).
- **Account-Seite**: Daten bearbeiten, **mehrere Adressen** verwalten,
  **Treuepunkte** ansehen, Konto löschen.

## 🔐 Login & Sicherheit
- **Echtes Login mit Session** (`/login`): „Angemeldet bleiben", Session-Timeout,
  Login-Sperre nach mehreren Fehlversuchen, gebannte Konten werden abgewiesen.
- **Zwei-Faktor (2FA/MFA)**: per **E-Mail-Code** oder **Authenticator-App (TOTP)**,
  inkl. **Backup-Codes**.
- **Passwort vergessen** (`/passwort-vergessen`): Reset per Code.
- Trennung von **Konto** (bleibt gespeichert) und **Session** (Logout beendet nur die Session).

## 🍽️ Restaurants & Speisekarte
- **Zentrale Datenquelle** unter `src/lib/data` (eine Quelle für alle Seiten).
- **Restaurant-Liste** (`/restaurants`): Suche, Küchen-Filter, Sortierung,
  **Favoriten-Filter**, **Vegetarisch-Filter**, **„Jetzt geöffnet"-Filter**,
  Durchschnitts-Bewertung aus Reviews.
- **Detailseite**: echte Speisekarte mit **Mengen-Auswahl**, **Veggie-Tags & Allergenen**,
  **Öffnungszeiten** (geöffnet/geschlossen), Favoriten-Herz.
- **Bewertungen**: nur möglich, wenn man dort **nachweislich bestellt** hat;
  Bewertung im Kopf ist anklickbar (springt zu den Reviews).

## 🛒 Warenkorb (`/cart`)
- Mengen ändern, entfernen, Zwischensumme/Liefergebühr/Gesamt.
- **Nach Restaurant gruppiert**, Mindestbestellwert-Check, bleibt nach Neuladen erhalten.

## 🧾 Checkout (`/checkout`)
- Lieferadresse + **Karten-Vorschau** der Adresse, **Auswahl unter mehreren Adressen**.
- **Lieferzeit**: sofort oder **geplant** (Uhrzeit wählen).
- **Zahlungsarten**: PayPal, Barzahlung, **Kreditkarte mit Validierung**
  (Luhn-Prüfung, Ablaufdatum, CVV, Kartentyp-Erkennung).
- **Trinkgeld** (0/5/10/15 %), **Gutscheincodes** (`LIEFERINO10`, `WILLKOMMEN5`, `GRATIS`).
- **Treuepunkte einlösen** (100 Punkte = 5 € Rabatt).
- **Bestellnummer**, voraussichtlicher **Liefertermin**, **Live-Lieferstatus**,
  **Browser-Benachrichtigung** bei Status-Wechsel, E-Mail-Bestätigung.
- **Rechnung als PDF / zum Drucken**.

## 🚚 Live-Tracking (`/tracking`)
- Zeigt den Liefer-Fortschritt der letzten Bestellung (aktualisiert sich automatisch).

## ⭐ Treuepunkte
- 1 Punkt je 1 € Bestellwert; im Account sichtbar, im Checkout einlösbar.

## 📍 Mehrere Adressen
- Im Account verwaltbar (hinzufügen/löschen), im Checkout auswählbar.

## 📦 Bestellverlauf (`/bestellungen`)
- Alle früheren Bestellungen mit Bestellnummer, **Detailansicht**, „🔁 Nochmal bestellen",
  **🚚 Verfolgen** (Live-Tracking) und **🧾 Rechnung** (PDF/Druck).

## 🛠️ Admin (`/admin`) — bleibt auf Deutsch
- Durch einen **Sicherheitsschlüssel geschützt** (Passwort ist im Code hinterlegt – nicht hier dokumentiert).
- **Editier-Modus-Schalter** (bleibt während der Sitzung erhalten).
- **Tabs**:
  - **Übersicht**: Kennzahlen + **Umsatz-Diagramm** (7 Tage).
  - **Lieferanten**: Restaurants **aktivieren/deaktivieren** (deaktivierte sind für normale Nutzer unsichtbar).
  - **Bewertungen**: bearbeiten/löschen.
  - **Nutzer**: **bannen/entbannen**, **MFA zurücksetzen**.

## 🌍 Mehrsprachigkeit
- **Sprach-Umschalter** im Menü, Auswahl bleibt gespeichert (`src/lib/i18n.js`).
- Sprachen: 🇩🇪 Deutsch, 🇬🇧 English, 🇪🇸 Español, 🇷🇺 Русский, 🇯🇵 日本語, 🏛️ Latina,
  🟪 **Enchanting** (Minecraft-Verzauberungstisch – braucht die Schrift `static/fonts/enchantment.ttf`).
- Übersetzt sind aktuell Navigation, Startseite und Restaurant-Liste (weitere Seiten folgen).
- Admin ist bewusst **nicht** übersetzt.

## 🌙 Dark-/Light-Mode
- Runder Schalter unten links auf jeder Seite, Auswahl bleibt gespeichert.

## 🥚 Easter Eggs
- **Jufalls-Jumpscare**: pro Sekunde 1-zu-10000-Chance auf ein kurzes GIF + Sound (global).
- **Konami-Code** (↑ ↑ ↓ ↓ ← → ← → B A)
- **„Drachenlord"** in die Suche tippen
- **`pizzapizzapizza`** in die Suche → Konfetti + Geheim-Gutschein `PIZZAPARTY` (25 %).
- **Geburtstag**: am Geburtstag (laut Konto) → Konfetti + Code `GEBURTSTAG` (20 %).
- **Hero-Titel 10× klicken** → Konfetti-Überraschung.
- **`foodcursor`** in die Suche → Emoji-Spur folgt der Maus (umschaltbar).
- **`schnee`** / **`winter`** in die Suche → fallender Saison-Effekt (umschaltbar).
- **`dragonpizza`** in die Suche → schaltet ein **geheimes Restaurant** frei.

---

## 🚨 Wo das Backend weiterarbeiten muss (Checkliste)

Aktuell läuft alles im Frontend (Daten im `localStorage`). Damit es „echt" wird,
muss das Backend folgende Punkte übernehmen:

### Konten & Login
- [ ] **Registrierung** speichern: `POST /api/auth/register` (Nutzer in Datenbank).
- [ ] **Login** serverseitig prüfen: `POST /api/auth/login` (kein Klartext-Vergleich im Browser).
- [ ] **Passwörter hashen** (z. B. bcrypt) – niemals im Klartext speichern.
- [ ] **Sessions/Token** serverseitig verwalten (z. B. JWT), inkl. Ablauf.
- [ ] **Login-Sperre / Rate-Limiting** serverseitig (Frontend-Sperre ist umgehbar).

### E-Mail (Service `src/lib/services/email.js` + Route `src/routes/api/email`)
- [ ] **Echten Versand** über das Gmail-Konto (App-Passwort in `.env`, siehe `EMAIL-SETUP.md`).
- [ ] **Verifizierungs-/Reset-Codes** serverseitig erzeugen, senden und prüfen.
- [ ] **Bestellbestätigung + Liefertermin** per E-Mail verschicken.

### Zwei-Faktor (MFA, `src/lib/services/mfa.js`)
- [ ] **TOTP-Secret** serverseitig erzeugen + sicher speichern, Codes serverseitig prüfen.
- [ ] **E-Mail-2FA-Codes** serverseitig.
- [ ] **Backup-Codes** serverseitig verwalten (Verbrauch zählen).

### Adresse
- [ ] **Adress-Prüfung** über einen eigenen Endpunkt proxen (statt direkt OpenStreetMap; Rate-Limit).

### Restaurants & Bewertungen
- [ ] **Restaurant-/Speisekarten-Daten** aus der Datenbank ausliefern.
- [ ] **Aktiv/Deaktiviert-Status** speichern: `PATCH /api/admin/restaurants/:slug`;
      normalen Nutzern nur **aktive** Restaurants ausliefern.
- [ ] **Bewertungen** zentral speichern/laden (`GET/POST /api/bewertungen`) und „nur nach
      Bestellung" serverseitig erzwingen.
- [ ] **Favoriten** pro Nutzer serverseitig speichern.

### Bestellungen & Zahlung
- [ ] **Bestellungen** speichern + ans Restaurant melden: `POST /api/bestellungen`
      (mit Bestellnummer, Artikeln, Summe, Trinkgeld, Liefertermin).
- [ ] **Echte Zahlung** über einen Anbieter (z. B. Stripe/PayPal) – Kartendaten NIE
      selbst speichern/übertragen (nur Token).
- [ ] **Gutscheincodes** serverseitig prüfen/einlösen (Missbrauch/Mehrfachnutzung verhindern).

### Admin
- [ ] **Echte Admin-Authentifizierung** statt des Frontend-Schlüssels.
- [ ] **Geschützte Admin-Endpunkte**: Nutzerliste, Bannen/Entbannen, MFA-Reset,
      Restaurant-Verwaltung, Bewertungs-Moderation.

### Sonstiges
- [ ] **Konto löschen** serverseitig umsetzen.
- [ ] Optional: Bestellverlauf, Treuepunkte usw. an die Datenbank koppeln.
