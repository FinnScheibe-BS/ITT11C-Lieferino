<script>
  import { holeRestaurants } from '$lib/api/restaurantService.js';
  import { aktiveRestaurants } from '$lib/stores/lieferanten.js';
  import { bewertungen } from '$lib/stores/bewertungen.js';
  import { t } from '$lib/utils/i18n.js';
  import { onMount } from 'svelte';
  
  let gewaehlterTyp = $state('alle');
  let sortierung = $state('standard');
  let suche = $state('');
  
  // Holt die Daten live aus der API, sobald die Seite lädt
  onMount(async () => {
    const daten = await holeRestaurants();
    if (daten && daten.length > 0) {
      $aktiveRestaurants = daten;
    }
  });
  
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
      .sort((a, b) => {
        if (sortierung === 'bewertung') return anzeigeBewertung(b) - anzeigeBewertung(a);
        if (sortierung === 'minbestell-auf') return a.minBestell - b.minBestell;
        return 0;
      })
  );
  
  let typen = $derived([...new Set($aktiveRestaurants.map((r) => r.typ))]);
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
  </div>

  <p class="treffer">{$t('rest.found').replace('{n}', gefilterteRestaurants.length)}</p>

  <div class="grid">
    {#each gefilterteRestaurants as r}
      <a href="/restaurant/{r.slug}" class="restaurant-card">
        <div class="card-bild">
          <span class="emoji-bild">{r.emoji}</span>
        </div>

        <div class="card-info">
          <h3>{r.name}</h3>
          <div class="card-meta">
            <span class="tag">{r.typ}</span>
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
    align-items: stretch;
  }

  .restaurant-card {
    display: flex !important;
    flex-direction: column !important;
    align-items: stretch !important;
    justify-content: flex-start !important;
    text-decoration: none;
    color: inherit;
    overflow: hidden;
    border-radius: 16px;
    background: #1c1710;
    border: 1px solid rgba(230, 168, 0, 0.15);
  }

  .card-bild {
    position: relative;
    width: 100%;
    height: 222px;
    flex: 0 0 222px !important;
    display: flex !important;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    background: radial-gradient(circle at 50% 40%, rgba(230, 168, 0, 0.12), rgba(0, 0, 0, 0) 70%);
  }

  /* Blur reduziert von 40% auf 25% */
  .card-bild::after {
    content: '';
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    height: 30%;
    backdrop-filter: blur(6px);
    -webkit-backdrop-filter: blur(6px);
    background: linear-gradient(
      to bottom,
      rgba(28, 23, 16, 0) 0%,
      rgba(28, 23, 16, 0.5) 55%,
      rgba(28, 23, 16, 0.92) 100%
    );
    z-index: 1;
    pointer-events: none;
  }

  /* Trennlinie angepasst an neue Blur-Höhe */
  .card-bild::before {
    content: '';
    position: absolute;
    left: 16px;
    right: 16px;
    bottom: 30%;
    height: 1px;
    background: linear-gradient(
      to right,
      rgba(255, 255, 255, 0) 0%,
      rgba(255, 255, 255, 0.4) 50%,
      rgba(255, 255, 255, 0) 100%
    );
    box-shadow: 0 1px 4px rgba(255, 255, 255, 0.15);
    z-index: 2;
    pointer-events: none;
  }

  .emoji-bild {
    font-size: 6.2rem;
    line-height: 1;
    filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.4));
    transition: transform 0.3s ease;
  }

  .restaurant-card:hover .emoji-bild {
    transform: scale(1.08);
  }

  /* Textbereich höher positioniert (margin von -14px auf -6px reduziert) */
  .card-info {
    width: 100%;
    flex: 1 1 auto !important;
    min-height: 64px;
    box-sizing: border-box;
    padding: 6px 16px 18px;
    margin-top: -6px;
    display: flex !important;
    flex-direction: column !important;
    align-items: stretch !important;
    justify-content: flex-start !important;
    gap: 3px;
    position: relative;
    z-index: 3;
  }

  .card-info h3 {
    width: 100%;
    margin: 0 !important;
    padding: 0 !important;
    font-size: 1.1rem;
    font-weight: 600;
    color: #fff !important;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.3;
    min-height: 1.3em;
    visibility: visible !important;
    opacity: 1 !important;
    display: -webkit-box !important;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .card-meta {
    display: flex;
    justify-content: flex-start;
    align-items: center;
    gap: 6px;
    flex-shrink: 0;
  }

  .tag {
    display: inline-flex;
    align-items: center;
    box-sizing: border-box;
    background: rgba(230, 168, 0, 0.2);
    color: #f9c932 !important;
    padding: 3px 10px;
    line-height: 1.3;
    border-radius: 12px;
    font-size: 0.72rem;
    font-weight: 600;
    border: 1px solid rgba(230, 168, 0, 0.25);
    white-space: nowrap;
    flex-shrink: 0;
  }

  .lieferzeit {
    color: rgba(255, 255, 255, 0.8) !important;
    font-size: 0.75rem;
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