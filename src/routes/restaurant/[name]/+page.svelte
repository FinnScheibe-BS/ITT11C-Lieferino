<script>
  import { page } from '$app/stores';
  import { 
    deaktivierteLieferanten, 
    findeRestaurant,
    ladeRestaurants,
    aktiveRestaurants 
  } from '$lib/stores/lieferanten.js';
  import { ladeSpeisen, speisen } from '$lib/stores/speisen.js';
  import Speisekarte from '$lib/components/restaurant/Speisekarte.svelte';
  import RestaurantHeader from '$lib/components/restaurant/RestaurantHeader.svelte';
  import Bewertungen from '$lib/components/restaurant/Bewertungen.svelte';
  import { onMount } from 'svelte';

  let slug = $derived($page.params.name);
  let restaurant = $state(null);
  let gerichte = $derived(restaurant?.id ? ($speisen[restaurant.id] ?? []) : []);
  let ladezustand = $state('loading');

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

  // Restaurant finden oder nachladen
  onMount(async () => {
    ladezustand = 'loading';
    
    // Erst im Store suchen
    restaurant = findeRestaurant(slug);
    
    // Wenn nicht gefunden, alle Restaurants nachladen
    if (!restaurant) {
      console.log('Restaurant nicht im Store, lade nach...', slug);
      const alle = await ladeRestaurants();
      restaurant = alle.find(r => r.slug === slug);
      
      // Fallback: Nach Name suchen falls Slug nicht matcht
      if (!restaurant) {
        restaurant = alle.find(r => 
          erstelleSlug(r.name) === slug || 
          r.name.toLowerCase().includes(slug.replace(/-/g, ' '))
        );
      }
    }
    
    // Gerichte laden wenn Restaurant gefunden
    if (restaurant?.id) {
      await ladeSpeisen(restaurant.id);
    }
    
    ladezustand = restaurant ? 'ready' : 'not_found';
  });

  function erstelleSlug(name) {
    return name
      .toLowerCase()
      .replace(/ä/g, 'ae')
      .replace(/ö/g, 'oe')
      .replace(/ü/g, 'ue')
      .replace(/ß/g, 'ss')
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }
</script>

{#if ladezustand === 'loading'}
  <div class="seite">
    <div class="ladezustand karte">
      <p>🍽️ Restaurant wird geladen...</p>
    </div>
  </div>
{:else if restaurant}
  <div class="seite">
    <a href="/restaurants" class="zurueck">⬅️ Zu allen Restaurants</a>

    <RestaurantHeader {restaurant} {slug} {geoeffnet} />

    <Speisekarte
      restaurant={{ ...restaurant, speisekarte: gerichte }}
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
      <p>Das Restaurant "{slug}" existiert nicht.</p>
      <a href="/restaurants" class="zurueck">⬅️ Zu allen Restaurants</a>
    </div>
  </div>
{/if}

<style>
  .ladezustand {
    padding: 48px 24px;
    text-align: center;
    opacity: 0.7;
  }
</style>