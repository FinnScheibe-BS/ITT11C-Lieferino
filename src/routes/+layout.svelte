<script>
  import { warenkorb } from '$lib/stores/cart.js';
  import { theme, themeWechseln } from '$lib/stores/theme.js';
  import { onMount, onDestroy } from 'svelte';

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

  onMount(() => {
    // Timer starten, sobald die Seite geladen ist.
    tickInterval = setInterval(tick, TICK_DAUER_MS);
  });

  onDestroy(() => {
    // Timer wieder aufräumen, damit er nicht im Hintergrund weiterläuft.
    clearInterval(tickInterval);
  });
</script>

<!-- 🎲 Overlay des Easter-Eggs. Liegt über allem (z-index sehr hoch). -->
{#if zeigeJumpscare}
  <div class="jumpscare-overlay">
    <img src={jumpscareGif} alt="" class="jumpscare-gif" />
  </div>
{/if}
<!-- Audio-Element ist immer im DOM, wird aber nur bei Bedarf abgespielt. -->
<audio bind:this={audioElement} src={jumpscareSound} preload="auto"></audio>

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
      <a href="/">🏠 Home</a>
      <a href="/restaurants">🍔 Restaurants</a>
      <a href="/cart" class="cart-link">
        🛒 Warenkorb
        {#if anzahl > 0}
          <span class="cart-badge">{anzahl}</span>
        {/if}
      </a>
      <a href="/account">👤 Account</a>

      <div class="nav-impressum">
        <h4>Impressum</h4>
        <p>Lieferino GmbH<br>Musterstraße 12<br>12345 Stadt</p>
      </div>
    </div>
  </div>
</div>

<!-- 🌙/☀️ Dark-/Light-Mode-Schalter: fest unten links auf JEDER Seite -->
<button class="theme-switch" onclick={themeWechseln} title="Hell/Dunkel umschalten" aria-label="Hell/Dunkel umschalten">
  {$theme === 'dark' ? '☀️' : '🌙'}
</button>

<div class="page-content">
  <slot />
</div>

<style>
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
     Wir setzen die Farben global, sobald am <html> data-theme="dark" steht
     (das macht der Theme-Store). :global(...) ist nötig, weil diese Elemente
     außerhalb dieser Komponente liegen.
     ===================================================================== */
  :global(html[data-theme='dark']) {
    background: #1a1a1f;
  }
  :global(html[data-theme='dark'] body) {
    background: #1a1a1f;
    color: #e8e8ea;
  }
  /* Helle Karten/Boxen im Dark-Mode abdunkeln */
  :global(html[data-theme='dark'] .login-box),
  :global(html[data-theme='dark'] .karte),
  :global(html[data-theme='dark'] .zeile),
  :global(html[data-theme='dark'] .block),
  :global(html[data-theme='dark'] .gericht) {
    background: #26262e !important;
    border-color: #3a3a44 !important;
    color: #e8e8ea;
  }
  :global(html[data-theme='dark'] .zusammenfassung) {
    background: #222229 !important;
    border-color: #3a3a44 !important;
  }
  /* Eingabefelder im Dark-Mode */
  :global(html[data-theme='dark'] input),
  :global(html[data-theme='dark'] select) {
    background: #2e2e38 !important;
    color: #e8e8ea !important;
    border-color: #44444f !important;
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