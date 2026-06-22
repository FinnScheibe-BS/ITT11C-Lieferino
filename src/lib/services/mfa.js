// =====================================================================
// 🔐 MFA / 2-FAKTOR-AUTHENTIFIZIERUNG
// Unterstützt zwei Methoden:
//   1. "email"  – ein Code wird per E-Mail geschickt (nutzt email.js)
//   2. "totp"   – ein rotierender Code aus einer Authenticator-App
//                 (z.B. Google Authenticator), nach RFC 6238.
//
// Die TOTP-Prüfung läuft hier im Frontend nur zu Demo-Zwecken. Im echten
// Betrieb übernimmt das Backend die Prüfung (siehe großer Hinweis unten).
// =====================================================================

// Das Alphabet für Base32 (so werden Authenticator-Secrets kodiert).
const BASE32 = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567';

// 🔑 Erzeugt ein zufälliges Secret (16 Zeichen Base32) für die Authenticator-App.
export function generiereSecret() {
  const zufall = crypto.getRandomValues(new Uint8Array(16));
  let secret = '';
  for (const b of zufall) secret += BASE32[b % 32];
  return secret;
}

// 🆘 Erzeugt Backup-Codes (Format "A1B2-C3D4"), falls die App verloren geht.
// Jeder Code ist nur EINMAL gültig.
export function generiereBackupCodes(anzahl = 6) {
  const codes = [];
  for (let i = 0; i < anzahl; i++) {
    const bytes = crypto.getRandomValues(new Uint8Array(4));
    const hex = Array.from(bytes)
      .map((b) => b.toString(16).padStart(2, '0'))
      .join('')
      .toUpperCase();
    codes.push(`${hex.slice(0, 4)}-${hex.slice(4, 8)}`);
  }
  return codes;
}

// 📲 Baut die "otpauth://"-Adresse, die eine Authenticator-App versteht.
// Daraus kann die App (oder ein QR-Code) das Konto einrichten.
export function otpauthUri(secret, email) {
  const label = encodeURIComponent(`Lieferino:${email}`);
  return `otpauth://totp/${label}?secret=${secret}&issuer=Lieferino`;
}

// Wandelt einen Base32-Text in echte Bytes um (für die HMAC-Berechnung).
function base32Decode(text) {
  const clean = text.toUpperCase().replace(/[^A-Z2-7]/g, '');
  let bits = 0;
  let wert = 0;
  const bytes = [];
  for (const zeichen of clean) {
    wert = (wert << 5) | BASE32.indexOf(zeichen);
    bits += 5;
    if (bits >= 8) {
      bits -= 8;
      bytes.push((wert >> bits) & 0xff);
    }
  }
  return new Uint8Array(bytes);
}

// Berechnet den 6-stelligen Code für einen bestimmten Zeitabschnitt (HOTP-Kern).
async function berechneCode(keyBytes, zaehler) {
  // Der Zähler muss als 8-Byte-Zahl (Big-Endian) vorliegen.
  const buffer = new ArrayBuffer(8);
  const view = new DataView(buffer);
  view.setUint32(0, Math.floor(zaehler / 2 ** 32), false);
  view.setUint32(4, zaehler >>> 0, false);

  // HMAC-SHA1 über den Schlüssel + Zähler (Web Crypto erledigt die Mathematik).
  const key = await crypto.subtle.importKey(
    'raw',
    keyBytes,
    { name: 'HMAC', hash: 'SHA-1' },
    false,
    ['sign']
  );
  const signatur = new Uint8Array(await crypto.subtle.sign('HMAC', key, buffer));

  // "Dynamic Truncation": aus der Signatur die 6 Ziffern herausrechnen.
  const offset = signatur[19] & 0x0f;
  const zahl =
    ((signatur[offset] & 0x7f) << 24) |
    ((signatur[offset + 1] & 0xff) << 16) |
    ((signatur[offset + 2] & 0xff) << 8) |
    (signatur[offset + 3] & 0xff);
  return (zahl % 1000000).toString().padStart(6, '0');
}

// ✅ Prüft, ob der eingegebene TOTP-Code stimmt.
// Wir erlauben ein Zeitfenster von ±30 Sekunden, falls die Uhren leicht abweichen.
export async function pruefeTotp(secret, eingabe) {
  const key = base32Decode(secret);
  const jetztAbschnitt = Math.floor(Date.now() / 1000 / 30);
  for (const fenster of [-1, 0, 1]) {
    const code = await berechneCode(key, jetztAbschnitt + fenster);
    if (code === eingabe.trim()) return true;
  }
  return false;
}

// ---------------------------------------------------------------------
// 🚨🚨🚨 ACHTUNG BACKEND-TEAM – MFA MUSS SERVERSEITIG ABGESICHERT WERDEN! 🚨🚨🚨
// ---------------------------------------------------------------------
// Aktuell (nur Frontend/Demo):
//   - Das TOTP-Secret liegt beim Nutzer im localStorage und die Prüfung
//     passiert im Browser. Das ist NICHT sicher – jeder könnte es auslesen.
//
// WAS DAS BACKEND BAUEN MUSS:
//   1. Beim Einrichten von TOTP das Secret SERVERSEITIG erzeugen und sicher
//      (verschlüsselt) speichern – niemals dauerhaft im Browser.
//   2. Endpunkt zum Aktivieren:   POST /api/mfa/totp/setup   -> liefert otpauth-URI
//      Endpunkt zum Bestätigen:   POST /api/mfa/totp/verify  { code }
//   3. Beim Login NACH korrektem Passwort den zweiten Faktor verlangen und
//      den Code SERVERSEITIG prüfen:  POST /api/mfa/verify { methode, code }
//   4. Für die E-Mail-Methode: Code im Backend erzeugen, per Mail senden und prüfen.
//   5. Außerdem: Backup-Codes anbieten, falls die App verloren geht.
// ---------------------------------------------------------------------
