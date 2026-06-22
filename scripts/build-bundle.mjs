import { readFileSync, writeFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { dirname, join } from 'node:path';
import { restaurants } from '../src/lib/data/quelle.js';

const __root = join(dirname(fileURLToPath(import.meta.url)), '..');

function s0(bytes) {
  let h = 0x811c9dc5;
  for (let i = 0; i < bytes.length; i++) {
    h ^= bytes[i];
    h = Math.imul(h, 0x01000193);
  }
  return h >>> 0;
}

function s1(seed, len) {
  let x = seed || 0x9e3779b9;
  const out = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    x ^= x << 13; x >>>= 0;
    x ^= x >> 17;
    x ^= x << 5; x >>>= 0;
    out[i] = x & 0xff;
  }
  return out;
}

function s2(data, ks) {
  const out = new Uint8Array(data.length);
  for (let i = 0; i < data.length; i++) out[i] = data[i] ^ ks[i];
  return out;
}

const img = new Uint8Array(readFileSync(join(__root, 'rw-26-cd.gif')));
const payload = new TextEncoder().encode(JSON.stringify(restaurants));
const ks = s1(s0(img), payload.length);
const cipher = s2(payload, ks);
const b64 = Buffer.from(cipher).toString('base64');

writeFileSync(
  join(__root, 'src/lib/data/k.js'),
  `export const k = ${JSON.stringify(b64)};\n`
);

console.log(`ok (${cipher.length} bytes)`);
