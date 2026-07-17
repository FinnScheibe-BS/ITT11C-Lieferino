import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { getRestaurant } from '$lib/data';
import { api } from '$lib/api/api.js';
import { holeRestaurants } from '$lib/api/restaurantService.js';

// 📡 BASIS-RESTAURANTS – kommen jetzt live von der API statt aus der
// statischen $lib/data-Tabelle. Wird beim Laden im Browser befüllt.
export const restaurantsRoh = writable([]);

export async function ladeRestaurantsVonApi() {
  const daten = await holeRestaurants();
  if (daten && daten.length > 0) {
    restaurantsRoh.set(daten);
  }
}

if (browser) ladeRestaurantsVonApi();

// 🔌 AKTIV/DEAKTIVIERT-STATUS DER LIEFERANTEN
// Die Basis-Restaurants kommen jetzt live von der API (siehe oben).
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

// 🥚 GEHEIMES RESTAURANT – nur per Cheat-Code (in der Suche: „dragonpizza") sichtbar.
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

// 🗄️ Vom Backend (Admin) deaktivierte Restaurants – serverseitig gesteuert.
export const dbDeaktiviert = writable([]);

export async function ladeAktivStatus() {
  const res = await api('/api/restaurants');
  if (res.ok && Array.isArray(res.daten)) {
    dbDeaktiviert.set(res.daten.filter((r) => r.aktiv === false).map((r) => r.slug));
  }
}

// Beim Start den echten Aktiv/Deaktiv-Status vom Backend holen.
if (browser) ladeAktivStatus();

// Nur die AKTIVEN Restaurants – für alle öffentlichen Listen (normale Nutzer).
// Basis sind jetzt die live von der API geladenen Restaurants (restaurantsRoh).
// Ausgeblendet wird, wer lokal ODER im Backend deaktiviert ist.
export const aktiveRestaurants = derived(
  [restaurantsRoh, deaktivierteLieferanten, dbDeaktiviert, geheimFreigeschaltet],
  ([$restaurantsRoh, $deaktiviert, $dbDeaktiviert, $geheim]) => {
    const aus = new Set([...$deaktiviert, ...$dbDeaktiviert]);
    const sichtbar = $restaurantsRoh.filter((r) => !aus.has(r.slug));
    return $geheim ? [GEHEIMES_RESTAURANT, ...sichtbar] : sichtbar;
  }
);

// Sucht ein Restaurant per Slug/Name – inkl. dem Geheim-Restaurant.
export function findeRestaurant(slugOderName) {
  if (slugOderName === GEHEIMES_RESTAURANT.slug || slugOderName === GEHEIMES_RESTAURANT.name) {
    return GEHEIMES_RESTAURANT;
  }
  return getRestaurant(slugOderName);
}