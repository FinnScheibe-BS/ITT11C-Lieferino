// $lib/stores/speisen.js
import { writable } from 'svelte/store';
import { holeGerichte } from '$lib/api/gerichtService.js';

export const speisen = writable({});

export async function ladeSpeisen(restaurantId) {
  const gerichte = await holeGerichte(restaurantId);
  
  // Mapping der API-Felder auf App-Felder
  const gemappteGerichte = gerichte.map(g => ({
    id: g.id,
    name: g.name,
    beschreibung: '', // API hat keine Beschreibung - optional hinzufügen
    preis: g.preis,
    veg: g.vegetarisch,
    vegan: g.vegan,
    allergene: [], // API hat keine Allergene - optional hinzufügen
    kategorie_id: g.kategorie_id
  }));
  
  speisen.update(s => ({ ...s, [restaurantId]: gemappteGerichte }));
  return gemappteGerichte;
}