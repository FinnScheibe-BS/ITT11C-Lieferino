// lib/stores/cart.js
import { writable, derived } from 'svelte/store';
import { holeWarenkorb, speichereWarenkorb } from '$lib/api/cartApi.js';
import { eingeloggt } from '$lib/stores/auth.js';

// Der reaktive Store für die Artikel
export const warenkorb = writable([]);

// Der reaktive Store für die Gesamtsumme
export const gesamtSumme = derived(warenkorb, ($warenkorb) => {
  return $warenkorb.reduce((summe, artikel) => summe + (artikel.preis * artikel.menge), 0);
});

// Sobald sich der Login-Status ändert, laden wir den Warenkorb vom Server
eingeloggt.subscribe(async (istEingeloggt) => {
  if (istEingeloggt) {
    const res = await holeWarenkorb();
    if (res.ok && res.daten) {
      warenkorb.set(res.daten);
    }
  } else {
    warenkorb.set([]); // Leeren bei Logout
  }
});

// Hilfsfunktion zum Speichern auf dem Server
async function syncMitServer(aktuellerWarenkorb) {
  await speichereWarenkorb(aktuellerWarenkorb);
}

// Menge ändern (+1 / -1)
export function aendereMenge(index, veraenderung) {
  warenkorb.update((aktuell) => {
    if (!aktuell[index]) return aktuell;
    
    const neueMenge = aktuell[index].menge + veraenderung;
    if (neueMenge <= 0) {
      // Artikel entfernen, wenn Menge auf 0 sinkt
      aktuell.splice(index, 1);
    } else {
      aktuell[index].menge = neueMenge;
    }
    
    syncMitServer(aktuell);
    return [...aktuell];
  });
}

// Artikel komplett löschen
export function entferneArtikel(index) {
  warenkorb.update((aktuell) => {
    aktuell.splice(index, 1);
    syncMitServer(aktuell);
    return [...aktuell];
  });
}

// Artikel hinzufügen (z. B. von einer Restaurant-Speisekarte aus)
export function fuegeHinzu(artikel) {
  warenkorb.update((aktuell) => {
    const existiert = aktuell.find(item => item.id === artikel.id);
    if (existiert) {
      existiert.menge += 1;
    } else {
      aktuell.push({ ...artikel, menge: 1 });
    }
    syncMitServer(aktuell);
    return [...aktuell];
  });
}