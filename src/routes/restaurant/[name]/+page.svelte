<script>
  import { page } from '$app/stores';
  import { getRestaurant } from '$lib/data';
  import { zumWarenkorb } from '$lib/stores/cart.js';

  // Den Slug aus der URL holen (z.B. "luigis-pizzeria") und das Restaurant suchen.
  let slug = $derived($page.params.name);
  let restaurant = $derived(getRestaurant(slug));

  // Kleine Bestätigung ("Pizza hinzugefügt"), die kurz eingeblendet wird.
  let hinweis = $state('');
  let hinweisTimer;

  function hinzufuegen(gericht) {
    zumWarenkorb(gericht, restaurant.name);
    hinweis = `${gericht.name} wurde hinzugefügt ✅`;
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
        <h1>{restaurant.name}</h1>
        <p class="meta">
          ⭐ {restaurant.bewertung} · ⏱️ {restaurant.lieferzeit} · Mindestbestellwert {restaurant.minBestell}€
        </p>
        <p class="beschreibung">{restaurant.beschreibung}</p>
      </div>
    </div>

    <h2>📋 Speisekarte</h2>
    <div class="speisekarte">
      {#each restaurant.speisekarte as gericht}
        <div class="gericht">
          <div class="gericht-info">
            <h3>{gericht.name}</h3>
            <p class="gericht-desc">{gericht.beschreibung}</p>
            <span class="preis">{gericht.preis.toFixed(2)}€</span>
          </div>
          <button class="add-btn" onclick={() => hinzufuegen(gericht)}>+ Hinzufügen</button>
        </div>
      {/each}
    </div>
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

  /* Toast unten zentriert */
  .toast { position: fixed; bottom: 30px; left: 50%; transform: translateX(-50%); background: #34c759; color: white; padding: 14px 24px; border-radius: 30px; font-weight: bold; box-shadow: 0 6px 20px rgba(0,0,0,0.2); z-index: 99998; }
</style>
