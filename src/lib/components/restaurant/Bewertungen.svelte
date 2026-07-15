<script>
  import { bewertungen, bewertungHinzufuegen, ladeBewertungen } from '$lib/stores/bewertungen.js';
  import { api, getToken } from '$lib/api/api.js';

  let { restaurant, slug } = $props();

  let meineBestellungen = $state([]);
  let bestellungenGeladen = $state(false);

  $effect(() => {
    if (bestellungenGeladen) return;

    bestellungenGeladen = true;
    meineBestellungen = JSON.parse(localStorage.getItem('lieferino_bestellungen') || '[]');

    if (getToken()) {
      api('/api/orders').then((res) => {
        if (res.ok && Array.isArray(res.daten)) {
          meineBestellungen = res.daten.map((o) => ({ ...o, artikel: o.positionen || [] }));
        }
      });
    }
  });

  $effect(() => {
    if (slug) ladeBewertungen(slug);
  });

  let reviews = $derived($bewertungen[slug] || []);

  let hatHierBestellt = $derived(
    meineBestellungen.some((b) =>
      (b.artikel || []).some((a) => a.restaurant === restaurant?.name)
    )
  );

  let neuName = $state('');
  let neuSterne = $state(5);
  let neuText = $state('');
  let bewertungFehler = $state('');

  async function bewertungAbschicken(e) {
    e.preventDefault();

    if (!hatHierBestellt) return;
    if (neuName.trim() === '' || neuText.trim() === '') return;

    bewertungHinzufuegen(slug, {
      name: neuName,
      sterne: neuSterne,
      text: neuText
    });

    neuName = '';
    neuText = '';
    neuSterne = 5;
  }
</script>

<h2 id="bewertungen">⭐ Bewertungen ({reviews.length})</h2>

{#if hatHierBestellt}
  <form class="review-form" onsubmit={bewertungAbschicken}>
    <input type="text" placeholder="Dein Name" bind:value={neuName} required />

    <select bind:value={neuSterne}>
      <option value={5}>⭐⭐⭐⭐⭐ (5)</option>
      <option value={4}>⭐⭐⭐⭐ (4)</option>
      <option value={3}>⭐⭐⭐ (3)</option>
      <option value={2}>⭐⭐ (2)</option>
      <option value={1}>⭐ (1)</option>
    </select>

    <textarea placeholder="Wie war dein Essen?" bind:value={neuText} required></textarea>

    {#if bewertungFehler}
      <p class="review-fehler">{bewertungFehler}</p>
    {/if}

    <button type="submit">Bewertung abschicken</button>
  </form>
{:else}
  <p class="review-sperre">🔒 Du kannst dieses Restaurant bewerten, sobald du hier etwas bestellt hast.</p>
{/if}

{#if reviews.length === 0}
  <p class="keine-reviews">Noch keine Bewertungen – sei die/der Erste! 🌟</p>
{:else}
  <div class="reviews">
    {#each reviews as review}
      <article class="review">
        <div class="review-kopf">
          <strong>{review.name}</strong>
          <span>{'⭐'.repeat(review.sterne)}</span>
        </div>
        <p>{review.text}</p>
      </article>
    {/each}
  </div>
{/if}