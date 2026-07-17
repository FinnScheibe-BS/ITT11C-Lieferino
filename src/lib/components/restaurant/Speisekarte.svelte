<script>
  import { zumWarenkorb } from '$lib/stores/cart.js';

  let { restaurant, onHinweis } = $props();

  let mengen = $state({});

  function menge(id) {
    return mengen[id] ?? 0;
  }

  function aendere(id, delta) {
    const neu = Math.max(0, menge(id) + delta);
    mengen = { ...mengen, [id]: neu };
  }

  function hinzufuegen(gericht) {
    aendere(gericht.id, 1);
    zumWarenkorb(gericht, restaurant.name, 1);

    if (onHinweis) {
      onHinweis`${gericht.name} hinzugefügt`;
    }
  }

  // Diese Funktion ignoriert bild_url und baut die URL sauber über den Namen auf
  function bildUrl(gericht) {
    if (!gericht || !gericht.name) return '';
    
    // 1. Dateiendung aus der alten bild_url auslesen (default: png)
    let endung = 'png'; 
    if (gericht.bild_url) {
      const match = gericht.bild_url.match(/\.([a-zA-Z0-9]+)$/);
      if (match) {
        endung = match[1];
      }
    }
    
    // 2. Dateiname aus dem echten Namen und der Endung zusammenbauen
    const dateiname = `${gericht.name}.${endung}`;
    
    // 3. encodeURIComponent wandelt Leerzeichen in %20 und Sonderzeichen (z.B. &) perfekt um
    return `http://172.30.4.90:8080/uploads/gerichte/${encodeURIComponent(dateiname)}`;
  }
</script>

<h2>Speisekarte</h2>

<div class="speisekarte">
  {#each restaurant.speisekarte as gericht (gericht.id)}
    <article class="gericht">
      <div class="gericht-bild-container">
        <img
          class="gericht-bild"
          src={bildUrl(gericht)}
          alt={gericht.name}
          loading="lazy"
        />
      </div>

      <div class="gericht-text">
        <div class="gericht-titelzeile">
          <h3>{gericht.name}</h3>
          {#if gericht.veg || gericht.vegetarisch}<span class="veg-punkt" title="vegetarisch">🌱</span>{/if}
        </div>

        {#if gericht.beschreibung}
          <p class="gericht-desc">{gericht.beschreibung}</p>
        {/if}

        {#if gericht.allergene && gericht.allergene.length > 0}
          <p class="allergene">{gericht.allergene.join(' · ')}</p>
        {/if}
      </div>

      <div class="gericht-aktion">
        <span class="preis">{gericht.preis.toFixed(2)}€</span>

        {#if menge(gericht.id) === 0}
          <button class="plus-btn" onclick={() => hinzufuegen(gericht)} aria-label="Hinzufügen">
            +
          </button>
        {:else}
          <div class="stepper">
            <button onclick={() => aendere(gericht.id, -1)} aria-label="Weniger">−</button>
            <span>{menge(gericht.id)}</span>
            <button
              onclick={() => {
                aendere(gericht.id, 1);
                zumWarenkorb(gericht, restaurant.name, 1);
              }}
              aria-label="Mehr"
            >+</button>
          </div>
        {/if}
      </div>
    </article>
  {/each}
</div>

<style>
  h2 {
    margin: 0 0 24px;
    font-size: 1.8rem;
    font-weight: 700;
  }

  .speisekarte {
    display: flex;
    flex-direction: column;
    gap: 20px;
    padding: 8px 0;
  }

  .gericht {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 24px;
    padding: 20px;
    border-radius: 14px;
    background: rgba(255, 248, 220, 0.03);
    border: 1px solid rgba(230, 168, 0, 0.15);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }

  .gericht:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(230, 168, 0, 0.12);
  }

  .gericht-bild-container {
    width: 140px;
    height: 140px;
    flex-shrink: 0;
    border-radius: 10px;
    overflow: hidden;
    background: rgba(255, 248, 220, 0.06);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .gericht-bild {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }

  .gericht-text {
    flex: 1;
    min-width: 0;
  }

  .gericht-titelzeile {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 8px;
  }

  .gericht-titelzeile h3 {
    margin: 0;
    font-size: 1.35rem;
    font-weight: 700;
    line-height: 1.3;
  }

  .veg-punkt {
    font-size: 1.2rem;
    line-height: 1;
    flex-shrink: 0;
  }

  .gericht-desc {
    color: rgba(245, 240, 232, 0.75);
    line-height: 1.5;
    margin: 0 0 10px;
    font-size: 1.05rem;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .allergene {
    color: rgba(245, 240, 232, 0.45);
    font-size: 0.85rem;
    margin: 0;
    line-height: 1.4;
  }

  .gericht-aktion {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    justify-content: center;
    gap: 14px;
    padding-left: 24px;
    border-left: 1px solid rgba(230, 168, 0, 0.12);
    min-width: 130px;
  }

  .preis {
    color: #f9c932;
    font-weight: 700;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
    font-size: 1.3rem;
    white-space: nowrap;
  }

  .plus-btn {
    width: 42px;
    height: 42px;
    min-width: 42px;
    padding: 0 !important;
    border-radius: 50% !important;
    font-size: 1.4rem !important;
    line-height: 1 !important;
    background: transparent !important;
    border: 1px solid rgba(230, 168, 0, 0.35) !important;
    color: #f9c932 !important;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .plus-btn:hover {
    background: rgba(230, 168, 0, 0.15) !important;
    transform: scale(1.05);
  }

  .stepper {
    display: flex;
    align-items: center;
    gap: 10px;
    font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  }

  .stepper button {
    width: 36px;
    height: 36px;
    min-width: 36px;
    padding: 0 !important;
    border-radius: 50% !important;
    font-size: 1.1rem !important;
    line-height: 1 !important;
    background: transparent !important;
    border: 1px solid rgba(230, 168, 0, 0.35) !important;
    color: #f9c932 !important;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .stepper button:hover {
    background: rgba(230, 168, 0, 0.15) !important;
    transform: scale(1.05);
  }

  .stepper span {
    min-width: 24px;
    text-align: center;
    font-weight: 700;
    font-size: 1.1rem;
  }

  :global(html[data-theme='light']) .gericht {
    background: rgba(26, 18, 0, 0.02);
    border-color: rgba(26, 18, 0, 0.1);
  }

  :global(html[data-theme='light']) .gericht:hover {
    box-shadow: 0 4px 16px rgba(26, 18, 0, 0.08);
  }

  :global(html[data-theme='light']) .gericht-desc {
    color: rgba(26, 18, 0, 0.75);
  }

  :global(html[data-theme='light']) .allergene {
    color: rgba(26, 18, 0, 0.45);
  }

  @media (max-width: 650px) {
    .gericht {
      flex-direction: row;
      align-items: flex-start;
      gap: 16px;
      padding: 16px;
    }

    .gericht-bild-container {
      width: 90px;
      height: 90px;
    }

    .gericht-text {
      flex: 1;
    }

    .gericht-titelzeile h3 {
      font-size: 1.15rem;
    }

    .gericht-desc {
      font-size: 0.95rem;
      -webkit-line-clamp: 3;
    }

    .gericht-aktion {
      border-left: none;
      padding-left: 0;
      min-width: auto;
      flex-direction: column;
      align-self: stretch;
      justify-content: space-between;
    }

    .preis {
      font-size: 1.1rem;
    }

    .plus-btn {
      width: 36px !important;
      height: 36px !important;
      min-width: 36px !important;
      font-size: 1.2rem !important;
    }

    .stepper button {
      width: 32px !important;
      height: 32px !important;
      min-width: 32px !important;
    }
  }
</style>