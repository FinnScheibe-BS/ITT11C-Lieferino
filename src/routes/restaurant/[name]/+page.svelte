<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { zumWarenkorb } from '$lib/stores/cart.js';
  import { favoriten, toggleFavorit } from '$lib/stores/favoriten.js';
  import { bewertungen, bewertungHinzufuegen } from '$lib/stores/bewertungen.js';
  import { deaktivierteLieferanten, findeRestaurant } from '$lib/stores/lieferanten.js';

  // Den Slug aus der URL holen und das Restaurant suchen (inkl. Geheim-Restaurant).
  // Deaktivierte Lieferanten behandeln wir für normale Nutzer wie "nicht gefunden".
  let slug = $derived($page.params.name);
  let restaurant = $derived($deaktivierteLieferanten.includes(slug) ? undefined : findeRestaurant(slug));

  // Frühere Bestellungen laden, um zu prüfen, ob man hier schon bestellt hat.
  let meineBestellungen = $state([]);
  onMount(() => {
    meineBestellungen = JSON.parse(localStorage.getItem('lieferino_bestellungen') || '[]');
  });

  // Darf nur bewerten, wer hier nachweislich schon bestellt hat.
  let hatHierBestellt = $derived(
    meineBestellungen.some((b) =>
      (b.artikel || []).some((a) => a.restaurant === restaurant?.name)
    )
  );

  // Reviews dieses Restaurants aus dem Store ziehen.
  let reviews = $derived($bewertungen[slug] || []);

  // Durchschnittsbewertung aus den Kundenreviews (sonst Basis-Bewertung).
  let durchschnitt = $derived.by(() => {
    if (reviews.length === 0) return restaurant?.bewertung ?? 0;
    const summe = reviews.reduce((s, r) => s + r.sterne, 0);
    return summe / reviews.length;
  });

  // Ist das Restaurant gerade geöffnet? Vergleich der aktuellen Uhrzeit
  // mit oeffnetUm/schliesstUm (Format "HH:MM").
  let geoeffnet = $derived.by(() => {
    if (!restaurant?.oeffnetUm) return true;
    const jetzt = new Date();
    const minutenJetzt = jetzt.getHours() * 60 + jetzt.getMinutes();
    const [oh, om] = restaurant.oeffnetUm.split(':').map(Number);
    const [sh, sm] = restaurant.schliesstUm.split(':').map(Number);
    return minutenJetzt >= oh * 60 + om && minutenJetzt < sh * 60 + sm;
  });

  // Mengen-Auswahl pro Gericht (Standard 1).
  let mengen = $state({});
  function menge(id) {
    return mengen[id] ?? 1;
  }

  function aendere(id, delta) {
    const neu = Math.max(1, menge(id) + delta);
    mengen = { ...mengen, [id]: neu };
  }

  // Eingabefelder für eine neue Bewertung.
  let neuName = $state('');
  let neuSterne = $state(5);
  let neuText = $state('');

  function bewertungAbschicken(e) {
    e.preventDefault();
    if (!hatHierBestellt) return;
    if (neuName.trim() === '' || neuText.trim() === '') return;

    bewertungHinzufuegen(slug, { name: neuName, sterne: neuSterne, text: neuText });

    neuName = '';
    neuText = '';
    neuSterne = 5;
  }

  // Kleine Bestätigung, die kurz eingeblendet wird.
  let hinweis = $state('');
  let hinweisTimer;

  function hinzufuegen(gericht) {
    const anzahl = menge(gericht.id);
    zumWarenkorb(gericht, restaurant.name, anzahl);
    hinweis = `${anzahl}× ${gericht.name} hinzugefügt ✅`;
    mengen = { ...mengen, [gericht.id]: 1 };

    clearTimeout(hinweisTimer);
    hinweisTimer = setTimeout(() => (hinweis = ''), 2000);
  }
</script>

{#if restaurant}
  <div class="seite">
    <a href="/restaurants" class="zurueck">⬅️ Zu allen Restaurants</a>

    <!-- Kopfbereich des Restaurants -->
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
          <a href="#bewertungen" class="bewertung-link">⭐ {durchschnitt.toFixed(1)} ({reviews.length} Bewertungen)</a>
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

    <h2>📋 Speisekarte</h2>
    <div class="speisekarte">
      {#each restaurant.speisekarte as gericht}
        <article class="gericht">
          <div class="gericht-info">
            <h3>
              {gericht.name}
              {#if gericht.veg}<span class="veg-tag">🌱 vegetarisch</span>{/if}
            </h3>

            <p class="gericht-desc">{gericht.beschreibung}</p>

            {#if gericht.allergene && gericht.allergene.length > 0}
              <p class="allergene">Enthält: {gericht.allergene.join(', ')}</p>
            {/if}

            <span class="preis">{gericht.preis.toFixed(2)}€</span>
          </div>

          <div class="gericht-aktion">
            <div class="stepper" aria-label="Menge auswählen">
              <button onclick={() => aendere(gericht.id, -1)} aria-label="Weniger">−</button>
              <span>{menge(gericht.id)}</span>
              <button onclick={() => aendere(gericht.id, 1)} aria-label="Mehr">+</button>
            </div>
            <button class="add-btn" onclick={() => hinzufuegen(gericht)}>+ Hinzufügen</button>
          </div>
        </article>
      {/each}
    </div>

    <h2 id="bewertungen">⭐ Bewertungen ({reviews.length})</h2>

    {#if hatHierBestellt}
      <form class="review-form" onsubmit={bewertungAbschicken}>
        <input type="text" placeholder="Dein Name" bind:value={neuName} required />
        <select bind:value={neuSterne}>
          <option value={5}>⭐⭐⭐⭐⭐ (5)</option>
          <option value={4}>⭐⭐⭐⭐ (4)</option>
          <option value={3}>⭐⭐⭐ (3)</option>
          <option value={2}>⭐⭐ (2)</option>
          <option value={1}>⭐ (1)</option>
        </select>
        <textarea placeholder="Wie war dein Essen?" bind:value={neuText} required></textarea>
        <button type="submit">Bewertung abschicken</button>
      </form>
    {:else}
      <p class="review-sperre">🔒 Du kannst dieses Restaurant bewerten, sobald du hier etwas bestellt hast.</p>
    {/if}

    {#if reviews.length === 0}
      <p class="keine-reviews">Noch keine Bewertungen – sei die/der Erste! 🌟</p>
    {:else}
      <div class="reviews">
        {#each reviews as review}
          <article class="review">
            <div class="review-kopf">
              <strong>{review.name}</strong>
              <span>{'⭐'.repeat(review.sterne)}</span>
            </div>
            <p>{review.text}</p>
          </article>
        {/each}
      </div>
    {/if}
  </div>

  {#if hinweis}
    <div class="toast">{hinweis}</div>
  {/if}
{:else}
  <div class="seite">
    <div class="leerzustand karte">
      <h1>😕 Restaurant nicht gefunden</h1>
      <a href="/restaurants" class="zurueck">⬅️ Zu allen Restaurants</a>
    </div>
  </div>
{/if}

<style>
  .seite {
    max-width: 920px;
    margin: 0 auto;
    padding: 28px 20px 64px;
  }

  .zurueck {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 18px;
    color: #f9c932;
    font-weight: 700;
    text-decoration: none;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  .zurueck:hover {
    color: #fff;
  }

  .kopf {
    display: flex;
    gap: 22px;
    align-items: center;
    padding: 28px;
    margin: 0 0 34px;
    position: relative;
    overflow: hidden;
  }

  .kopf::before {
    content: '';
    position: absolute;
    inset: -40%;
    background:
      radial-gradient(circle at 20% 20%, rgba(230, 168, 0, 0.22), transparent 34%),
      radial-gradient(circle at 80% 80%, rgba(249, 201, 50, 0.10), transparent 32%);
    pointer-events: none;
  }

  .kopf .emoji {
    position: relative;
    z-index: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 104px;
    height: 104px;
    border-radius: 24px;
    background: linear-gradient(135deg, rgba(230, 168, 0, 0.22), rgba(122, 80, 0, 0.18));
    border: 1px solid rgba(230, 168, 0, 0.28);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.28);
    font-size: 4.2rem;
    flex-shrink: 0;
  }

  .kopf-inhalt {
    position: relative;
    z-index: 1;
    min-width: 0;
  }

  .kopf h1 {
    display: flex;
    align-items: center;
    gap: 10px;
    margin: 0;
    line-height: 1.05;
  }

  .meta {
    display: flex;
    flex-wrap: wrap;
    gap: 7px;
    align-items: center;
    color: rgba(245, 240, 232, 0.72);
    font-size: 0.94rem;
    margin: 10px 0 6px;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  .bewertung-link {
    color: #f9c932;
    text-decoration: underline;
    text-decoration-style: dotted;
    text-underline-offset: 4px;
  }

  .oeffnung {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;
    margin: 6px 0 12px;
    font-size: 0.92rem;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  .status {
    font-weight: 800;
  }

  .status.auf {
    color: #63d471;
  }

  .status.zu {
    color: #ff6b6b;
  }

  .zeiten {
    color: rgba(245, 240, 232, 0.52);
    font-size: 0.84rem;
  }

  .beschreibung {
    max-width: 64ch;
    color: rgba(245, 240, 232, 0.78);
    margin: 0;
    line-height: 1.55;
  }

  h2 {
    margin: 30px 0 14px;
  }

  .speisekarte,
  .reviews {
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  .gericht {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 18px;
    padding: 18px;
  }

  .gericht-info {
    min-width: 0;
  }

  .gericht-info h3 {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;
    margin: 0 0 6px;
  }

  .gericht-desc {
    color: rgba(245, 240, 232, 0.68);
    font-size: 0.9rem;
    line-height: 1.45;
    margin: 0 0 8px;
  }

  .allergene {
    color: rgba(245, 240, 232, 0.42);
    font-size: 0.78rem;
    margin: 0 0 8px;
  }

  .preis {
    display: inline-flex;
    color: #f9c932;
    font-weight: 800;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  .veg-tag {
    display: inline-flex;
    align-items: center;
    border-radius: 999px;
    padding: 3px 9px;
    background: rgba(99, 212, 113, 0.12);
    border: 1px solid rgba(99, 212, 113, 0.28);
    color: #8ff09a;
    font-size: 0.72rem;
    font-weight: 700;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  .gericht-aktion {
    display: flex;
    flex-direction: column;
    gap: 9px;
    align-items: flex-end;
    flex-shrink: 0;
  }

  .stepper {
    display: flex;
    align-items: center;
    gap: 9px;
    padding: 5px;
    border-radius: 999px;
    background: rgba(255, 248, 220, 0.06);
    border: 1px solid rgba(230, 168, 0, 0.18);
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  .stepper button {
    width: 28px;
    height: 28px;
    min-width: 28px;
    padding: 0 !important;
    border-radius: 50% !important;
    font-size: 1rem !important;
    line-height: 1 !important;
  }

  .stepper span {
    min-width: 18px;
    color: #f5f0e8;
    text-align: center;
    font-weight: 800;
  }

  .add-btn {
    white-space: nowrap;
  }

  .herz {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    padding: 0 !important;
    border-radius: 50% !important;
    background: rgba(255, 248, 220, 0.08) !important;
    border: 1px solid rgba(230, 168, 0, 0.22) !important;
    box-shadow: none !important;
    color: inherit !important;
    font-size: 1.2rem !important;
  }

  .herz:hover {
    transform: scale(1.06) !important;
    opacity: 1 !important;
  }

  .review-form {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 18px;
  }

  .review-form textarea {
    min-height: 92px;
    resize: vertical;
  }

  .review-sperre,
  .keine-reviews {
    padding: 16px;
    border-radius: 16px;
    border: 1px solid rgba(230, 168, 0, 0.18);
    background: rgba(255, 248, 220, 0.05);
    color: rgba(245, 240, 232, 0.68);
  }

  .review {
    padding: 16px;
  }

  .review-kopf {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 7px;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  .review p {
    margin: 0;
    color: rgba(245, 240, 232, 0.72);
    line-height: 1.5;
  }

  .toast {
    position: fixed;
    left: 50%;
    bottom: 30px;
    transform: translateX(-50%);
    z-index: 99998;
    padding: 14px 24px;
    border-radius: 999px;
    background: linear-gradient(135deg, #e6a800, #b87c00);
    color: #1a0f00;
    font-weight: 800;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
    box-shadow: 0 10px 34px rgba(0, 0, 0, 0.38), 0 0 18px rgba(230, 168, 0, 0.25);
  }

  .leerzustand {
    padding: 32px;
    text-align: center;
  }

  .leerzustand .zurueck {
    margin: 18px 0 0;
  }

  :global(html[data-theme='light']) .meta,
  :global(html[data-theme='light']) .beschreibung,
  :global(html[data-theme='light']) .gericht-desc,
  :global(html[data-theme='light']) .review p {
    color: rgba(26, 18, 0, 0.72);
  }

  :global(html[data-theme='light']) .zeiten,
  :global(html[data-theme='light']) .allergene {
    color: rgba(26, 18, 0, 0.45);
  }

  :global(html[data-theme='light']) .stepper span {
    color: #1a1200;
  }

  :global(html[data-theme='light']) .review-sperre,
  :global(html[data-theme='light']) .keine-reviews {
    background: rgba(255, 252, 235, 0.80);
    color: rgba(26, 18, 0, 0.66);
  }

  @media (max-width: 680px) {
    .seite {
      padding: 20px 14px 54px;
    }

    .kopf {
      align-items: flex-start;
      padding: 22px;
      gap: 16px;
    }

    .kopf .emoji {
      width: 76px;
      height: 76px;
      border-radius: 18px;
      font-size: 3rem;
    }

    .gericht {
      align-items: stretch;
      flex-direction: column;
    }

    .gericht-aktion {
      flex-direction: row;
      justify-content: space-between;
      align-items: center;
      width: 100%;
    }
  }
</style>