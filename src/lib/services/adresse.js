// =====================================================================
// 🏠 ADRESS-PRÜFUNG
// Prüft, ob eine eingegebene Adresse tatsächlich existiert. Dafür nutzen wir
// den kostenlosen OpenStreetMap-Dienst "Nominatim". Das funktioniert direkt
// im Frontend per fetch() – ohne API-Schlüssel.
// =====================================================================

// Prüft eine Adresse und gibt zurück, ob sie gefunden wurde.
// Erwartet ein Objekt mit { strasse, hausnummer, plz, ort }.
export async function pruefeAdresse({ strasse, hausnummer, plz, ort }) {
  // Wir bauen aus den Einzelteilen eine Such-Adresse zusammen.
  const suchtext = `${strasse} ${hausnummer}, ${plz} ${ort}`;

  // encodeURIComponent sorgt dafür, dass Leerzeichen & Umlaute in der URL
  // korrekt codiert werden.
  const url = `https://nominatim.openstreetmap.org/search?format=json&limit=1&countrycodes=de&q=${encodeURIComponent(
    suchtext
  )}`;

  try {
    const antwort = await fetch(url, {
      headers: {
        // Nominatim möchte gerne wissen, wer anfragt (höfliche Nutzung).
        'Accept-Language': 'de'
      }
    });
    const ergebnisse = await antwort.json();

    // Wenn die Liste mindestens einen Treffer hat, existiert die Adresse.
    if (Array.isArray(ergebnisse) && ergebnisse.length > 0) {
      return { gefunden: true, treffer: ergebnisse[0].display_name };
    }
    return { gefunden: false };
  } catch (fehler) {
    // Falls der Dienst gerade nicht erreichbar ist, blockieren wir den Nutzer
    // nicht, sondern melden den Fehler getrennt zurück.
    console.error('Adressprüfung fehlgeschlagen:', fehler);
    return { gefunden: false, fehler: true };
  }
}

// ---------------------------------------------------------------------
// 🚨🚨🚨 HINWEIS FÜRS BACKEND-TEAM 🚨🚨🚨
// ---------------------------------------------------------------------
// Aktuell fragt das Frontend Nominatim direkt an. Das ist für ein Schul-
// projekt okay, hat aber zwei Schwächen:
//   1. Nominatim erlaubt nur wenige Anfragen pro Sekunde (Rate-Limit).
//   2. Man kann den Anbieter nicht zentral steuern/wechseln.
//
// WAS DAS BACKEND SPÄTER ÜBERNEHMEN SOLLTE:
//   - Einen eigenen Endpunkt anbieten, z.B.  POST /api/adresse/pruefen
//     der die Anfrage an einen Geocoding-Dienst weiterleitet (Nominatim,
//     Google Maps Geocoding o.ä.) und das Ergebnis zurückgibt.
//   - Dann ruft das Frontend nur noch den eigenen Endpunkt auf.
// ---------------------------------------------------------------------
