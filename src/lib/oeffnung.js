// 🕒 Prüft, ob ein Restaurant gerade geöffnet ist (anhand oeffnetUm/schliesstUm "HH:MM").
export function istGeoeffnet(restaurant) {
  if (!restaurant?.oeffnetUm || !restaurant?.schliesstUm) return true;
  const jetzt = new Date();
  const minutenJetzt = jetzt.getHours() * 60 + jetzt.getMinutes();
  const [oh, om] = restaurant.oeffnetUm.split(':').map(Number);
  const [sh, sm] = restaurant.schliesstUm.split(':').map(Number);
  return minutenJetzt >= oh * 60 + om && minutenJetzt < sh * 60 + sm;
}
