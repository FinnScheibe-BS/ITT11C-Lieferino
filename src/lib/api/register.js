// lib/api/register.js
// Direkte Anbindung an das Backend für die Registrierung.
// TODO: IP ist aktuell hart codiert (Dev-Backend). Für Produktion besser
// über eine Umgebungsvariable lösen, z.B. PUBLIC_API_BASE_URL in einer .env
// und hier: import { PUBLIC_API_BASE_URL } from '$env/static/public';
const API_BASE = 'http://172.30.4.90:8080';

/**
 * Registriert einen neuen Nutzer beim Backend.
 * @param {{ email: string, passwort: string, vorname: string, nachname: string, geburtsdatum: string }} daten
 * @returns {Promise<{ ok: boolean, status: number, daten: any, offline: boolean }>}
 */
export async function registriere(daten) {
  try {
    const antwort = await fetch(`${API_BASE}/api/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(daten)
    });

    let inhalt = null;
    try {
      inhalt = await antwort.json();
    } catch {
      // Backend hat evtl. keinen Body zurückgeschickt (z.B. bei 204)
      inhalt = null;
    }

    return {
      ok: antwort.ok,
      status: antwort.status,
      daten: inhalt,
      offline: false
    };
  } catch (err) {
    // fetch wirft z.B. bei DNS-/Verbindungsfehlern, Backend nicht erreichbar etc.
    console.error('Registrierung: Backend nicht erreichbar', err);
    return { ok: false, status: 0, daten: null, offline: true };
  }
}