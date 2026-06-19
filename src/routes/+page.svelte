<script>
  // Eure Restaurant-Datenbank mit deinem Bildlink für Luigi's Pizzeria
  let restaurants = [
    { 
      name: "Luigi's Pizzeria", 
      typ: "pasta", 
      preis: 12, 
      bild: "https://th.bing.com/th/id/OIP.HfHopD9EcXAjSi7G-1M9BgHaFj?w=210&h=180&c=7&r=0&o=7&pid=1.7&rm=3", 
      beschreibung: "Steinofenpizza & Pasta" 
    },
    { 
      name: "Burger Heaven", 
      typ: "burger", 
      preis: 15, 
      bild: "🍔", // Hier könnt ihr später einen weiteren Bildlink oder Pfad eintragen
      beschreibung: "Premium Burger & Fries" 
    },
    { 
      name: "Sushi Sakura", 
      typ: "sushi", 
      preis: 22, 
      bild: "🍣", 
      beschreibung: "Frische Sushi-Platten" 
    },
    { 
      name: "Green & Fresh", 
      typ: "salat", 
      preis: 9, 
      bild: "🌱", 
      beschreibung: "Gesunde Salat-Bowls" 
    },
    { 
      name: "Mamma Mia", 
      typ: "pasta", 
      preis: 14, 
      bild: "🍝", 
      beschreibung: "Echte italienische Pasta" 
    },
    { 
      name: "Burger Kingz", 
      typ: "burger", 
      preis: 10, 
      bild: "🍟", 
      beschreibung: "Günstige, schnelle Burger" 
    }
  ];

  let gewaehlterTyp = $state("alle"); 
  let sortierung = $state("standard"); 

  let gefilterteRestaurants = $derived(
    restaurants
      .filter(r => gewaehlterTyp === "alle" || r.typ === gewaehlterTyp)
      .sort((a, b) => {
        if (sortierung === "preis-auf") return a.preis - b.preis;
        if (sortierung === "preis-ab") return b.preis - a.preis;
        return 0; 
      })
  );
</script>

<div class="page-container">
  <h2>Finde dein Lieblingsessen</h2>

  <div class="controls">
    <div class="control-group">
      <label for="filter">Kategorie:</label>
      <select id="filter" bind:value={gewaehlterTyp}>
        <option value="alle">🍕 Alles anzeigen</option>
        <option value="burger">🍔 Burger</option>
        <option value="pasta">🍝 Pasta</option>
        <option value="sushi">🍣 Sushi</option>
      </select>
    </div>

    <div class="control-group">
      <label for="sort">Sortieren nach:</label>
      <select id="sort" bind:value={sortierung}>
        <option value="standard">Standard</option>
        <option value="preis-auf">Preis: Günstig zuerst</option>
        <option value="preis-ab">Preis: Teuer zuerst</option>
      </select>
    </div>
  </div>

  <div class="restaurant-list">
    {#each gefilterteRestaurants as restaurant}
      <div class="restaurant-row">
        <div class="restaurant-media">
          {#if restaurant.bild.startsWith('http')}
            <img src={restaurant.bild} alt={restaurant.name} class="restaurant-img" />
          {:else}
            <span class="emoji">{restaurant.bild}</span>
          {/if}
        </div>
        
        <div class="details">
          <h3>{restaurant.name}</h3>
          <p>{restaurant.beschreibung}</p>
        </div>
        <div class="preis-tag">
          Ab {restaurant.preis}€
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .page-container {
    max-width: 800px;
    margin: 40px auto;
    font-family: sans-serif;
    padding: 0 20px;
  }

  h2 {
    color: #333;
    margin-bottom: 25px;
  }

  .controls {
    display: flex;
    gap: 20px;
    background: #f5f5f5;
    padding: 15px;
    border-radius: 12px;
    margin-bottom: 30px;
  }

  .control-group {
    display: flex;
    flex-direction: column;
    gap: 5px;
    flex: 1;
  }

  select {
    padding: 10px;
    border-radius: 8px;
    border: 1px solid #ccc;
    font-size: 1rem;
    background: white;
  }

  .restaurant-list {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }

  .restaurant-row {
    display: flex;
    align-items: center;
    background: white;
    padding: 20px;
    border-radius: 12px;
    border: 1px solid #eee;
    box-shadow: 0 4px 12px rgba(0,0,0,0.02);
  }

  .restaurant-media {
    width: 70px;
    height: 70px;
    margin-right: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .restaurant-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 10px;
  }

  .emoji {
    font-size: 2.5rem;
  }

  .details {
    flex: 1;
  }

  .details h3 {
    margin: 0 0 5px 0;
    font-size: 1.2rem;
  }

  .details p {
    margin: 0;
    color: #666;
    font-size: 0.9rem;
  }

  .preis-tag {
    font-weight: bold;
    background: #673ab7;
    color: white;
    padding: 8px 14px;
    border-radius: 20px;
    font-size: 0.9rem;
  }
</style>