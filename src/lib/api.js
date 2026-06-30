import { browser } from '$app/environment';

// 🌐 API-CLIENT fürs Go-Backend
// Basis-URL des Backends. In der Entwicklung läuft das Backend (Gin) auf Port 8090.
export const API_BASE = 'http://localhost:8090';

// --- JWT-Token (kommt vom Login/Registrierung) ---
export function getToken() {
  return browser ? localStorage.getItem('lieferino_token') : null;
}
export function setzeToken(t) {
  if (browser) localStorage.setItem('lieferino_token', t);
}
export function loescheToken() {
  if (browser) localStorage.removeItem('lieferino_token');
}

// Zentraler fetch-Helfer: hängt JSON-Header + JWT-Token an und liefert
// ein einheitliches Ergebnis { ok, status, daten }.
export async function api(pfad, { method = 'GET', body } = {}) {
  const headers = { 'Content-Type': 'application/json' };
  const t = getToken();
  if (t) headers['Authorization'] = 'Bearer ' + t;

  try {
    const res = await fetch(API_BASE + pfad, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined
    });
    let daten = null;
    try {
      daten = await res.json();
    } catch {
      /* keine JSON-Antwort */
    }
    return { ok: res.ok, status: res.status, daten };
  } catch (fehler) {
    // Backend nicht erreichbar
    return { ok: false, status: 0, daten: null, offline: true };
  }
}
