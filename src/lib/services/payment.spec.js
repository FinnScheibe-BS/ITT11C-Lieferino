import { describe, test, expect } from 'vitest';
import { pruefeKartennummer, pruefeAblauf, pruefeCvv, kartenTyp } from './payment.js';

describe('Kreditkarten-Prüfung', () => {
  test('gültige Test-Kartennummer (Luhn)', () => {
    expect(pruefeKartennummer('4242 4242 4242 4242')).toBe(true);
  });

  test('ungültige Kartennummer (Luhn schlägt fehl)', () => {
    expect(pruefeKartennummer('4242 4242 4242 4241')).toBe(false);
  });

  test('Kartentyp wird erkannt', () => {
    expect(kartenTyp('4242424242424242')).toBe('Visa');
    expect(kartenTyp('5500005555555559')).toBe('Mastercard');
  });

  test('abgelaufenes Datum ist ungültig', () => {
    expect(pruefeAblauf('01/20')).toBe(false);
  });

  test('zukünftiges Datum ist gültig', () => {
    expect(pruefeAblauf('12/35')).toBe(true);
  });

  test('CVV mit 3 Ziffern ist gültig', () => {
    expect(pruefeCvv('123', 'Visa')).toBe(true);
    expect(pruefeCvv('12', 'Visa')).toBe(false);
  });
});
