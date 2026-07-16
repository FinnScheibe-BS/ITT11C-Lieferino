<script>
  import { favoriten, toggleFavorit } from '$lib/stores/favoriten.js';
  import { bewertungen } from '$lib/stores/bewertungen.js';

  let { restaurant, slug, geoeffnet } = $props();

  let reviews = $derived($bewertungen[slug] || []);

  let durchschnitt = $derived.by(() => {
    if (reviews.length === 0) return restaurant?.bewertung ?? 0;

    const summe = reviews.reduce((gesamt, review) => gesamt + review.sterne, 0);
    return summe / reviews.length;
  });
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
        {durchschnitt.toFixed(1)} <span class="stern">★</span> ({reviews.length})
      </a>
      <span class="trenner">·</span>
      <span>{restaurant.lieferzeit}</span>
      <span class="trenner">·</span>
      <span>Min. {restaurant.minBestell}€</span>
    </p>

    <p class="oeffnung">
      <span class="status" class:auf={geoeffnet} class:zu={!geoeffnet}>
        <span class="dot"></span>
        {geoeffnet ? 'Jetzt geöffnet' : 'Geschlossen'}
      </span>
      <span class="zeiten">{restaurant.oeffnetUm}–{restaurant.schliesstUm} Uhr</span>
    </p>

    <p class="beschreibung">{restaurant.beschreibung}</p>
  </div>
</section>

<style>
  .kopf {
    display: flex;
    gap: 20px;
    align-items: flex-start;
    padding: 28px;
  }

  .emoji {
    font-size: 2.4rem;
    line-height: 1;
    flex-shrink: 0;
  }

  .kopf-inhalt {
    display: flex;
    flex-direction: column;
    gap: 10px;
    min-width: 0;
  }

  h1 {
    display: flex;
    align-items: center;
    gap: 12px;
    margin: 0;
    font-size: 1.7rem;
    line-height: 1.2;
  }

  .herz {
    background: none;
    border: none;
    font-size: 1.3rem;
    padding: 4px;
    cursor: pointer;
    line-height: 1;
  }

  .meta {
    display: flex;
    align-items: center;
    gap: 10px;
    margin: 0;
    color: rgba(245, 240, 232, 0.68);
    font-size: 0.92rem;
  }

  .bewertung-link {
    color: #f9c932;
    font-weight: 700;
    text-decoration: none;
  }

  .bewertung-link:hover {
    text-decoration: underline;
  }

  .stern {
    font-size: 0.85em;
  }

  .trenner {
    opacity: 0.4;
  }

  .oeffnung {
    display: flex;
    align-items: center;
    gap: 10px;
    margin: 2px 0 0;
    font-size: 0.88rem;
  }

  .status {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-weight: 600;
  }

  .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    display: inline-block;
  }

  .status.auf {
    color: #8ff09a;
  }
  .status.auf .dot {
    background: #4ade63;
    box-shadow: 0 0 6px rgba(74, 222, 99, 0.6);
  }

  .status.zu {
    color: #ff8f8f;
  }
  .status.zu .dot {
    background: #ff5c5c;
    box-shadow: 0 0 6px rgba(255, 92, 92, 0.5);
  }

  .zeiten {
    color: rgba(245, 240, 232, 0.45);
  }

  .beschreibung {
    margin: 6px 0 0;
    color: rgba(245, 240, 232, 0.6);
    font-size: 0.92rem;
    line-height: 1.5;
  }

  :global(html[data-theme='light']) .meta {
    color: rgba(26, 18, 0, 0.68);
  }
  :global(html[data-theme='light']) .zeiten {
    color: rgba(26, 18, 0, 0.45);
  }
  :global(html[data-theme='light']) .beschreibung {
    color: rgba(26, 18, 0, 0.62);
  }

  @media (max-width: 680px) {
    .kopf {
      padding: 20px;
      gap: 14px;
    }
    h1 {
      font-size: 1.4rem;
    }
  }
</style>