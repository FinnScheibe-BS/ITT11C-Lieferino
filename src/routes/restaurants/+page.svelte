<script>
  import { aktiveRestaurants } from '$lib/stores/lieferanten.js';
  import { favoriten, toggleFavorit } from '$lib/stores/favoriten.js';
  import { bewertungen } from '$lib/stores/bewertungen.js';
  import { t } from '$lib/utils/i18n.js';
  import { istGeoeffnet } from '$lib/utils/oeffnung.js';

  let gewaehlterTyp = $state('alle');
  let sortierung = $state('standard');
  let suche = $state('');
  let nurFavoriten = $state(false);
  let nurVeg = $state(false);
  let nurGeoeffnet = $state(false);

  function anzeigeBewertung(r) {
    const reviews = $bewertungen[r.slug] || [];
    if (reviews.length === 0) return r.bewertung;
    return reviews.reduce((s, b) => s + b.sterne, 0) / reviews.length;
  }

  let gefilterteRestaurants = $derived(
    $aktiveRestaurants
      .slice()
      .filter((r) => gewaehlterTyp === 'alle' || r.typ === gewaehlterTyp)
      .filter((r) => r.name.toLowerCase().includes(suche.toLowerCase()))
      .filter((r) => !nurFavoriten || $favoriten.includes(r.slug))
      .filter((r) => !nurVeg || r.speisekarte.some((g) => g.veg))
      .filter((r) => !nurGeoeffnet || istGeoeffnet(r))
      .sort((a, b) => {
        if (sortierung === 'bewertung') return anzeigeBewertung(b) - anzeigeBewertung(a);
        if (sortierung === 'minbestell-auf') return a.minBestell - b.minBestell;
        return 0;
      })
  );

  let typen = $derived([...new Set($aktiveRestaurants.map((r) => r.typ))]);

  function herzKlick(event, slug) {
    event.preventDefault();
    event.stopPropagation();
    toggleFavorit(slug);
  }
</script>

<div class="seite">
  <div class="hero-box karte">
    <h1>{$t('rest.title')}</h1>
    <p>{$t('rest.subtitle')}</p>
  </div>

  <div class="filter-bar karte">
    <input type="search" placeholder={$t('common.search_placeholder')} bind:value={suche} class="suchfeld" />

    <select bind:value={gewaehlterTyp}>
      <option value="alle">{$t('common.all_cuisines')}</option>
      {#each typen as typ}
        <option value={typ}>{typ}</option>
      {/each}
    </select>

    <select bind:value={sortierung}>
      <option value="standard">{$t('rest.sort_standard')}</option>
      <option value="bewertung">{$t('rest.sort_rating')}</option>
      <option value="minbestell-auf">{$t('rest.sort_minorder')}</option>
    </select>

    <button class="fav-filter" class:aktiv={nurFavoriten} onclick={() => (nurFavoriten = !nurFavoriten)}>
      {nurFavoriten ? $t('rest.only_favs') : $t('rest.all')}
    </button>

    <button class="fav-filter veg" class:aktiv={nurVeg} onclick={() => (nurVeg = !nurVeg)}>
      {$t('rest.veg')}
    </button>

    <button class="fav-filter offen" class:aktiv={nurGeoeffnet} onclick={() => (nurGeoeffnet = !nurGeoeffnet)}>
      🟢 Jetzt geöffnet
    </button>
  </div>

  <p class="treffer">{$t('rest.found').replace('{n}', gefilterteRestaurants.length)}</p>

  <div class="grid">
    {#each gefilterteRestaurants as r}
      <a href="/restaurant/{r.slug}" class="restaurant-card karte">
        <span class="emoji-bild">{r.emoji}</span>

        <!-- Herz & Badge oben -->
        <button class="herz" onclick={(e) => herzKlick(e, r.slug)} aria-label="Favorit">
          {$favoriten.includes(r.slug) ? '❤️' : '🤍'}
        </button>
        <span class="rating-badge">⭐ {anzeigeBewertung(r).toFixed(1)}</span>

        <!-- Text unten -->
        <div class="card-info">
          <p>{r.beschreibung}</p>
          <h3>{r.name}</h3>
          <div class="card-meta">
            <span class="tag">{r.typ}</span>
            <span class="lieferzeit">⏱️ {r.lieferzeit}</span>
          </div>
        </div>
      </a>
    {/each}
  </div>

  {#if gefilterteRestaurants.length === 0}
    <p class="leer karte">{$t('rest.none')}</p>
  {/if}
</div>

<style>
  .seite {
    max-width: 1200px;
    margin: 0 auto;
    padding: 24px 20px 48px;
  }

  .hero-box {
    text-align: center;
    padding: 48px 24px;
    margin-bottom: 24px;
    background: linear-gradient(135deg, rgba(230, 168, 0, 0.15), rgba(184, 124, 0, 0.10));
  }

  .hero-box h1 {
    margin: 0 0 12px;
    font-size: clamp(1.8rem, 5vw, 2.5rem);
    font-weight: 700;
    letter-spacing: -0.03em;
  }

  .hero-box p {
    margin: 0;
    opacity: 0.85;
    font-size: 1.05rem;
  }

  .filter-bar {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    margin-bottom: 16px;
    padding: 16px;
    align-items: center;
  }

  .suchfeld {
    flex: 1;
    min-width: 200px;
  }

  .filter-bar select {
    min-width: 140px;
    cursor: pointer;
  }

  .fav-filter {
    padding: 9px 16px;
    border-radius: 10px;
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.18s ease;
    white-space: nowrap;
  }

  .fav-filter:hover {
    transform: scale(1.02);
  }

  .fav-filter.aktiv {
    border-color: rgba(230, 168, 0, 0.8) !important;
    box-shadow: 0 0 0 3px rgba(230, 168, 0, 0.15) !important;
  }

  .fav-filter.veg.aktiv {
    border-color: #4caf50 !important;
    box-shadow: 0 0 0 3px rgba(76, 175, 80, 0.2) !important;
  }

  .fav-filter.offen.aktiv {
    border-color: #2e9e4f !important;
    box-shadow: 0 0 0 3px rgba(46, 158, 79, 0.2) !important;
  }

  .treffer {
    color: rgba(245, 240, 232, 0.6);
    font-size: 0.85rem;
    margin: 0 0 20px;
    padding-left: 4px;
  }

  :global(html[data-theme='light']) .treffer {
    color: rgba(26, 18, 0, 0.5);
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 24px;
  }

  /* ─── Card Basis ──────────────────────────────────────────────────── */
  .restaurant-card {
    aspect-ratio: 4 / 3;
    display: block;
    text-decoration: none;
    color: inherit;
    position: relative;
    overflow: hidden;
    --textbox-h: 90px; /* zentrale Stellschraube für Höhe des Textbereichs */
  }

  /* ─── Emoji Bild (zentriert im Bereich ÜBER der Textbox) ─────────── */
  .emoji-bild {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: var(--textbox-h);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 6.5rem;
    filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.4));
    transition: transform 0.3s ease;
    z-index: 1;
  }

  .restaurant-card:hover .emoji-bild {
    transform: scale(1.08);
  }

  /* ─── Herz Button ─────────────────────────────────────────────────── */
  .herz {
    position: absolute;
    top: 16px;
    left: 16px;
    background: rgba(255, 255, 255, 0.15);
    border: none;
    border-radius: 50%;
    width: 42px;
    height: 42px;
    font-size: 1.3rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.18s ease;
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    z-index: 4;
    padding: 0;
  }

  .herz:hover {
    background: rgba(255, 255, 255, 0.25);
    transform: scale(1.12);
  }

  /* ─── Rating Badge ────────────────────────────────────────────────── */
  .rating-badge {
    position: absolute;
    top: 16px;
    right: 16px;
    background: rgba(0, 0, 0, 0.75);
    color: #ffd700;
    padding: 5px 12px;
    border-radius: 20px;
    font-size: 0.85rem;
    font-weight: 700;
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    z-index: 4;
    font-family: 'Geist Sans', -apple-system, sans-serif !important;
  }

  /* ─── Card Info Text ──────────────────────────────────────────────── */
  .card-info {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: var(--textbox-h);
  box-sizing: border-box;
  overflow: hidden;
  padding: 14px 18px 16px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  z-index: 3;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.65) 0%, rgba(0, 0, 0, 0) 100%);
  }

  .card-info h3 {
    margin: 0;
    font-size: 1.15rem;
    font-weight: 600;
    color: #fff !important;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.3;
  }

  .card-info p {
    margin: 0;
    color: rgba(255, 220, 100, 0.9) !important;
    font-size: 0.82rem;
    line-height: 1.2;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .card-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 8px;
    margin-top: 2px;
  }

  .tag {
    background: rgba(230, 168, 0, 0.2);
    color: #f9c932 !important;
    padding: 4px 12px;
    border-radius: 14px;
    font-size: 0.75rem;
    font-weight: 600;
    border: 1px solid rgba(230, 168, 0, 0.25);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 60%;
  }

  .lieferzeit {
    color: rgba(255, 255, 255, 0.8) !important;
    font-size: 0.78rem;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .leer {
    text-align: center;
    padding: 48px 24px;
    color: rgba(245, 240, 232, 0.6);
    font-size: 1.05rem;
  }

  :global(html[data-theme='light']) .leer {
    color: rgba(26, 18, 0, 0.55);
  }
</style>