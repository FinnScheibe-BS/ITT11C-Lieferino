const API_BASE_URL = 'http://172.30.4.90:8080/api';

function generiereSlug(name) {
    return name.toLowerCase().replace(/[^a-z0-9 ]/g, '').trim().replace(/\s+/g, '-');
}

function holeEmoji(nationalitaet) {
    const nat = nationalitaet.toLowerCase();
    if (nat.includes('deutsch')) return '🍺';
    if (nat.includes('asiatisch') || nat.includes('fusion')) return '🥢';
    if (nat.includes('amerikanisch')) return '🍔';
    if (nat.includes('französisch')) return '🥐';
    if (nat.includes('italienisch')) return '🍕';
    return '🍽️';
}

export async function holeRestaurants() {
    try {
        const response = await fetch(`${API_BASE_URL}/restaurants`);
        if (!response.ok) throw new Error('Fehler beim Laden');
        const apiDaten = await response.json();
        
        return apiDaten.map(r => ({
            id: r.id,
            name: r.name,
            slug: generiereSlug(r.name),
            beschreibung: r.adresse,
            typ: r.nationalitaet.trim(),
            emoji: holeEmoji(r.nationalitaet),
            lieferzeit: "30-40 Min",
            minBestell: 10,
            bewertung: 4.5,
            speisekarte: []
        }));
    } catch (error) {
        console.error(error);
        return [];
    }
}