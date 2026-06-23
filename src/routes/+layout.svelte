<script>import "../app.css";
import { warenkorb } from '$lib/stores/cart.js';
import { theme, themeWechseln } from '$lib/stores/theme.js';
import { t, sprache, setzeSprache, SPRACHEN } from '$lib/i18n.js';
import { onMount, onDestroy } from 'svelte';
import Drachenlord from '$lib/Drachenlord.svelte';
import FunOverlay from '$lib/FunOverlay.svelte';

// Die Asset-Dateien liegen in src/sets. Vite verwandelt diese Imports
// automatisch in eine gültige URL, die im Browser geladen werden kann.
import jumpscareGif from '../sets/44b67d5e479d46f672031fb9ee0229cf.gif';
import jumpscareSound from '../sets/myinstants.mp3';

// Der '?' Operator sorgt dafür: Wenn warenkorb undefined ist,
// stürzt es nicht ab, sondern gibt 0 zurück.
let anzahl = $derived($warenkorb?.length ?? 0);

// =====================================================================
// 🎲 EASTER-EGG / "JUMPSCARE" 🎲
// Solange man auf der Seite ist, läuft ein "Tick". Bei JEDEM Tick gibt es
// eine 1-zu-10000-Chance, dass das GIF zusammen mit dem Sound abgespielt wird.
// Weil das hier im +layout.svelte steckt, läuft es auf JEDER Unterseite mit.
// =====================================================================

// Wie oft ein Tick passiert (in Millisekunden). 1000 = ein Tick pro Sekunde.
// Diesen Wert kann man kleiner machen, damit das Easter-Egg häufiger "würfelt".
const TICK_DAUER_MS = 1000;

// Die Chance: 1 von 10000 pro Tick.
const CHANCE = 10000;

let zeigeJumpscare = $state(false);   // Steuert das Anzeigen des Overlays
let tickInterval;                     // Merkt sich den Timer, damit wir ihn stoppen können
let audioElement;                     // Verweis auf das <audio>-Element im HTML

// Wird bei jedem Tick aufgerufen und "würfelt".
function tick() {
  // Math.random() gibt eine Zahl zwischen 0 und 1. Mal CHANCE und abgerundet
  // ergibt das eine Zufallszahl von 0 bis CHANCE-1. Genau eine davon (die 0)
  // löst den Jumpscare aus -> exakt 1-zu-10000-Chance.
  const wuerfel = Math.floor(Math.random() * CHANCE);
  if (wuerfel === 0) {
    loeseJumpscareAus();
  }
}

function loeseJumpscareAus() {
  zeigeJumpscare = true;

  // Sound von vorne starten und abspielen.
  if (audioElement) {
    audioElement.currentTime = 0;
    // .play() kann vom Browser blockiert werden, wenn der Nutzer die Seite
    // noch nie angeklickt hat. Mit .catch() verhindern wir einen Absturz.
    audioElement.play().catch(() => {});
  }

  // Nach 4 Sekunden blenden wir das Overlay automatisch wieder aus.
  setTimeout(() => {
    zeigeJumpscare = false;
  }, 4000);
}

// 🎮 KONAMI-CODE: ↑ ↑ ↓ ↓ ← → ← → B A  → löst den Jumpscare aus.
const KONAMI = [
  'ArrowUp', 'ArrowUp', 'ArrowDown', 'ArrowDown',
  'ArrowLeft', 'ArrowRight', 'ArrowLeft', 'ArrowRight', 'b', 'a'
];
let konamiFortschritt = 0;
function konamiTaste(e) {
  const taste = e.key.length === 1 ? e.key.toLowerCase() : e.key;
  if (taste === KONAMI[konamiFortschritt]) {
    konamiFortschritt += 1;
    if (konamiFortschritt === KONAMI.length) {
      konamiFortschritt = 0;
      loeseJumpscareAus(); // 😈 Überraschung!
    }
  } else {
    // Bei falscher Taste von vorne (außer es ist der erste Schritt).
    konamiFortschritt = taste === KONAMI[0] ? 1 : 0;
  }
}

onMount(() => {
  // Timer starten, sobald die Seite geladen ist.
  tickInterval = setInterval(tick, TICK_DAUER_MS);
  window.addEventListener('keydown', konamiTaste);
});

onDestroy(() => {
  // Timer wieder aufräumen, damit er nicht im Hintergrund weiterläuft.
  clearInterval(tickInterval);
  if (typeof window !== 'undefined') window.removeEventListener('keydown', konamiTaste);
});</script>

<!-- 🐉 Verstecktes „Drachenlord"-Easter-Egg (global) -->
<Drachenlord></Drachenlord>
<!-- 🎉 Fun-Modus (Emoji-Cursor + Saison-Effekt) -->
<FunOverlay></FunOverlay>

<!-- 🎲 Overlay des Easter-Eggs. Liegt über allem (z-index sehr hoch). -->
{#if zeigeJumpscare}
  <div class="jumpscare-overlay">
    <img src="{jumpscareGif}" alt class="jumpscare-gif" />
  </div>
{/if}
<!-- Audio-Element ist immer im DOM, wird aber nur bei Bedarf abgespielt. -->
<audio bind:this="{audioElement}" src="{jumpscareSound}" preload="auto"></audio>

<div class="nav-container">
  <input type="checkbox" id="menu-toggle" class="menu-checkbox" />
  <label for="menu-toggle" class="nav-burger-btn">
    ☰
    {#if anzahl > 0}
      <span class="cart-badge-burger">{anzahl}</span>
    {/if}
  </label>

  <div class="nav-dropdown-balken">
    <div class="button-umrundung"></div>

    <div class="nav-links-wrapper">
      <a href="/">{$t('nav.home')}</a>
      <a href="/restaurants">{$t('nav.restaurants')}</a>
      <a href="/cart" class="cart-link">
        {$t('nav.cart')}
        {#if anzahl > 0}
          <span class="cart-badge">{anzahl}</span>
        {/if}
      </a>
      <a href="/bestellungen">{$t('nav.orders')}</a>
      <a href="/account">{$t('nav.account')}</a>
      <a href="/login">{$t('nav.login')}</a>
      <a href="/admin">{$t('nav.admin')}</a>

      <!-- 🌍 Sprach-Umschalter (mit Icon der aktuellen Sprache davor) -->
      <div class="sprach-wahl">
        {#each SPRACHEN as s}
          {#if s.code === $sprache && s.icon}
            <img src="{s.icon}" alt class="sprach-icon" />
          {/if}
        {/each}
        <select value="{$sprache}" onchange="{(e)" => setzeSprache(e.target.value)} aria-label="Sprache wählen">
          {#each SPRACHEN as s}
            <option value="{s.code}">{s.flag} {s.label}</option>
          {/each}
        </select>
      </div>

      <div class="nav-impressum">
        <h4>Impressum</h4>
        <p>Lieferino GmbH<br />Musterstraße 12<br />12345 Stadt</p>
      </div>
    </div>
  </div>
</div>

<!-- 🌙/☀️ Dark-/Light-Mode-Schalter: fest unten links auf JEDER Seite -->
<button class="theme-switch" onclick="{themeWechseln}" title="Hell/Dunkel umschalten" aria-label="Hell/Dunkel umschalten">
  {$theme === 'dark' ? '☀️' : '🌙'}
</button>

<div class="page-content">
  <slot></slot>
</div>

<style>
  /* 🌍 Sprach-Umschalter im Menü */
  .sprach-wahl { margin-top: 10px; display: flex; align-items: center; gap: 8px; }
  .sprach-wahl select { flex: 1; padding: 8px; border-radius: 8px; border: 1px solid #ccc; font-size: 0.9rem; cursor: pointer; }
  .sprach-icon { width: 26px; height: 26px; object-fit: contain; border-radius: 4px; image-rendering: pixelated; }

  /* 🟪 Minecraft-Verzauberungstisch (Standard Galactic Alphabet)
     Die Schrift-Datei muss unter static/fonts/enchantment.ttf liegen.
     Fehlt sie, fällt der Text einfach auf normale (englische) Schrift zurück. */
  @font-face {
    font-family: 'Enchanting';
    src: url('/fonts/enchantment.ttf') format('truetype');
    font-display: swap;
  }
  /* Auf ALLE Elemente anwenden + !important, damit die seiteneigenen
     "font-family: sans-serif"-Regeln überschrieben werden. Zahlen/Emojis,
     die die Schrift nicht hat, fallen automatisch auf sans-serif zurück. */
  :global(html[data-sga='true']),
  :global(html[data-sga='true'] body),
  :global(html[data-sga='true'] body *) {
    font-family: 'Enchanting', sans-serif !important;
  }
  /* Gezielt für die Navigation (überschreibt deren eigenes !important). */
  :global(html[data-sga='true'] .nav-links-wrapper a),
  :global(html[data-sga='true'] .nav-burger-btn),
  :global(html[data-sga='true'] .nav-impressum) {
    font-family: 'Enchanting', sans-serif !important;
  }

  /* 🌙/☀️ Theme-Schalter: runder Button unten links */
  .theme-switch {
    position: fixed;
    bottom: 20px;
    left: 20px;
    width: 50px;
    height: 50px;
    border-radius: 50%;
    border: none;
    background: #673ab7;
    color: white;
    font-size: 1.4rem;
    cursor: pointer;
    z-index: 100001;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  /* =====================================================================
     🌙 DARK-MODE
     Sobald am <html> data-theme="dark" steht (das macht der Theme-Store),
     ziehen sich diese Farben durch ALLE Seiten. :global(...) ist nötig, weil
     die Elemente außerhalb dieser Layout-Komponente liegen. !important sorgt
     dafür, dass wir die seiteneigenen (hellen) Styles sicher übersteuern.
     ===================================================================== */

  /* Grundfläche: Hintergrund + Standard-Textfarbe */
  :global(html[data-theme='dark']),
  :global(html[data-theme='dark'] body) {
    background: #16161b !important;
    color: #e8e8ea !important;
  }
  :global(html[data-theme='dark'] .page-content) {
    background: #16161b;
  }

  /* Helle Karten/Boxen/Flächen abdunkeln (Klassen aus allen Seiten) */
  :global(html[data-theme='dark'] .login-box),
  :global(html[data-theme='dark'] .karte),
  :global(html[data-theme='dark'] .restaurant-card),
  :global(html[data-theme='dark'] .netflix-card),
  :global(html[data-theme='dark'] .zeile),
  :global(html[data-theme='dark'] .block),
  :global(html[data-theme='dark'] .gericht),
  :global(html[data-theme='dark'] .review),
  :global(html[data-theme='dark'] .review-form),
  :global(html[data-theme='dark'] .bestellung),
  :global(html[data-theme='dark'] .tabelle),
  :global(html[data-theme='dark'] .fav-filter) {
    background: #24242c !important;
    border-color: #3a3a44 !important;
    color: #e8e8ea !important;
  }

  /* Leicht abgesetzte Flächen (Zusammenfassung, Kopfbereiche, Tabellenkopf) */
  :global(html[data-theme='dark'] .zusammenfassung),
  :global(html[data-theme='dark'] .kopf),
  :global(html[data-theme='dark'] .zeile.kopf),
  :global(html[data-theme='dark'] .filter-bar) {
    background: #1e1e25 !important;
    border-color: #3a3a44 !important;
  }

  /* Überschriften + Standardtext hell halten */
  :global(html[data-theme='dark'] h1),
  :global(html[data-theme='dark'] h2),
  :global(html[data-theme='dark'] h3),
  :global(html[data-theme='dark'] h4),
  :global(html[data-theme='dark'] .seite),
  :global(html[data-theme='dark'] .welcome-container),
  :global(html[data-theme='dark'] .gesamt),
  :global(html[data-theme='dark'] .zwischensumme),
  :global(html[data-theme='dark'] .rechnung-zeile.gesamt) {
    color: #f1f1f4 !important;
  }

  /* Eingabefelder + Buttons mit hellem Hintergrund */
  :global(html[data-theme='dark'] input),
  :global(html[data-theme='dark'] select),
  :global(html[data-theme='dark'] textarea) {
    background: #2e2e38 !important;
    color: #e8e8ea !important;
    border-color: #44444f !important;
  }

  /* Helle "Chips"/Tags lesbar machen */
  :global(html[data-theme='dark'] .tag),
  :global(html[data-theme='dark'] .chip) {
    background: #3a2f4d !important;
    color: #d3bdf0 !important;
  }

  /* Emoji-Platzhalter etwas dunkler tönen */
  :global(html[data-theme='dark'] .emoji-bild) {
    background: linear-gradient(135deg, #2a2333, #1f1b29) !important;
  }

  /* 🎲 EASTER-EGG OVERLAY: deckt den kompletten Bildschirm ab */
  .jumpscare-overlay {
    position: fixed;
    inset: 0; /* Kurzform für top/right/bottom/left = 0 */
    background: #000;
    z-index: 999999; /* Über wirklich allem, auch über der Navigation */
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .jumpscare-gif {
    width: 100%;
    height: 100%;
    object-fit: cover; /* Füllt den ganzen Bildschirm aus */
  }

  .nav-container {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    z-index: 99999;
  }

  .menu-checkbox {
    display: none !important;
  }

  .nav-burger-btn {
    position: fixed !important;
    top: 20px;
    right: 4px;
    width: 40px;
    height: 36px;
    z-index: 100002;
    background: transparent !important;
    border: none !important;
    cursor: pointer;
    font-size: 26px;
    color: white !important;
    display: flex !important;
    align-items: center;
    justify-content: center;
    user-select: none;
    transition: right 0.3s ease;
  }

  .cart-badge-burger {
    position: absolute;
    top: -2px;
    right: -2px;
    background: #ff3b30;
    color: white;
    font-size: 11px;
    font-weight: 700;
    line-height: 1;
    padding: 2px 5px;
    border-radius: 10px;
    font-family: sans-serif;
  }

  .nav-dropdown-balken {
    position: fixed !important;
    top: 0;
    bottom: 0;
    right: -215px;
    width: 240px;
    height: 100vh;
    background: #673ab7 !important;
    z-index: 100000;
    transition: right 0.3s cubic-bezier(0.25, 1, 0.5, 1), box-shadow 0.3s ease;
    box-shadow: none;
  }

  .button-umrundung {
    position: absolute !important;
    top: 14px;
    left: -18px;
    width: 35px;
    height: 48px;
    background: #673ab7 !important;
    border-radius: 12px 0 0 12px;
    z-index: 100001;
    transition: opacity 0.2s ease;
  }

  .menu-checkbox:checked ~ .nav-dropdown-balken {
    right: 0 !important;
    box-shadow: -8px 0 25px rgba(0, 0, 0, 0.4) !important;
  }

  .menu-checkbox:checked ~ .nav-dropdown-balken .button-umrundung {
    opacity: 0 !important;
  }

  .menu-checkbox:checked ~ .nav-burger-btn {
    right: 185px !important;
  }

  .nav-links-wrapper {
    display: flex !important;
    flex-direction: column !important;
    height: 100%;
    padding-top: 90px;
    box-sizing: border-box;
    position: relative;
    z-index: 100003;
  }

  .nav-links-wrapper a {
    color: white !important;
    text-decoration: none !important;
    padding: 15px 25px !important;
    font-size: 18px !important;
    font-family: sans-serif !important;
    font-weight: 700 !important;
    display: block !important;
    transition: background 0.2s;
    white-space: nowrap;
  }

  .nav-links-wrapper a:hover {
    background: rgba(255, 255, 255, 0.15) !important;
  }

  .cart-link {
    display: flex !important;
    align-items: center;
    gap: 8px;
  }

  .cart-badge {
    background: white;
    color: #673ab7;
    font-size: 0.75rem;
    font-weight: 800;
    padding: 2px 8px;
    border-radius: 12px;
    line-height: 1.3;
  }

  .nav-impressum {
    margin-top: auto;
    padding: 25px;
    color: rgba(255, 255, 255, 0.7);
    font-family: sans-serif;
    font-size: 0.8rem;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
  }

  .nav-impressum h4 {
    margin: 0 0 5px 0;
    color: white;
  }

  .nav-impressum p {
    margin: 0;
    line-height: 1.3;
  }

  .page-content {
    padding: 20px;
  }
</style>
