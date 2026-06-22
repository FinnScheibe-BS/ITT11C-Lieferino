import { describe, test, expect, beforeEach } from 'vitest';
import { get } from 'svelte/store';
import { warenkorb, gesamtSumme, zumWarenkorb, aendereMenge, leereWarenkorb } from './cart.js';

describe('Warenkorb-Store', () => {
  beforeEach(() => leereWarenkorb());

  test('fügt ein Gericht mit Menge hinzu', () => {
    zumWarenkorb({ name: 'Pizza', preis: 8.5 }, 'Luigi', 2);
    const inhalt = get(warenkorb);
    expect(inhalt).toHaveLength(1);
    expect(inhalt[0].menge).toBe(2);
    expect(get(gesamtSumme)).toBe(17);
  });

  test('gleiches Gericht erhöht nur die Menge', () => {
    zumWarenkorb({ name: 'Pizza', preis: 8.5 }, 'Luigi');
    zumWarenkorb({ name: 'Pizza', preis: 8.5 }, 'Luigi');
    expect(get(warenkorb)).toHaveLength(1);
    expect(get(warenkorb)[0].menge).toBe(2);
  });

  test('Menge auf 0 entfernt den Artikel', () => {
    zumWarenkorb({ name: 'Pizza', preis: 8.5 }, 'Luigi');
    aendereMenge(0, -1);
    expect(get(warenkorb)).toHaveLength(0);
  });
});
