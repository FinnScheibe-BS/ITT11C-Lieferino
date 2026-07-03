import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { api, getToken } from '$lib/api/api.js';

// ❤️ FAVORITEN-STORE
// Speichert die Slugs der Lieblings-Restaurants. Eingeloggt: zusätzlich im
// Backend (Datenbank), damit sie geräteübergreifend erhalten bleiben.

function ladeStart() {
  if (!browser) return [];
  try {
    return JSON.parse(localStorage.getItem('lieferino_favoriten') || '[]');
  } catch {
    return [];
  }
}

export const favoriten = writable(ladeStart());

// Jede Änderung speichern.
if (browser) {
  favoriten.subscribe((liste) => {
    localStorage.setItem('lieferino_favoriten', JSON.stringify(liste));
  });
}

// Schaltet ein Restaurant als Favorit an/aus.
export function toggleFavorit(slug) {
  let nunFavorit;
  favoriten.update((liste) => {
    if (liste.includes(slug)) {
      nunFavorit = false;
      return liste.filter((s) => s !== slug); // war Favorit -> entfernen
    }
    nunFavorit = true;
    return [...liste, slug]; // war kein Favorit -> hinzufügen
  });

  // 🗄️ Wenn eingeloggt: Änderung auch im Backend (DB) speichern.
  if (getToken()) {
    api('/api/favorites/' + slug, { method: nunFavorit ? 'POST' : 'DELETE' });
  }
}

// 🗄️ Lädt die Favoriten aus dem Backend (wird nach dem Login aufgerufen).
export async function ladeFavoritenVomBackend() {
  if (!getToken()) return;
  const res = await api('/api/favorites');
  if (res.ok && Array.isArray(res.daten)) {
    favoriten.set(res.daten);
  }
}

// Beim Start (eingeloggt) gleich die Favoriten aus der DB holen.
if (browser) {
  ladeFavoritenVomBackend();
}
