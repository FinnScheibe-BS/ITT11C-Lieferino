<script>
  import { onMount, onDestroy } from 'svelte';
  import { warenkorb, gesamtSumme, leereWarenkorb } from '$lib/stores/cart.js';
  import { sendeBestellBestaetigung } from '$lib/services/email.js';

  const LIEFERGEBUEHR = 2.49;

  // Zahlungsart-Auswahl
  let zahlungsart = $state('paypal');

  // Lieferadresse + E-Mail laden wir aus den gespeicherten Nutzerdaten.
  let user = $state(null);

  // Nach dem Abschicken zeigen wir eine Erfolgsmeldung statt des Formulars.
  let bestellt = $state(false);
  let liefertermin = $state('');

  // 🎟️ GUTSCHEINE: gültige Codes und was sie bewirken.
  const GUTSCHEINE = {
    LIEFERINO10: { typ: 'prozent', wert: 10, text: '10% Rabatt' },
    WILLKOMMEN5: { typ: 'betrag', wert: 5, text: '5€ Rabatt' },
    GRATIS: { typ: 'lieferung', wert: 0, text: 'Gratis Lieferung' }
  };
  let gutscheinCode = $state('');
  let aktiverGutschein = $state(null); // der eingelöste Gutschein (oder null)
  let gutscheinFehler = $state('');

  function gutscheinEinloesen() {
    const code = gutscheinCode.trim().toUpperCase();
    if (GUTSCHEINE[code]) {
      aktiverGutschein = { code, ...GUTSCHEINE[code] };
      gutscheinFehler = '';
    } else {
      aktiverGutschein = null;
      gutscheinFehler = 'Dieser Gutscheincode ist leider ungültig. 🙁';
    }
  }

  // Liefergebühr: 0€, wenn ein "Gratis Lieferung"-Gutschein aktiv ist.
  let lieferkosten = $derived(aktiverGutschein?.typ === 'lieferung' ? 0 : LIEFERGEBUEHR);

  // Rabatt-Betrag je nach Gutschein-Typ.
  let rabatt = $derived.by(() => {
    if (!aktiverGutschein) return 0;
    if (aktiverGutschein.typ === 'prozent') return ($gesamtSumme * aktiverGutschein.wert) / 100;
    if (aktiverGutschein.typ === 'betrag') return Math.min(aktiverGutschein.wert, $gesamtSumme);
    return 0;
  });

  // Endsumme = Zwischensumme − Rabatt + Lieferkosten (nie unter 0).
  let endsumme = $derived(Math.max(0, $gesamtSumme - rabatt + lieferkosten));

  // 🚚 LIVE-LIEFERSTATUS: Die Bestellung durchläuft mehrere Phasen.
  const STATUS_PHASEN = ['Bestellung erhalten', 'Wird zubereitet', 'Unterwegs', 'Geliefert'];
  let statusIndex = $state(0);
  let statusTimer;

  onMount(() => {
    const gespeichert = localStorage.getItem('lieferino_user');
    if (gespeichert) user = JSON.parse(gespeichert);
  });

  onDestroy(() => clearInterval(statusTimer));

  // Lässt den Status alle 3 Sekunden eine Phase weiterspringen.
  function starteLieferstatus() {
    statusIndex = 0;
    statusTimer = setInterval(() => {
      if (statusIndex < STATUS_PHASEN.length - 1) {
        statusIndex += 1;
      } else {
        clearInterval(statusTimer);
      }
    }, 3000);
  }

  // Berechnet einen voraussichtlichen Liefertermin: jetzt + ca. 40 Minuten.
  function berechneLiefertermin() {
    const jetzt = new Date();
    jetzt.setMinutes(jetzt.getMinutes() + 40);
    // Schöne Uhrzeit im Format HH:MM.
    return jetzt.toLocaleTimeString('de-DE', { hour: '2-digit', minute: '2-digit' });
  }

  async function bestellungAbschicken() {
    liefertermin = berechneLiefertermin();

    // 🚨 BACKEND-HINWEIS: Diese Bestellung wird aktuell NICHT an einen Server
    // geschickt, sondern nur lokal verarbeitet. Das Backend muss hier später
    // einen Endpunkt anbieten, z.B.  POST /api/bestellungen
    //   Body: { kunde, artikel, summe, zahlungsart, liefertermin }
    // der die Bestellung in der Datenbank speichert und ans Restaurant meldet.

    // Bestellung in den lokalen Bestellverlauf legen (für die Verlauf-Seite).
    const bestellung = {
      datum: new Date().toISOString(),
      artikel: $warenkorb,
      summe: endsumme,
      gutschein: aktiverGutschein?.code || null,
      zahlungsart,
      liefertermin
    };
    const verlauf = JSON.parse(localStorage.getItem('lieferino_bestellungen') || '[]');
    verlauf.unshift(bestellung);
    localStorage.setItem('lieferino_bestellungen', JSON.stringify(verlauf));

    // E-Mail-Bestätigung anstoßen (echter Versand erfolgt später übers Backend).
    if (user?.email) {
      await sendeBestellBestaetigung(user.email, bestellung, liefertermin);
    }

    leereWarenkorb();
    bestellt = true;
    starteLieferstatus(); // 🚚 Live-Status starten
  }
</script>

<div class="seite">
  {#if bestellt}
    <!-- ✅ Erfolgsansicht nach der Bestellung -->
    <div class="erfolg">
      <h1>🎉 Bestellung erfolgreich!</h1>
      <p>Vielen Dank für deine Bestellung bei Lieferino.</p>
      <p class="liefertermin">🚚 Voraussichtliche Lieferung bis <strong>{liefertermin} Uhr</strong></p>

      <!-- 🚚 LIVE-LIEFERSTATUS: zeigt die aktuelle Phase an -->
      <div class="status-tracker">
        {#each STATUS_PHASEN as phase, i}
          <div class="status-schritt" class:erreicht={i <= statusIndex} class:aktuell={i === statusIndex}>
            <div class="status-punkt">{i < statusIndex ? '✅' : i === statusIndex ? '🔄' : '⬜'}</div>
            <span>{phase}</span>
          </div>
        {/each}
      </div>

      {#if user?.email}
        <p class="mail-info">📧 Eine Bestätigung wurde an {user.email} geschickt.</p>
      {/if}
      <a href="/bestellungen" class="btn sekundaer">Meine Bestellungen</a>
      <a href="/" class="btn">Zurück zur Startseite</a>
    </div>
  {:else if $warenkorb.length === 0}
    <!-- Kein Inhalt zum Bestellen -->
    <div class="leer">
      <h1>🧾 Kasse</h1>
      <p>Dein Warenkorb ist leer.</p>
      <a href="/restaurants" class="btn">Restaurants entdecken</a>
    </div>
  {:else}
    <h1>🧾 Kasse</h1>

    <!-- 🏠 Lieferadresse -->
    <section class="block">
      <h2>Lieferadresse</h2>
      {#if user}
        <p>
          {user.vorname} {user.nachname}<br />
          {user.strasse} {user.hausnummer}<br />
          {user.plz} {user.ort}
        </p>
      {:else}
        <p class="warnung">⚠️ Keine Adresse gefunden. Bitte zuerst registrieren.</p>
      {/if}
    </section>

    <!-- 💳 Zahlungsart -->
    <section class="block">
      <h2>Zahlungsart</h2>
      <div class="zahlung-optionen">
        <label class:aktiv={zahlungsart === 'paypal'}>
          <input type="radio" bind:group={zahlungsart} value="paypal" /> 🅿️ PayPal
        </label>
        <label class:aktiv={zahlungsart === 'karte'}>
          <input type="radio" bind:group={zahlungsart} value="karte" /> 💳 Kreditkarte
        </label>
        <label class:aktiv={zahlungsart === 'bar'}>
          <input type="radio" bind:group={zahlungsart} value="bar" /> 💶 Barzahlung
        </label>
      </div>
    </section>

    <!-- 🎟️ Gutscheincode -->
    <section class="block">
      <h2>Gutscheincode</h2>
      <div class="gutschein-zeile">
        <input type="text" placeholder="z.B. LIEFERINO10" bind:value={gutscheinCode} />
        <button class="einloesen" onclick={gutscheinEinloesen}>Einlösen</button>
      </div>
      {#if aktiverGutschein}
        <p class="gutschein-ok">✅ „{aktiverGutschein.code}" aktiv: {aktiverGutschein.text}</p>
      {:else if gutscheinFehler}
        <p class="warnung">{gutscheinFehler}</p>
      {/if}
      <p class="gutschein-tipp">💡 Versuch's mal mit <strong>LIEFERINO10</strong>, <strong>WILLKOMMEN5</strong> oder <strong>GRATIS</strong>.</p>
    </section>

    <!-- 📦 Bestellübersicht -->
    <section class="block">
      <h2>Deine Bestellung</h2>
      {#each $warenkorb as artikel}
        <div class="uebersicht-zeile">
          <span>{artikel.menge}× {artikel.name}</span>
          <span>{(artikel.preis * artikel.menge).toFixed(2)}€</span>
        </div>
      {/each}
      <div class="uebersicht-zeile">
        <span>Zwischensumme</span>
        <span>{$gesamtSumme.toFixed(2)}€</span>
      </div>
      <!-- Rabatt nur anzeigen, wenn ein Gutschein greift -->
      {#if rabatt > 0}
        <div class="uebersicht-zeile rabatt">
          <span>Rabatt ({aktiverGutschein.code})</span>
          <span>−{rabatt.toFixed(2)}€</span>
        </div>
      {/if}
      <div class="uebersicht-zeile">
        <span>Liefergebühr</span>
        <span>{lieferkosten === 0 ? 'Gratis 🎉' : lieferkosten.toFixed(2) + '€'}</span>
      </div>
      <div class="uebersicht-zeile gesamt">
        <span>Gesamt</span>
        <span>{endsumme.toFixed(2)}€</span>
      </div>
    </section>

    <button class="btn voll" onclick={bestellungAbschicken} disabled={!user}>
      Jetzt kostenpflichtig bestellen ✅
    </button>
  {/if}
</div>

<style>
  .seite { max-width: 640px; margin: 0 auto; padding: 20px; font-family: sans-serif; }
  h1 { margin-bottom: 20px; }

  .block { background: white; border: 1px solid #eee; border-radius: 16px; padding: 20px; margin-bottom: 16px; }
  .block h2 { margin: 0 0 12px; font-size: 1.1rem; }

  .zahlung-optionen { display: flex; flex-direction: column; gap: 10px; }
  .zahlung-optionen label { border: 1px solid #ddd; border-radius: 12px; padding: 12px; cursor: pointer; }
  .zahlung-optionen label.aktiv { border-color: #673ab7; background: #faf7ff; }

  .uebersicht-zeile { display: flex; justify-content: space-between; margin-bottom: 8px; color: #555; }
  .uebersicht-zeile.gesamt { font-weight: bold; color: #222; border-top: 1px solid #ddd; padding-top: 10px; font-size: 1.15rem; }

  .warnung { color: #d97706; font-weight: 600; }

  .btn { display: inline-block; background: #673ab7; color: white; text-decoration: none; padding: 14px 22px; border-radius: 12px; font-weight: bold; border: none; cursor: pointer; font-size: 1rem; }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .voll { width: 100%; }

  .erfolg, .leer { text-align: center; padding: 40px 20px; }
  .erfolg .liefertermin { font-size: 1.1rem; margin: 16px 0; }
  .erfolg .mail-info { color: #777; }
  .erfolg .btn, .leer .btn { margin-top: 20px; }
  .btn.sekundaer { background: #f1f1f1; color: #333; margin-right: 10px; }

  /* 🎟️ Gutschein */
  .gutschein-zeile { display: flex; gap: 10px; }
  .gutschein-zeile input { flex: 1; padding: 11px; border: 1px solid #ddd; border-radius: 10px; }
  .einloesen { background: #673ab7; color: white; border: none; padding: 0 18px; border-radius: 10px; font-weight: bold; cursor: pointer; }
  .gutschein-ok { color: #34c759; font-weight: 600; margin: 10px 0 0; }
  .gutschein-tipp { color: #999; font-size: 0.82rem; margin: 8px 0 0; }
  .rabatt { color: #34c759; }

  /* 🚚 Live-Lieferstatus */
  .status-tracker { display: flex; justify-content: space-between; gap: 6px; max-width: 460px; margin: 24px auto; }
  .status-schritt { flex: 1; text-align: center; opacity: 0.4; transition: opacity 0.3s ease; }
  .status-schritt.erreicht { opacity: 1; }
  .status-schritt .status-punkt { font-size: 1.5rem; margin-bottom: 6px; }
  .status-schritt span { font-size: 0.78rem; display: block; }
  .status-schritt.aktuell span { font-weight: bold; color: #673ab7; }
</style>
