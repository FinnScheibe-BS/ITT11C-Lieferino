<script>
  import { warenkorb, gesamtSumme, aendereMenge, entferneArtikel } from '$lib/stores/cart.js';
  import { getRestaurant } from '$lib/data';

  // 🏠 Liefergebühr als feste Pauschale (kann später angepasst werden).
  const LIEFERGEBUEHR = 2.49;

  // Höchster Mindestbestellwert aller Restaurants, die im Warenkorb liegen.
  // (Wir nehmen den höchsten, damit alle Restaurants bedient werden können.)
  let mindestbestellwert = $derived.by(() => {
    let max = 0;
    for (const artikel of $warenkorb) {
      const r = getRestaurant(artikel.restaurant);
      if (r && r.minBestell > max) max = r.minBestell;
    }
    return max;
  });

  // Ist der Mindestbestellwert erreicht?
  let mindestErreicht = $derived($gesamtSumme >= mindestbestellwert);
</script>

<div class="seite">
  <h1>🛒 Dein Warenkorb</h1>

  {#if $warenkorb.length === 0}
    <!-- Leerer Warenkorb -->
    <div class="leer">
      <p>😋 Dein Warenkorb ist noch leer.</p>
      <a href="/restaurants" class="btn">Jetzt Restaurants entdecken</a>
    </div>
  {:else}
    <div class="liste">
      {#each $warenkorb as artikel, index}
        <div class="zeile">
          <div class="info">
            <h3>{artikel.name}</h3>
            <p class="herkunft">von {artikel.restaurant}</p>
            <span class="preis">{artikel.preis.toFixed(2)}€</span>
          </div>

          <!-- Mengen-Steuerung: - [Anzahl] + -->
          <div class="menge">
            <button onclick={() => aendereMenge(index, -1)} aria-label="Weniger">−</button>
            <span>{artikel.menge}</span>
            <button onclick={() => aendereMenge(index, 1)} aria-label="Mehr">+</button>
          </div>

          <!-- Zwischensumme pro Artikel (Preis × Menge) -->
          <div class="zwischensumme">{(artikel.preis * artikel.menge).toFixed(2)}€</div>

          <button class="loeschen" onclick={() => entferneArtikel(index)} aria-label="Entfernen">🗑️</button>
        </div>
      {/each}
    </div>

    <!-- Zusammenfassung -->
    <div class="zusammenfassung">
      <div class="rechnung-zeile">
        <span>Zwischensumme</span>
        <span>{$gesamtSumme.toFixed(2)}€</span>
      </div>
      <div class="rechnung-zeile">
        <span>Liefergebühr</span>
        <span>{LIEFERGEBUEHR.toFixed(2)}€</span>
      </div>
      <div class="rechnung-zeile gesamt">
        <span>Gesamt</span>
        <span>{($gesamtSumme + LIEFERGEBUEHR).toFixed(2)}€</span>
      </div>

      <!-- Hinweis, falls der Mindestbestellwert noch nicht erreicht ist -->
      {#if !mindestErreicht}
        <p class="warnung">
          ⚠️ Mindestbestellwert {mindestbestellwert.toFixed(2)}€ noch nicht erreicht
          (es fehlen {(mindestbestellwert - $gesamtSumme).toFixed(2)}€).
        </p>
      {/if}

      <!-- Weiter zur Kasse: nur aktiv, wenn der Mindestbestellwert erreicht ist -->
      <a
        href="/checkout"
        class="btn checkout-btn"
        class:gesperrt={!mindestErreicht}
        aria-disabled={!mindestErreicht}
      >
        Zur Kasse 🧾
      </a>
    </div>
  {/if}
</div>

<style>
  .seite { max-width: 760px; margin: 0 auto; padding: 20px; font-family: sans-serif; }
  h1 { margin-bottom: 20px; }

  .leer { text-align: center; padding: 40px; color: #777; }

  .liste { display: flex; flex-direction: column; gap: 10px; }
  .zeile { display: grid; grid-template-columns: 1fr auto auto auto; align-items: center; gap: 16px; background: white; border: 1px solid #eee; border-radius: 14px; padding: 14px; }
  .info h3 { margin: 0 0 2px; }
  .herkunft { color: #999; font-size: 0.8rem; margin: 0 0 4px; }
  .preis { color: #673ab7; font-weight: bold; font-size: 0.9rem; }

  .menge { display: flex; align-items: center; gap: 10px; }
  .menge button { width: 30px; height: 30px; border-radius: 50%; border: 1px solid #ddd; background: #f7f7f7; font-size: 1.1rem; cursor: pointer; }
  .menge span { min-width: 20px; text-align: center; font-weight: bold; }

  .zwischensumme { font-weight: bold; min-width: 70px; text-align: right; }
  .loeschen { background: none; border: none; cursor: pointer; font-size: 1.1rem; }

  .zusammenfassung { margin-top: 24px; background: #faf7ff; border: 1px solid #ede7f6; border-radius: 16px; padding: 20px; }
  .rechnung-zeile { display: flex; justify-content: space-between; margin-bottom: 8px; color: #555; }
  .rechnung-zeile.gesamt { font-size: 1.2rem; font-weight: bold; color: #222; border-top: 1px solid #ddd; padding-top: 10px; margin-top: 4px; }

  .warnung { color: #d97706; font-size: 0.9rem; font-weight: 600; margin: 12px 0; }

  .btn { display: inline-block; background: #673ab7; color: white; text-decoration: none; padding: 14px 22px; border-radius: 12px; font-weight: bold; text-align: center; }
  .checkout-btn { width: 100%; box-sizing: border-box; margin-top: 10px; }
  /* Gesperrter Button: ausgegraut und nicht klickbar */
  .gesperrt { opacity: 0.5; pointer-events: none; }
</style>
