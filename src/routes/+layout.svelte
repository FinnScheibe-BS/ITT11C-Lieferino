<div class="nav-container">
  <input type="checkbox" id="menu-toggle" class="menu-checkbox" />

  <label for="menu-toggle" class="nav-burger-btn">
    ☰
  </label>

  <div class="nav-dropdown-balken">
    <div class="button-umrundung"></div>

    <div class="nav-links-wrapper">
      <a href="/">🏠 Home</a>
      <a href="/restaurants">🍔 Restaurants</a>
      <a href="/cart">🛒 Warenkorb</a>
      <a href="/login">👤 Login</a>
    </div>
  </div>
</div>

<div class="page-content">
  <slot />
</div>

<style>
  /* Container fixieren */
  .nav-container {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    z-index: 10000;
  }

  .menu-checkbox {
    display: none;
  }

  /* BURGER BUTTON (Kompakter & ohne Eigenschatten für perfekten Übergang) */
  .nav-burger-btn {
    position: fixed;
    top: 20px;
    right: 20px; 
    width: 40px;  /* 👈 Etwas schmaler */
    height: 36px; /* 👈 Etwas flacher */
    z-index: 10003;
    background: #673ab7; 
    border: none;
    border-radius: 6px; /* Leicht angepasster Radius für die kleinere Größe */
    cursor: pointer;
    font-size: 18px; /* Schriftgröße leicht verringert */
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    user-select: none;
    box-shadow: none; /* 👈 Kein eigener Schatten mehr, damit er mit dem Hintergrund verschmilzt */
    transition: right 0.3s ease;
  }

  /* DER LILA BALKEN (Seitenwand) */
  .nav-dropdown-balken {
    position: fixed;
    top: 0;
    bottom: 0;
    right: -225px; 
    width: 240px; 
    height: 100vh; 
    background: #673ab7; 
    z-index: 10001;
    transition: right 0.3s cubic-bezier(0.25, 1, 0.5, 1), box-shadow 0.3s ease;
    box-shadow: none;
  }

  /* DIE UMRUNDUNG (Exakt auf die neue kompakte Button-Größe kalibriert) */
  .button-umrundung {
    position: absolute;
    top: 14px;      
    left: -48px;    /* Genau austariert für die neue Button-Breite */
    width: 60px;    
    height: 48px;   /* Umschließt den kleineren Button perfekt gleichmäßig */
    background: #673ab7; 
    border-radius: 12px 0 0 12px; 
    z-index: 10002;
    transition: opacity 0.2s ease;
  }

  /* WENN GEÖFFNET */
  .menu-checkbox:checked ~ .nav-dropdown-balken {
    right: 0; 
    box-shadow: -8px 0 25px rgba(0, 0, 0, 0.4); 
  }

  .menu-checkbox:checked ~ .nav-dropdown-balken .button-umrundung {
    opacity: 0; 
  }

  /* Da der Button kleiner ist, wandert er beim Öffnen ein Stück weniger weit nach links */
  .menu-checkbox:checked ~ .nav-burger-btn {
    right: 180px; 
  }

  /* LINKS IM MENÜ */
  .nav-links-wrapper {
    display: flex;
    flex-direction: column;
    padding-top: 90px; 
    position: relative;
    z-index: 10004;
  }

  .nav-links-wrapper a {
    color: white;
    text-decoration: none;
    padding: 15px 25px;
    font-size: 18px;
    font-family: sans-serif;
    font-weight: 700; 
    transition: background 0.2s;
    white-space: nowrap;
  }

  .nav-links-wrapper a:hover {
    background: rgba(255, 255, 255, 0.15);
  }

  .page-content {
    padding: 20px;
  }
</style>