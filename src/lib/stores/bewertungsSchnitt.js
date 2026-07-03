import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { api } from '$lib/api/api.js';

// ⭐ DURCHSCHNITTS-BEWERTUNGEN (live aus der DB)
// Hält pro Restaurant-Slug { schnitt, anzahl } aus den echten Bewertungen.
// So zeigen die Kacheln die "lebendige" Sterne-Bewertung aller Nutzer.

export const bewertungsSchnitt = writable({});

export async function ladeBewertungsSchnitt() {
  const res = await api('/api/reviews/schnitt');
  if (res.ok && res.daten && typeof res.daten === 'object') {
    bewertungsSchnitt.set(res.daten);
  }
}

// Beim Start einmal laden (öffentlich, kein Login nötig).
if (browser) ladeBewertungsSchnitt();
