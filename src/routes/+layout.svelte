<script>
  import '../app.css';
  import 'geist-svelte/font/sans';
  import 'geist-svelte/font/mono';

  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { onMount } from 'svelte';

  import { warenkorb } from '$lib/stores/cart.js';
  import { theme, themeWechseln } from '$lib/stores/theme.js';
  import { t, sprache, setzeSprache, SPRACHEN } from '$lib/utils/i18n.js';
  import { eingeloggt, logout } from '$lib/stores/auth.js';

  let anzahl = $state(0);
  let warenkorbSumme = $state(0);
  let menuOffen = $state(false);

  let istHomepage = $derived($page.route.id === '/' || $page.url.pathname === '/');

  function menuSchliessen() {
    menuOffen = false;
  }

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
    const unsub = warenkorb.subscribe((v) => {
      anzahl = v?.length ?? 0;
      berechneSumme();
    });

    return () => {
      unsub();
    };
  });
</script>

<!-- ░░░ NAVIGATION ░░░ -->

{#if menuOffen}
  <div class="nav-backdrop" onclick={menuSchliessen} aria-hidden="true"></div>
{/if}

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

<nav class="nav-drawer" class:drawer-offen={menuOffen} aria-label="Hauptnavigation">
  <div class="drawer-inner">
    <div class="drawer-logo">
      <span class="logo-ring">🍕</span>
      <span class="logo-name">Lieferino</span>
    </div>

    <div class="nav-links">
      <a href="/" onclick={menuSchliessen} class="nav-link">
        {$t('nav.home')}
      </a>
      <a href="/restaurants" onclick={menuSchliessen} class="nav-link">
        {$t('nav.restaurants')}
      </a>
      <a href="/cart" onclick={menuSchliessen} class="nav-link cart-link-item">
        {$t('nav.cart')}
        {#if anzahl > 0}
          <span class="nav-cart-badge">{anzahl}</span>
        {/if}
      </a>
      <a href="/bestellungen" onclick={menuSchliessen} class="nav-link">
        {$t('nav.orders')}
      </a>
      <a href="/account" onclick={menuSchliessen} class="nav-link">
        {$t('nav.account')}
      </a>
      {#if !$eingeloggt}
        <a href="/login" onclick={menuSchliessen} class="nav-link">
          {$t('nav.login')}
        </a>
      {:else}
        <button type="button" onclick={ausloggen} class="nav-link nav-logout">
          <span class="nav-icon">🚪</span>
          {$t('nav.logout')}
        </button>
      {/if}
      <a href="/admin" onclick={menuSchliessen} class="nav-link">
        {$t('nav.admin')}
      </a>
    </div>

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

    <div class="drawer-impressum">
      <p class="impressum-title">Impressum</p>
      <p>Lieferino GmbH<br />Musterstraße 12<br />12345 Stadt</p>
    </div>
  </div>
</nav>

<a href="/cart" class="cart-fab" title="Warenkorb" aria-label="Zum Warenkorb">
  🛒
  {#if anzahl > 0}
    <span class="cart-fab-badge">{anzahl}</span>
  {/if}
</a>

<div class="page-content">
  <slot />
</div>

<style>
  :root {
    --gold-100: #fff8e7;
    --gold-200: #fde8a0;
    --gold-300: #f9c932;
    --gold-400: #e6a800;
    --gold-500: #b87c00;
    --gold-600: #7a5000;
    --card-radius: 18px;
    --tile-label-h: 72px;
    --nav-width: 260px;
    --transition: 0.32s cubic-bezier(0.4, 0, 0.2, 1);
  }

  :global(*, *::before, *::after) { box-sizing: border-box; margin: 0; padding: 0; }

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

  :global(html[data-theme='light']),
  :global(html[data-theme='light'] body) {
    background:
      radial-gradient(ellipse 80% 60% at 20% 10%, rgba(230,168,0,0.10) 0%, transparent 60%),
      radial-gradient(ellipse 60% 50% at 80% 80%, rgba(184,124,0,0.07) 0%, transparent 55%),
      #faf6ee !important;
    color: #1a1200 !important;
  }

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

  :global(.kachel .kachel-info),
  :global(.restaurant-card .card-info),
  :global(.netflix-card .card-info) {
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
  }

  :global(.karte),
  :global(.block),
  :global(.login-box),
  :global(.review),
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
  :global(html[data-theme='light'] .login-box) {
    background: rgba(255, 252, 235, 0.80) !important;
    border-color: rgba(230, 168, 0, 0.35) !important;
    color: #1a1200 !important;
  }

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
    font-family: 'Geist Sans', -apple-system, sans-serif !important;
  }

  :global(button:not(.nav-burger):not(.theme-btn):hover),
  :global(.btn:hover) {
    opacity: 0.88 !important;
    transform: scale(0.98) !important;
  }

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
    font-family: 'Geist Sans', -apple-system, sans-serif !important;
  }

  :global(input:focus),
  :global(select:focus),
  :global(textarea:focus) {
    border-color: rgba(230, 168, 0, 0.80) !important;
    box-shadow: 0 0 0 3px rgba(230,168,0,0.15) !important;
  }

  :global(h1), :global(h2), :global(h3), :global(h4) {
    font-family: 'Geist Sans', -apple-system, sans-serif !important;
  }

  :global(a) { color: #f9c932; text-decoration: none; transition: color 0.15s; }
  :global(a:hover) { color: #fff; }

  .page-content {
    padding: 24px 20px 48px;
    max-width: 1200px;
    margin: 0 auto;
  }

  .nav-backdrop {
    position: fixed;
    inset: 0;
    z-index: 9998;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    -webkit-backdrop-filter: blur(4px);
  }

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
  }

  .nav-burger:hover {
    border-color: rgba(230, 168, 0, 0.70) !important;
    box-shadow: 0 4px 20px rgba(230,168,0,0.2) !important;
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
  }

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
    display: flex !important;
    align-items: center;
    justify-content: center;
    transition: box-shadow 0.2s, border-color 0.2s, transform 0.15s !important;
  }

  .theme-btn:hover {
    border-color: rgba(230, 168, 0, 0.70) !important;
    box-shadow: 0 4px 20px rgba(0,0,0,0.5), 0 0 16px rgba(230,168,0,0.25) !important;
    transform: scale(1.06) !important;
  }

  .cart-fab {
    position: fixed;
    bottom: 22px;
    right: 22px;
    width: 48px;
    height: 48px;
    border-radius: 50% !important;
    border: 1px solid rgba(230, 168, 0, 0.35) !important;
    background: rgba(20, 12, 0, 0.70) !important;
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    color: #f9c932 !important;
    font-size: 1.25rem;
    cursor: pointer;
    z-index: 10001;
    display: flex !important;
    align-items: center;
    justify-content: center;
    text-decoration: none !important;
    transition: box-shadow 0.2s, border-color 0.2s, transform 0.15s !important;
  }

  .cart-fab:hover {
    border-color: rgba(230, 168, 0, 0.70) !important;
    box-shadow: 0 4px 20px rgba(0,0,0,0.5), 0 0 16px rgba(230,168,0,0.25) !important;
    transform: scale(1.06) !important;
  }

  .cart-fab-badge {
    position: absolute;
    top: -4px;
    right: -4px;
    background: linear-gradient(135deg, #e6a800, #b87c00);
    color: #1a0f00;
    font-size: 11px;
    font-weight: 800;
    min-width: 18px;
    height: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 4px;
    border-radius: 10px;
    box-shadow: 0 2px 6px rgba(0,0,0,0.35);
  }

  :global(::-webkit-scrollbar) { width: 6px; height: 6px; }
  :global(::-webkit-scrollbar-track) { background: transparent; }
  :global(::-webkit-scrollbar-thumb) {
    background: rgba(230, 168, 0, 0.30);
    border-radius: 3px;
  }
  :global(::-webkit-scrollbar-thumb:hover) {
    background: rgba(230, 168, 0, 0.55);
  }
</style>