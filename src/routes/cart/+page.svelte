<script>
  import { cart } from "$lib/stores/cart";
</script>

<h1>Warenkorb</h1>

{#each $cart as item}
  <p>{item.name} - {item.price}</p>
{/each}