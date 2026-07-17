<script>
  import { zumWarenkorb } from '$lib/stores/cart.js';

  let { restaurant, onHinweis } = $props();

  let mengen = $state({});
  let bildFehler = $state({});

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
      onHinweis(`${gericht.name} hinzugefügt`);
    }
  }

  function bildUrl(gericht) {
    return `http://172.30.4.90:8080/uploads/gerichte/${encodeURIComponent(gericht.name)}.jpeg`;
  }

  function markiereBildFehler(id) {
    bildFehler = { ...bildFehler, [id]: true };
  }
</script>

<h2>🍽️ Speisekarte</h2>

<div class="speisekarte">
  {#each restaurant.speisekarte as gericht}
    <article class="gericht">
      {#if !bildFehler[gericht.id]}
        <img
          class="gericht-bild"
          src={bildUrl(gericht)}
          alt={gericht.name}
          loading="lazy"
          onerror={() => markiereBildFehler(gericht.id)}
        />
      {:else}
        <div class="gericht-bild gericht-bild-platzhalter">🍽️</div>
      {/if}

      <div class="gericht-text">
        <div class="gericht-titelzeile">
          <h3>{gericht.name}</h3>
          {#if gericht.veg}<span class="veg-punkt" title="vegetarisch">🌱</span>{/if}
        </div>

        <p class="gericht-desc">{gericht.beschreibung}</p>

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
  margin: 0 0 16px;
}

.speisekarte {
  display: flex;
  flex-direction: column;
}

.gericht {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 28px 4px;
  border-bottom: 1px solid rgba(230, 168, 0, 0.1);
  transition: background-color 0.15s ease;
}

.gericht:first-child {
  padding-top: 4px;
}

.gericht:last-child {
  border-bottom: none;
  padding-bottom: 4px;
}

.gericht:hover {
  background: rgba(230, 168, 0, 0.04);
}

.gericht-bild {
  width: 88px;
  height: 88px;
  min-width: 88px;
  border-radius: 12px;
  object-fit: cover;
  background: rgba(255, 248, 220, 0.06);
}

.gericht-bild-platzhalter {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.8rem;
  border: 1px solid rgba(230, 168, 0, 0.15);
}

.gericht-text {
  flex: 1;
  min-width: 0;
}

.gericht-titelzeile {
  display: flex;
  align-items: center;
  gap: 8px;
}

.gericht-titelzeile h3 {
  margin: 0;
  font-weight: 600;
}

.veg-punkt {
  font-size: 0.9rem;
  line-height: 1;
}

.gericht-desc {
  color: rgba(245, 240, 232, 0.65);
  line-height: 1.55;
  margin: 6px 0 0;
  max-width: 46ch;
}

.allergene {
  color: rgba(245, 240, 232, 0.4);
  font-size: 0.9rem;
  margin: 8px 0 0;
}

.gericht-aktion {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
}

.preis {
  color: #f9c932;
  font-weight: 700;
  font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
  white-space: nowrap;
}

.plus-btn {
  width: 36px;
  height: 36px;
  min-width: 36px;
  padding: 0 !important;
  border-radius: 50% !important;
  font-size: 1.2rem !important;
  line-height: 1 !important;
  background: transparent !important;
  border: 1px solid rgba(230, 168, 0, 0.35) !important;
  color: #f9c932 !important;
}

.plus-btn:hover {
  background: rgba(230, 168, 0, 0.12) !important;
}

.stepper {
  display: flex;
  align-items: center;
  gap: 10px;
  font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
}

.stepper button {
  width: 32px;
  height: 32px;
  min-width: 32px;
  padding: 0 !important;
  border-radius: 50% !important;
  font-size: 1.05rem !important;
  line-height: 1 !important;
  background: transparent !important;
  border: 1px solid rgba(230, 168, 0, 0.35) !important;
  color: #f9c932 !important;
}

.stepper button:hover {
  background: rgba(230, 168, 0, 0.12) !important;
}

.stepper span {
  min-width: 18px;
  text-align: center;
  font-weight: 700;
}

:global(html[data-theme='light']) .gericht {
  border-bottom-color: rgba(26, 18, 0, 0.08);
}

:global(html[data-theme='light']) .gericht:hover {
  background: rgba(230, 168, 0, 0.05);
}

:global(html[data-theme='light']) .gericht-desc {
  color: rgba(26, 18, 0, 0.62);
}

:global(html[data-theme='light']) .allergene {
  color: rgba(26, 18, 0, 0.4);
}

@media (max-width: 680px) {
  .gericht {
    gap: 14px;
    padding: 22px 4px;
  }

  .gericht-bild {
    width: 68px;
    height: 68px;
    min-width: 68px;
  }

  .gericht-desc {
    max-width: none;
  }
}
</style>