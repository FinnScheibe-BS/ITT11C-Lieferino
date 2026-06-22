import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// ❤️ FAVORITEN-STORE
// Speichert die Slugs der Lieblings-Restaurants im localStorage.

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
  favoriten.update((liste) => {
    if (liste.includes(slug)) {
      return liste.filter((s) => s !== slug); // war Favorit -> entfernen
    }
    return [...liste, slug]; // war kein Favorit -> hinzufügen
  });
}
