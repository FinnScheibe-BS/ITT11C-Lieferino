<script>
  import { restaurants } from "$lib/data/restaurants";
</script>

<h1>Restaurants</h1>

{#each restaurants as r}
  <div>
    <h2>{r.name}</h2>
    <p>⭐ {r.rating}</p>
    <p>⏱ {r.deliveryTime}</p>
  </div>
{/each}