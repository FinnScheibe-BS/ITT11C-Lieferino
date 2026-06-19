<script>
  // Eure Restaurant-Datenbank
  let restaurants = [
    { name: "Luigi's Pizzeria", typ: "pasta", preis: 12, emoji: "🍕", beschreibung: "Steinofenpizza & Pasta" },
    { name: "Burger Heaven", typ: "burger", preis: 15, emoji: "🍔", beschreibung: "Premium Burger & Fries" },
    { name: "Sushi Sakura", typ: "sushi", preis: 22, emoji: "🍣", beschreibung: "Frische Sushi-Platten" },
    { name: "Green & Fresh", typ: "salat", preis: 9, emoji: "🌱", beschreibung: "Gesunde Salat-Bowls" },
    { name: "Mamma Mia", typ: "pasta", preis: 14, emoji: "🍝", beschreibung: "Echte italienische Pasta" },
    { name: "Burger Kingz", typ: "burger", preis: 10, emoji: "🍟", beschreibung: "Günstige, schnelle Burger" }
  ];

  // 1. Variablen, die sich ändern können, bekommen ein $state()
  let gewaehlterTyp = $state("alle"); 
  let sortierung = $state("standard"); 

  // 2. Das alte "$:" wird in Svelte 5 durch ein stabiles $derived() ersetzt
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