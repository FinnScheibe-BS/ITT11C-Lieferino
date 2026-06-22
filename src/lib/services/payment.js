// =====================================================================
// 💳 ZAHLUNGS-SERVICE (Kreditkarten-Prüfung)
// Prüft im Frontend, ob eine eingegebene Kreditkarte gültig AUSSIEHT.
// Achtung: Das ist nur eine FORMAT-Prüfung – ob die Karte echt gedeckt ist,
// kann nur ein echter Zahlungsanbieter über das Backend feststellen.
// =====================================================================

// Erkennt den Kartentyp anhand der ersten Ziffern.
export function kartenTyp(nummer) {
  const n = nummer.replace(/\s/g, '');
  if (/^4/.test(n)) return 'Visa';
  if (/^5[1-5]/.test(n)) return 'Mastercard';
  if (/^3[47]/.test(n)) return 'American Express';
  return 'Unbekannt';
}

// ✅ Luhn-Algorithmus: die Standard-Prüfsumme für Kreditkartennummern.
// Jede zweite Ziffer (von rechts) wird verdoppelt; ist die Summe durch 10
// teilbar, ist die Nummer formal gültig.
export function pruefeKartennummer(nummer) {
  const ziffern = nummer.replace(/\s/g, '');
  // Muss 13–19 Ziffern haben und nur aus Zahlen bestehen.
  if (!/^\d{13,19}$/.test(ziffern)) return false;

  let summe = 0;
  let verdoppeln = false;
  // Von rechts nach links durchgehen.
  for (let i = ziffern.length - 1; i >= 0; i--) {
    let z = parseInt(ziffern[i], 10);
    if (verdoppeln) {
      z *= 2;
      if (z > 9) z -= 9; // z.B. 16 -> 1+6 = 7
    }
    summe += z;
    verdoppeln = !verdoppeln;
  }
  return summe % 10 === 0;
}

// 📅 Prüft das Ablaufdatum im Format MM/JJ (z.B. "07/27").
// Gültig, wenn das Datum in der Zukunft liegt.
export function pruefeAblauf(ablauf) {
  const treffer = /^(\d{2})\/(\d{2})$/.exec(ablauf.trim());
  if (!treffer) return false;

  const monat = parseInt(treffer[1], 10);
  const jahr = 2000 + parseInt(treffer[2], 10);
  if (monat < 1 || monat > 12) return false;

  const jetzt = new Date();
  // Letzter gültiger Tag = letzter Tag des Ablaufmonats.
  const ablaufDatum = new Date(jahr, monat, 0, 23, 59, 59);
  return ablaufDatum >= jetzt;
}

// 🔢 Prüft den CVV/Sicherheitscode (3 Ziffern, bei Amex 4).
export function pruefeCvv(cvv, typ) {
  const laenge = typ === 'American Express' ? 4 : 3;
  return new RegExp(`^\\d{${laenge}}$`).test(cvv.trim());
}

// Fügt beim Tippen automatisch Leerzeichen ein (4er-Blöcke) für die Anzeige.
export function formatiereNummer(nummer) {
  return nummer
    .replace(/\s/g, '')
    .replace(/(\d{4})(?=\d)/g, '$1 ')
    .trim();
}

// Gesamt-Prüfung: gibt zurück, ob die Karte formal gültig ist + welche Felder fehlen.
export function pruefeKarte({ nummer, ablauf, cvv }) {
  const typ = kartenTyp(nummer);
  const fehler = {
    nummer: !pruefeKartennummer(nummer),
    ablauf: !pruefeAblauf(ablauf),
    cvv: !pruefeCvv(cvv, typ)
  };
  const gueltig = !fehler.nummer && !fehler.ablauf && !fehler.cvv;
  return { gueltig, fehler, typ };
}

// ---------------------------------------------------------------------
// 🚨🚨🚨 ACHTUNG BACKEND-TEAM – ECHTE ZAHLUNG NUR ÜBER BACKEND! 🚨🚨🚨
// ---------------------------------------------------------------------
// Diese Prüfungen sagen NUR, ob die Eingaben formal stimmen. Eine echte
// Abbuchung darf NIEMALS im Frontend passieren und Kartendaten dürfen nicht
// gespeichert werden!
//
// WAS DAS BACKEND BAUEN MUSS:
//   - Einen Zahlungsanbieter anbinden (z.B. Stripe oder PayPal).
//   - Das Frontend schickt nur ein sicheres Token an:  POST /api/zahlung
//   - Der Anbieter prüft Deckung, führt die Zahlung aus und meldet das Ergebnis.
//   - Kartennummer/CVV werden NIE in unserer Datenbank gespeichert (PCI-DSS).
// ---------------------------------------------------------------------
