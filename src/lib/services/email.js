// =====================================================================
// 📧 E-MAIL-SERVICE
// Hier sind alle Funktionen rund um das Versenden von E-Mails gebündelt.
// Aktuell läuft alles im Frontend (zum Testen). Der echte E-Mail-Versand
// MUSS aber vom Backend übernommen werden (siehe große Hinweise unten).
// =====================================================================

// Erzeugt einen 6-stelligen Verifizierungscode als Text, z.B. "048213".
// padStart(6, "0") sorgt dafür, dass auch kleine Zahlen 6 Stellen haben.
export function generiereCode() {
  return Math.floor(Math.random() * 1000000)
    .toString()
    .padStart(6, '0');
}

// ---------------------------------------------------------------------
// 🚨🚨🚨 ACHTUNG BACKEND-TEAM – HIER MUSS ETWAS GEBAUT WERDEN! 🚨🚨🚨
// ---------------------------------------------------------------------
// WAS PASSIERT JETZT (nur Test/Frontend):
//   - Wir erzeugen den Code im Browser und geben ihn nur in der Konsole aus.
//   - Es wird KEINE echte E-Mail verschickt.
//
// WAS DAS BACKEND BAUEN MUSS:
//   1. Einen Endpunkt bereitstellen, z.B.  POST /api/email/verify-code
//        Body: { "email": "kunde@gmail.com" }
//   2. Im Backend einen Code erzeugen (NICHT im Frontend, sonst kann man
//      ihn in den Entwicklertools auslesen!) und in der Datenbank speichern.
//   3. Den Code per echtem Mailversand über das Gmail-Konto verschicken.
//      Empfohlen: Nodemailer + Gmail SMTP  (smtp.gmail.com, Port 587)
//      ODER die Gmail-API mit einem App-Passwort.
//   4. Einen zweiten Endpunkt zum Prüfen, z.B.  POST /api/email/check-code
//        Body: { "email": "...", "code": "048213" }
//        Antwort: { "gueltig": true/false }
//
// Sobald der Endpunkt steht, hier nur noch den fetch() aktivieren:
//
//   await fetch('/api/email/verify-code', {
//     method: 'POST',
//     headers: { 'Content-Type': 'application/json' },
//     body: JSON.stringify({ email })
//   });
// ---------------------------------------------------------------------
export async function sendeVerifizierungsEmail(email, code) {
  console.log(`📧 [TEST] Verifizierungscode für ${email}: ${code}`);
  // Im Testbetrieb geben wir den Code zusätzlich zurück, damit die UI ihn
  // anzeigen kann. SOBALD das Backend angebunden ist, wird hier nichts mehr
  // zurückgegeben (der Nutzer bekommt den Code dann per echter E-Mail).
  return { erfolg: true, testCode: code };
}

// ---------------------------------------------------------------------
// 🚨🚨🚨 ACHTUNG BACKEND-TEAM – AUCH HIER ANBINDUNG NÖTIG! 🚨🚨🚨
// ---------------------------------------------------------------------
// Diese Funktion soll später bei einer Bestellung eine Bestätigungs-Mail
// mit Bestellhinweis und voraussichtlichem Liefertermin verschicken.
//
// WAS DAS BACKEND BAUEN MUSS:
//   - Endpunkt z.B.  POST /api/email/order-confirmation
//       Body: { "email": "...", "bestellung": {...}, "liefertermin": "..." }
//   - Im Backend die echte E-Mail über das Gmail-Konto verschicken.
// ---------------------------------------------------------------------
export async function sendeBestellBestaetigung(email, bestellung, liefertermin) {
  console.log(`📦 [TEST] Bestellbestätigung an ${email} – Liefertermin: ${liefertermin}`, bestellung);
  return { erfolg: true };
}
