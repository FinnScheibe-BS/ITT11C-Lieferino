<script>
  import { page } from '$app/stores';
  import { findeRestaurant, ladeRestaurants } from '$lib/stores/lieferanten.js';
  import { ladeSpeisen, speisen } from '$lib/stores/speisen.js';
  import Speisekarte from '$lib/components/restaurant/Speisekarte.svelte';
  import Bewertungen from '$lib/components/restaurant/Bewertungen.svelte';
  import { onMount } from 'svelte';

  let slug = $derived($page.params.name);
  let restaurant = $state(null);
  let gerichte = $derived(restaurant?.id ? ($speisen[restaurant.id] ?? []) : []);
  let ladezustand = $state('loading');

  // Ort aus Adresse (letzter Teil nach Komma)
  let ort = $derived(
    restaurant?.adresse
      ? restaurant.adresse.split(',').pop().trim()
      : ''
  );

  // Öffnungsstatus berechnen
  let zeitInfo = $derived.by(() => {
    if (!restaurant?.oeffnetUm || !restaurant?.schliesstUm) {
      return { offen: true, text: 'Öffnungszeiten unbekannt', klasse: 'offen' };
    }

    const jetzt = new Date();
    const minutenJetzt = jetzt.getHours() * 60 + jetzt.getMinutes();
    
    const [oh, om] = restaurant.oeffnetUm.split(':').map(Number);
    const [sh, sm] = restaurant.schliesstUm.split(':').map(Number);
    
    const oeffnetMinuten = oh * 60 + om;
    const schliesstMinuten = sh * 60 + sm;
    
    const istOffen = minutenJetzt >= oeffnetMinuten && minutenJetzt < schliesstMinuten;
    
    if (istOffen) {
      return {
        offen: true,
        text: `Offen bis ${restaurant.schliesstUm} Uhr`,
        klasse: 'offen'
      };
    } else {
      let diffMinuten;
      
      if (minutenJetzt < oeffnetMinuten) {
        diffMinuten = oeffnetMinuten - minutenJetzt;
      } else {
        const minutenBisMitternacht = 1440 - minutenJetzt;
        diffMinuten = minutenBisMitternacht + oeffnetMinuten;
      }
      
      let zeitText;
      if (diffMinuten < 60) {
        zeitText = `in ${diffMinuten} Min`;
      } else if (diffMinuten < 1440) {
        const stunden = Math.floor(diffMinuten / 60);
        const min = diffMinuten % 60;
        zeitText = min > 0 ? `in ${stunden}h ${min}min` : `in ${stunden}h`;
      } else {
        zeitText = `morgen um ${restaurant.oeffnetUm}`;
      }
      
      return {
        offen: false,
        text: `Geschlossen – Öffnet ${zeitText}`,
        klasse: 'geschlossen'
      };
    }
  });

  let hinweis = $state('');
  let hinweisTimer;

  onMount(async () => {
    ladezustand = 'loading';
    restaurant = findeRestaurant(slug);
    if (!restaurant) {
      const alle = await ladeRestaurants();
      restaurant = alle.find(r => r.slug === slug);
      if (!restaurant) {
        restaurant = alle.find(r =>
          r.name.toLowerCase().includes(slug.replace(/-/g, ' '))
        );
      }
    }
    if (restaurant?.id) {
      await ladeSpeisen(restaurant.id);
    }
    ladezustand = restaurant ? 'ready' : 'not_found';
  });
</script>

{#if ladezustand === 'loading'}
  <div class="seite">
    <div class="ladezustand karte">
      <p>🍽️ Restaurant wird geladen...</p>
    </div>
  </div>
{:else if restaurant}
  <div class="seite">
    <div class="header-bereich">
      <a href="/restaurants" class="zurueck-link">⬅️ Zurück zu allen Restaurants</a>
      <h1>{restaurant.name}</h1>
      <div class="meta-info">
        <span class="tag">{ort}</span>
        <span class="status {zeitInfo.klasse}">
          <span class="punkt"></span>
          <span class="status-text">{zeitInfo.text}</span>
        </span>
      </div>
    </div>

    <div class="inhalt">
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
  </div>

  {#if hinweis}
    <div class="toast">{hinweis}</div>
  {/if}
{:else}
  <div class="seite">
    <div class="header-bereich">
      <a href="/restaurants" class="zurueck-link">⬅️ Zurück zu allen Restaurants</a>
      <h1>😕 Restaurant nicht gefunden</h1>
      <p>Das Restaurant "{slug}" existiert nicht.</p>
      <a href="/restaurants" class="button">Zur Übersicht</a>
    </div>
  </div>
{/if}

<style>
  .seite {
    max-width: 1200px;
    margin: 0 auto;
    padding: 24px 20px 48px;
  }

  .header-bereich {
    text-align: center;
    padding: 64px 24px 40px;
    margin-bottom: 24px;
  }

  .zurueck-link {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    color: rgba(245, 240, 232, 0.7);
    text-decoration: none;
    font-size: 0.95rem;
    margin-bottom: 24px;
    transition: color 0.3s ease;
  }

  .zurueck-link:hover {
    color: #f9c932;
  }

  .header-bereich h1 {
    margin: 0 0 16px;
    font-size: clamp(3rem, 7vw, 4.5rem);
    font-weight: 700;
    letter-spacing: -0.03em;
    color: #f9c932;
    text-shadow: 0 2px 12px rgba(230, 168, 0, 0.4);
  }

  .meta-info {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
  }

  .tag {
    display: inline-flex;
    align-items: center;
    background: rgba(230, 168, 0, 0.2);
    color: #f9c932;
    padding: 6px 14px;
    line-height: 1.3;
    border-radius: 12px;
    font-size: 0.85rem;
    font-weight: 600;
    border: 1px solid rgba(230, 168, 0, 0.25);
  }

  .status {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 6px 14px;
    border-radius: 12px;
    font-size: 0.85rem;
    font-weight: 600;
  }

  .status.offen {
    background: rgba(76, 175, 80, 0.15);
    color: #4caf50;
    border: 1px solid rgba(76, 175, 80, 0.3);
  }

  .status.geschlossen {
    background: rgba(244, 67, 54, 0.15);
    color: #f44336;
    border: 1px solid rgba(244, 67, 54, 0.3);
  }

  .punkt {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    animation: pulsieren 2s infinite;
    flex-shrink: 0;
  }

  .status.offen .punkt {
    background: #4caf50;
    box-shadow: 0 0 8px rgba(76, 175, 80, 0.6);
  }

  .status.geschlossen .punkt {
    background: #f44336;
    box-shadow: 0 0 8px rgba(244, 67, 54, 0.6);
  }

  .status-text {
    white-space: nowrap;
  }

  @keyframes pulsieren {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.4; }
  }

  .inhalt {
    display: flex;
    flex-direction: column;
    gap: 32px;
  }

  .ladezustand {
    padding: 48px 24px;
    text-align: center;
    opacity: 0.7;
  }

  .button {
    display: inline-block;
    background: linear-gradient(135deg, #e6a800 0%, #f9c932 100%);
    color: #1c1710;
    padding: 12px 28px;
    border-radius: 12px;
    font-weight: 600;
    text-decoration: none;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
  }

  .button:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(230, 168, 0, 0.3);
  }

  .toast {
    position: fixed;
    bottom: 24px;
    left: 50%;
    transform: translateX(-50%);
    background: #1c1710;
    border: 1px solid rgba(230, 168, 0, 0.3);
    color: #f9c932;
    padding: 12px 24px;
    border-radius: 12px;
    font-weight: 500;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
    z-index: 1000;
    animation: slideUp 0.3s ease;
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateX(-50%) translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateX(-50%) translateY(0);
    }
  }

  :global(html[data-theme='light']) .zurueck-link {
    color: rgba(26, 18, 0, 0.6);
  }

  :global(html[data-theme='light']) .toast {
    background: #fff;
    color: #1c1710;
    border-color: rgba(230, 168, 0, 0.2);
  }
</style>