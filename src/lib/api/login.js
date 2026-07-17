// lib/api/login.js
// Direkte Anbindung an das Backend für den Login.
// TODO: IP ist aktuell hart codiert (Dev-Backend). Für Produktion besser
// über eine Umgebungsvariable lösen, z.B. PUBLIC_API_BASE_URL in einer .env
// und hier: import { PUBLIC_API_BASE_URL } from '$env/static/public';
const API_BASE = 'http://172.30.4.90:8080';

/**
 * Meldet einen Nutzer mit E-Mail + Passwort beim Backend an.
 * Heißt bewusst nicht "login", weil $lib/stores/auth.js schon eine
 * login()-Funktion exportiert (setzt den eingeloggt-Store) - sonst
 * Namenskollision beim Import in derselben Datei.
 * @param {{ email: string, passwort: string }} daten
 * @returns {Promise<{ ok: boolean, status: number, daten: any, offline: boolean }>}
 */
export async function anmelden(daten) {
  try {
    const antwort = await fetch(`${API_BASE}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(daten)
    });

    let inhalt = null;
    try {
      inhalt = await antwort.json();
    } catch {
      inhalt = null;
    }

    return {
      ok: antwort.ok,
      status: antwort.status,
      daten: inhalt,
      offline: false
    };
  } catch (err) {
    console.error('Login: Backend nicht erreichbar', err);
    return { ok: false, status: 0, daten: null, offline: true };
  }
}