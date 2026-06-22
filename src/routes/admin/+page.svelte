<script>
  import { onMount } from 'svelte';
  import { restaurants } from '$lib/data';

  // 🛠️ ADMIN-DASHBOARD
  // Zeigt eine Übersicht über Restaurants und die (lokal gespeicherten)
  // Bestellungen. Reine Anzeige – echtes Verwalten muss das Backend liefern.

  let bestellungen = $state([]);

  onMount(() => {
    bestellungen = JSON.parse(localStorage.getItem('lieferino_bestellungen') || '[]');
  });

  // Kennzahlen automatisch aus den Bestellungen berechnen.
  let umsatz = $derived(bestellungen.reduce((s, b) => s + b.summe, 0));
  let verkaufteArtikel = $derived(
    bestellungen.reduce((s, b) => s + b.artikel.reduce((a, art) => a + art.menge, 0), 0)
  );

  // 📊 Umsatz der letzten 7 Tage für das Balken-Diagramm.
  let umsatzProTag = $derived.by(() => {
    const tage = [];
    const heute = new Date();
    // Wir bauen 7 leere Tage (von vor 6 Tagen bis heute).
    for (let i = 6; i >= 0; i--) {
      const d = new Date(heute);
      d.setDate(heute.getDate() - i);
      const schluessel = d.toISOString().slice(0, 10); // YYYY-MM-DD
      tage.push({ schluessel, label: d.toLocaleDateString('de-DE', { weekday: 'short' }), summe: 0 });
    }
    // Bestellungen den Tagen zuordnen.
    for (const b of bestellungen) {
      const schluessel = (b.datum || '').slice(0, 10);
      const tag = tage.find((t) => t.schluessel === schluessel);
      if (tag) tag.summe += b.summe;
    }
    return tage;
  });

  // Höchster Tagesumsatz (für die Balkenhöhe), mind. 1 um Division durch 0 zu vermeiden.
  let maxTag = $derived(Math.max(1, ...umsatzProTag.map((t) => t.summe)));
</script>

<div class="seite">
  <h1>🛠️ Admin-Dashboard</h1>

  <!-- 🚨🚨🚨 HINWEIS FÜRS BACKEND-TEAM 🚨🚨🚨 -->
  <!--
    Diese Seite zeigt nur Daten an, die lokal im Browser liegen. Damit ein
    echtes Admin-Panel funktioniert, muss das Backend folgendes bereitstellen:
      - GET  /api/admin/bestellungen   (alle Bestellungen aller Kunden)
      - GET/POST/PUT/DELETE /api/admin/restaurants  (Restaurants verwalten)
      - GET/POST/PUT/DELETE /api/admin/gerichte      (Speisekarte verwalten)
    Außerdem braucht es eine Zugriffskontrolle (nur Admins dürfen hier rein).
  -->

  <!-- Kennzahlen -->
  <div class="kacheln">
    <div class="kachel">
      <span class="zahl">{restaurants.length}</span>
      <span class="label">Restaurants</span>
    </div>
    <div class="kachel">
      <span class="zahl">{bestellungen.length}</span>
      <span class="label">Bestellungen</span>
    </div>
    <div class="kachel">
      <span class="zahl">{verkaufteArtikel}</span>
      <span class="label">Verkaufte Artikel</span>
    </div>
    <div class="kachel">
      <span class="zahl">{umsatz.toFixed(2)}€</span>
      <span class="label">Umsatz</span>
    </div>
  </div>

  <!-- 📊 Umsatz-Diagramm (letzte 7 Tage) -->
  <h2>📊 Umsatz (letzte 7 Tage)</h2>
  <div class="chart">
    {#each umsatzProTag as tag}
      <div class="balken-spalte">
        <span class="balken-wert">{tag.summe > 0 ? tag.summe.toFixed(0) + '€' : ''}</span>
        <div class="balken" style="height: {(tag.summe / maxTag) * 100}%"></div>
        <span class="balken-label">{tag.label}</span>
      </div>
    {/each}
  </div>

  <!-- Restaurant-Tabelle -->
  <h2>🍽️ Restaurants</h2>
  <div class="tabelle">
    <div class="zeile kopf">
      <span>Restaurant</span>
      <span>Typ</span>
      <span>Bewertung</span>
      <span>Gerichte</span>
    </div>
    {#each restaurants as r}
      <div class="zeile">
        <span>{r.emoji} {r.name}</span>
        <span>{r.typ}</span>
        <span>⭐ {r.bewertung}</span>
        <span>{r.speisekarte.length}</span>
      </div>
    {/each}
  </div>
</div>

<style>
  .seite { max-width: 900px; margin: 0 auto; padding: 20px; font-family: sans-serif; }
  h1 { margin-bottom: 24px; }
  h2 { margin: 30px 0 14px; }

  .kacheln { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 16px; }
  .kachel { background: #673ab7; color: white; border-radius: 16px; padding: 24px; text-align: center; }
  .kachel .zahl { display: block; font-size: 2rem; font-weight: bold; }
  .kachel .label { opacity: 0.85; font-size: 0.9rem; }

  /* 📊 Balken-Diagramm */
  .chart { display: flex; align-items: flex-end; gap: 10px; height: 200px; background: white; border: 1px solid #eee; border-radius: 14px; padding: 20px; }
  .balken-spalte { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: flex-end; height: 100%; gap: 6px; }
  .balken-wert { font-size: 0.72rem; color: #673ab7; font-weight: bold; }
  .balken { width: 100%; max-width: 46px; background: linear-gradient(180deg, #7e57c2, #673ab7); border-radius: 8px 8px 0 0; min-height: 3px; transition: height 0.3s ease; }
  .balken-label { font-size: 0.78rem; color: #888; }

  .tabelle { background: white; border: 1px solid #eee; border-radius: 14px; overflow: hidden; }
  .zeile { display: grid; grid-template-columns: 2fr 1fr 1fr 1fr; padding: 14px 18px; border-bottom: 1px solid #f0f0f0; align-items: center; }
  .zeile:last-child { border-bottom: none; }
  .zeile.kopf { background: #faf7ff; font-weight: bold; }
</style>
