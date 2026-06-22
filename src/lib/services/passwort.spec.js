import { describe, test, expect } from 'vitest';
import { pruefePasswortStaerke } from './passwort.js';

describe('Passwort-Stärke', () => {
  test('schwaches Passwort ist nicht sicher', () => {
    const r = pruefePasswortStaerke('abc');
    expect(r.istSicher).toBe(false);
  });

  test('starkes Passwort erfüllt alle Regeln', () => {
    const r = pruefePasswortStaerke('Abcdef1!');
    expect(r.istSicher).toBe(true);
    expect(r.score).toBe(5);
  });

  test('ohne Sonderzeichen reicht es trotzdem (score 4)', () => {
    const r = pruefePasswortStaerke('Abcdef12');
    expect(r.score).toBe(4);
    expect(r.istSicher).toBe(true);
  });
});
