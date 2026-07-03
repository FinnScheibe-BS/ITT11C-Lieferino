<script>
  import { page } from '$app/stores';
  import { deaktivierteLieferanten, findeRestaurant } from '$lib/stores/lieferanten.js';
  import Speisekarte from '$lib/components/restaurant/Speisekarte.svelte';
  import RestaurantHeader from '$lib/components/restaurant/RestaurantHeader.svelte';
  import Bewertungen from '$lib/components/restaurant/Bewertungen.svelte';

  let slug = $derived($page.params.name);
  let restaurant = $derived($deaktivierteLieferanten.includes(slug) ? undefined : findeRestaurant(slug));

  let geoeffnet = $derived.by(() => {
    if (!restaurant?.oeffnetUm) return true;

    const jetzt = new Date();
    const minutenJetzt = jetzt.getHours() * 60 + jetzt.getMinutes();

    const [oh, om] = restaurant.oeffnetUm.split(':').map(Number);
    const [sh, sm] = restaurant.schliesstUm.split(':').map(Number);

    return minutenJetzt >= oh * 60 + om && minutenJetzt < sh * 60 + sm;
  });

  let hinweis = $state('');
  let hinweisTimer;
</script>

{#if restaurant}
  <div class="seite">
    <a href="/restaurants" class="zurueck">⬅️ Zu allen Restaurants</a>

    <RestaurantHeader {restaurant} {slug} {geoeffnet} />

    <Speisekarte
      {restaurant}
      onHinweis={(text) => {
        hinweis = text;
        clearTimeout(hinweisTimer);
        hinweisTimer = setTimeout(() => (hinweis = ''), 2000);
      }}
    />

    <Bewertungen {restaurant} {slug} />
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