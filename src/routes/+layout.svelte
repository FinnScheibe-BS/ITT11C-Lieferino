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
  /* Globaler Container, der das Menü rechts fixiert */
  .nav-container {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    z-index: 99999; /* Extrem hoch, damit es über allem liegt */
  }

  /* Versteckt die Checkbox, die wir für den Klick-Status nutzen */
  .menu-checkbox {
    display: none !important;
  }

  /* BURGER BUTTON (Das weiße Icon, das auf dem Balken sitzt) */
  .nav-burger-btn {
    position: fixed !important;
    top: 20px;
    right: 4px; /* Sitzt voll auf dem schmalen Balken */
    width: 40px;  
    height: 36px; 
    z-index: 100002; /* Liegt ganz oben */
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

  /* DER LILA BALKEN (Die Seitenwand) */
  .nav-dropdown-balken {
    position: fixed !important;
    top: 0;
    bottom: 0;
    right: -215px; /* Versteckt die Wand, lässt ca. 25px als Streifen stehen */
    width: 240px; 
    height: 100vh; 
    background: #673ab7 !important; /* Lila */
    z-index: 100000;
    transition: right 0.3s cubic-bezier(0.25, 1, 0.5, 1), box-shadow 0.3s ease;
    box-shadow: none;
  }

  /* DIE UMRUNDUNG (Wölbung links neben dem Balken auf Button-Höhe) */
  .button-umrundung {
    position: absolute !important;
    top: 14px;      
    left: -18px; /* Klebt dezent am Balken */
    width: 35px;    
    height: 48px;   
    background: #673ab7 !important; 
    border-radius: 12px 0 0 12px; 
    z-index: 100001; /* Zwischen Wand und Icon */
    transition: opacity 0.2s ease;
  }

  /* WENN DIE CHECKBOX GEKLICKT IST (Menü öffnet sich) */
  .menu-checkbox:checked ~ .nav-dropdown-balken {
    right: 0 !important; 
    box-shadow: -8px 0 25px rgba(0, 0, 0, 0.4) !important; 
  }

  /* Blendet die Ausbuchtung aus, wenn die Wand ganz da ist */
  .menu-checkbox:checked ~ .nav-dropdown-balken .button-umrundung {
    opacity: 0 !important; 
  }

  /* Das Icon wandert mit der Wand nach links */
  .menu-checkbox:checked ~ .nav-burger-btn {
    right: 185px !important; 
  }

  /* DIE LINKS IM MENÜ (Mit Flexbox für das Impressum) */
  .nav-links-wrapper {
    display: flex !important;
    flex-direction: column !important;
    height: 100%; /* Füllt die gesamte Höhe des Balkens aus */
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

  /* Drückt das Impressum automatisch ganz an das untere Ende */
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

  /* Hauptinhalt der Seite */
  .page-content {
    padding: 20px;
  }
</style>