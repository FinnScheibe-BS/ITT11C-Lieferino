import { writable, derived, get } from 'svelte/store';
import { browser } from '$app/environment';
import { api } from '$lib/api/api.js';
import { holeRestaurants } from '$lib/api/restaurantService.js';

// 📡 BASIS-RESTAURANTS – live von der API
export const restaurantsRoh = writable([]);

// ✅ Funktion richtig benannt
export async function ladeRestaurants() {
  const daten = await holeRestaurants();
  if (daten && daten.length > 0) {
    restaurantsRoh.set(daten);
  }
  return daten; // ← Return für direkten Aufruf
}

if (browser) ladeRestaurants();

// 🔌 DEAKTIVIERTE LIEFERANTEN
function ladeStart() {
  if (!browser) return [];
  try {
    return JSON.parse(localStorage.getItem('lieferino_deaktivierte_lieferanten') || '[]');
  } catch {
    return [];
  }
}

export const deaktivierteLieferanten = writable(ladeStart());

if (browser) {
  deaktivierteLieferanten.subscribe((liste) => {
    localStorage.setItem('lieferino_deaktivierte_lieferanten', JSON.stringify(liste));
  });
}

export function setzeLieferantAktiv(slug, aktiv) {
  deaktivierteLieferanten.update((liste) => {
    if (aktiv) {
      return liste.filter((s) => s !== slug);
    }
    return liste.includes(slug) ? liste : [...liste, slug];
  });
}

// 🥚 GEHEIMES RESTAURANT
const GEHEIMES_RESTAURANT = {
  slug: 'drachen-grill',
  name: '🐉 Drachen-Grill (geheim)',
  typ: 'geheim',
  emoji: '🐲',
  bewertung: 5.0,
  lieferzeit: '6-6 min',
  minBestell: 0,
  oeffnetUm: '00:00',
  schliesstUm: '23:59',
  beschreibung: 'Das legendäre Geheim-Restaurant. Nur für Eingeweihte. 🔥',
  geheim: true,
  speisekarte: [
    { id: 1, name: 'Drachenpizza 🐉', preis: 6.66, beschreibung: 'Höllisch scharf', veg: false, allergene: ['Gluten'] },
    { id: 2, name: 'Lava-Wings', preis: 9.99, beschreibung: 'Mit Geheimsauce', veg: false, allergene: [] },
    { id: 3, name: 'Goldener Apfel', preis: 4.2, beschreibung: 'Gibt +2 Herzen ❤️', veg: true, allergene: [] }
  ]
};

function ladeGeheim() {
  if (!browser) return false;
  return localStorage.getItem('lieferino_geheim') === 'true';
}
export const geheimFreigeschaltet = writable(ladeGeheim());
if (browser) {
  geheimFreigeschaltet.subscribe((v) => localStorage.setItem('lieferino_geheim', String(v)));
}
export function geheimFreischalten() {
  geheimFreigeschaltet.set(true);
}

// 🗄️ BACKEND-DEAKTIVIERTE
export const dbDeaktiviert = writable([]);

export async function ladeAktivStatus() {
  const res = await api('/api/restaurants');
  if (res.ok && Array.isArray(res.daten)) {
    dbDeaktiviert.set(res.daten.filter((r) => r.aktiv === false).map((r) => r.slug));
  }
}

if (browser) ladeAktivStatus();

// ✅ NUR AKTIVE RESTAURANTS
export const aktiveRestaurants = derived(
  [restaurantsRoh, deaktivierteLieferanten, dbDeaktiviert, geheimFreigeschaltet],
  ([$restaurantsRoh, $deaktiviert, $dbDeaktiviert, $geheim]) => {
    const aus = new Set([...$deaktiviert, ...$dbDeaktiviert]);
    const sichtbar = $restaurantsRoh.filter((r) => !aus.has(r.slug));
    return $geheim ? [GEHEIMES_RESTAURANT, ...sichtbar] : sichtbar;
  }
);

// ✅ FINDE RESTAURANT – sucht jetzt im Store, nicht in static data
export function findeRestaurant(slugOderName) {
  // Geheim-Restaurant
  if (slugOderName === GEHEIMES_RESTAURANT.slug || slugOderName === GEHEIMES_RESTAURANT.name) {
    return GEHEIMES_RESTAURANT;
  }
  
  // 🔍 Im Store suchen (mit get() für synchronen Zugriff)
  const alle = get(restaurantsRoh);
  return alle.find(r => 
    r.slug === slugOderName || 
    r.name.toLowerCase() === slugOderName.toLowerCase()
  ) || null;
}

// Helper: Slug erstellen (für Fallback-Suche)
export function erstelleSlug(name) {
  return name
    .toLowerCase()
    .replace(/ä/g, 'ae')
    .replace(/ö/g, 'oe')
    .replace(/ü/g, 'ue')
    .replace(/ß/g, 'ss')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
}