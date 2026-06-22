// =====================================================================
// 🍽️ ZENTRALE DATENQUELLE (Single Source of Truth)
// Hier stehen ALLE Restaurants samt Speisekarte an EINER Stelle.
// Früher gab es mehrere getrennte Listen – die sind jetzt zusammengeführt.
//
// WICHTIG: Diese Datei wird NICHT direkt von der App geladen, sondern dient
// als Vorlage. Wenn du hier etwas änderst, danach einmal den Generator laufen
// lassen:   node scripts/build-bundle.mjs
// =====================================================================

export const restaurants = [
  {
    slug: 'luigis-pizzeria',
    name: "Luigi's Pizzeria",
    typ: 'italienisch',
    emoji: '🍕',
    bewertung: 4.9,
    lieferzeit: '25-35 min',
    minBestell: 10,
    beschreibung: 'Steinofenpizza & frische Pasta aus Italien',
    speisekarte: [
      { id: 1, name: 'Pizza Margherita', preis: 8.5, beschreibung: 'Tomate, Mozzarella, Basilikum' },
      { id: 2, name: 'Pizza Salami', preis: 9.5, beschreibung: 'Tomate, Mozzarella, Salami' },
      { id: 3, name: 'Pasta Carbonara', preis: 10.0, beschreibung: 'Speck, Ei, Parmesan' },
      { id: 4, name: 'Tiramisu', preis: 4.5, beschreibung: 'Hausgemachtes Dessert' }
    ]
  },
  {
    slug: 'burger-heaven',
    name: 'Burger Heaven',
    typ: 'amerikanisch',
    emoji: '🍔',
    bewertung: 4.8,
    lieferzeit: '20-30 min',
    minBestell: 15,
    beschreibung: 'Premium Burger & knusprige Fries',
    speisekarte: [
      { id: 1, name: 'Classic Cheeseburger', preis: 9.0, beschreibung: 'Rindfleisch, Cheddar, Salat' },
      { id: 2, name: 'Bacon Burger', preis: 11.0, beschreibung: 'Doppelt Speck, BBQ-Sauce' },
      { id: 3, name: 'Pommes groß', preis: 4.0, beschreibung: 'Knusprig mit Meersalz' },
      { id: 4, name: 'Milkshake', preis: 4.5, beschreibung: 'Vanille oder Schoko' }
    ]
  },
  {
    slug: 'sushi-sakura',
    name: 'Sushi Sakura',
    typ: 'asiatisch',
    emoji: '🍣',
    bewertung: 4.7,
    lieferzeit: '30-40 min',
    minBestell: 20,
    beschreibung: 'Frische Sushi-Platten & japanische Klassiker',
    speisekarte: [
      { id: 1, name: 'Sushi Set (12 St.)', preis: 14.0, beschreibung: 'Gemischte Auswahl' },
      { id: 2, name: 'California Roll', preis: 7.5, beschreibung: 'Surimi, Avocado, Gurke' },
      { id: 3, name: 'Miso Suppe', preis: 3.5, beschreibung: 'Mit Tofu und Algen' },
      { id: 4, name: 'Edamame', preis: 4.0, beschreibung: 'Gedämpfte Sojabohnen' }
    ]
  },
  {
    slug: 'taco-loco',
    name: 'Taco Loco',
    typ: 'mexikanisch',
    emoji: '🌮',
    bewertung: 4.6,
    lieferzeit: '20-30 min',
    minBestell: 12,
    beschreibung: 'Würzige Tacos, Burritos & Nachos',
    speisekarte: [
      { id: 1, name: 'Taco Trio', preis: 8.0, beschreibung: '3 Tacos nach Wahl' },
      { id: 2, name: 'Burrito Grande', preis: 9.5, beschreibung: 'Reis, Bohnen, Hähnchen' },
      { id: 3, name: 'Nachos mit Käse', preis: 6.0, beschreibung: 'Mit Jalapeños' },
      { id: 4, name: 'Guacamole', preis: 3.5, beschreibung: 'Frisch zubereitet' }
    ]
  },
  {
    slug: 'mamma-mia',
    name: 'Mamma Mia',
    typ: 'italienisch',
    emoji: '🍝',
    bewertung: 4.5,
    lieferzeit: '25-35 min',
    minBestell: 10,
    beschreibung: 'Echte italienische Pasta wie bei Nonna',
    speisekarte: [
      { id: 1, name: 'Spaghetti Bolognese', preis: 9.0, beschreibung: 'Klassische Hackfleischsauce' },
      { id: 2, name: 'Lasagne', preis: 10.5, beschreibung: 'Mit Béchamel überbacken' },
      { id: 3, name: 'Bruschetta', preis: 5.0, beschreibung: 'Geröstetes Brot mit Tomaten' }
    ]
  },
  {
    slug: 'el-sabor',
    name: 'El Sabor',
    typ: 'spanisch',
    emoji: '🥘',
    bewertung: 4.3,
    lieferzeit: '35-45 min',
    minBestell: 25,
    beschreibung: 'Tapas & Paella mit spanischem Flair',
    speisekarte: [
      { id: 1, name: 'Paella Mixta', preis: 13.0, beschreibung: 'Für 1 Person, Meeresfrüchte & Huhn' },
      { id: 2, name: 'Patatas Bravas', preis: 5.5, beschreibung: 'Mit scharfer Sauce' },
      { id: 3, name: 'Tortilla Española', preis: 6.0, beschreibung: 'Kartoffel-Omelett' }
    ]
  },
  {
    slug: 'curry-house',
    name: 'Curry House',
    typ: 'asiatisch',
    emoji: '🍛',
    bewertung: 4.2,
    lieferzeit: '30-40 min',
    minBestell: 15,
    beschreibung: 'Aromatische Currys aus Indien',
    speisekarte: [
      { id: 1, name: 'Chicken Tikka Masala', preis: 11.0, beschreibung: 'Cremiges Tomaten-Curry' },
      { id: 2, name: 'Veggie Korma', preis: 9.5, beschreibung: 'Mildes Gemüse-Curry' },
      { id: 3, name: 'Naan Brot', preis: 3.0, beschreibung: 'Frisch aus dem Ofen' }
    ]
  },
  {
    slug: 'green-and-fresh',
    name: 'Green & Fresh',
    typ: 'vegetarisch',
    emoji: '🥗',
    bewertung: 4.4,
    lieferzeit: '15-25 min',
    minBestell: 8,
    beschreibung: 'Gesunde Salat-Bowls & Smoothies',
    speisekarte: [
      { id: 1, name: 'Caesar Bowl', preis: 8.5, beschreibung: 'Salat, Croutons, Parmesan' },
      { id: 2, name: 'Falafel Bowl', preis: 9.0, beschreibung: 'Mit Hummus und Gemüse' },
      { id: 3, name: 'Mango Smoothie', preis: 4.5, beschreibung: 'Frisch püriert' }
    ]
  }
];
