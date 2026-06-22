<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { getRestaurant } from '$lib/data';
  import { zumWarenkorb } from '$lib/stores/cart.js';
  import { favoriten, toggleFavorit } from '$lib/stores/favoriten.js';
  import { bewertungen, bewertungHinzufuegen } from '$lib/stores/bewertungen.js';
  import { deaktivierteLieferanten } from '$lib/stores/lieferanten.js';

  // Den Slug aus der URL holen (z.B. "luigis-pizzeria") und das Restaurant suchen.
  // Deaktivierte Lieferanten behandeln wir für normale Nutzer wie "nicht gefunden".
  let slug = $derived($page.params.name);
  let restaurant = $derived($deaktivierteLieferanten.includes(slug) ? undefined : getRestaurant(slug));

  // Frühere Bestellungen laden, um zu prüfen, ob man hier schon bestellt hat.
  let meineBestellungen = $state([]);
  onMount(() => {
    meineBestellungen = JSON.parse(localStorage.getItem('lieferino_bestellungen') || '[]');
  });

  // ✅ Darf nur bewerten, wer hier nachweislich schon bestellt hat.
  let hatHierBestellt = $derived(
    meineBestellungen.some((b) =>
      (b.artikel || []).some((a) => a.restaurant === restaurant?.name)
    )
  );

  // Reviews dieses Restaurants aus dem Store ziehen.
  let reviews = $derived($bewertungen[slug] || []);

  // ⭐ Durchschnittsbewertung aus den Kundenreviews (sonst Basis-Bewertung).
  let durchschnitt = $derived.by(() => {
    if (reviews.length === 0) return restaurant?.bewertung ?? 0;
    const summe = reviews.reduce((s, r) => s + r.sterne, 0);
    return summe / reviews.length;
  });

  // 🕒 Ist das Restaurant gerade geöffnet? Vergleich der aktuellen Uhrzeit
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
    if (!hatHierBestellt) return; // Schutz: nur nach Bestellung
    if (neuName.trim() === '' || neuText.trim() === '') return;
    bewertungHinzufuegen(slug, { name: neuName, sterne: neuSterne, text: neuText });
    // Felder zurücksetzen
    neuName = '';
    neuText = '';
    neuSterne = 5;
  }

  // Kleine Bestätigung ("Pizza hinzugefügt"), die kurz eingeblendet wird.
  let hinweis = $state('');
  let hinweisTimer;

  function hinzufuegen(gericht) {
    const anzahl = menge(gericht.id);
    zumWarenkorb(gericht, restaurant.name, anzahl);
    hinweis = `${anzahl}× ${gericht.name} hinzugefügt ✅`;
    mengen = { ...mengen, [gericht.id]: 1 }; // Menge zurücksetzen
    // Alten Timer löschen, damit sich die Meldungen nicht überlagern.
    clearTimeout(hinweisTimer);
    hinweisTimer = setTimeout(() => (hinweis = ''), 2000);
  }
</script>

{#if restaurant}
  <div class="seite">
    <a href="/restaurants" class="zurueck">⬅️ Zu allen Restaurants</a>

    <!-- Kopfbereich des Restaurants -->
    <div class="kopf">
      <span class="emoji">{restaurant.emoji}</span>
      <div>
        <h1>
          {restaurant.name}
          <!-- ❤️ Favoriten-Herz -->
          <button class="herz" onclick={() => toggleFavorit(slug)} aria-label="Favorit">
            {$favoriten.includes(slug) ? '❤️' : '🤍'}
          </button>
        </h1>
        <p class="meta">
          <!-- Klick auf die Bewertung springt zu den Reviews -->
          <a href="#bewertungen" class="bewertung-link">⭐ {durchschnitt.toFixed(1)} ({reviews.length} Bewertungen)</a>
          · ⏱️ {restaurant.lieferzeit} · Min. {restaurant.minBestell}€
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
    </div>

    <h2>📋 Speisekarte</h2>
    <div class="speisekarte">
      {#each restaurant.speisekarte as gericht}
        <div class="gericht">
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
            <!-- Mengen-Stepper -->
            <div class="stepper">
              <button onclick={() => aendere(gericht.id, -1)} aria-label="Weniger">−</button>
              <span>{menge(gericht.id)}</span>
              <button onclick={() => aendere(gericht.id, 1)} aria-label="Mehr">+</button>
            </div>
            <button class="add-btn" onclick={() => hinzufuegen(gericht)}>+ Hinzufügen</button>
          </div>
        </div>
      {/each}
    </div>

    <!-- ⭐ BEWERTUNGEN -->
    <h2 id="bewertungen">⭐ Bewertungen ({reviews.length})</h2>

    <!-- Formular nur, wenn man hier schon bestellt hat -->
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

    <!-- Liste der vorhandenen Bewertungen -->
    {#if reviews.length === 0}
      <p class="keine-reviews">Noch keine Bewertungen – sei die/der Erste! 🌟</p>
    {:else}
      <div class="reviews">
        {#each reviews as review}
          <div class="review">
            <div class="review-kopf">
              <strong>{review.name}</strong>
              <span>{'⭐'.repeat(review.sterne)}</span>
            </div>
            <p>{review.text}</p>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <!-- Kurzer Hinweis nach dem Hinzufügen -->
  {#if hinweis}
    <div class="toast">{hinweis}</div>
  {/if}
{:else}
  <!-- Falls die URL kein gültiges Restaurant trifft -->
  <div class="seite">
    <h1>😕 Restaurant nicht gefunden</h1>
    <a href="/restaurants" class="zurueck">⬅️ Zu allen Restaurants</a>
  </div>
{/if}

<style>
  .seite { max-width: 800px; margin: 0 auto; padding: 20px; font-family: sans-serif; }
  .zurueck { color: #673ab7; text-decoration: none; font-weight: 600; }

  .kopf { display: flex; gap: 20px; align-items: center; background: #f3e5f5; padding: 24px; border-radius: 20px; margin: 16px 0 30px; }
  .kopf .emoji { font-size: 4.5rem; }
  .kopf h1 { margin: 0; }
  .meta { color: #555; font-size: 0.9rem; margin: 6px 0; }
  .beschreibung { color: #777; margin: 0; }

  .speisekarte { display: flex; flex-direction: column; gap: 12px; }
  .gericht { display: flex; justify-content: space-between; align-items: center; gap: 12px; background: white; border: 1px solid #eee; border-radius: 14px; padding: 16px; }
  .gericht-info h3 { margin: 0 0 4px; }
  .gericht-desc { color: #777; font-size: 0.85rem; margin: 0 0 6px; }
  .preis { font-weight: bold; color: #673ab7; }
  .add-btn { background: #673ab7; color: white; border: none; padding: 12px 18px; border-radius: 12px; font-weight: bold; cursor: pointer; white-space: nowrap; }
  .add-btn:hover { background: #5a2da3; }

  /* Öffnungsstatus */
  .oeffnung { margin: 4px 0; font-size: 0.9rem; }
  .status { font-weight: 700; }
  .status.auf { color: #2e9e4f; }
  .status.zu { color: #dc3545; }
  .zeiten { color: #888; font-size: 0.82rem; }

  /* Veggie-Tag + Allergene */
  .veg-tag { background: #e8f5e9; color: #2e7d32; font-size: 0.7rem; padding: 2px 8px; border-radius: 10px; font-weight: 600; margin-left: 6px; vertical-align: middle; }
  .allergene { color: #999; font-size: 0.78rem; margin: 0 0 6px; }

  /* Mengen-Stepper + Aktion */
  .gericht-aktion { display: flex; flex-direction: column; gap: 8px; align-items: flex-end; }
  .stepper { display: flex; align-items: center; gap: 8px; }
  .stepper button { width: 28px; height: 28px; border-radius: 50%; border: 1px solid #ddd; background: #f7f7f7; font-size: 1rem; cursor: pointer; }
  .stepper span { min-width: 18px; text-align: center; font-weight: bold; }

  /* Toast unten zentriert */
  .toast { position: fixed; bottom: 30px; left: 50%; transform: translateX(-50%); background: #34c759; color: white; padding: 14px 24px; border-radius: 30px; font-weight: bold; box-shadow: 0 6px 20px rgba(0,0,0,0.2); z-index: 99998; }

  /* ❤️ Herz-Button neben dem Namen */
  .herz { background: none; border: none; font-size: 1.4rem; cursor: pointer; vertical-align: middle; }

  /* ⭐ Bewertungen */
  .review-form { display: flex; flex-direction: column; gap: 10px; background: white; border: 1px solid #eee; border-radius: 14px; padding: 16px; margin-bottom: 16px; }
  .review-form input, .review-form select, .review-form textarea { padding: 10px; border: 1px solid #ddd; border-radius: 10px; font-size: 0.95rem; font-family: inherit; }
  .review-form textarea { min-height: 70px; resize: vertical; }
  .review-form button { background: #673ab7; color: white; border: none; padding: 12px; border-radius: 10px; font-weight: bold; cursor: pointer; }

  .bewertung-link { color: inherit; text-decoration: underline; text-decoration-style: dotted; cursor: pointer; }
  .review-sperre { background: #f7f7f7; border: 1px solid #eee; border-radius: 12px; padding: 14px; color: #777; }
  .keine-reviews { color: #777; }
  .reviews { display: flex; flex-direction: column; gap: 10px; }
  .review { background: white; border: 1px solid #eee; border-radius: 12px; padding: 14px; }
  .review-kopf { display: flex; justify-content: space-between; margin-bottom: 6px; }
  .review p { margin: 0; color: #555; }
</style>
