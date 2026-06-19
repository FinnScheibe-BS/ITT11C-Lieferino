<script>
  import { page } from '$app/stores';
  // Pfad ist jetzt absolut korrekt auf dein 'stores/cart.js' ausgerichtet
  import { warenkorb } from '$lib/stores/cart.js';

  // Wir greifen auf die ID aus der URL zu (z.B. "luigi")
  let restaurantId = $page.params.name;

  let speisekarte = [
    { id: 1, name: "Pizza", preis: 8.50 },
    { id: 2, name: "Pasta", preis: 10.00 }
  ];

  function inWarenkorb(gericht) {
    warenkorb.update(n => [...n, gericht]);
  }
</script>

<h1>Restaurant ID: {restaurantId}</h1>

{#each speisekarte as gericht}
  <div style="margin-bottom: 10px; border: 1px solid #eee; padding: 10px;">
    {gericht.name} - {gericht.preis}€
    <button on:click={() => inWarenkorb(gericht)}>+ Hinzufügen</button>
  </div>
{/each}

<a href="/">Zurück zur Übersicht</a>