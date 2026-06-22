<script>
  import { onMount } from 'svelte';
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

  onMount(() => {
    const gespeichert = localStorage.getItem('lieferino_user');
    if (gespeichert) user = JSON.parse(gespeichert);
  });

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
      summe: $gesamtSumme + LIEFERGEBUEHR,
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
  }
</script>

<div class="seite">
  {#if bestellt}
    <!-- ✅ Erfolgsansicht nach der Bestellung -->
    <div class="erfolg">
      <h1>🎉 Bestellung erfolgreich!</h1>
      <p>Vielen Dank für deine Bestellung bei Lieferino.</p>
      <p class="liefertermin">🚚 Voraussichtliche Lieferung bis <strong>{liefertermin} Uhr</strong></p>
      {#if user?.email}
        <p class="mail-info">📧 Eine Bestätigung wurde an {user.email} geschickt.</p>
      {/if}
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
        <span>Liefergebühr</span>
        <span>{LIEFERGEBUEHR.toFixed(2)}€</span>
      </div>
      <div class="uebersicht-zeile gesamt">
        <span>Gesamt</span>
        <span>{($gesamtSumme + LIEFERGEBUEHR).toFixed(2)}€</span>
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
</style>
