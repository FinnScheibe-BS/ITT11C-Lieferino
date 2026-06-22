<script>
  // 🍽️ Restaurants kommen jetzt aus der zentralen Quelle (siehe $lib/data).
  import { restaurants } from '$lib/data';
  import { favoriten, toggleFavorit } from '$lib/stores/favoriten.js';

  // Variablen, die sich ändern können, bekommen ein $state().
  let gewaehlterTyp = $state('alle');
  let sortierung = $state('standard');
  let suche = $state(''); // 🔍 Freitext-Suche nach Name
  let nurFavoriten = $state(false); // ❤️ Nur Favoriten anzeigen?

  // $derived berechnet die gefilterte + sortierte Liste automatisch neu,
  // sobald sich Filter, Sortierung, Suchtext oder Favoriten ändern.
  let gefilterteRestaurants = $derived(
    restaurants
      .filter((r) => gewaehlterTyp === 'alle' || r.typ === gewaehlterTyp)
      // toLowerCase() macht die Suche groß-/kleinschreibungs-egal.
      .filter((r) => r.name.toLowerCase().includes(suche.toLowerCase()))
      // Favoriten-Filter: nur Restaurants, deren Slug in der Favoritenliste steht.
      .filter((r) => !nurFavoriten || $favoriten.includes(r.slug))
      .sort((a, b) => {
        if (sortierung === 'bewertung') return b.bewertung - a.bewertung;
        if (sortierung === 'minbestell-auf') return a.minBestell - b.minBestell;
        return 0;
      })
  );

  // Wir sammeln alle vorkommenden Küchen-Typen für das Dropdown (ohne Dopplungen).
  let typen = [...new Set(restaurants.map((r) => r.typ))];

  // Herz-Klick: verhindert, dass der Link dahinter ausgelöst wird.
  function herzKlick(event, slug) {
    event.preventDefault();
    event.stopPropagation();
    toggleFavorit(slug);
  }
</script>

<div class="seite">
  <div class="hero-box">
    <h1>🍔 Alle Restaurants</h1>
    <p>Stöbere durch unsere Auswahl und finde dein Lieblingsessen</p>
  </div>

  <!-- 🔍 Such- und Filterleiste -->
  <div class="filter-bar">
    <input type="search" placeholder="🔍 Restaurant suchen…" bind:value={suche} class="suchfeld" />

    <select bind:value={gewaehlterTyp}>
      <option value="alle">🌍 Alle Küchen</option>
      {#each typen as typ}
        <option value={typ}>{typ}</option>
      {/each}
    </select>

    <select bind:value={sortierung}>
      <option value="standard">Sortierung: Standard</option>
      <option value="bewertung">Beste Bewertung</option>
      <option value="minbestell-auf">Niedrigster Mindestbestellwert</option>
    </select>

    <!-- ❤️ Umschalter: nur Favoriten anzeigen -->
    <button class="fav-filter" class:aktiv={nurFavoriten} onclick={() => (nurFavoriten = !nurFavoriten)}>
      {nurFavoriten ? '❤️ Nur Favoriten' : '🤍 Alle'}
    </button>
  </div>

  <!-- Treffer-Anzahl -->
  <p class="treffer">{gefilterteRestaurants.length} Restaurant(s) gefunden</p>

  <div class="grid">
    {#each gefilterteRestaurants as r}
      <a href="/restaurant/{r.slug}" class="karte">
        <div class="emoji-bild">
          <span class="emoji-gross">{r.emoji}</span>
          <span class="rating-badge">⭐ {r.bewertung}</span>
          <!-- ❤️ Favoriten-Herz oben links -->
          <button class="herz" onclick={(e) => herzKlick(e, r.slug)} aria-label="Favorit">
            {$favoriten.includes(r.slug) ? '❤️' : '🤍'}
          </button>
        </div>
        <h3>{r.name}</h3>
        <p class="desc">{r.beschreibung}</p>
        <div class="footer">
          <span class="tag">{r.typ}</span>
          <span class="min">⏱️ {r.lieferzeit}</span>
        </div>
      </a>
    {/each}
  </div>

  <!-- Falls die Suche/der Filter nichts findet -->
  {#if gefilterteRestaurants.length === 0}
    <p class="leer">😕 Keine Restaurants gefunden. Versuche einen anderen Filter.</p>
  {/if}
</div>

<style>
  .seite { max-width: 1100px; margin: 0 auto; padding: 20px; font-family: sans-serif; }
  .hero-box { background: #673ab7; color: white; padding: 40px 20px; border-radius: 24px; text-align: center; margin-bottom: 30px; }
  .hero-box h1 { margin: 0 0 8px; }
  .hero-box p { margin: 0; opacity: 0.9; }

  .filter-bar { display: flex; flex-wrap: wrap; gap: 12px; margin-bottom: 12px; }
  .filter-bar .suchfeld { flex: 1; min-width: 200px; }
  .filter-bar input, .filter-bar select { padding: 11px; border: 1px solid #ddd; border-radius: 10px; font-size: 0.95rem; }
  .fav-filter { padding: 11px 16px; border: 1px solid #ddd; border-radius: 10px; background: white; cursor: pointer; font-size: 0.95rem; }
  .fav-filter.aktiv { border-color: #e0245e; background: #fff0f5; }
  .herz { position: absolute; top: 8px; left: 8px; background: rgba(255,255,255,0.85); border: none; border-radius: 50%; width: 34px; height: 34px; font-size: 1.1rem; cursor: pointer; }

  .treffer { color: #777; font-size: 0.9rem; margin: 0 0 16px; }

  .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 18px; }
  .karte { background: white; border: 1px solid #eee; border-radius: 16px; padding: 14px; text-decoration: none; color: inherit; box-shadow: 0 4px 12px rgba(0,0,0,0.05); transition: transform 0.15s ease; }
  .karte:hover { transform: translateY(-4px); }

  .emoji-bild { height: 130px; border-radius: 12px; display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #f3e5f5, #ede7f6); position: relative; }
  .emoji-gross { font-size: 3.5rem; }
  .rating-badge { position: absolute; top: 8px; right: 8px; background: rgba(0,0,0,0.7); color: white; padding: 3px 8px; border-radius: 20px; font-size: 0.8rem; }

  .karte h3 { margin: 12px 0 4px; }
  .desc { color: #777; font-size: 0.85rem; margin: 0 0 10px; }
  .footer { display: flex; justify-content: space-between; align-items: center; }
  .tag { background: #f3e5f5; color: #673ab7; padding: 3px 10px; border-radius: 20px; font-size: 0.78rem; }
  .min { color: #999; font-size: 0.8rem; }

  .leer { text-align: center; color: #777; margin-top: 30px; }
</style>
