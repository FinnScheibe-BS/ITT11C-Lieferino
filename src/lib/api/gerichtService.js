// $lib/api/gerichtService.js
const BASE_URL = 'http://localhost:8080/api';
//const BASE_URL = 'http://172.30.4.90:8080/api';
export async function holeGerichte(restaurantId) {
  try {
    const response = await fetch(`${BASE_URL}/gerichte`);
    if (!response.ok) throw new Error('Netzwerkfehler');
    
    const alleGerichte = await response.json();
    
    // Filtere Gerichte nach restaurant_id
    return alleGerichte.filter(g => g.restaurant_id === restaurantId);
  } catch (error) {
    console.error('Fehler beim Laden der Gerichte:', error);
    return [];
  }
}

export async function holeAlleGerichte() {
  try {
    const response = await fetch(`${BASE_URL}/gerichte`);
    if (!response.ok) throw new Error('Netzwerkfehler');
    return await response.json();
  } catch (error) {
    console.error('Fehler beim Laden aller Gerichte:', error);
    return [];
  }
}