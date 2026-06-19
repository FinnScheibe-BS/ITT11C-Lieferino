<script>
  import { warenkorb } from '$lib/stores/cart.js';
  
  // Der '?' Operator sorgt dafür: Wenn warenkorb undefined ist, 
  // stürzt es nicht ab, sondern gibt 0 zurück.
  let anzahl = $derived($warenkorb?.length ?? 0);
</script>

<div class="nav-container">
  <input type="checkbox" id="menu-toggle" class="menu-checkbox" />
  <label for="menu-toggle" class="nav-burger-btn">
    ☰
    {#if anzahl > 0}
      <span class="cart-badge-burger">{anzahl}</span>
    {/if}
  </label>

  <div class="nav-dropdown-balken">
    <div class="button-umrundung"></div>

    <div class="nav-links-wrapper">
      <a href="/">🏠 Home</a>
      <a href="/restaurants">🍔 Restaurants</a>
      <a href="/cart" class="cart-link">
        🛒 Warenkorb
        {#if anzahl > 0}
          <span class="cart-badge">{anzahl}</span>
        {/if}
      </a>
      <a href="/account">👤 Account</a>

      <div class="nav-impressum">
        <h4>Impressum</h4>
        <p>Lieferino GmbH<br>Musterstraße 12<br>12345 Stadt</p>
      </div>
    </div>
  </div>
</div>

<div class="page-content">
  <slot />
</div>

<style>
  .nav-container {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    z-index: 99999;
  }

  .menu-checkbox {
    display: none !important;
  }

  .nav-burger-btn {
    position: fixed !important;
    top: 20px;
    right: 4px;
    width: 40px;
    height: 36px;
    z-index: 100002;
    background: transparent !important;
    border: none !important;
    cursor: pointer;
    font-size: 26px;
    color: white !important;
    display: flex !important;
    align-items: center;
    justify-content: center;
    user-select: none;
    transition: right 0.3s ease;
  }

  .cart-badge-burger {
    position: absolute;
    top: -2px;
    right: -2px;
    background: #ff3b30;
    color: white;
    font-size: 11px;
    font-weight: 700;
    line-height: 1;
    padding: 2px 5px;
    border-radius: 10px;
    font-family: sans-serif;
  }

  .nav-dropdown-balken {
    position: fixed !important;
    top: 0;
    bottom: 0;
    right: -215px;
    width: 240px;
    height: 100vh;
    background: #673ab7 !important;
    z-index: 100000;
    transition: right 0.3s cubic-bezier(0.25, 1, 0.5, 1), box-shadow 0.3s ease;
    box-shadow: none;
  }

  .button-umrundung {
    position: absolute !important;
    top: 14px;
    left: -18px;
    width: 35px;
    height: 48px;
    background: #673ab7 !important;
    border-radius: 12px 0 0 12px;
    z-index: 100001;
    transition: opacity 0.2s ease;
  }

  .menu-checkbox:checked ~ .nav-dropdown-balken {
    right: 0 !important;
    box-shadow: -8px 0 25px rgba(0, 0, 0, 0.4) !important;
  }

  .menu-checkbox:checked ~ .nav-dropdown-balken .button-umrundung {
    opacity: 0 !important;
  }

  .menu-checkbox:checked ~ .nav-burger-btn {
    right: 185px !important;
  }

  .nav-links-wrapper {
    display: flex !important;
    flex-direction: column !important;
    height: 100%;
    padding-top: 90px;
    box-sizing: border-box;
    position: relative;
    z-index: 100003;
  }

  .nav-links-wrapper a {
    color: white !important;
    text-decoration: none !important;
    padding: 15px 25px !important;
    font-size: 18px !important;
    font-family: sans-serif !important;
    font-weight: 700 !important;
    display: block !important;
    transition: background 0.2s;
    white-space: nowrap;
  }

  .nav-links-wrapper a:hover {
    background: rgba(255, 255, 255, 0.15) !important;
  }

  .cart-link {
    display: flex !important;
    align-items: center;
    gap: 8px;
  }

  .cart-badge {
    background: white;
    color: #673ab7;
    font-size: 0.75rem;
    font-weight: 800;
    padding: 2px 8px;
    border-radius: 12px;
    line-height: 1.3;
  }

  .nav-impressum {
    margin-top: auto;
    padding: 25px;
    color: rgba(255, 255, 255, 0.7);
    font-family: sans-serif;
    font-size: 0.8rem;
    border-top: 1px solid rgba(255, 255, 255, 0.1);
  }

  .nav-impressum h4 {
    margin: 0 0 5px 0;
    color: white;
  }

  .nav-impressum p {
    margin: 0;
    line-height: 1.3;
  }

  .page-content {
    padding: 20px;
  }
</style>