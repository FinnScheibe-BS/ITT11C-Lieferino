<script>
  import "../app.css";
  import 'geist-svelte/font/sans';
  import 'geist-svelte/font/mono';
  import { warenkorb } from '$lib/stores/cart.js';
  import { theme, themeWechseln } from '$lib/stores/theme.js';
  import { t, sprache, setzeSprache, SPRACHEN } from '$lib/i18n.js';
  import { eingeloggt, logout } from '$lib/stores/auth.js';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import Drachenlord from '$lib/Drachenlord.svelte';
  import FunOverlay from '$lib/FunOverlay.svelte';

  import jumpscareGif from '../sets/44b67d5e479d46f672031fb9ee0229cf.gif';
  import jumpscareSound from '../sets/myinstants.mp3';

  let anzahl = $state(0);
  let warenkorbSumme = $state(0);
  
  // Prüfen, ob wir auf der Homepage sind
  let istHomepage = $derived($page.route.id === '/' || $page.url.pathname === '/');

  const TICK_DAUER_MS = 1000;
  const CHANCE = 10000;

  let zeigeJumpscare = $state(false);
  let menuOffen = $state(false);
  let tickInterval;
  let audioElement;

  function tick() {
    const wuerfel = Math.floor(Math.random() * CHANCE);
    if (wuerfel === 0) loeseJumpscareAus();
  }

  function loeseJumpscareAus() {
    zeigeJumpscare = true;
    if (audioElement) {
      audioElement.currentTime = 0;
      audioElement.play().catch(() => {});
    }
    setTimeout(() => { zeigeJumpscare = false; }, 4000);
  }

  const KONAMI = ['ArrowUp','ArrowUp','ArrowDown','ArrowDown','ArrowLeft','ArrowRight','ArrowLeft','ArrowRight','b','a'];
  let konamiFortschritt = 0;
  function konamiTaste(e) {
    const taste = e.key.length === 1 ? e.key.toLowerCase() : e.key;
    if (taste === KONAMI[konamiFortschritt]) {
      konamiFortschritt += 1;
      if (konamiFortschritt === KONAMI.length) {
        konamiFortschritt = 0;
        loeseJumpscareAus();
      }
    } else {
      konamiFortschritt = taste === KONAMI[0] ? 1 : 0;
    }
  }

  function menuSchliessen() { menuOffen = false; }

  function ausloggen() {
    logout();
    menuSchliessen();
    goto('/');
  }

  function berechneSumme() {
    if ($warenkorb && Array.isArray($warenkorb)) {
      warenkorbSumme = $warenkorb.reduce((sum, item) => sum + (item.preis ?? 0), 0);
    } else {
      warenkorbSumme = 0;
    }
  }

  onMount(() => {
    const unsub = warenkorb.subscribe(v => {
      anzahl = v?.length ?? 0;
      berechneSumme();
    });
    tickInterval = setInterval(tick, TICK_DAUER_MS);
    window.addEventListener('keydown', konamiTaste);
    return () => {
      unsub();
      clearInterval(tickInterval);
      window.removeEventListener('keydown', konamiTaste);
    };
  });
</script>

<!-- Easter Eggs -->
<Drachenlord />
<FunOverlay />

<!-- Jumpscare Overlay -->
{#if zeigeJumpscare}
  <div class="jumpscare-overlay">
    <img src={jumpscareGif} alt="" class="jumpscare-gif" />
  </div>
{/if}
<audio bind:this={audioElement} src={jumpscareSound} preload="auto"></audio>

<!-- ░░░ NAVIGATION ░░░ -->

<!-- Backdrop zum Schließen -->
{#if menuOffen}
  <div class="nav-backdrop" onclick={menuSchliessen} aria-hidden="true"></div>
{/if}

<!-- Burger Button -->
<button
  class="nav-burger"
  class:menu-offen={menuOffen}
  onclick={() => (menuOffen = !menuOffen)}
  aria-label="Menü öffnen"
  aria-expanded={menuOffen}
>
  <span class="burger-bar"></span>
  <span class="burger-bar"></span>
  <span class="burger-bar"></span>
  {#if anzahl > 0}
    <span class="burger-badge">{anzahl}</span>
  {/if}
</button>

<!-- Side Drawer -->
<nav class="nav-drawer" class:drawer-offen={menuOffen} aria-label="Hauptnavigation">
  <div class="drawer-inner">
    <div class="drawer-logo">
      <span class="logo-ring">🍕</span>
      <span class="logo-name">Lieferino</span>
    </div>

    <div class="nav-links">
      <a href="/" onclick={menuSchliessen} class="nav-link">
        <span class="nav-icon">⌂</span>
        {$t('nav.home')}
      </a>
      <a href="/restaurants" onclick={menuSchliessen} class="nav-link">
        <span class="nav-icon">🍽</span>
        {$t('nav.restaurants')}
      </a>
      <a href="/cart" onclick={menuSchliessen} class="nav-link cart-link-item">
        <span class="nav-icon">🛒</span>
        {$t('nav.cart')}
        {#if anzahl > 0}
          <span class="nav-cart-badge">{anzahl}</span>
        {/if}
      </a>
      <a href="/bestellungen" onclick={menuSchliessen} class="nav-link">
        <span class="nav-icon">📦</span>
        {$t('nav.orders')}
      </a>
      <a href="/account" onclick={menuSchliessen} class="nav-link">
        <span class="nav-icon">👤</span>
        {$t('nav.account')}
      </a>
      {#if !$eingeloggt}
        <a href="/login" onclick={menuSchliessen} class="nav-link">
          <span class="nav-icon">🔑</span>
          {$t('nav.login')}
        </a>
      {:else}
        <button type="button" onclick={ausloggen} class="nav-link nav-logout">
          <span class="nav-icon">🚪</span>
          {$t('nav.logout')}
        </button>
      {/if}
      <a href="/admin" onclick={menuSchliessen} class="nav-link">
        <span class="nav-icon">⚙</span>
        {$t('nav.admin')}
      </a>
    </div>

    <!-- Sprache -->
    <div class="sprach-block">
      {#each SPRACHEN as s}
        {#if s.code === $sprache && s.icon}
          <img src={s.icon} alt="" class="sprach-icon" />
        {/if}
      {/each}
      <select
        value={$sprache}
        onchange={(e) => setzeSprache(e.target.value)}
        aria-label="Sprache wählen"
        class="sprach-select"
      >
        {#each SPRACHEN as s}
          <option value={s.code}>{s.flag} {s.label}</option>
        {/each}
      </select>
    </div>

    <!-- Impressum -->
    <div class="drawer-impressum">
      <p class="impressum-title">Impressum</p>
      <p>Lieferino GmbH<br />Musterstraße 12<br />12345 Stadt</p>
    </div>
  </div>
</nav>

<!-- Theme Toggle -->
<button
  class="theme-btn"
  onclick={themeWechseln}
  title="Hell/Dunkel"
  aria-label="Hell/Dunkel umschalten"
>
  {$theme === 'dark' ? '☀️' : '🌙'}
</button>

<!-- ░░░ HERO BANNER (nur auf Homepage) ░░░ -->
{#if istHomepage}
  <section class="hero-banner">
    <div class="hero-bg-effects">
      <div class="hero-glow hero-glow-1"></div>
      <div class="hero-glow hero-glow-2"></div>
      <div class="hero-glow hero-glow-3"></div>
      <div class="hero-grid"></div>
    </div>
    
    <div class="hero-content">
      <div class="hero-badge">
        <span class="badge-icon">✨</span>
        <span>Premium Food Delivery</span>
      </div>
      
      <h1 class="hero-title">
        <span class="title-line">Dein Geschmack.</span>
        <span class="title-line title-gold">Unsere Leidenschaft.</span>
      </h1>
      
      <p class="hero-subtitle">
        Entdecke exklusive Restaurants und lass dir 
        <span class="highlight">kulinarische Meisterwerke</span> direkt vor die Tür bringen.
      </p>
      
      <div class="hero-features">
        <div class="feature-item">
          <span class="feature-icon">🚀</span>
          <span>Schnelle Lieferung</span>
        </div>
        <div class="feature-item">
          <span class="feature-icon">⭐</span>
          <span>Premium Qualität</span>
        </div>
        <div class="feature-item">
          <span class="feature-icon">🔒</span>
          <span>Sichere Zahlung</span>
        </div>
      </div>
      
      <div class="hero-actions">
        <a href="/restaurants" class="hero-btn hero-btn-primary">
          <span>Restaurants entdecken</span>
          <span class="btn-arrow">→</span>
        </a>
        <a href="/bestellungen" class="hero-btn hero-btn-secondary">
          <span>Bestellungen ansehen</span>
        </a>
      </div>
      
      <div class="hero-stats">
        <div class="stat-item">
          <span class="stat-value">500+</span>
          <span class="stat-label">Restaurants</span>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-item">
          <span class="stat-value">30min</span>
          <span class="stat-label">Ø Lieferzeit</span>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-item">
          <span class="stat-value">10k+</span>
          <span class="stat-label">Kunden</span>
        </div>
      </div>

      <!-- ░░░ WARENKORB BOX ░░░ -->
      <a href="/cart" class="hero-cart-box">
        <div class="cart-box-icon">
          <span class="cart-icon">🛒</span>
          {#if anzahl > 0}
            <span class="cart-count">{anzahl}</span>
          {/if}
        </div>
        <div class="cart-box-info">
          <span class="cart-label">Dein Warenkorb</span>
          <span class="cart-total">{warenkorbSumme.toFixed(2)} €</span>
        </div>
        <span class="cart-arrow">→</span>
      </a>
    </div>
    
    <div class="hero-decorative">
      <div class="food-floating food-1">🍕</div>
      <div class="food-floating food-2">🍔</div>
      <div class="food-floating food-3">🍣</div>
      <div class="food-floating food-4">🥗</div>
      <div class="food-floating food-5">🍰</div>
    </div>

    <!-- Weicher Fade-Out am Ende -->
    <div class="hero-fade-out"></div>
    
    <div class="hero-scroll-indicator">
      <span>Scrollen</span>
      <div class="scroll-arrow">↓</div>
    </div>
  </section>
{:else}
  <!-- Einfacher Header auf allen anderen Seiten -->
  <section class="simple-header">
    <h1>Willkommen bei Lieferino</h1>
  </section>
{/if}

<!-- Seiteninhalt -->
<div class="page-content">
  <slot />
</div>

<style>
  /* ─── CSS-Tokens: Gold-Palette ─────────────────────────────────────── */
  :root {
    --gold-100: #fff8e7;
    --gold-200: #fde8a0;
    --gold-300: #f9c932;
    --gold-400: #e6a800;
    --gold-500: #b87c00;
    --gold-600: #7a5000;
    --blur-bg: rgba(255, 248, 220, 0.55);
    --blur-dark: rgba(30, 20, 5, 0.60);
    --card-radius: 18px;
    --tile-label-h: 72px;
    --nav-width: 260px;
    --transition: 0.32s cubic-bezier(0.4, 0, 0.2, 1);
  }

  /* ─── Globale Basis ─────────────────────────────────────────────────── */
  :global(*, _::before,_ ::after) { box-sizing: border-box; margin: 0; padding: 0; }

  :global(html) {
    background: #0d0d0d;
    color: #f5f0e8;
    font-family: 'Libre Caslon Text', Georgia, serif;
    -webkit-font-smoothing: antialiased;
  }

  :global(body) {
    background:
      radial-gradient(ellipse 80% 60% at 20% 10%, rgba(230,168,0,0.18) 0%, transparent 60%),
      radial-gradient(ellipse 60% 50% at 80% 80%, rgba(184,124,0,0.12) 0%, transparent 55%),
      #0d0d0d;
    min-height: 100vh;
  }

  /* ─── Hell-Modus ────────────────────────────────────────────────────── */
  :global(html[data-theme='light']),
  :global(html[data-theme='light'] body) {
    background:
      radial-gradient(ellipse 80% 60% at 20% 10%, rgba(230,168,0,0.10) 0%, transparent 60%),
      radial-gradient(ellipse 60% 50% at 80% 80%, rgba(184,124,0,0.07) 0%, transparent 55%),
      #faf6ee !important;
    color: #1a1200 !important;
  }

  /* ─── Apple-Kacheln (global, damit alle Seiten profitieren) ─────────── */

  :global(.kachel),
  :global(.restaurant-card),
  :global(.netflix-card) {
    position: relative;
    overflow: hidden;
    border-radius: var(--card-radius) !important;
    border: 1px solid rgba(230, 168, 0, 0.25) !important;
    box-shadow:
      0 0 0 0.5px rgba(230,168,0,0.15),
      0 8px 32px rgba(0,0,0,0.45) !important;
    background-color: #1a1200;
    transition: transform var(--transition), box-shadow var(--transition);
    cursor: pointer;
  }

  :global(.kachel:hover),
  :global(.restaurant-card:hover),
  :global(.netflix-card:hover) {
    transform: scale(1.025) translateY(-2px);
    box-shadow:
      0 0 0 1px rgba(230,168,0,0.5),
      0 16px 48px rgba(0,0,0,0.55),
      0 0 24px rgba(230,168,0,0.15) !important;
  }

  :global(.kachel img),
  :global(.restaurant-card img),
  :global(.netflix-card img),
  :global(.kachel .emoji-bild),
  :global(.restaurant-card .emoji-bild) {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 4rem;
    background: linear-gradient(135deg, #2a1f00, #1a1200) !important;
  }

  :global(.kachel::after),
  :global(.restaurant-card::after),
  :global(.netflix-card::after) {
    content: '';
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    height: var(--tile-label-h);
    backdrop-filter: blur(20px) saturate(1.6) brightness(0.85);
    -webkit-backdrop-filter: blur(20px) saturate(1.6) brightness(0.85);
    background: linear-gradient(
      to top,
      rgba(20, 12, 0, 0.82) 0%,
      rgba(30, 18, 0, 0.60) 55%,
      transparent 100%
    );
    border-top: 0.5px solid rgba(230, 168, 0, 0.18);
    border-radius: 0 0 var(--card-radius) var(--card-radius);
    z-index: 2;
    pointer-events: none;
  }

  :global(.kachel .kachel-info),
  :global(.restaurant-card .card-info),
  :global(.netflix-card .card-info),
  :global(.kachel h3),
  :global(.restaurant-card h3),
  :global(.netflix-card h3),
  :global(.kachel .name),
  :global(.restaurant-card .name) {
    position: absolute !important;
    bottom: 0 !important;
    left: 0 !important;
    right: 0 !important;
    height: var(--tile-label-h) !important;
    display: flex !important;
    flex-direction: column !important;
    justify-content: center !important;
    padding: 10px 14px !important;
    z-index: 3 !important;
    color: #fff !important;
    background: none !important;
    font-size: 0.95rem !important;
    font-weight: 600 !important;
    letter-spacing: -0.01em;
  }

  :global(.kachel .kachel-info span),
  :global(.restaurant-card .card-info span),
  :global(.netflix-card .card-info span),
  :global(.kachel p),
  :global(.restaurant-card p) {
    color: rgba(255, 220, 100, 0.85) !important;
    font-size: 0.78rem !important;
    font-weight: 400 !important;
    margin-top: 2px !important;
  }

  /* ─── Helle Karten / allgemeine Surfaces ────────────────────────────── */
  :global(.karte),
  :global(.block),
  :global(.login-box),
  :global(.review),
  :global(.review-form),
  :global(.bestellung),
  :global(.gericht),
  :global(.zeile) {
    background: rgba(255, 248, 220, 0.06) !important;
    border: 1px solid rgba(230, 168, 0, 0.18) !important;
    border-radius: var(--card-radius) !important;
    backdrop-filter: blur(16px) saturate(1.4);
    -webkit-backdrop-filter: blur(16px) saturate(1.4);
    color: #f5f0e8 !important;
    box-shadow: 0 4px 24px rgba(0,0,0,0.3) !important;
  }

  :global(html[data-theme='light'] .karte),
  :global(html[data-theme='light'] .block),
  :global(html[data-theme='light'] .login-box),
  :global(html[data-theme='light'] .review),
  :global(html[data-theme='light'] .bestellung),
  :global(html[data-theme='light'] .gericht),
  :global(html[data-theme='light'] .zeile) {
    background: rgba(255, 252, 235, 0.80) !important;
    border-color: rgba(230, 168, 0, 0.35) !important;
    color: #1a1200 !important;
  }

  /* ─── Buttons global ────────────────────────────────────────────────── */
  :global(button:not(.nav-burger):not(.theme-btn)),
  :global(.btn),
  :global(input[type='submit']) {
    background: linear-gradient(135deg, #e6a800, #b87c00) !important;
    color: #1a0f00 !important;
    border: none !important;
    border-radius: 12px !important;
    font-weight: 700 !important;
    font-size: 0.9rem !important;
    padding: 10px 20px !important;
    cursor: pointer !important;
    box-shadow: 0 2px 12px rgba(230,168,0,0.35) !important;
    transition: opacity 0.15s, transform 0.15s !important;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif !important;
  }

  :global(button:not(.nav-burger):not(.theme-btn):hover),
  :global(.btn:hover) {
    opacity: 0.88 !important;
    transform: scale(0.98) !important;
  }

  /* ─── Inputs global ─────────────────────────────────────────────────── */
  :global(input:not([type='checkbox']):not([type='radio'])),
  :global(select),
  :global(textarea) {
    background: rgba(255, 248, 220, 0.08) !important;
    border: 1px solid rgba(230, 168, 0, 0.30) !important;
    border-radius: 10px !important;
    color: #f5f0e8 !important;
    padding: 9px 14px !important;
    font-size: 0.9rem !important;
    outline: none !important;
    transition: border-color 0.18s !important;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif !important;
  }

  :global(input:focus),
  :global(select:focus),
  :global(textarea:focus) {
    border-color: rgba(230, 168, 0, 0.80) !important;
    box-shadow: 0 0 0 3px rgba(230,168,0,0.15) !important;
  }

  :global(html[data-theme='light'] input:not([type='checkbox']):not([type='radio'])),
  :global(html[data-theme='light'] select),
  :global(html[data-theme='light'] textarea) {
    background: rgba(255, 250, 235, 0.85) !important;
    color: #1a1200 !important;
  }

  :global(html[data-theme='light'] input::placeholder),
  :global(html[data-theme='light'] textarea::placeholder) {
    color: rgba(26, 18, 0, 0.45) !important;
  }

  /* ─── Überschriften ─────────────────────────────────────────────────── */
  :global(h1), :global(h2), :global(h3), :global(h4) {
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif !important;
  }
  :global(h1) { font-size: clamp(1.8rem, 5vw, 3rem); font-weight: 700; letter-spacing: -0.03em; color: #fff; }
  :global(h2) { font-size: clamp(1.3rem, 3vw, 1.9rem); font-weight: 600; letter-spacing: -0.02em; color: #f5d87c; }
  :global(h3) { font-size: 1.1rem; font-weight: 600; color: #f5d87c; }
  :global(h4) { font-size: 0.9rem; font-weight: 500; color: rgba(245, 216, 124, 0.7); }

  :global(html[data-theme='light'] h1) { color: #1a0f00; }
  :global(html[data-theme='light'] h2) { color: #7a5000; }
  :global(html[data-theme='light'] h3) { color: #7a5000; }
  :global(html[data-theme='light'] h4) { color: #7a5000; }

  :global(strong), :global(b), :global(.font-bold) {
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif !important;
    font-weight: 700;
  }

  /* ─── Links ─────────────────────────────────────────────────────────── */
  :global(a) { color: #f9c932; text-decoration: none; transition: color 0.15s; }
  :global(a:hover) { color: #fff; }
  :global(html[data-theme='light'] a) { color: #9a6600; }
  :global(html[data-theme='light'] a:hover) { color: #7a5000; }

  /* ─── Seiteninhalt ──────────────────────────────────────────────────── */
  .page-content {
    padding: 24px 20px 48px;
    max-width: 1200px;
    margin: 0 auto;
  }

  /* ─── JUMPSCARE ─────────────────────────────────────────────────────── */
  .jumpscare-overlay {
    position: fixed;
    inset: 0;
    background: #000;
    z-index: 999999;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .jumpscare-gif { width: 100%; height: 100%; object-fit: cover; }

  /* ─── Nav Backdrop ──────────────────────────────────────────────────── */
  .nav-backdrop {
    position: fixed;
    inset: 0;
    z-index: 9998;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    -webkit-backdrop-filter: blur(4px);
  }

  /* ─── Burger Button ─────────────────────────────────────────────────── */
  .nav-burger {
    position: fixed;
    top: 18px;
    right: 18px;
    z-index: 10001;
    width: 44px;
    height: 44px;
    border-radius: 14px !important;
    background: rgba(20, 12, 0, 0.70) !important;
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border: 1px solid rgba(230, 168, 0, 0.35) !important;
    box-shadow: 0 4px 16px rgba(0,0,0,0.4) !important;
    display: flex !important;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 5px;
    cursor: pointer;
    padding: 0 !important;
    transition: var(--transition) !important;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif !important;
  }

  .nav-burger:hover {
    border-color: rgba(230, 168, 0, 0.70) !important;
    box-shadow: 0 4px 20px rgba(230,168,0,0.2) !important;
    transform: none !important;
    opacity: 1 !important;
  }

  .burger-bar {
    display: block;
    width: 20px;
    height: 2px;
    background: #f9c932;
    border-radius: 2px;
    transition: var(--transition);
  }

  .menu-offen .burger-bar:nth-child(1) {
    transform: translateY(7px) rotate(45deg);
  }
  .menu-offen .burger-bar:nth-child(2) {
    opacity: 0;
    transform: scaleX(0);
  }
  .menu-offen .burger-bar:nth-child(3) {
    transform: translateY(-7px) rotate(-45deg);
  }

  .burger-badge {
    position: absolute;
    top: -4px;
    right: -4px;
    background: #e6a800;
    color: #1a0f00;
    font-size: 10px;
    font-weight: 800;
    padding: 2px 5px;
    border-radius: 8px;
    line-height: 1;
    font-family: 'Geist Sans', -apple-system, sans-serif !important;
  }

  /* ─── Drawer ────────────────────────────────────────────────────────── */
  .nav-drawer {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    width: var(--nav-width);
    z-index: 9999;
    transform: translateX(100%);
    transition: transform var(--transition);
    background: rgba(14, 9, 0, 0.82);
    backdrop-filter: blur(40px) saturate(1.8);
    -webkit-backdrop-filter: blur(40px) saturate(1.8);
    border-left: 1px solid rgba(230, 168, 0, 0.20);
    box-shadow: -12px 0 48px rgba(0, 0, 0, 0.60);
  }

  .drawer-offen {
    transform: translateX(0);
  }

  .drawer-inner {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: 80px 0 32px;
  }

  .drawer-logo {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 0 20px 28px;
    border-bottom: 1px solid rgba(230, 168, 0, 0.12);
    margin-bottom: 12px;
  }

  .logo-ring {
    width: 36px;
    height: 36px;
    border-radius: 10px;
    background: linear-gradient(135deg, #e6a800, #7a5000);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.2rem;
    box-shadow: 0 2px 8px rgba(230,168,0,0.35);
  }

  .logo-name {
    font-size: 1.1rem;
    font-weight: 700;
    color: #f9c932;
    letter-spacing: -0.02em;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif !important;
  }

  .nav-links {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow-y: auto;
    padding: 4px 10px;
  }

  .nav-link {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 14px;
    border-radius: 12px;
    color: rgba(245, 240, 232, 0.85) !important;
    font-size: 0.95rem;
    font-weight: 500;
    text-decoration: none;
    transition: background 0.18s, color 0.18s;
    margin-bottom: 2px;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif !important;
  }

  .nav-link:hover {
    background: rgba(230, 168, 0, 0.12);
    color: #f9c932 !important;
  }

  .nav-icon {
    font-size: 1.1rem;
    width: 22px;
    text-align: center;
    flex-shrink: 0;
  }

  .cart-link-item { position: relative; }

  .nav-cart-badge {
    margin-left: auto;
    background: linear-gradient(135deg, #e6a800, #b87c00);
    color: #1a0f00;
    font-size: 0.72rem;
    font-weight: 800;
    padding: 2px 8px;
    border-radius: 10px;
    line-height: 1.4;
    font-family: 'Geist Sans', -apple-system, sans-serif !important;
  }

  .sprach-block {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 14px 20px;
    border-top: 1px solid rgba(230, 168, 0, 0.10);
    margin-top: 8px;
  }

  .sprach-icon {
    width: 24px;
    height: 24px;
    object-fit: contain;
    border-radius: 4px;
    image-rendering: pixelated;
    flex-shrink: 0;
  }

  .sprach-select {
    flex: 1;
    background: rgba(255, 248, 220, 0.06) !important;
    border: 1px solid rgba(230, 168, 0, 0.25) !important;
    border-radius: 8px !important;
    color: #f5f0e8 !important;
    padding: 7px 10px !important;
    font-size: 0.85rem !important;
    cursor: pointer;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif !important;
  }

  .drawer-impressum {
    padding: 16px 20px;
    border-top: 1px solid rgba(230, 168, 0, 0.10);
    font-size: 0.75rem;
    color: rgba(245, 240, 232, 0.40);
    line-height: 1.5;
  }

  .impressum-title {
    font-weight: 600;
    color: rgba(245, 240, 232, 0.60);
    margin-bottom: 4px;
    font-size: 0.8rem;
  }

  /* ─── Theme Button ──────────────────────────────────────────────────── */
  .theme-btn {
    position: fixed;
    bottom: 22px;
    left: 22px;
    width: 48px;
    height: 48px;
    border-radius: 50% !important;
    border: 1px solid rgba(230, 168, 0, 0.35) !important;
    background: rgba(20, 12, 0, 0.70) !important;
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    color: #f9c932;
    font-size: 1.25rem;
    cursor: pointer;
    z-index: 10001;
    box-shadow: 0 4px 16px rgba(0,0,0,0.4), 0 0 0 0 rgba(230,168,0,0);
    display: flex !important;
    align-items: center;
    justify-content: center;
    transition: box-shadow 0.2s, border-color 0.2s, transform 0.15s !important;
  }

  .theme-btn:hover {
    border-color: rgba(230, 168, 0, 0.70) !important;
    box-shadow: 0 4px 20px rgba(0,0,0,0.5), 0 0 16px rgba(230,168,0,0.25) !important;
    transform: scale(1.06) !important;
    opacity: 1 !important;
  }

  /* ─── SGA Easter Egg ────────────────────────────────────────────────── */
  :global(html[data-sga='true']),
  :global(html[data-sga='true'] body),
  :global(html[data-sga='true'] body *) {
    font-family: 'Enchanting', sans-serif !important;
  }

  @font-face {
    font-family: 'Enchanting';
    src: url('/fonts/enchantment.ttf') format('truetype');
    font-display: swap;
  }

  /* ─── Scrollbar ─────────────────────────────────────────────────────── */
  :global(::-webkit-scrollbar) { width: 6px; height: 6px; }
  :global(::-webkit-scrollbar-track) { background: transparent; }
  :global(::-webkit-scrollbar-thumb) {
    background: rgba(230, 168, 0, 0.30);
    border-radius: 3px;
  }
  :global(::-webkit-scrollbar-thumb:hover) {
    background: rgba(230, 168, 0, 0.55);
  }

  /* ════════════════════════════════════════════════════════════════════
     ─── HERO BANNER STYLES ─────────────────────────────────────────────
     ════════════════════════════════════════════════════════════════════ */

  .hero-banner {
    position: relative;
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: visible;
    padding: 120px 24px 80px;
  }

  /* Hintergrund-Effekte */
  .hero-bg-effects {
    position: absolute;
    inset: 0;
    pointer-events: none;
    overflow: hidden;
  }

  .hero-glow {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.5;
    animation: glow-pulse 8s ease-in-out infinite;
  }

  .hero-glow-1 {
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, rgba(230,168,0,0.4) 0%, transparent 70%);
    top: -200px;
    left: -100px;
    animation-delay: 0s;
  }

  .hero-glow-2 {
    width: 500px;
    height: 500px;
    background: radial-gradient(circle, rgba(184,124,0,0.3) 0%, transparent 70%);
    bottom: -150px;
    right: -100px;
    animation-delay: 2s;
  }

  .hero-glow-3 {
    width: 400px;
    height: 400px;
    background: radial-gradient(circle, rgba(249,201,50,0.25) 0%, transparent 70%);
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    animation-delay: 4s;
  }

  .hero-grid {
    position: absolute;
    inset: 0;
    background-image: 
      linear-gradient(rgba(230,168,0,0.03) 1px, transparent 1px),
      linear-gradient(90deg, rgba(230,168,0,0.03) 1px, transparent 1px);
    background-size: 60px 60px;
    mask-image: radial-gradient(ellipse at center, black 0%, transparent 70%);
    -webkit-mask-image: radial-gradient(ellipse at center, black 0%, transparent 70%);
  }

  @keyframes glow-pulse {
    0%, 100% { opacity: 0.4; transform: scale(1); }
    50% { opacity: 0.6; transform: scale(1.1); }
  }

  /* Haupt-Inhalt */
  .hero-content {
    position: relative;
    z-index: 10;
    text-align: center;
    max-width: 900px;
    animation: fade-in-up 1s ease-out;
  }

  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(40px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  /* Badge */
  .hero-badge {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 8px 20px;
    background: rgba(230,168,0,0.12);
    border: 1px solid rgba(230,168,0,0.35);
    border-radius: 50px;
    font-size: 0.85rem;
    font-weight: 600;
    color: #f9c932;
    margin-bottom: 24px;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  :global(html[data-theme='light']) .hero-badge {
    background: rgba(230,168,0,0.18);
    color: #7a5000;
    border-color: rgba(122,80,0,0.35);
  }

  .badge-icon {
    animation: sparkle 2s ease-in-out infinite;
  }

  @keyframes sparkle {
    0%, 100% { transform: scale(1) rotate(0deg); }
    50% { transform: scale(1.2) rotate(10deg); }
  }

  /* Titel */
  .hero-title {
    font-size: clamp(2.5rem, 8vw, 5rem);
    font-weight: 800;
    line-height: 1.1;
    letter-spacing: -0.04em;
    margin-bottom: 20px;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  .title-line {
    display: block;
  }

  .title-gold {
    background: linear-gradient(135deg, #f9c932 0%, #e6a800 50%, #fde8a0 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  :global(html[data-theme='light']) .title-gold {
    background: linear-gradient(135deg, #7a5000 0%, #b87c00 50%, #e6a800 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  /* Untertitel */
  .hero-subtitle {
    font-size: clamp(1.1rem, 2.5vw, 1.4rem);
    color: rgba(245,240,232,0.85);
    line-height: 1.6;
    margin-bottom: 32px;
    max-width: 700px;
    margin-left: auto;
    margin-right: auto;
    font-family: 'Libre Caslon Text', Georgia, serif;
  }

  :global(html[data-theme='light']) .hero-subtitle {
    color: rgba(26,18,0,0.75);
  }

  .hero-subtitle .highlight {
    color: #f9c932;
    font-weight: 600;
  }

  :global(html[data-theme='light']) .hero-subtitle .highlight {
    color: #b87c00;
  }

  /* Features */
  .hero-features {
    display: flex;
    justify-content: center;
    gap: 24px;
    flex-wrap: wrap;
    margin-bottom: 40px;
  }

  .feature-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 20px;
    background: rgba(255,248,220,0.06);
    border: 1px solid rgba(230,168,0,0.25);
    border-radius: 50px;
    font-size: 0.9rem;
    font-weight: 500;
    color: #f5f0e8;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
    transition: all var(--transition);
  }

  :global(html[data-theme='light']) .feature-item {
    background: rgba(255,252,235,0.7);
    border-color: rgba(230,168,0,0.35);
    color: #1a1200;
  }

  .feature-item:hover {
    background: rgba(230,168,0,0.15);
    border-color: rgba(230,168,0,0.5);
    transform: translateY(-2px);
  }

  .feature-icon {
    font-size: 1.2rem;
  }

  /* Actions / Buttons */
  .hero-actions {
    display: flex;
    justify-content: center;
    gap: 16px;
    flex-wrap: wrap;
    margin-bottom: 48px;
  }

  .hero-btn {
    display: inline-flex;
    align-items: center;
    gap: 10px;
    padding: 16px 32px;
    border-radius: 14px;
    font-size: 1rem;
    font-weight: 700;
    text-decoration: none;
    transition: all var(--transition);
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
    cursor: pointer;
  }

  .hero-btn-primary {
    background: linear-gradient(135deg, #e6a800 0%, #b87c00 100%);
    color: #1a0f00;
    box-shadow: 0 4px 24px rgba(230,168,0,0.4);
  }

  .hero-btn-primary:hover {
    transform: translateY(-3px) scale(1.02);
    box-shadow: 0 8px 32px rgba(230,168,0,0.55);
  }

  .hero-btn-secondary {
    background: rgba(255,248,220,0.08);
    color: #f9c932;
    border: 1px solid rgba(230,168,0,0.35);
  }

  :global(html[data-theme='light']) .hero-btn-secondary {
    background: rgba(255,252,235,0.7);
    color: #7a5000;
    border-color: rgba(122,80,0,0.35);
  }

  .hero-btn-secondary:hover {
    background: rgba(230,168,0,0.15);
    border-color: rgba(230,168,0,0.5);
    transform: translateY(-3px);
  }

  .btn-arrow {
    transition: transform var(--transition);
  }

  .hero-btn-primary:hover .btn-arrow {
    transform: translateX(4px);
  }

  /* Stats */
  .hero-stats {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 0;
    padding: 24px 32px;
    background: rgba(255,248,220,0.04);
    border: 1px solid rgba(230,168,0,0.2);
    border-radius: 20px;
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    max-width: 500px;
    margin: 0 auto 32px;
  }

  :global(html[data-theme='light']) .hero-stats {
    background: rgba(255,252,235,0.6);
    border-color: rgba(230,168,0,0.3);
  }

  .stat-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 0 16px;
  }

  .stat-value {
    font-size: 1.5rem;
    font-weight: 800;
    color: #f9c932;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  :global(html[data-theme='light']) .stat-value {
    color: #7a5000;
  }

  .stat-label {
    font-size: 0.8rem;
    color: rgba(245,240,232,0.6);
    margin-top: 4px;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  :global(html[data-theme='light']) .stat-label {
    color: rgba(26,18,0,0.5);
  }

  .stat-divider {
    width: 1px;
    height: 40px;
    background: rgba(230,168,0,0.3);
  }

  /* ════════════════════════════════════════════════════════════════════
     ─── WARENKORB BOX ──────────────────────────────────────────────────
     ════════════════════════════════════════════════════════════════════ */

  .hero-cart-box {
    display: inline-flex;
    align-items: center;
    gap: 16px;
    padding: 14px 24px;
    background: rgba(255,248,220,0.08);
    border: 1px solid rgba(230,168,0,0.35);
    border-radius: 16px;
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
    text-decoration: none;
    transition: all var(--transition);
    animation: fade-in-up 1s ease-out 0.3s both;
    cursor: pointer;
    position: relative;
    overflow: hidden;
  }

  .hero-cart-box::before {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, rgba(230,168,0,0.1) 0%, transparent 50%, rgba(230,168,0,0.1) 100%);
    opacity: 0;
    transition: opacity var(--transition);
  }

  .hero-cart-box:hover::before {
    opacity: 1;
  }

  .hero-cart-box:hover {
    background: rgba(230,168,0,0.15);
    border-color: rgba(230,168,0,0.6);
    transform: translateY(-4px);
    box-shadow: 0 8px 32px rgba(230,168,0,0.25);
  }

  :global(html[data-theme='light']) .hero-cart-box {
    background: rgba(255,252,235,0.7);
    border-color: rgba(122,80,0,0.35);
  }

  :global(html[data-theme='light']) .hero-cart-box:hover {
    background: rgba(230,168,0,0.12);
    border-color: rgba(122,80,0,0.5);
  }

  .cart-box-icon {
    position: relative;
    width: 48px;
    height: 48px;
    background: linear-gradient(135deg, #e6a800, #b87c00);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 16px rgba(230,168,0,0.35);
    flex-shrink: 0;
  }

  .cart-icon {
    font-size: 1.4rem;
    filter: drop-shadow(0 2px 4px rgba(0,0,0,0.2));
  }

  .cart-count {
    position: absolute;
    top: -6px;
    right: -6px;
    background: #fff;
    color: #1a0f00;
    font-size: 0.7rem;
    font-weight: 800;
    padding: 2px 6px;
    border-radius: 8px;
    font-family: 'Geist Sans', -apple-system, sans-serif;
    box-shadow: 0 2px 8px rgba(0,0,0,0.2);
  }

  .cart-box-info {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    text-align: left;
    z-index: 1;
  }

  .cart-label {
    font-size: 0.75rem;
    color: rgba(245,240,232,0.6);
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  :global(html[data-theme='light']) .cart-label {
    color: rgba(26,18,0,0.5);
  }

  .cart-total {
    font-size: 1.25rem;
    font-weight: 700;
    color: #f9c932;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
    letter-spacing: -0.02em;
  }

  :global(html[data-theme='light']) .cart-total {
    color: #7a5000;
  }

  .cart-arrow {
    font-size: 1.2rem;
    color: #f9c932;
    opacity: 0;
    transform: translateX(-8px);
    transition: all var(--transition);
    margin-left: 8px;
  }

  :global(html[data-theme='light']) .cart-arrow {
    color: #b87c00;
  }

  .hero-cart-box:hover .cart-arrow {
    opacity: 1;
    transform: translateX(0);
  }

  /* Leerzustand wenn Warenkorb leer */
  .hero-cart-box.empty {
    opacity: 0.7;
  }

  .hero-cart-box.empty:hover {
    opacity: 1;
  }

  /* Dekorativ schwebende Elemente */
  .hero-decorative {
    position: absolute;
    inset: 0;
    pointer-events: none;
    overflow: hidden;
    z-index: 5;
  }

  .food-floating {
    position: absolute;
    font-size: 3rem;
    opacity: 0.15;
    animation: float 6s ease-in-out infinite;
    filter: drop-shadow(0 0 20px rgba(230,168,0,0.3));
  }

  :global(html[data-theme='light']) .food-floating {
    opacity: 0.08;
  }

  .food-1 { top: 15%; left: 10%; animation-delay: 0s; }
  .food-2 { top: 25%; right: 15%; animation-delay: 1s; }
  .food-3 { bottom: 30%; left: 12%; animation-delay: 2s; }
  .food-4 { bottom: 20%; right: 10%; animation-delay: 3s; }
  .food-5 { top: 50%; left: 5%; animation-delay: 4s; }

  @keyframes float {
    0%, 100% { transform: translateY(0px) rotate(0deg); }
    50% { transform: translateY(-20px) rotate(5deg); }
  }

  /* ════════════════════════════════════════════════════════════════════
     ─── FADE-OUT AM ENDE DES HERO ──────────────────────────────────────
     ════════════════════════════════════════════════════════════════════ */

  .hero-fade-out {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 180px;
    background: linear-gradient(
      to bottom,
      transparent 0%,
      transparent 20%,
      #0d0d0d 100%
    );
    z-index: 20;
    pointer-events: none;
  }

  :global(html[data-theme='light']) .hero-fade-out {
    background: linear-gradient(
      to bottom,
      transparent 0%,
      transparent 20%,
      #faf6ee 100%
    );
  }

  /* Scroll Indicator */
  .hero-scroll-indicator {
    position: absolute;
    bottom: 30px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    color: rgba(245,240,232,0.5);
    font-size: 0.8rem;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
    animation: bounce 2s ease-in-out infinite;
    z-index: 25;
  }

  :global(html[data-theme='light']) .hero-scroll-indicator {
    color: rgba(26,18,0,0.4);
  }

  .scroll-arrow {
    font-size: 1.2rem;
    color: #f9c932;
  }

  @keyframes bounce {
    0%, 100% { transform: translateX(-50%) translateY(0); }
    50% { transform: translateX(-50%) translateY(8px); }
  }

  /* ════════════════════════════════════════════════════════════════════
     ─── SIMPLE HEADER (für alle anderen Seiten) ────────────────────────
     ════════════════════════════════════════════════════════════════════ */

  .simple-header {
    padding: 100px 24px 40px;
    text-align: center;
    background:
      radial-gradient(ellipse 80% 60% at 20% 10%, rgba(230,168,0,0.12) 0%, transparent 60%),
      radial-gradient(ellipse 60% 50% at 80% 80%, rgba(184,124,0,0.08) 0%, transparent 55%);
  }

  .simple-header h1 {
    font-size: clamp(2rem, 6vw, 3.5rem);
    font-weight: 700;
    letter-spacing: -0.03em;
    background: linear-gradient(135deg, #f9c932 0%, #e6a800 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin-bottom: 16px;
  }

  :global(html[data-theme='light']) .simple-header {
    background:
      radial-gradient(ellipse 80% 60% at 20% 10%, rgba(230,168,0,0.08) 0%, transparent 60%),
      radial-gradient(ellipse 60% 50% at 80% 80%, rgba(184,124,0,0.05) 0%, transparent 55%);
  }

  :global(html[data-theme='light']) .simple-header h1 {
    background: linear-gradient(135deg, #7a5000 0%, #b87c00 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  /* Responsive Anpassungen */
  @media (max-width: 768px) {
    .hero-banner {
      padding: 100px 16px 60px;
    }

    .hero-features {
      gap: 12px;
    }

    .feature-item {
      padding: 10px 16px;
      font-size: 0.85rem;
    }

    .hero-actions {
      flex-direction: column;
      align-items: center;
    }

    .hero-btn {
      width: 100%;
      max-width: 280px;
      justify-content: center;
    }

    .hero-stats {
      flex-direction: column;
      gap: 16px;
      padding: 20px;
    }

    .stat-divider {
      width: 40px;
      height: 1px;
    }

    .food-floating {
      font-size: 2rem;
    }

    /* Warenkorb Box Mobile */
    .hero-cart-box {
      width: 100%;
      max-width: 320px;
      justify-content: center;
      padding: 12px 18px;
    }

    .cart-box-info {
      align-items: center;
      text-align: center;
    }

    .hero-fade-out {
      height: 120px;
    }

    .simple-header {
      padding: 80px 16px 30px;
    }
  }
</style>