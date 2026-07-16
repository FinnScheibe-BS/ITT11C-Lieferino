<script>
  import { holeRestaurants } from '$lib/api/restaurantService.js';
  import { aktiveRestaurants } from '$lib/stores/lieferanten.js';
  import { bewertungen } from '$lib/stores/bewertungen.js';
  import { t } from '$lib/utils/i18n.js';
  import { onMount } from 'svelte';
<<<<<<< HEAD
 
  let gewaehlterTyp = $state('alle');
  let sortierung = $state('standard');
  let suche = $state('');
 
=======

  let gewaehlterTyp = $state('alle');
  let sortierung = $state('standard');
  let suche = $state('');

>>>>>>> 407a6ebd4f3f43b182c45ab53f596dbe0b754e17
  // Holt die Daten live aus der API, sobald die Seite lädt
  onMount(async () => {
    const daten = await holeRestaurants();
    if (daten && daten.length > 0) {
<<<<<<< HEAD
      $aktiveRestaurants = daten;
=======
      $aktiveRestaurants = daten; 
>>>>>>> 407a6ebd4f3f43b182c45ab53f596dbe0b754e17
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
      <a href="/restaurant/{r.slug}" class="restaurant-card karte">
        <span class="emoji-bild">{r.emoji}</span>
<<<<<<< HEAD
 
=======

>>>>>>> 407a6ebd4f3f43b182c45ab53f596dbe0b754e17
        <!-- Text unten -->
        <div class="card-info">
          <p>{r.beschreibung}</p>
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
<<<<<<< HEAD
 
=======

>>>>>>> 407a6ebd4f3f43b182c45ab53f596dbe0b754e17
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
<<<<<<< HEAD
 
=======

>>>>>>> 407a6ebd4f3f43b182c45ab53f596dbe0b754e17
  .restaurant-card {
    aspect-ratio: 4 / 3;
    display: block;
    text-decoration: none;
    color: inherit;
    position: relative;
    overflow: hidden;
    --textbox-h: 100px; /* Etwas höher für besseren Textplatz */
  }
<<<<<<< HEAD
 
=======

>>>>>>> 407a6ebd4f3f43b182c45ab53f596dbe0b754e17
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
<<<<<<< HEAD
 
=======

>>>>>>> 407a6ebd4f3f43b182c45ab53f596dbe0b754e17
  .card-info {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: var(--textbox-h);
    box-sizing: border-box;
    padding: 12px 16px;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    gap: 6px;
    z-index: 3;
    background: linear-gradient(to top, rgba(0, 0, 0, 0.85) 0%, rgba(0, 0, 0, 0.6) 60%, rgba(0, 0, 0, 0) 100%);
  }
 
  .card-info h3 {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 600;
    color: #fff !important;
    white-space: normal; /* ✅ Umbrechen erlaubt */
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.3;
    display: -webkit-box;
    -webkit-line-clamp: 2; /* Max 2 Zeilen */
    -webkit-box-orient: vertical;
  }
 
  .card-info p {
    margin: 0;
    color: rgba(255, 220, 100, 0.85) !important;
    font-size: 0.78rem;
    line-height: 1.25;
    white-space: normal; /* ✅ Umbrechen erlaubt */
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2; /* Max 2 Zeilen */
    -webkit-box-orient: vertical;
  }
 
  .card-meta {
    display: flex;
    justify-content: flex-start;
    align-items: center;
    gap: 8px;
    margin-top: 4px;
  }
 
  .tag {
    background: rgba(230, 168, 0, 0.2);
    color: #f9c932 !important;
    padding: 4px 10px;
    border-radius: 12px;
    font-size: 0.72rem;
    font-weight: 600;
    border: 1px solid rgba(230, 168, 0, 0.25);
    white-space: nowrap;
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
 