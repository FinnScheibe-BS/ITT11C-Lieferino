export function berechneDurchschnitt(reviews, fallback = 0) {
  if (!reviews || reviews.length === 0) return fallback;

  const summe = reviews.reduce((gesamt, review) => gesamt + review.sterne, 0);
  return summe / reviews.length;
}