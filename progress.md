# 📋 Lieferino – Was die App alles kann

Vollständige Übersicht über die Funktionen von **Frontend** (SvelteKit), **Backend**
(Go / Gin) und **Datenbank** (PostgreSQL) – jeweils mit einer kurzen Erklärung, wie es
funktioniert.

> **Aufbau:** Das Frontend spricht über eine REST-API (`http://localhost:8090`) mit dem
> Go-Backend, das alle Daten in einer PostgreSQL-Datenbank speichert. Backend + DB laufen
> in Docker. Ist das Backend aus, funktioniert die App im Browser-Modus (localStorage) weiter.

---

## 👤 Registrierung & E-Mail-Verifizierung
**Was:** Konto anlegen in 3 Schritten (Login-Daten → persönliche Daten → Adresse). Jede
E-Mail ist nur **einmal** verwendbar. Nach dem Anlegen muss die E-Mail per **6-stelligem
Code** bestätigt werden.
**Wie:** Das Frontend prüft schon in Schritt 1 über `POST /api/auth/email-check`, ob die
E-Mail frei ist (sonst „E-Mail bereits registriert"). Beim Absenden legt
`POST /api/auth/register` den Nutzer in der DB an (E-Mail-Spalte ist eindeutig) und schickt
den Code per E-Mail. `POST /api/auth/verify-email` prüft den Code. Ein Konto ist erst
nutzbar, wenn E-Mail bestätigt **und** MFA eingerichtet ist. Bei der Registrierung müssen
außerdem **AGB und Datenschutz** per Pflicht-Häkchen akzeptiert werden (Links zu `/agb` und
`/datenschutz`).

## 🔐 Login mit Zwei-Faktor (MFA) – Pflicht
**Was:** Anmeldung ist **zweistufig**: Passwort, danach ein 6-stelliger Code aus einer
**Authenticator-App** (TOTP, z. B. Google Authenticator/Authy).
**Wie:** `POST /api/auth/login` prüft das Passwort serverseitig (bcrypt-Hash). Stimmt es,
fordert der Server den zweiten Faktor an; `POST /api/auth/mfa/verify` prüft den TOTP-Code und
gibt erst dann ein Zugangs-Token (JWT) aus. Die MFA wird bei der Registrierung über einen
**QR-Code** eingerichtet (`/mfa/setup` + `/mfa/enable`).

## 🔑 Passwort zurücksetzen
**Was:** „Passwort vergessen" → Einmal-Code per E-Mail → Code eingeben → neues Passwort
wählen + bestätigen.
**Wie:** `POST /api/auth/reset-anfordern` schickt den Code (verrät aus Datenschutzgründen
nicht, ob die E-Mail existiert), `/reset-code` prüft ihn, `/reset-neu` setzt das neue
Passwort (wieder gegen die Passwort-Regeln geprüft, bcrypt-gehasht gespeichert).

## 🛡️ Sicherheit
**Was:** Mehrere Schutzmechanismen gegen Missbrauch.
**Wie:**
- **Passwörter** werden nur als **bcrypt-Hash** gespeichert, nie im Klartext.
- **Starke Passwort-Policy** (≥10 Zeichen, Groß-/Kleinbuchstabe, Zahl, Sonderzeichen) – mit
  Live-Checkliste im Frontend und identischer Prüfung im Backend.
- **Rate-Limiting** auf Login/Registrierung/E-Mail (pro IP) gegen Brute-Force/Spam.
- **Account-Sperre** nach zu vielen Fehlversuchen (temporär), plus Admin-Sperre.
- **CORS-Whitelist** (nur das eigene Frontend darf zugreifen) + Security-Header.
- **MFA ist Pflicht** für jeden Konto-Zugang.

## 🍽️ Verkäufer & Speisekarte
**Was:** Restaurant-Liste mit Suche, Küchen-Filter, Sortierung (u. a. nach Bewertung),
Favoriten-/Vegetarisch-/„Jetzt geöffnet"-Filter; Detailseite mit Speisekarte, Mengenwahl,
Öffnungszeiten.
**Wie:** Alle **Verkäufer und Produkte liegen in der Datenbank** (`GET /api/restaurants`).
Der Admin kann Restaurants **aktiv/deaktiv** schalten – deaktivierte sind auf der Seite
sofort ausgeblendet (Frontend liest den Status vom Backend).

## 🛒 Warenkorb & 🧾 Checkout
**Was:** Warenkorb (Mengen, Zwischensumme, nach Restaurant gruppiert), Checkout mit
Lieferadresse (+ Kartenvorschau), Lieferzeit, Zahlungsart (PayPal/Bar/Kreditkarte mit
Luhn-Validierung), Trinkgeld, Gutschein, Treuepunkte, Rechnung als PDF.
**Wie – wichtig (Manipulationsschutz):** Beim Bestellen (`POST /api/orders`) prüft das
**Backend jeden Artikel + Preis gegen die Datenbank**. Wer im Browser den Preis ändert,
wird abgelehnt. Auch **Gutschein, Liefergebühr, Mindestbestellwert und Treuepunkte** werden
**serverseitig** berechnet/erzwungen – die Endsumme kommt autoritativ vom Server.

## 🚚 Bestell-Status & Live-Tracking
**Was:** Bestellung durchläuft: *Bestellung erhalten → Wird zubereitet → Unterwegs →
Geliefert*. Checkout-Erfolgsseite und Tracking-Seite zeigen denselben Status.
**Wie:** Der Status wird **serverseitig aus der vergangenen Zeit** seit der Bestellung
berechnet (`GET /api/orders/:nummer/status`) – manipulationssicher und ohne Hintergrund-Job.
Beide Seiten fragen denselben Endpunkt ab, daher immer konsistent.

## ⭐ Bewertungen
**Was:** Sterne + Text pro Restaurant, Durchschnitt auf den Kacheln.
**Wie:** Bewertungen liegen in der DB. `POST /api/restaurants/:slug/reviews` erlaubt eine
Bewertung **nur, wenn man dort nachweislich bestellt hat** (serverseitig geprüft).
`GET /api/reviews/schnitt` liefert Durchschnitt + Anzahl je Restaurant.

## ❤️ Favoriten & ⭐ Treuepunkte
**Was:** Lieblings-Restaurants markieren; Treuepunkte sammeln (1 Punkt je 1 €) und im
Checkout einlösen (100 Punkte = 5 € Rabatt).
**Wie:** Beides liegt **pro Konto in der Datenbank** (`/api/favorites`, Punktestand am
Nutzer). Punkte werden beim Bestellen serverseitig gutgeschrieben/abgezogen – nicht im
Browser manipulierbar.

## 📦 Bestellverlauf & Profil
**Was:** Frühere Bestellungen mit Details, „Nochmal bestellen", Verfolgen, Rechnung.
Profil + mehrere Lieferadressen verwalten.
**Wie:** `GET /api/orders` (Bestellungen) und `GET/PUT /api/me` (Profil/Adressen) – alles in
der DB gespeichert, damit nichts verloren geht.

## 📧 E-Mail-Versand
**Was:** Verifizierungs-Code, Passwort-Reset-Code und **Bestellbestätigung** kommen per Mail.
**Wie:** Das Backend verschickt echte E-Mails über ein Gmail-Konto (SMTP). Die
Bestellbestätigung wird direkt beim Anlegen der Bestellung im Hintergrund verschickt.

## 🛠️ Admin-Bereich (`/admin`) — bleibt auf Deutsch
**Was:** Dashboard mit Kennzahlen (Umsatz, Bestellungen, Artikel, Nutzer, aktive
Restaurants), Verwaltung von Restaurants (aktiv/deaktiv), Bewertungen (löschen) und Nutzern
(sperren/entsperren, MFA zurücksetzen).
**Wie:** Zugang nur für Admin-Konten – geprüft über eine **Rollen-Middleware** im Backend
(`/api/admin/...`). Wer kein Admin ist, bekommt 403.

> ### 🔑 Admin-Zugang zum Testen (Klartext)
> - **E-Mail:** `admin@lieferino.de`
> - **Passwort:** `Admin!Lieferino2026`
> - **MFA-Setup-Key (TOTP):** `DTKYE5RLIBA5NDVWAPVWXXVCZ6WX4U64`
>
> **So einloggen:** E-Mail + Passwort eingeben → beim MFA-Schritt den Setup-Key in einer
> Authenticator-App hinterlegen (z. B. Google Authenticator → „Setup-Schlüssel eingeben")
> und den angezeigten 6-stelligen Code eintippen.
> *(Wer Admin wird, steuert die Umgebungsvariable `ADMIN_EMAIL`; Standard ist
> `admin@lieferino.de`.)*

## 🌍 Mehrsprachigkeit
**Was:** Umschalter im Menü; Sprachen: 🇩🇪 Deutsch, 🇬🇧 English, 🇪🇸 Español, 🇷🇺 Русский,
🇯🇵 日本語, 🏛️ Latina, 🟪 Enchanting (Minecraft-Schrift).
**Wie:** Zentrales Wörterbuch (`src/lib/i18n.js`), Auswahl bleibt im localStorage. Übersetzt
sind Navigation, Startseite, Restaurant-Liste, Login/MFA/Passwort-Reset, Account und Tracking.
Admin bleibt bewusst deutsch. *(Checkout + Registrierungsformular sind noch überwiegend deutsch.)*

## 🌙 Dark-/Light-Mode
**Was:** Umschalter unten links auf jeder Seite, Auswahl bleibt gespeichert.
**Wie:** Über ein `data-theme`-Attribut am `<html>`; Farben passen sich per CSS-Variablen an.

## 🥚 Easter Eggs
Konami-Code (↑↑↓↓←→←→BA), „Drachenlord"/`pizzapizzapizza`/`dragonpizza`/`foodcursor`/
`schnee` in die Suche tippen, Geburtstags-Konfetti, Hero-Titel 10× klicken, Zufalls-Jumpscare.

---

## 🧱 Technik-Stack
- **Frontend:** SvelteKit (Svelte 5), Vite.
- **Backend:** Go mit Gin, GORM, JWT, bcrypt, TOTP.
- **Datenbank:** PostgreSQL (in Docker).
- **Start:** `docker compose up -d --build` (Backend + DB) + `npm run dev` (Frontend),
  oder die Start-Skripte (`start-linux.sh`, `start-mac.command`, `start-windows.bat`).
  Optional: Kubernetes-Aufbau mit hochverfügbarer Patroni-Datenbank (siehe `k8s/README.md`).

## 🔭 Noch offen / optional
- Checkout + Registrierungsformular vollständig übersetzen.
- Echte Bezahlung über einen Zahlungsanbieter (aktuell simuliert).
- Bewertungen bearbeiten/löschen auch für Nutzer (Admin kann bereits moderieren).
