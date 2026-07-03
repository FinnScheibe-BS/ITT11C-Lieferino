import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { api, getToken } from '$lib/api/api.js';

// ⭐ TREUEPUNKTE
// Der ECHTE Punktestand liegt jetzt SERVERSEITIG am Konto (manipulationssicher).
// Sammeln/Einlösen passiert beim Bestellen im Backend. Der localStorage-Wert ist
// nur ein Cache für die Anzeige (nicht eingeloggt / offline).

export const PUNKTE_PRO_EURO = 1;
export const EINLOESE_SCHRITT = 100; // 100 Punkte ...
export const EINLOESE_WERT = 5; // ... = 5 € Rabatt

function ladeStart() {
  if (!browser) return 0;
  return Number(localStorage.getItem('lieferino_treuepunkte') || '0');
}

export const treuepunkte = writable(ladeStart());

if (browser) {
  treuepunkte.subscribe((p) => localStorage.setItem('lieferino_treuepunkte', String(p)));
}

// 🗄️ Echten Punktestand vom Backend (Konto) laden.
export async function ladeTreuepunkte() {
  if (!getToken()) return;
  const res = await api('/api/me');
  if (res.ok && res.daten && typeof res.daten.treuepunkte === 'number') {
    treuepunkte.set(res.daten.treuepunkte);
  }
}

// Beim Start (eingeloggt) gleich den echten Stand holen.
if (browser) ladeTreuepunkte();
