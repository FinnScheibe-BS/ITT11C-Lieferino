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
      onHinweis(`${anzahl}× ${gericht.name} hinzugefügt ✅`);
    }

    mengen = { ...mengen, [gericht.id]: 1 };
  }
</script>

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