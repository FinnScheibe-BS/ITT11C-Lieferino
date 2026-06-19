<script>
  let name = "";
  let loggedIn = false;

  function login() {
    loggedIn = true;
  }

  function logout() {
    loggedIn = false;
    name = "";
  }
</script>

<h1>Account 👤</h1>

{#if !loggedIn}
  <input placeholder="Name eingeben" bind:value={name} />
  <button on:click={login}>Login</button>
{:else}
  <p>Hallo {name} 👋</p>
  <button on:click={logout}>Logout</button>
{/if}