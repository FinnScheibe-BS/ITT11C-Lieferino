import { describe, test, expect, vi, afterEach } from 'vitest';
import { generiereSecret, otpauthUri, generiereBackupCodes, pruefeTotp } from './mfa.js';

afterEach(() => vi.restoreAllMocks());

describe('MFA', () => {
  test('Secret hat 16 Base32-Zeichen', () => {
    expect(generiereSecret()).toMatch(/^[A-Z2-7]{16}$/);
  });

  test('otpauthUri enthält Secret und Issuer', () => {
    const uri = otpauthUri('ABC123', 'kunde@beispiel.de');
    expect(uri).toContain('secret=ABC123');
    expect(uri).toContain('issuer=Lieferino');
  });

  test('Backup-Codes haben das richtige Format', () => {
    const codes = generiereBackupCodes(3);
    expect(codes).toHaveLength(3);
    expect(codes[0]).toMatch(/^[0-9A-F]{4}-[0-9A-F]{4}$/);
  });

  test('TOTP stimmt mit RFC-6238-Testvektor überein', async () => {
    // Bei T=59s ist der Zähler 1 -> erwarteter 6-stelliger Code: 287082
    vi.spyOn(Date, 'now').mockReturnValue(59000);
    const ok = await pruefeTotp('GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ', '287082');
    expect(ok).toBe(true);
  });

  test('falscher TOTP-Code wird abgelehnt', async () => {
    vi.spyOn(Date, 'now').mockReturnValue(59000);
    const ok = await pruefeTotp('GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ', '000000');
    expect(ok).toBe(false);
  });
});
