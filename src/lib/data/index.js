import g from '../../../rw-26-cd.gif?inline';
import { k } from './k.js';

function b(s) {
  const bin = atob(s);
  const out = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
  return out;
}

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

const c = b(k);
const img = b(g.slice(g.indexOf(',') + 1));
const ks = s1(s0(img), c.length);
const o = new Uint8Array(c.length);
for (let i = 0; i < c.length; i++) o[i] = c[i] ^ ks[i];

export const restaurants = JSON.parse(new TextDecoder().decode(o));

export function getRestaurant(slugOderName) {
  return restaurants.find(
    (r) => r.slug === slugOderName || r.name === slugOderName
  );
}
