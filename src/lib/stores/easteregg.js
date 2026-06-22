import { writable } from 'svelte/store';

// 🥚 Steuert das versteckte „Drachenlord"-Easter-Egg.
// Wird true gesetzt, wenn jemand in der Suche „Drachenlord" eingibt.
export const drachenAktiv = writable(false);

export function drachenlordAusloesen() {
  drachenAktiv.set(true);
}
