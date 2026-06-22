import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// 🔑 AUTH-STORE (Anmelde-Status / Session)
// WICHTIG: Wir trennen das KONTO (lieferino_user – bleibt dauerhaft gespeichert)
// von der SESSION (lieferino_session – sagt nur, ob man gerade eingeloggt ist).
// Beim Logout löschen wir NUR die Session, nicht das Konto.

function ladeSession() {
  if (!browser) return false;
  return !!localStorage.getItem('lieferino_session');
}

export const eingeloggt = writable(ladeSession());

// Meldet den Nutzer an (Session setzen).
export function login() {
  if (browser) localStorage.setItem('lieferino_session', new Date().toISOString());
  eingeloggt.set(true);
}

// Meldet den Nutzer ab – das Konto bleibt erhalten!
export function logout() {
  if (browser) localStorage.removeItem('lieferino_session');
  eingeloggt.set(false);
}

// Prüft, ob überhaupt ein Konto angelegt ist (für "Willkommen zurück"-Anzeige).
export function hatKonto() {
  if (!browser) return false;
  return !!localStorage.getItem('lieferino_user');
}
