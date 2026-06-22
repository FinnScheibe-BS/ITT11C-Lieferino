import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { restaurants } from '$lib/data';

// 🗑️ VERWALTUNG GELÖSCHTER LIEFERANTEN
// Die Basis-Restaurants kommen aus der (unveränderlichen) zentralen Quelle.
// Damit der Admin Lieferanten "löschen" kann, merken wir uns hier die Slugs
// der ausgeblendeten Lieferanten im localStorage. Öffentliche Listen filtern
// diese dann heraus.
//
// 🚨 BACKEND-HINWEIS: Das ist ein reines Frontend-Ausblenden. Echtes Löschen/
// Verwalten von Restaurants gehört ins Backend (DELETE /api/admin/restaurants/:id).

function ladeStart() {
  if (!browser) return [];
  try {
    return JSON.parse(localStorage.getItem('lieferino_geloeschte_lieferanten') || '[]');
  } catch {
    return [];
  }
}

export const geloeschteLieferanten = writable(ladeStart());

if (browser) {
  geloeschteLieferanten.subscribe((liste) => {
    localStorage.setItem('lieferino_geloeschte_lieferanten', JSON.stringify(liste));
  });
}

export function loescheLieferant(slug) {
  geloeschteLieferanten.update((l) => (l.includes(slug) ? l : [...l, slug]));
}

export function stelleLieferantWiederHer(slug) {
  geloeschteLieferanten.update((l) => l.filter((s) => s !== slug));
}

// Nur die NICHT gelöschten Restaurants – für alle öffentlichen Listen.
export const aktiveRestaurants = derived(geloeschteLieferanten, ($geloescht) =>
  restaurants.filter((r) => !$geloescht.includes(r.slug))
);
