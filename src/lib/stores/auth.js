import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// 🔑 AUTH-STORE (Anmelde-Status / Session + Sicherheit)
// WICHTIG: Wir trennen das KONTO (lieferino_user – bleibt dauerhaft gespeichert)
// von der SESSION (lieferino_session – sagt nur, ob man gerade eingeloggt ist).
// Beim Logout löschen wir NUR die Session, nicht das Konto.

const SESSION_KEY = 'lieferino_session';
const FEHL_KEY = 'lieferino_loginversuche';

// Wie lange eine Session ohne "angemeldet bleiben" gilt (in Minuten).
const TIMEOUT_MIN = 30;
// Nach so vielen Fehlversuchen wird der Login kurz gesperrt.
const MAX_VERSUCHE = 5;
const SPERRE_MIN = 1;

// Prüft, ob eine gültige (nicht abgelaufene) Session existiert.
export function istSessionGueltig() {
  if (!browser) return false;
  try {
    const s = JSON.parse(localStorage.getItem(SESSION_KEY) || 'null');
    if (!s) return false;
    // ablauf === null bedeutet "angemeldet bleiben" (läuft nicht ab).
    if (s.ablauf && Date.now() > s.ablauf) {
      localStorage.removeItem(SESSION_KEY);
      return false;
    }
    return true;
  } catch {
    return false;
  }
}

export const eingeloggt = writable(istSessionGueltig());

// Meldet den Nutzer an. Mit angemeldetBleiben=true läuft die Session nicht ab.
export function login(angemeldetBleiben = false) {
  const ablauf = angemeldetBleiben ? null : Date.now() + TIMEOUT_MIN * 60 * 1000;
  if (browser) {
    localStorage.setItem(SESSION_KEY, JSON.stringify({ ts: Date.now(), ablauf }));
  }
  resetFehlversuche(); // erfolgreicher Login -> Zähler zurücksetzen
  eingeloggt.set(true);
}

// Meldet den Nutzer ab – das Konto bleibt erhalten!
export function logout() {
  if (browser) {
    localStorage.removeItem(SESSION_KEY);
    localStorage.removeItem('lieferino_token'); // API-Token (JWT) auch entfernen
  }
  eingeloggt.set(false);
}

// Prüft, ob überhaupt ein Konto angelegt ist (für "Willkommen zurück"-Anzeige).
export function hatKonto() {
  if (!browser) return false;
  return !!localStorage.getItem('lieferino_user');
}

// ---------------------------------------------------------------------
// 🔒 Login-Sperre nach zu vielen Fehlversuchen (Brute-Force-Schutz)
// 🚨 BACKEND-HINWEIS: Diese Sperre läuft nur im Browser und ist leicht zu
// umgehen. Echter Brute-Force-Schutz (Rate-Limiting) MUSS serverseitig sein.
// ---------------------------------------------------------------------

// Liefert den aktuellen Status: gesperrt? wie lange noch? wie viele Versuche?
export function loginStatus() {
  if (!browser) return { gesperrt: false, restSek: 0, versuche: 0 };
  const d = JSON.parse(localStorage.getItem(FEHL_KEY) || 'null') || { versuche: 0, sperreBis: 0 };
  if (d.sperreBis && Date.now() < d.sperreBis) {
    return { gesperrt: true, restSek: Math.ceil((d.sperreBis - Date.now()) / 1000), versuche: d.versuche };
  }
  return { gesperrt: false, restSek: 0, versuche: d.versuche };
}

// Zählt einen Fehlversuch und sperrt ggf.
export function registriereFehlversuch() {
  if (!browser) return;
  const d = JSON.parse(localStorage.getItem(FEHL_KEY) || 'null') || { versuche: 0, sperreBis: 0 };
  d.versuche += 1;
  if (d.versuche >= MAX_VERSUCHE) {
    d.sperreBis = Date.now() + SPERRE_MIN * 60 * 1000;
    d.versuche = 0;
  }
  localStorage.setItem(FEHL_KEY, JSON.stringify(d));
}

export function resetFehlversuche() {
  if (browser) localStorage.removeItem(FEHL_KEY);
}
