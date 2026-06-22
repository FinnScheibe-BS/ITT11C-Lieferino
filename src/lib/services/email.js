// =====================================================================
// 📧 E-MAIL-SERVICE (Frontend-Seite)
// Erzeugt Codes und ruft die Server-Route /api/email auf, die die echte
// E-Mail verschickt. Solange in der .env keine Gmail-Zugangsdaten stehen,
// läuft alles im "Test-Modus": es wird nichts verschickt und der Code wird
// in der App direkt angezeigt (damit man weiterkommt).
//
// 👉 Echten Versand aktivieren: siehe EMAIL-SETUP.md
// =====================================================================

// Erzeugt einen 6-stelligen Verifizierungscode als Text, z.B. "048213".
export function generiereCode() {
  return Math.floor(Math.random() * 1000000)
    .toString()
    .padStart(6, '0');
}

// Ruft die Server-Route auf, die die Mail verschickt.
// Gibt { gesendet, testModus? } zurück.
async function sendeMail(an, betreff, text) {
  try {
    const antwort = await fetch('/api/email', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ an, betreff, text })
    });
    return await antwort.json();
  } catch {
    // Server nicht erreichbar -> Test-Modus.
    return { gesendet: false, testModus: true };
  }
}

// Verschickt den Verifizierungscode.
// Im Test-Modus geben wir den Code als testCode zurück, damit die UI ihn
// anzeigen kann. Bei echtem Versand bleibt testCode leer.
export async function sendeVerifizierungsEmail(email, code) {
  const ergebnis = await sendeMail(
    email,
    'Dein Lieferino-Verifizierungscode',
    `Hallo!\n\nDein Verifizierungscode lautet: ${code}\n\nViele Grüße\nDein Lieferino-Team`
  );
  return {
    erfolg: true,
    testCode: ergebnis.gesendet ? '' : code,
    testModus: !ergebnis.gesendet
  };
}

// Verschickt eine Bestellbestätigung mit Liefertermin.
export async function sendeBestellBestaetigung(email, bestellung, liefertermin) {
  const zeilen = (bestellung.artikel || [])
    .map((a) => `- ${a.menge}x ${a.name}`)
    .join('\n');
  await sendeMail(
    email,
    'Deine Lieferino-Bestellung',
    `Vielen Dank für deine Bestellung!\n\n${zeilen}\n\nSumme: ${bestellung.summe?.toFixed(2)}€\nVoraussichtliche Lieferung bis: ${liefertermin} Uhr\n\nGuten Appetit!\nDein Lieferino-Team`
  );
  return { erfolg: true };
}
