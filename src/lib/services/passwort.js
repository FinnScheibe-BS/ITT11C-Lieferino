// =====================================================================
// 🔒 PASSWORT-SICHERHEIT (IT-Sicherheit)
// Diese Funktion prüft im Frontend, ob ein Passwort "stark genug" ist,
// und gibt eine Bewertung zurück, die wir dem Nutzer anzeigen können.
// =====================================================================

// 🔐 Diese Regeln MÜSSEN mit dem Backend übereinstimmen
// (backend/handlers/passwort.go), sonst sieht der Nutzer "grün", wird aber
// trotzdem abgelehnt. Daher: mind. 10 Zeichen UND alle 5 Regeln.
export const MIN_LAENGE = 10;

// Prüft ein Passwort gegen mehrere Regeln und gibt ein Ergebnis-Objekt zurück.
// score = wie viele Regeln erfüllt sind (0 bis 5).
// stufeKey = i18n-Schlüssel für die Stärke (Komponente übersetzt via $t).
export function pruefePasswortStaerke(passwort) {
  // Jede Regel ist true/false. So sieht der Nutzer genau, was noch fehlt.
  const regeln = {
    laenge: passwort.length >= MIN_LAENGE, // mindestens 10 Zeichen
    grossbuchstabe: /[A-Z]/.test(passwort), // mind. ein Großbuchstabe
    kleinbuchstabe: /[a-z]/.test(passwort), // mind. ein Kleinbuchstabe
    zahl: /[0-9]/.test(passwort), // mind. eine Zahl
    sonderzeichen: /[^A-Za-z0-9]/.test(passwort) // mind. ein Sonderzeichen
  };

  // Wir zählen, wie viele Regeln erfüllt sind.
  const score = Object.values(regeln).filter(Boolean).length;

  // Ein passender i18n-Schlüssel + Farbe zur Stärke.
  let stufeKey = 'pw.very_weak';
  let farbe = '#ff3b30'; // rot
  if (score === 3) {
    stufeKey = 'pw.medium';
    farbe = '#ff9500'; // orange
  } else if (score === 4) {
    stufeKey = 'pw.strong';
    farbe = '#34c759'; // grün
  } else if (score === 5) {
    stufeKey = 'pw.very_strong';
    farbe = '#248a3d'; // dunkelgrün
  }

  // "Sicher genug" = ALLE Regeln erfüllt (gleiche Schwelle wie das Backend).
  const istSicher = score === 5;

  return { regeln, score, stufeKey, farbe, istSicher };
}

// ---------------------------------------------------------------------
// 🚨🚨🚨 WICHTIGER HINWEIS FÜRS BACKEND-TEAM 🚨🚨🚨
// ---------------------------------------------------------------------
// Mit dem Kollegen abgesprochen: In diesem Schulprojekt dürfen E-Mail und
// Passwort im Klartext gespeichert/übertragen werden (KEINE Zertifikate).
//
// WICHTIG ZU WISSEN (für die echte Welt / falls das je live geht):
//   - Passwörter NIEMALS im Klartext speichern! Im Backend muss das Passwort
//     mit einem sicheren Hash-Verfahren gespeichert werden (z.B. bcrypt).
//   - Die obige Prüfung läuft im Frontend und ist nur "Komfort". Das Backend
//     MUSS die gleiche Mindeststärke nochmal prüfen, denn Frontend-Checks
//     kann man umgehen.
// ---------------------------------------------------------------------
