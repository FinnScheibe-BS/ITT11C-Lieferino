import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { api, getToken } from '$lib/api.js';

// ⭐ BEWERTUNGEN-STORE
// Eingeloggt: Reviews kommen aus dem Backend (DB) – so sehen alle dieselben.
// Nicht eingeloggt / Backend aus: Fallback auf localStorage.
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

// 🗄️ Lädt die Bewertungen eines Restaurants aus dem Backend (öffentlich).
export async function ladeBewertungen(slug) {
  const res = await api('/api/restaurants/' + slug + '/reviews');
  if (res.ok && Array.isArray(res.daten)) {
    bewertungen.update((daten) => ({ ...daten, [slug]: res.daten }));
  }
}

// Fügt eine neue Bewertung hinzu. Eingeloggt -> ans Backend (DB), das auch
// serverseitig prüft, ob man dort bestellt hat. Gibt { ok, status } zurück.
export async function bewertungHinzufuegen(slug, neueBewertung) {
  if (getToken()) {
    const res = await api('/api/restaurants/' + slug + '/reviews', {
      method: 'POST',
      body: {
        name: neueBewertung.name,
        sterne: neueBewertung.sterne,
        text: neueBewertung.text,
        restaurantName: neueBewertung.restaurantName
      }
    });
    if (res.ok) await ladeBewertungen(slug); // echte Liste vom Server holen
    return { ok: res.ok, status: res.status, daten: res.daten };
  }

  // Fallback (nicht eingeloggt): lokal speichern.
  bewertungen.update((daten) => {
    const liste = daten[slug] || [];
    return {
      ...daten,
      [slug]: [{ ...neueBewertung, datum: new Date().toISOString() }, ...liste]
    };
  });
  return { ok: true };
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

// ✅ Bewertungen laufen jetzt über das Backend:
//   GET  /api/restaurants/:slug/reviews      (Liste, öffentlich)
//   POST /api/restaurants/:slug/reviews      (neu, nur nach Bestellung – geschützt)
// Löschen/Bearbeiten (oben) sind aktuell noch lokal (kommt mit dem Admin-Backend).
