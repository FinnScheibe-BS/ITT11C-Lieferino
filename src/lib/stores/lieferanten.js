// $lib/stores/lieferanten.js
import { writable } from 'svelte/store';
import { holeRestaurants } from '$lib/api/restaurantService.js';

export const aktiveRestaurants = writable([]);
export const deaktivierteLieferanten = [];

// Hilfsfunktion zum Erstellen eines Slugs
export function nameToSlug(name) {
  return name
    .toLowerCase()
    .replace(/ä/g, 'ae')
    .replace(/ö/g, 'oe')
    .replace(/ü/g, 'ue')
    .replace(/ß/g, 'ss')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
}

// Lade-Funktion
export async function ladeRestaurants() {
  try {
    const daten = await holeRestaurants();
    if (daten && daten.length > 0) {
      const mitSlugs = daten.map(r => ({
        ...r,
        slug: r.slug || nameToSlug(r.name)
      }));
      aktiveRestaurants.set(mitSlugs);
      return mitSlugs;
    }
  } catch (error) {
    console.error('Fehler beim Laden der Restaurants:', error);
  }
  return [];
}

// ❌ FALSCH: aktuelleRestaurants (mit "l")
// ✅ RICHTIG: aktiveRestaurants (ohne "l")
export function findeRestaurant(slug) {
  let gefunden = undefined;
  
  aktiveRestaurants.subscribe(restaurants => {
    gefunden = restaurants.find(r => r.slug === slug);
  })();
  
  return gefunden;
}

// Alternative ohne Subscription
export function findeRestaurantSync(slug) {
  let ergebnis = undefined;
  aktiveRestaurants.update(r => {
    ergebnis = r.find(restaurant => restaurant.slug === slug);
    return r;
  });
  return ergebnis;
}