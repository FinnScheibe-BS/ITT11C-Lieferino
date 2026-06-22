import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// 🌙/☀️ THEME-STORE (Dark-/Light-Mode)
// Merkt sich die Auswahl im localStorage, damit sie nach dem Neuladen erhalten bleibt.

function ladeStart() {
  if (!browser) return 'light';
  return localStorage.getItem('lieferino_theme') || 'light';
}

export const theme = writable(ladeStart());

// Bei jeder Änderung: speichern UND am <html>-Element setzen, damit unser
// globales CSS (siehe +layout.svelte) das passende Farbschema anwenden kann.
if (browser) {
  theme.subscribe((wert) => {
    localStorage.setItem('lieferino_theme', wert);
    document.documentElement.setAttribute('data-theme', wert);
  });
}

// Schaltet zwischen hell und dunkel um.
export function themeWechseln() {
  theme.update((t) => (t === 'light' ? 'dark' : 'light'));
}
