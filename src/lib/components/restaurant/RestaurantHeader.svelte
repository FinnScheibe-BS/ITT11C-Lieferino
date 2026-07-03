<script>
  import { favoriten, toggleFavorit } from '$lib/stores/favoriten.js';

  let { restaurant, slug, durchschnitt, reviews, geoeffnet } = $props();
</script>

<section class="kopf karte">
  <span class="emoji">{restaurant.emoji}</span>

  <div class="kopf-inhalt">
    <h1>
      {restaurant.name}
      <button class="herz" onclick={() => toggleFavorit(slug)} aria-label="Favorit">
        {$favoriten.includes(slug) ? '❤️' : '🤍'}
      </button>
    </h1>

    <p class="meta">
      <a href="#bewertungen" class="bewertung-link">
        ⭐ {durchschnitt.toFixed(1)} ({reviews.length} Bewertungen)
      </a>
      <span>·</span>
      <span>⏱️ {restaurant.lieferzeit}</span>
      <span>·</span>
      <span>Min. {restaurant.minBestell}€</span>
    </p>

    <p class="oeffnung">
      {#if geoeffnet}
        <span class="status auf">🟢 Jetzt geöffnet</span>
      {:else}
        <span class="status zu">🔴 Geschlossen</span>
      {/if}
      <span class="zeiten">({restaurant.oeffnetUm}–{restaurant.schliesstUm} Uhr)</span>
    </p>

    <p class="beschreibung">{restaurant.beschreibung}</p>
  </div>
</section>