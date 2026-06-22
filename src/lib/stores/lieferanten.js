import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { restaurants } from '$lib/data';

// 🔌 AKTIV/DEAKTIVIERT-STATUS DER LIEFERANTEN
// Die Basis-Restaurants kommen aus der (unveränderlichen) zentralen Quelle.
// Der Admin kann Lieferanten "deaktivieren" – wir merken uns die Slugs der
// deaktivierten Lieferanten im localStorage. Deaktivierte Lieferanten sind für
// normale Nutzer (öffentliche Listen) unsichtbar.
//
// 🚨🚨🚨 BACKEND-HINWEIS 🚨🚨🚨
// Der Aktiv/Deaktiviert-Status liegt aktuell nur lokal im Browser. Das Backend
// MUSS diesen Status speichern und ausliefern, z.B.:
//   PATCH /api/admin/restaurants/:slug   Body: { aktiv: true|false }
//   GET   /api/restaurants               liefert nur AKTIVE für normale Nutzer
// Dann lädt das Frontend den Status vom Server statt aus dem localStorage.

function ladeStart() {
  if (!browser) return [];
  try {
    return JSON.parse(localStorage.getItem('lieferino_deaktivierte_lieferanten') || '[]');
  } catch {
    return [];
  }
}

// Liste der Slugs, die deaktiviert sind.
export const deaktivierteLieferanten = writable(ladeStart());

if (browser) {
  deaktivierteLieferanten.subscribe((liste) => {
    localStorage.setItem('lieferino_deaktivierte_lieferanten', JSON.stringify(liste));
  });
}

// Setzt ein Restaurant auf aktiv (true) oder deaktiviert (false).
export function setzeLieferantAktiv(slug, aktiv) {
  deaktivierteLieferanten.update((liste) => {
    if (aktiv) {
      return liste.filter((s) => s !== slug); // aus der Deaktiviert-Liste entfernen
    }
    return liste.includes(slug) ? liste : [...liste, slug]; // hinzufügen
  });
}

// Prüft, ob ein Lieferant aktiv ist.
export function istAktiv(slug, deaktivierteListe) {
  return !deaktivierteListe.includes(slug);
}

// Nur die AKTIVEN Restaurants – für alle öffentlichen Listen (normale Nutzer).
export const aktiveRestaurants = derived(deaktivierteLieferanten, ($deaktiviert) =>
  restaurants.filter((r) => !$deaktiviert.includes(r.slug))
);
