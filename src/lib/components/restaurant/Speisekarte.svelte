<script>
  import { zumWarenkorb } from '$lib/stores/cart.js';

  let { restaurant, onHinweis } = $props();

  let mengen = $state({});

  function menge(id) {
    return mengen[id] ?? 1;
  }

  function aendere(id, delta) {
    const neu = Math.max(1, menge(id) + delta);
    mengen = { ...mengen, [id]: neu };
  }

  function hinzufuegen(gericht) {
    const anzahl = menge(gericht.id);

    zumWarenkorb(gericht, restaurant.name, anzahl);

    if (onHinweis) {
      onHinweis(`${anzahl}× ${gericht.name} hinzugefügt`);
    }

    mengen = { ...mengen, [gericht.id]: 1 };
  }
</script>

<h2>Speisekarte</h2>

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

        <button class="add-btn" onclick={() => hinzufuegen(gericht)}>Hinzufügen</button>
      </div>
    </article>
  {/each}
</div>

<style>
h2 {
  margin: 0 0 18px;
  font-size: 1.25rem;
}

.speisekarte {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.gericht {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
  padding: 24px;
  border-radius: 16px;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.gericht:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.gericht-info {
  min-width: 0;
}

.gericht-info h3 {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
  margin: 0 0 8px;
  font-size: 1.05rem;
}

.gericht-desc {
  color: rgba(245, 240, 232, 0.68);
  font-size: 0.9rem;
  line-height: 1.5;
  margin: 0 0 10px;
}

.allergene {
  color: rgba(245, 240, 232, 0.42);
  font-size: 0.78rem;
  margin: 0 0 10px;
}

.preis {
  display: inline-flex;
  color: #f9c932;
  font-weight: 800;
  font-size: 1.05rem;
  font-family: 'Geist Sans', -apple-system, 'SF Pro Display', sans-serif;
}

.veg-tag {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  padding: 3px 10px;
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
  gap: 12px;
  align-items: flex-end;
  flex-shrink: 0;
}

.stepper {
  display: flex;
  align-items: center;
  gap: 10px;
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

:global(html[data-theme='light']) .gericht-desc {
  color: rgba(26, 18, 0, 0.72);
}

:global(html[data-theme='light']) .allergene {
  color: rgba(26, 18, 0, 0.45);
}

:global(html[data-theme='light']) .stepper span {
  color: #1a1200;
}

@media (max-width: 680px) {
  .gericht {
    align-items: stretch;
    flex-direction: column;
    padding: 20px;
  }

  .gericht-aktion {
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    width: 100%;
  }
}
</style>