// lib/api/cartApi.js
import { getToken } from '$lib/api/api.js';

const API_BASE = 'http://172.30.4.90:8080';

/**
 * Holt den aktuellen Warenkorb des Benutzers vom Server.
 */
export async function holeWarenkorb() {
  const token = getToken();
  if (!token) return { ok: false, status: 401, daten: [] };

  try {
    const antwort = await fetch(`${API_BASE}/api/cart`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (!antwort.ok) return { ok: false, status: antwort.status, daten: [] };
    const daten = await antwort.json();
    return { ok: true, status: antwort.status, daten };
  } catch (err) {
    console.error('Warenkorb konnte nicht geladen werden:', err);
    return { ok: false, status: 0, daten: [], offline: true };
  }
}

/**
 * Aktualisiert den Warenkorb auf dem Server (sendet das gesamte Array).
 */
export async function speichereWarenkorb(artikelListe) {
  const token = getToken();
  if (!token) return { ok: false, status: 401 };

  try {
    const antwort = await fetch(`${API_BASE}/api/cart`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(artikelListe)
    });

    return { ok: antwort.ok, status: antwort.status };
  } catch (err) {
    console.error('Warenkorb-Synchronisation fehlgeschlagen:', err);
    return { ok: false, status: 0, offline: true };
  }
}