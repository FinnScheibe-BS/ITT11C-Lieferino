import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// ⭐ BEWERTUNGEN-STORE
// Speichert Kunden-Reviews pro Restaurant (nach Slug) im localStorage.
// Aufbau: { "luigis-pizzeria": [ { name, sterne, text, datum }, ... ] }

function ladeStart() {
  if (!browser) return {};
  try {
    return JSON.parse(localStorage.getItem('lieferino_bewertungen') || '{}');
  } catch {
    return {};
  }
}

export const bewertungen = writable(ladeStart());

if (browser) {
  bewertungen.subscribe((daten) => {
    localStorage.setItem('lieferino_bewertungen', JSON.stringify(daten));
  });
}

// Fügt eine neue Bewertung für ein Restaurant hinzu.
export function bewertungHinzufuegen(slug, neueBewertung) {
  bewertungen.update((daten) => {
    const liste = daten[slug] || [];
    return {
      ...daten,
      [slug]: [{ ...neueBewertung, datum: new Date().toISOString() }, ...liste]
    };
  });
}

// 🗑️ Löscht eine Bewertung eines Restaurants (per Index in dessen Liste).
export function bewertungLoeschen(slug, index) {
  bewertungen.update((daten) => {
    const liste = [...(daten[slug] || [])];
    liste.splice(index, 1);
    return { ...daten, [slug]: liste };
  });
}

// ✏️ Bearbeitet Text und/oder Sterne einer Bewertung.
export function bewertungBearbeiten(slug, index, patch) {
  bewertungen.update((daten) => {
    const liste = [...(daten[slug] || [])];
    if (liste[index]) liste[index] = { ...liste[index], ...patch };
    return { ...daten, [slug]: liste };
  });
}

// 🚨 BACKEND-HINWEIS: Bewertungen liegen aktuell nur lokal im Browser.
// Damit alle Nutzer dieselben Reviews sehen, muss das Backend sie speichern:
//   POST /api/bewertungen   (neue Bewertung)
//   GET  /api/bewertungen?restaurant=slug   (alle Bewertungen laden)
