import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { pathToFileURL } from 'node:url';

const r = process.cwd();
const m = await import(pathToFileURL(join(r, 'src/lib/data/quelle.js')).href);

function a(b) {
  let h = 0x811c9dc5;
  for (let i = 0; i < b.length; i++) {
    h ^= b[i];
    h = Math.imul(h, 0x01000193);
  }
  return h >>> 0;
}

function c(s, n) {
  let x = s || 0x9e3779b9;
  const o = new Uint8Array(n);
  for (let i = 0; i < n; i++) {
    x ^= x << 13; x >>>= 0;
    x ^= x >> 17;
    x ^= x << 5; x >>>= 0;
    o[i] = x & 0xff;
  }
  return o;
}

function d(p, k) {
  const o = new Uint8Array(p.length);
  for (let i = 0; i < p.length; i++) o[i] = p[i] ^ k[i];
  return o;
}

const g = new Uint8Array(readFileSync(join(r, 'rw-26-cd.gif')));
const p = new TextEncoder().encode(JSON.stringify(m.restaurants));
const e = d(p, c(a(g), p.length));

writeFileSync(
  join(r, 'src/lib/data/k.js'),
  `export const k = ${JSON.stringify(Buffer.from(e).toString('base64'))};\n`
);

console.log(`ok (${e.length} bytes)`);
