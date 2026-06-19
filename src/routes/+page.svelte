<div class="welcome-container">
  <h1>Willkommen bei Lieferino</h1>
  <p>Dein Lieblingsessen, nur wenige Klicks entfernt 🍕</p>

  <div class="restaurant-grid">
    
    <a href="/restaurants" class="restaurant-card">
      <div class="image-placeholder">
        🍕
      </div>
      <h3>Luigi's Pizzeria</h3>
      <p>Die beste Steinofenpizza der Stadt.</p>
    </a>

    <a href="/restaurants" class="restaurant-card">
      <div class="image-placeholder">
        🍔
      </div>
      <h3>Burger Heaven</h3>
      <p>Saftige Smash-Burger & frische Pommes.</p>
    </a>

    <a href="/restaurants" class="restaurant-card">
      <div class="image-placeholder">
        🍣
      </div>
      <h3>Sushi Sakura</h3>
      <p>Frisches Sushi und asiatische Spezialitäten.</p>
    </a>

    <a href="/restaurants" class="restaurant-card">
      <div class="image-placeholder">
        🌱
      </div>
      <h3>Green & Fresh</h3>
      <p>Leckere Bowls und gesunde Salate.</p>
    </a>

  </div>
</div>

<style>
  /* Haupt-Container zentrieren */
  .welcome-container {
    max-width: 1200px;
    margin: 0 auto;
    text-align: center;
    padding-top: 40px;
    font-family: sans-serif;
  }

  /* Überschrift "Mitte Oben" */
  h1 {
    font-size: 3rem;
    color: #333;
    margin-bottom: 10px;
    font-weight: 800;
  }

  p {
    font-size: 1.2rem;
    color: #666;
    margin-bottom: 50px;
  }

  /* Grid-Layout für die 4 Kästen */
  .restaurant-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 30px;
    padding: 0 20px;
  }

  /* Die einzelnen Kästen (Karten als Link) */
  .restaurant-card {
    background: white;
    border-radius: 16px;
    padding: 20px;
    box-shadow: 0 8px 20px rgba(0, 0, 0, 0.06);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
    border: 1px solid #eee;
    text-align: left; /* Text im Kasten linksbündig */
    display: block;   /* Macht den gesamten Link-Kasten klickbar */
    text-decoration: none; /* 👈 HIER: Verhindert hässliche Unterstreichungen */
    color: inherit;         /* 👈 HIER: Verhindert, dass die Schrift blau wird */
  }

  /* Hover-Effekt: Kasten hebt sich leicht an */
  .restaurant-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 12px 25px rgba(103, 58, 183, 0.15); /* Lila-Hauch beim Drüberschweben */
  }

  /* Platzhalter für die Restaurant-Bilder */
  .image-placeholder {
    width: 100%;
    height: 160px;
    background: #f3e5f5; 
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 3rem; 
    margin-bottom: 15px;
    overflow: hidden;
  }

  /* Automatische Anpassung für spätere echte Bilder */
  .image-placeholder img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .restaurant-card h3 {
    font-size: 1.4rem;
    color: #222;
    margin: 10px 0 5px 0;
  }

  .restaurant-card p {
    font-size: 0.95rem;
    color: #777;
    margin-bottom: 0;
  }
</style>