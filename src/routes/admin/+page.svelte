<script>
  import { onMount } from 'svelte';
  import { restaurants } from '$lib/data';

  // 🛠️ ADMIN-DASHBOARD
  // Zeigt eine Übersicht über Restaurants und die (lokal gespeicherten)
  // Bestellungen. Reine Anzeige – echtes Verwalten muss das Backend liefern.

  let bestellungen = $state([]);

  // 🔐 Sicherheitsschlüssel für den Admin-Bereich.
  // ⚠️ Steht hier absichtlich im Klartext (vom Team so gewünscht). Für einen
  // echten Schutz müsste die Prüfung im Backend mit einem gehashten Schlüssel
  // erfolgen – im Frontend kann man den Schlüssel im Quelltext auslesen.
  const ADMIN_KEY = 'Lieferino#2026_DevOps';

  let entsperrt = $state(false);
  let schluesselEingabe = $state('');
  let schluesselFehler = $state('');

  onMount(() => {
    bestellungen = JSON.parse(localStorage.getItem('lieferino_bestellungen') || '[]');
    // Schon in dieser Browser-Sitzung entsperrt?
    entsperrt = sessionStorage.getItem('lieferino_admin') === 'true';
  });

  function entsperren(e) {
    e.preventDefault();
    if (schluesselEingabe === ADMIN_KEY) {
      entsperrt = true;
      sessionStorage.setItem('lieferino_admin', 'true'); // bis zum Tab-Schließen merken
      schluesselFehler = '';
    } else {
      schluesselFehler = 'Falscher Sicherheitsschlüssel. 🔒';
    }
  }

  function adminSperren() {
    sessionStorage.removeItem('lieferino_admin');
    entsperrt = false;
    schluesselEingabe = '';
  }

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
  {#if !entsperrt}
    <!-- 🔐 Sperrbildschirm: Sicherheitsschlüssel verlangen -->
    <div class="sperre">
      <h1>🔐 Admin-Bereich</h1>
      <p>Dieser Bereich ist geschützt. Bitte gib den Sicherheitsschlüssel ein.</p>
      <form onsubmit={entsperren}>
        <input type="password" placeholder="Sicherheitsschlüssel" bind:value={schluesselEingabe} />
        {#if schluesselFehler}<p class="sperre-fehler">{schluesselFehler}</p>{/if}
        <button type="submit">Entsperren 🔓</button>
      </form>
    </div>
  {:else}
  <h1>🛠️ Admin-Dashboard</h1>
  <button class="sperren-btn" onclick={adminSperren}>🔒 Sperren</button>

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
  {/if}
</div>

<style>
  .seite { max-width: 900px; margin: 0 auto; padding: 20px; font-family: sans-serif; }

  /* 🔐 Sperrbildschirm */
  .sperre { max-width: 380px; margin: 60px auto; background: white; border: 1px solid #eee; border-radius: 20px; padding: 30px; text-align: center; box-shadow: 0 10px 30px rgba(0,0,0,0.06); }
  .sperre p { color: #777; font-size: 0.9rem; }
  .sperre form { display: flex; flex-direction: column; gap: 12px; margin-top: 16px; }
  .sperre input { padding: 12px; border: 1px solid #ddd; border-radius: 10px; font-size: 0.95rem; }
  .sperre button { padding: 13px; background: #673ab7; color: white; border: none; border-radius: 12px; font-weight: bold; cursor: pointer; }
  .sperre-fehler { color: #dc3545; font-weight: 600; font-size: 0.85rem; margin: 0; }
  .sperren-btn { float: right; background: #f1f1f1; border: none; border-radius: 10px; padding: 8px 14px; cursor: pointer; font-weight: 600; color: #555; }
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
