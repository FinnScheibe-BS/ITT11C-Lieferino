# 📧 E-Mail-Versand aktivieren (Anleitung)

Der komplette Code ist schon fertig. Damit **echte** E-Mails verschickt werden
(Verifizierungscode, Bestellbestätigung, Passwort-Reset), müssen nur noch die
Zugangsdaten hinterlegt werden.

> Solange das nicht eingerichtet ist, läuft alles im **Test-Modus**: Es wird
> nichts verschickt, der Code wird stattdessen direkt in der App angezeigt.
> Außerdem gibt es im Test-Modus überall einen **„Dev: überspringen"-Button**.

## Zugangsdaten (Lieferino-Konto)
- **E-Mail:** `lieferino5@gmail.com`
- **Konto-Passwort:** liegt lokal in der `.env`-Datei (wird aus Sicherheitsgründen
  NICHT in dieses öffentliche Repo geschrieben). Im Team bitte privat weitergeben.
  - ⚠️ Das normale Konto-Passwort funktioniert **NICHT** für den Programm-Versand!
    Google verlangt dafür ein **App-Passwort** (siehe unten).

## Schritt für Schritt

### 1. 2-Faktor-Authentifizierung im Gmail-Konto aktivieren
- Bei `lieferino5@gmail.com` anmelden → https://myaccount.google.com/security
- „Bestätigung in zwei Schritten" einschalten (ohne das gibt es keine App-Passwörter).

### 2. App-Passwort erzeugen
- https://myaccount.google.com/apppasswords öffnen
- App-Name z. B. „Lieferino" eingeben → **Erstellen**
- Google zeigt ein **16-stelliges Passwort** an (z. B. `abcd efgh ijkl mnop`).
- Dieses Passwort kopieren (die Leerzeichen kann man weglassen).

### 3. In die .env-Datei eintragen
Im Projekt-Ordner gibt es die Datei `.env` (wird **nicht** zu GitHub hochgeladen).
Dort eintragen:

```
GMAIL_USER=lieferino5@gmail.com
GMAIL_APP_PASSWORD=abcdefghijklmnop
```

### 4. Server neu starten
```
npm run dev
```
Fertig – ab jetzt verschickt die App echte E-Mails. ✅

## Wie funktioniert es im Code?
- Frontend → ruft die Server-Route **`/api/email`** auf
  ([src/routes/api/email/+server.js](src/routes/api/email/+server.js)).
- Die Route verschickt mit **Nodemailer** über Gmail.
- Der Service [src/lib/services/email.js](src/lib/services/email.js) wird von
  Registrierung, Login (2FA), Passwort-Reset und Bestellbestätigung genutzt.

## Wichtig für später (Produktion)
- Die `.env` niemals zu GitHub hochladen (ist bereits in `.gitignore`).
- Das App-Passwort sicher im Team teilen (nicht öffentlich posten).
- Die **„Dev: überspringen"-Buttons** erscheinen nur im Entwicklungsmodus
  (`npm run dev`) und tauchen in der echten Version automatisch NICHT auf.
