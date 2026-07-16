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

  onMount(async () => {
    console.log('🔍 START - Slug:', slug);
    console.log('📦 Store vor Load:', $aktiveRestaurants?.length ?? 0, 'Restaurants');
    
    try {
      // Erst im Store suchen
      restaurant = findeRestaurant(slug);
      console.log('✅ Im Store gefunden:', restaurant?.name ?? 'NEIN');
      
      // Wenn nicht gefunden, nachladen
      if (!restaurant) {
        console.log('⏳ Lade Restaurants nach...');
        const alle = await ladeRestaurants();
        console.log('📥 Nachgeladen:', alle.length, 'Restaurants');
        
        restaurant = alle.find(r => r.slug === slug);
        console.log('✅ Nach Laden gefunden:', restaurant?.name ?? 'NEIN');
        
        // Fallback: Slug generieren und vergleichen
        if (!restaurant) {
          const slugFromName = nameToSlug(slug.replace(/-/g, ' '));
          restaurant = alle.find(r => nameToSlug(r.name) === slug);
          console.log('✅ Fallback-Suche:', restaurant?.name ?? 'NEIN');
        }
      }
      
      // Gerichte laden wenn Restaurant gefunden
      if (restaurant?.id) {
        console.log('🍽️ Lade Gerichte für Restaurant ID:', restaurant.id);
        await ladeSpeisen(restaurant.id);
        console.log('✅ Gerichte geladen:', $speisen[restaurant.id]?.length ?? 0);
      }
      
      ladezustand = restaurant ? 'ready' : 'not_found';
      console.log('🏁 ENDE - Zustand:', ladezustand);
      
    } catch (error) {
      console.error('❌ FEHLER im onMount:', error);
      ladezustand = 'error';
    }
  });

  function nameToSlug(name) {
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
      <p class="debug">Slug: {slug}</p>
    </div>
  </div>
{:else if ladezustand === 'error'}
  <div class="seite">
    <div class="fehler-zustand karte">
      <h1>😕 Fehler beim Laden</h1>
      <p>Bitte prüfe die Browser-Console (F12)</p>
      <a href="/restaurants" class="zurueck">⬅️ Zu allen Restaurants</a>
    </div>
  </div>
{:else if restaurant}
  <!-- Normale Restaurant-Anzeige -->
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
  .ladezustand, .fehler-zustand {
    padding: 48px 24px;
    text-align: center;
  }
  .debug {
    font-size: 0.8rem;
    opacity: 0.5;
    margin-top: 8px;
  }
  .fehler-zustand {
    color: #ff6b6b;
  }
</style>