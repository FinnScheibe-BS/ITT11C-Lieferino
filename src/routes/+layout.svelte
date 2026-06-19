<script>
  let open = false;
</script>

<button on:click={() => open = !open}>
  TOGGLE
</button>

{#if open}
  <div style="position:fixed; top:0; right:0; width:200px; height:100vh; background:orange;">
    OPEN
  </div>
{/if}