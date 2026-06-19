<script>
  // Ändere 'cart' zu 'warenkorb'
  import { warenkorb } from "$lib/stores/cart";
</script>

<h1>Warenkorb</h1>

{#each $warenkorb as item}
  <p>{item.name} - {item.preis}€</p>
{/each}