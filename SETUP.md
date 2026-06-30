# 🛠️ Lieferino – Einrichtung & Start

Die App besteht aus **zwei Teilen**, die zusammen laufen:
- **Frontend** (SvelteKit) – die Website
- **Backend** (Go/Gin) + **PostgreSQL-Datenbank** – läuft in Docker

> Wenn das Backend läuft, werden Registrierung, Login, Profil/Adressen und
> Bestellungen in der **Datenbank** gespeichert (gehen also nicht verloren).
> Ist das Backend aus, funktioniert die App weiter über den Browser-Speicher.

---

## 🚀 Schnellstart (alles mit einem Klick)
Es gibt fertige Start-Skripte, die **alles** auf einmal hochfahren (Backend + Datenbank + Frontend)
und den Browser öffnen. **Voraussetzung:** Docker läuft + Node.js ist installiert.

| Betriebssystem | So startest du |
|---|---|
| 🐧 **Linux** | im Terminal: `./start-linux.sh` (einmalig evtl. `chmod +x start-linux.sh`) |
| 🍎 **macOS** | **Doppelklick** auf `start-mac.command` (oder im Terminal: `./start-mac.command`) |
| 🪟 **Windows** | **Doppelklick** auf `start-windows.bat` |

Danach läuft die App auf **http://localhost:5173** (Backend auf :8090).

**Beenden:** im Terminal `Strg+C`. Die Docker-Container laufen weiter – zum Stoppen:
```sh
cd backend
docker compose down
```

> Wer es lieber von Hand macht, findet die Einzelschritte unten.

---

## 0. Voraussetzungen
- **Node.js** (für das Frontend)
- **Docker** + **Docker Compose** (für Backend + Datenbank)

## 1. Frontend starten
Im Projekt-Hauptordner:
```sh
npm install
npm run dev
```
→ Website läuft auf **http://localhost:5173**

## 2. Backend + Datenbank starten
In einem zweiten Terminal:
```sh
cd backend
docker compose up --build
```
→ Backend läuft auf **http://localhost:8090** (Datenbank läuft intern im Container).
Test: http://localhost:8090/health sollte `{"status":"ok"}` zeigen.

> Eigener Docker-Projektname `lieferino` + Port 8090 → kollidiert nicht mit anderen Docker-Projekten.
> Zum Stoppen: `docker compose down` (Daten bleiben erhalten).

---

## 3. 📧 E-Mail-Versand aktivieren
Damit echte E-Mails (Verifizierungscode usw.) verschickt werden, braucht Gmail ein
**App-Passwort** (16-stellig) – das **normale Konto-Passwort funktioniert NICHT**.

**Schritt für Schritt:**
1. Beim Gmail-Konto `lieferino5@gmail.com` die **2-Faktor-Authentifizierung** aktivieren
   (https://myaccount.google.com/security).
2. Ein **App-Passwort** erzeugen: https://myaccount.google.com/apppasswords
   → 16-stelligen Code kopieren.
3. In die Datei **`backend/handlers/email.go`** eintragen:
   ```go
   const standardGmailAppPasswort = "abcdefghijklmnop"   // <- hier das App-Passwort
   ```
4. Backend neu bauen:
   ```sh
   cd backend
   docker compose up -d --build backend
   ```

Danach werden E-Mails wirklich verschickt. Solange das App-Passwort leer ist, läuft
der Versand im **Test-Modus** (der Code wird in der App angezeigt, damit man testen kann).

> ⚠️ Hinweis: Das Repo ist öffentlich – ein dort sichtbares App-Passwort wird von Google
> oft automatisch deaktiviert. Notfalls einfach neu erzeugen.

---

## 4. Nach erfolgreicher E-Mail-Einrichtung
Sobald der echte Versand läuft, kann der **„Dev: überspringen"-Button** in der
Verifizierung entfernt werden (er erscheint ohnehin nur im Entwicklungsmodus).

---

## Kurzfassung
| Was | Befehl | Adresse |
|---|---|---|
| Frontend | `npm install && npm run dev` | http://localhost:5173 |
| Backend + DB | `cd backend && docker compose up --build` | http://localhost:8090 |
| E-Mail | App-Passwort in `backend/handlers/email.go` + Backend neu bauen | – |
