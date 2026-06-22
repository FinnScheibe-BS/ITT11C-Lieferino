import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

// 🛒 WARENKORB-STORE
// Der Warenkorb merkt sich Artikel inkl. Menge. Damit er einen Seiten-Neuladen
// übersteht, speichern wir ihn zusätzlich im localStorage des Browsers.

// Beim Start versuchen wir, einen gespeicherten Warenkorb zu laden.
function ladeStart() {
  if (!browser) return []; // Auf dem Server gibt es kein localStorage.
  try {
    const gespeichert = localStorage.getItem('lieferino_warenkorb');
    return gespeichert ? JSON.parse(gespeichert) : [];
  } catch {
    return [];
  }
}

export const warenkorb = writable(ladeStart());

// Jede Änderung am Warenkorb wird automatisch im localStorage gesichert.
if (browser) {
  warenkorb.subscribe((inhalt) => {
    localStorage.setItem('lieferino_warenkorb', JSON.stringify(inhalt));
  });
}

// 📦 Fügt ein Gericht hinzu (Standard-Menge 1, optional mehr auf einmal).
// Ist es schon drin, erhöhen wir nur die Menge.
export function zumWarenkorb(gericht, restaurantName, menge = 1) {
  warenkorb.update((inhalt) => {
    // Wir erkennen ein gleiches Gericht am Namen + Restaurant.
    const vorhanden = inhalt.find(
      (i) => i.name === gericht.name && i.restaurant === restaurantName
    );
    if (vorhanden) {
      vorhanden.menge += menge;
      return [...inhalt];
    }
    return [...inhalt, { ...gericht, menge, restaurant: restaurantName }];
  });
}

// ➕➖ Erhöht oder senkt die Menge. Bei Menge 0 fliegt der Artikel raus.
export function aendereMenge(index, delta) {
  warenkorb.update((inhalt) => {
    inhalt[index].menge += delta;
    if (inhalt[index].menge <= 0) {
      inhalt.splice(index, 1); // Artikel entfernen
    }
    return [...inhalt];
  });
}

// 🗑️ Entfernt einen Artikel komplett.
export function entferneArtikel(index) {
  warenkorb.update((inhalt) => {
    inhalt.splice(index, 1);
    return [...inhalt];
  });
}

// 🧹 Leert den ganzen Warenkorb (z.B. nach erfolgreicher Bestellung).
export function leereWarenkorb() {
  warenkorb.set([]);
}

// 💶 Gesamtsumme als abgeleiteter Store: rechnet sich automatisch neu.
export const gesamtSumme = derived(warenkorb, ($warenkorb) =>
  $warenkorb.reduce((summe, i) => summe + i.preis * i.menge, 0)
);
