<script>
  import { goto } from '$app/navigation';
  import { warenkorb } from '$lib/stores/cart.js';
  import { holeBestellungen } from '$lib/api/bestellung.js';

  // 🧾 BESTELLVERLAUF
  // Liest die früheren Bestellungen aus dem localStorage (werden beim Checkout
  // gespeichert) und zeigt sie an.

  let bestellungen = $state([]);
  let offen = $state({}); // welche Bestellungen sind aufgeklappt?
  let druckBestellung = $state(null); // Bestellung, die gerade gedruckt wird

  // 🧾 Rechnung als PDF/Druck für eine bestimmte Bestellung.
  function rechnungDrucken(bestellung) {
    druckBestellung = bestellung;
    // Kurz warten, damit die Druck-Ansicht gerendert ist, dann drucken.
    setTimeout(() => window.print(), 50);
  }

  // Hinweis: onMount feuert in diesem Setup gerade nicht zuverlässig -> $effect
  // (läuft sicher beim ersten Rendern). initGeladen sorgt dafür, dass das nur 1x passiert.
  let initGeladen = false;
  $effect(() => {
    if (initGeladen) return;
    initGeladen = true;

    // Erst die lokale Kopie zeigen (schnell + offline).
    bestellungen = JSON.parse(localStorage.getItem('lieferino_bestellungen') || '[]');

    // 🗄️ Falls eingeloggt: echte Bestellungen aus dem Backend (DB) nachladen.
    (async () => {
      const res = await holeBestellungen();
      if (res.ok && Array.isArray(res.daten)) {
        // Backend liefert "positionen"; die Anzeige nutzt "artikel".
        bestellungen = res.daten.map((o) => ({ ...o, artikel: o.positionen || [] }));
      }
    })();
  });

  // 🔁 "Nochmal bestellen": legt die Artikel der alten Bestellung neu in den
  // Warenkorb und leitet zum Warenkorb weiter.
  function nochmalBestellen(bestellung) {
    warenkorb.set(bestellung.artikel.map((a) => ({ ...a })));
    goto('/cart');
  }

  // Wandelt das gespeicherte ISO-Datum in ein schönes deutsches Format um.
  function formatiere(iso) {
    return new Date(iso).toLocaleString('de-DE', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
</script>

<div class="seite">
  <h1>🧾 Deine Bestellungen</h1>

  {#if bestellungen.length === 0}
    <div class="leer">
      <p>Du hast noch nichts bestellt.</p>
      <a href="/restaurants" class="btn">Jetzt etwas bestellen</a>
    </div>
  {:else}
    <div class="liste">
      {#each bestellungen as bestellung}
        <div class="bestellung">
          <div class="kopf">
            <div>
              {#if bestellung.nummer}<span class="nummer">{bestellung.nummer}</span>{/if}
              <span class="datum">📅 {formatiere(bestellung.datum)}</span>
            </div>
            <span class="summe">{bestellung.summe.toFixed(2)}€</span>
          </div>

          <!-- 🔎 Detailansicht ein-/ausklappen -->
          <button class="details-toggle" onclick={() => (offen[bestellung.nummer || bestellung.datum] = !offen[bestellung.nummer || bestellung.datum])}>
            {offen[bestellung.nummer || bestellung.datum] ? '▲ Details ausblenden' : '▼ Details anzeigen'}
          </button>

          {#if offen[bestellung.nummer || bestellung.datum]}
            <!-- Volle Aufschlüsselung -->
            <div class="detail">
              {#each bestellung.artikel as a}
                <div class="detail-zeile">
                  <span>{a.menge}× {a.name}</span>
                  <span>{(a.preis * a.menge).toFixed(2)}€</span>
                </div>
              {/each}
              {#if bestellung.trinkgeld > 0}
                <div class="detail-zeile"><span>Trinkgeld</span><span>{bestellung.trinkgeld.toFixed(2)}€</span></div>
              {/if}
              {#if bestellung.gutschein}
                <div class="detail-zeile"><span>Gutschein</span><span>{bestellung.gutschein}</span></div>
              {/if}
              <div class="detail-zeile"><span>Liefertermin</span><span>{bestellung.liefertermin} Uhr</span></div>
            </div>
          {:else}
            <div class="artikel">
              {#each bestellung.artikel as a}
                <span class="chip">{a.menge}× {a.name}</span>
              {/each}
            </div>
          {/if}

          <div class="fuss">
            <span class="zahlung">Bezahlt mit: {bestellung.zahlungsart}</span>
            <div class="fuss-btns">
              <a href="/tracking" class="btn klein sekundaer">🚚 Verfolgen</a>
              <button class="btn klein sekundaer" onclick={() => rechnungDrucken(bestellung)}>🧾 Rechnung</button>
              <button class="btn klein" onclick={() => nochmalBestellen(bestellung)}>🔁 Nochmal bestellen</button>
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <!-- 🧾 Druckbare Rechnung (nur beim Drucken/PDF sichtbar) -->
  {#if druckBestellung}
    <div class="rechnung-druck">
      <h2>Lieferino – Rechnung</h2>
      <p>Bestellnummer: {druckBestellung.nummer || '—'}<br />
        Datum: {formatiere(druckBestellung.datum)}<br />
        Liefertermin: {druckBestellung.liefertermin} Uhr</p>
      <table>
        <tbody>
          {#each druckBestellung.artikel as a}
            <tr><td>{a.menge}× {a.name}</td><td>{(a.preis * a.menge).toFixed(2)}€</td></tr>
          {/each}
          {#if druckBestellung.trinkgeld > 0}<tr><td>Trinkgeld</td><td>{druckBestellung.trinkgeld.toFixed(2)}€</td></tr>{/if}
          <tr class="r-gesamt"><td>Gesamt</td><td>{druckBestellung.summe.toFixed(2)}€</td></tr>
        </tbody>
      </table>
      <p>Bezahlt mit: {druckBestellung.zahlungsart}</p>
      <p>Vielen Dank für deine Bestellung! 🍕</p>
    </div>
  {/if}
</div>

<style>
  /* 🧾 Druckbare Rechnung: am Bildschirm versteckt, nur beim Drucken sichtbar */
  .rechnung-druck { display: none; }
  @media print {
    :global(body *) { visibility: hidden !important; }
    .rechnung-druck, .rechnung-druck * { visibility: visible !important; }
    .rechnung-druck { display: block; position: absolute; top: 0; left: 0; width: 100%; padding: 20px; font-family: sans-serif; color: #000; }
    .rechnung-druck table { width: 100%; border-collapse: collapse; margin: 16px 0; }
    .rechnung-druck td { padding: 4px 0; border-bottom: 1px solid #eee; }
    .rechnung-druck .r-gesamt td { font-weight: bold; border-top: 2px solid #000; border-bottom: none; }
  }

  .fuss-btns { display: flex; gap: 8px; flex-wrap: wrap; }
  .btn.sekundaer { background: #f1f1f1; color: #333; }
  .seite { max-width: 760px; margin: 0 auto; padding: 20px; font-family: sans-serif; }
  h1 { margin-bottom: 20px; }

  .leer { text-align: center; padding: 40px; color: #777; }

  .liste { display: flex; flex-direction: column; gap: 14px; }
  .bestellung { background: white; border: 1px solid #eee; border-radius: 16px; padding: 18px; }
  .kopf { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
  .datum { color: #777; font-size: 0.9rem; }
  .nummer { display: inline-block; background: #f3e5f5; color: #673ab7; font-weight: bold; padding: 2px 10px; border-radius: 10px; font-size: 0.8rem; margin-right: 8px; letter-spacing: 1px; }
  .details-toggle { background: none; border: none; color: #673ab7; cursor: pointer; font-size: 0.85rem; padding: 0; margin-bottom: 10px; }
  .detail { background: #faf7ff; border-radius: 10px; padding: 12px; margin-bottom: 12px; }
  .detail-zeile { display: flex; justify-content: space-between; font-size: 0.88rem; color: #555; padding: 3px 0; }
  .summe { font-weight: bold; font-size: 1.1rem; color: #673ab7; }

  .artikel { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 14px; }
  .chip { background: #f3e5f5; color: #673ab7; padding: 4px 12px; border-radius: 20px; font-size: 0.85rem; }

  .fuss { display: flex; justify-content: space-between; align-items: center; }
  .zahlung { color: #999; font-size: 0.85rem; text-transform: capitalize; }

  .btn { display: inline-block; background: #673ab7; color: white; text-decoration: none; padding: 14px 22px; border-radius: 12px; font-weight: bold; border: none; cursor: pointer; }
  .btn.klein { padding: 10px 16px; font-size: 0.9rem; }
</style>