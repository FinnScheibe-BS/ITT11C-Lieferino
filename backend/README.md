# 🍕 Lieferino Backend (Go + Gin + PostgreSQL)

Das Backend zur Lieferino-App. Geschrieben in **Go** mit **Gin**, Datenbank **PostgreSQL**
über **GORM**, alles läuft in **Docker**.

## 🚀 Starten
Im Ordner `backend` ausführen:

```sh
docker compose up --build
```

Das startet zwei Container:
- **lieferino-db** – die PostgreSQL-Datenbank (nur intern erreichbar)
- **lieferino-backend** – das Go-Backend, erreichbar unter **http://localhost:8090**

> Hinweis: Eigener Projektname `lieferino` + Port **8090**, damit es nicht mit
> anderen Docker-Projekten kollidiert.

## 🔌 Endpunkte (Fundament)
| Methode | Pfad | Beschreibung |
|---|---|---|
| GET  | `/health` | Lebt das Backend? |
| POST | `/api/auth/register` | Registrierung `{ email, passwort, username, vorname, nachname, geburtsdatum }` |
| POST | `/api/auth/login` | Login `{ email, passwort }` → JWT-Token |
| POST | `/api/email` | E-Mail senden `{ an, betreff, text }` (Test-Modus ohne Gmail-Zugang) |

## 📧 E-Mail aktivieren
In der `docker-compose.yml` (oder einer `.env`) ein **Gmail-App-Passwort** setzen:
```
GMAIL_USER=lieferino5@gmail.com
GMAIL_APP_PASSWORD=<16-stelliges App-Passwort>
```
Ohne diese Werte läuft der E-Mail-Endpunkt im Test-Modus (sendet nichts).

## 🗄️ Datenbank
Die Tabellen werden beim Start automatisch angelegt (GORM AutoMigrate).
Aktuell: `users`, `addresses`. Weitere Modelle (Restaurants, Bestellungen, …) folgen.

## 🛠️ Lokale Entwicklung ohne Docker
Geht auch, braucht aber lokal installiertes Go + eine laufende PostgreSQL.
Werte über Umgebungsvariablen setzen (siehe `.env.example`), dann `go run .`.
