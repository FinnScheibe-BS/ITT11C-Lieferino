import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// ⭐ TREUEPUNKTE
// Pro ausgegebenem Euro gibt es 1 Punkt. Punkte kann man im Checkout einlösen
// (100 Punkte = 5 € Rabatt). Alles lokal im localStorage.
//
// 🚨 BACKEND-HINWEIS: Treuepunkte müssen später serverseitig (pro Konto) gespeichert
// werden, sonst kann man sie im Browser einfach manipulieren.

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

// Punkte für einen Bestellbetrag gutschreiben.
export function punkteSammeln(betrag) {
  treuepunkte.update((p) => p + Math.floor(betrag * PUNKTE_PRO_EURO));
}

// Punkte abziehen (beim Einlösen).
export function punkteAbziehen(anzahl) {
  treuepunkte.update((p) => Math.max(0, p - anzahl));
}
