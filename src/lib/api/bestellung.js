// lib/api/bestellung.js
// Direkte Anbindung ans Backend für den Bestellverlauf.
// TODO: IP ist aktuell hart codiert (Dev-Backend). Für Produktion besser
// über eine Umgebungsvariable lösen, z.B. PUBLIC_API_BASE_URL in einer .env
// und hier: import { PUBLIC_API_BASE_URL } from '$env/static/public';
import { getToken } from '$lib/api/api.js';

const API_BASE = 'http://172.30.4.90:8080';

/**
 * Holt die Bestellungen des eingeloggten Kunden vom Backend.
 * Braucht ein Login (Token), weil das Backend sonst nicht weiß,
 * welcher Kunde (FK_ID_Kunde) gemeint ist.
 * @returns {Promise<{ ok: boolean, status: number, daten: any, offline: boolean }>}
 */
export async function holeBestellungen() {
  const token = getToken();
  if (!token) {
    // Kein Login -> gar nicht erst anfragen, Backend würde eh ablehnen.
    return { ok: false, status: 401, daten: null, offline: false };
  }

  try {
    const antwort = await fetch(`${API_BASE}/api/orders`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    let inhalt = null;
    try {
      inhalt = await antwort.json();
    } catch {
      // Backend hat evtl. keinen Body zurückgeschickt (z.B. bei 204/401)
      inhalt = null;
    }

    return {
      ok: antwort.ok,
      status: antwort.status,
      daten: inhalt,
      offline: false
    };
  } catch (err) {
    console.error('Bestellungen: Backend nicht erreichbar', err);
    return { ok: false, status: 0, daten: null, offline: true };
  }
}