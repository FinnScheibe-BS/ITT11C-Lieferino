import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// 🎉 FUN-MODUS (Easter-Egg-Spielereien)
// emojiCursor: Food-Emojis folgen der Maus.
// saison: fallende Emojis (z.B. Schnee) über der ganzen Seite.
// Beides wird über versteckte Such-Codewörter an-/ausgeschaltet.

function lade(key) {
  if (!browser) return false;
  return localStorage.getItem(key) === 'true';
}

export const emojiCursor = writable(lade('lieferino_emojicursor'));
export const saison = writable(lade('lieferino_saison'));

if (browser) {
  emojiCursor.subscribe((v) => localStorage.setItem('lieferino_emojicursor', String(v)));
  saison.subscribe((v) => localStorage.setItem('lieferino_saison', String(v)));
}

export function toggleEmojiCursor() {
  emojiCursor.update((v) => !v);
}
export function toggleSaison() {
  saison.update((v) => !v);
}
