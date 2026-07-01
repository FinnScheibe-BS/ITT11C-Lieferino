<script>
  import { onMount, onDestroy } from 'svelte';
  import { warenkorb, gesamtSumme, leereWarenkorb } from '$lib/stores/cart.js';
  import { pruefeKarte, kartenTyp, formatiereNummer } from '$lib/services/payment.js';
  import { treuepunkte, punkteSammeln, punkteAbziehen, EINLOESE_SCHRITT, EINLOESE_WERT } from '$lib/stores/treue.js';
  import { api, getToken } from '$lib/api.js';

  const LIEFERGEBUEHR = 2.49;

  // Zahlungsart-Auswahl
  let zahlungsart = $state('paypal');

  // 💳 Kreditkarten-Eingaben
  let kartennummer = $state('');
  let kartenAblauf = $state('');
  let kartenCvv = $state('');
  let kartePruefung = $state(null); // Ergebnis der letzten Prüfung (oder null)

  // Live-Kartentyp (Visa/Mastercard/…) für die Anzeige.
  let aktKartenTyp = $derived(kartennummer ? kartenTyp(kartennummer) : '');

  // Formatiert die Nummer beim Tippen in 4er-Blöcke.
  function nummerEingabe(e) {
    kartennummer = formatiereNummer(e.target.value);
    kartePruefung = null; // alte Fehlermeldung zurücksetzen
  }

  // 💶 Trinkgeld (in Prozent der Zwischensumme).
  let trinkgeldProzent = $state(0);
  let trinkgeld = $derived(($gesamtSumme * trinkgeldProzent) / 100);

  // Lieferadresse + E-Mail laden wir aus den gespeicherten Nutzerdaten.
  let user = $state(null);

  // 🏠 Adresswahl: Nutzer kann mehrere Adressen haben (user.adressen).
  // Fallback: die einzelnen Felder aus der Registrierung.
  let gewaehlteAdresse = $state(0);
  let adressen = $derived(
    user?.adressen?.length
      ? user.adressen
      : user
        ? [{ label: 'Standard', strasse: user.strasse, hausnummer: user.hausnummer, plz: user.plz, ort: user.ort }]
        : []
  );
  let aktiveAdresse = $derived(adressen[gewaehlteAdresse] || adressen[0] || null);

  // 🗺️ Karten-Vorschau der gewählten Lieferadresse.
  let kartenUrl = $derived(
    aktiveAdresse
      ? `https://maps.google.com/maps?q=${encodeURIComponent(
          `${aktiveAdresse.strasse} ${aktiveAdresse.hausnummer}, ${aktiveAdresse.plz} ${aktiveAdresse.ort}`
        )}&z=15&output=embed`
      : ''
  );

  // ⏰ Geplante Lieferung: 'asap' (so schnell wie möglich) oder 'geplant'.
  let lieferOption = $state('asap');
  let geplanteZeit = $state(''); // Format "HH:MM"

  // ⭐ Treuepunkte einlösen?
  let punkteEinloesen = $state(false);
  let punkteRabatt = $derived(
    punkteEinloesen && $treuepunkte >= EINLOESE_SCHRITT ? EINLOESE_WERT : 0
  );

  // Nach dem Abschicken zeigen wir eine Erfolgsmeldung statt des Formulars.
  let bestellt = $state(false);
  let liefertermin = $state('');

  // 🎟️ GUTSCHEINE: gültige Codes und was sie bewirken.
  const GUTSCHEINE = {
    LIEFERINO10: { typ: 'prozent', wert: 10, text: '10% Rabatt' },
    WILLKOMMEN5: { typ: 'betrag', wert: 5, text: '5€ Rabatt' },
    GRATIS: { typ: 'lieferung', wert: 0, text: 'Gratis Lieferung' },
    // 🥚 Versteckte Easter-Egg-Codes (über die Suche zu finden)
    PIZZAPARTY: { typ: 'prozent', wert: 25, text: '25% Rabatt 🍕' },
    GEBURTSTAG: { typ: 'prozent', wert: 20, text: '20% Geburtstags-Rabatt 🎂' }
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

  // Endsumme = Zwischensumme − Rabatt − Punkte-Rabatt + Lieferkosten + Trinkgeld (nie unter 0).
  let endsumme = $derived(Math.max(0, $gesamtSumme - rabatt - punkteRabatt + lieferkosten + trinkgeld));

  // Eindeutige Bestellnummer (wird beim Abschicken erzeugt).
  let bestellnummer = $state('');
  // Die abgeschlossene Bestellung (für die Rechnung – der Warenkorb wird ja geleert).
  let letzteBestellung = $state(null);

  // 🧾 Rechnung als PDF: öffnet den Druck-Dialog (dort „Als PDF speichern").
  function rechnungDrucken() {
    window.print();
  }

  // 🚚 LIVE-LIEFERSTATUS: Die Bestellung durchläuft mehrere Phasen.
  const STATUS_PHASEN = ['Bestellung erhalten', 'Wird zubereitet', 'Unterwegs', 'Geliefert'];
  let statusIndex = $state(0);
  let statusTimer;

  onMount(() => {
    const gespeichert = localStorage.getItem('lieferino_user');
    if (gespeichert) user = JSON.parse(gespeichert);
  });

  onDestroy(() => clearInterval(statusTimer));

  // 🔔 Zeigt eine Browser-Benachrichtigung (falls der Nutzer erlaubt hat).
  function benachrichtige(text) {
    if (typeof Notification === 'undefined') return; // Browser kann das nicht
    if (Notification.permission === 'granted') {
      new Notification('Lieferino 🛵', { body: text });
    }
  }

  // Lässt den Status alle 3 Sekunden eine Phase weiterspringen + benachrichtigt.
  function starteLieferstatus() {
    statusIndex = 0;
    statusTimer = setInterval(() => {
      if (statusIndex < STATUS_PHASEN.length - 1) {
        statusIndex += 1;
        benachrichtige(`Status: ${STATUS_PHASEN[statusIndex]}`); // 🔔 bei jedem Wechsel
      } else {
        clearInterval(statusTimer);
      }
    }, 3000);
  }

  // Berechnet den Liefertermin: bei "geplant" die gewählte Zeit, sonst jetzt + 40 Min.
  function berechneLiefertermin() {
    if (lieferOption === 'geplant' && geplanteZeit) {
      return geplanteZeit; // bereits Format "HH:MM"
    }
    const jetzt = new Date();
    jetzt.setMinutes(jetzt.getMinutes() + 40);
    return jetzt.toLocaleTimeString('de-DE', { hour: '2-digit', minute: '2-digit' });
  }

  async function bestellungAbschicken() {
    // 💳 Wenn mit Kreditkarte bezahlt wird, zuerst die Karte prüfen.
    if (zahlungsart === 'karte') {
      kartePruefung = pruefeKarte({ nummer: kartennummer, ablauf: kartenAblauf, cvv: kartenCvv });
      if (!kartePruefung.gueltig) {
        return; // Fehler werden im Formular angezeigt – nicht weiter.
      }
    }

    // ⏰ Bei geplanter Lieferung muss eine Zeit gewählt sein.
    if (lieferOption === 'geplant' && !geplanteZeit) {
      return;
    }

    liefertermin = berechneLiefertermin();

    // Bestellnummer erzeugen, z.B. "LF-7F3A9C".
    bestellnummer =
      'LF-' + Math.floor(Math.random() * 0xffffff).toString(16).toUpperCase().padStart(6, '0');

    // 🚨 BACKEND-HINWEIS: Diese Bestellung wird aktuell NICHT an einen Server
    // geschickt, sondern nur lokal verarbeitet. Das Backend muss hier später
    // einen Endpunkt anbieten, z.B.  POST /api/bestellungen
    //   Body: { nummer, kunde, artikel, summe, zahlungsart, trinkgeld, liefertermin }
    // der die Bestellung in der Datenbank speichert und ans Restaurant meldet.
    // 🚨 Die KARTENDATEN dürfen dabei NIEMALS mitgespeichert/gesendet werden –
    //     echte Zahlung läuft über einen Anbieter (siehe payment.js).

    // Bestellung in den lokalen Bestellverlauf legen (für die Verlauf-Seite).
    const bestellung = {
      nummer: bestellnummer,
      datum: new Date().toISOString(),
      artikel: $warenkorb,
      zwischensumme: $gesamtSumme,
      trinkgeld,
      summe: endsumme,
      gutschein: aktiverGutschein?.code || null,
      zahlungsart,
      liefertermin,
      geplant: lieferOption === 'geplant',
      adresse: aktiveAdresse
    };
    const verlauf = JSON.parse(localStorage.getItem('lieferino_bestellungen') || '[]');
    verlauf.unshift(bestellung);
    localStorage.setItem('lieferino_bestellungen', JSON.stringify(verlauf));

    // 🗄️ Bestellung im Backend (DB) speichern – damit nichts verloren geht.
    // Nur, wenn ein Token vorhanden ist (eingeloggt). Best-effort.
    if (getToken()) {
      await api('/api/orders', {
        method: 'POST',
        body: {
          summe: endsumme,
          zwischensumme: $gesamtSumme,
          trinkgeld,
          gutschein: aktiverGutschein?.code || '',
          zahlungsart,
          liefertermin,
          positionen: $warenkorb.map((a) => ({
            name: a.name, preis: a.preis, menge: a.menge, restaurant: a.restaurant
          }))
        }
      });
      // 🏠 Die verwendete Lieferadresse zusätzlich im Profil (DB) speichern.
      if (aktiveAdresse) {
        const meRes = await api('/api/me');
        if (meRes.ok && meRes.daten) {
          const adressen = meRes.daten.adressen || [];
          const schon = adressen.some(
            (a) => a.strasse === aktiveAdresse.strasse && a.plz === aktiveAdresse.plz && a.hausnummer === aktiveAdresse.hausnummer
          );
          if (!schon) {
            adressen.push({ label: aktiveAdresse.label || 'Lieferadresse', ...aktiveAdresse });
            await api('/api/me', {
              method: 'PUT',
              body: {
                username: meRes.daten.username, vorname: meRes.daten.vorname,
                nachname: meRes.daten.nachname, geburtsdatum: meRes.daten.geburtsdatum,
                adressen
              }
            });
          }
        }
      }
    }

    // Für die Rechnung merken (der Warenkorb wird gleich geleert).
    letzteBestellung = bestellung;

    // ⭐ Treuepunkte: ggf. eingelöste Punkte abziehen, dann neue Punkte sammeln.
    if (punkteRabatt > 0) punkteAbziehen(EINLOESE_SCHRITT);
    punkteSammeln(endsumme);

    // Die Bestellbestätigungs-E-Mail verschickt das Backend automatisch beim
    // Anlegen der Bestellung (siehe POST /api/orders) – hier nichts nötig.

    // 🔔 Um Liefer-Benachrichtigungen bitten (falls noch nicht entschieden).
    if (typeof Notification !== 'undefined' && Notification.permission === 'default') {
      Notification.requestPermission();
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
      <p class="bestellnr">Bestellnummer: <strong>{bestellnummer}</strong></p>
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
      <button class="btn sekundaer" onclick={rechnungDrucken}>🧾 Rechnung (PDF / Drucken)</button>
      <a href="/tracking" class="btn sekundaer">🚚 Bestellung verfolgen</a>
      <a href="/" class="btn">Zurück zur Startseite</a>
    </div>

    <!-- 🧾 Druckbare Rechnung (nur beim Drucken/PDF sichtbar) -->
    {#if letzteBestellung}
      <div class="rechnung-druck">
        <h2>Lieferino – Rechnung</h2>
        <p>Bestellnummer: {letzteBestellung.nummer}<br />
          Datum: {new Date(letzteBestellung.datum).toLocaleString('de-DE')}<br />
          Liefertermin: {letzteBestellung.liefertermin} Uhr</p>
        {#if letzteBestellung.adresse}
          <p>Lieferadresse:<br />
            {user?.vorname} {user?.nachname}<br />
            {letzteBestellung.adresse.strasse} {letzteBestellung.adresse.hausnummer}, {letzteBestellung.adresse.plz} {letzteBestellung.adresse.ort}</p>
        {/if}
        <table>
          <tbody>
            {#each letzteBestellung.artikel as a}
              <tr><td>{a.menge}× {a.name}</td><td>{(a.preis * a.menge).toFixed(2)}€</td></tr>
            {/each}
            {#if letzteBestellung.trinkgeld > 0}<tr><td>Trinkgeld</td><td>{letzteBestellung.trinkgeld.toFixed(2)}€</td></tr>{/if}
            <tr class="r-gesamt"><td>Gesamt</td><td>{letzteBestellung.summe.toFixed(2)}€</td></tr>
          </tbody>
        </table>
        <p>Bezahlt mit: {letzteBestellung.zahlungsart}</p>
        <p>Vielen Dank für deine Bestellung! 🍕</p>
      </div>
    {/if}
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
      {#if aktiveAdresse}
        <!-- Auswahl, falls mehrere Adressen vorhanden sind -->
        {#if adressen.length > 1}
          <select bind:value={gewaehlteAdresse} class="adress-wahl">
            {#each adressen as a, i}
              <option value={i}>{a.label ? a.label + ': ' : ''}{a.strasse} {a.hausnummer}, {a.plz} {a.ort}</option>
            {/each}
          </select>
        {/if}
        <p>
          {user.vorname} {user.nachname}<br />
          {aktiveAdresse.strasse} {aktiveAdresse.hausnummer}<br />
          {aktiveAdresse.plz} {aktiveAdresse.ort}
        </p>
        <!-- 🗺️ Karten-Vorschau der Lieferadresse -->
        <iframe class="karte-vorschau" src={kartenUrl} title="Lieferadresse auf der Karte" loading="lazy"></iframe>
      {:else}
        <p class="warnung">⚠️ Keine Adresse gefunden. Bitte zuerst registrieren.</p>
      {/if}
    </section>

    <!-- ⏰ Lieferzeit -->
    <section class="block">
      <h2>Lieferzeit</h2>
      <div class="zeit-optionen">
        <label class:aktiv={lieferOption === 'asap'}>
          <input type="radio" bind:group={lieferOption} value="asap" /> 🚀 So schnell wie möglich
        </label>
        <label class:aktiv={lieferOption === 'geplant'}>
          <input type="radio" bind:group={lieferOption} value="geplant" /> 🕒 Zu einer bestimmten Uhrzeit
        </label>
      </div>
      {#if lieferOption === 'geplant'}
        <input type="time" bind:value={geplanteZeit} class="zeit-feld" />
        {#if !geplanteZeit}<p class="warnung">Bitte eine Uhrzeit wählen.</p>{/if}
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

      <!-- 💳 Kreditkarten-Formular (nur bei Auswahl "Kreditkarte") -->
      {#if zahlungsart === 'karte'}
        <div class="karte-formular">
          <div class="kf-feld">
            <label for="kn">Kartennummer {#if aktKartenTyp && aktKartenTyp !== 'Unbekannt'}<span class="kartentyp">{aktKartenTyp}</span>{/if}</label>
            <input id="kn" type="text" inputmode="numeric" maxlength="23" placeholder="1234 5678 9012 3456" value={kartennummer} oninput={nummerEingabe} />
            {#if kartePruefung?.fehler.nummer}<span class="kf-fehler">Ungültige Kartennummer</span>{/if}
          </div>
          <div class="kf-row">
            <div class="kf-feld">
              <label for="ka">Gültig bis (MM/JJ)</label>
              <input id="ka" type="text" maxlength="5" placeholder="07/27" bind:value={kartenAblauf} />
              {#if kartePruefung?.fehler.ablauf}<span class="kf-fehler">Ungültig / abgelaufen</span>{/if}
            </div>
            <div class="kf-feld">
              <label for="cv">CVV</label>
              <input id="cv" type="text" inputmode="numeric" maxlength="4" placeholder="123" bind:value={kartenCvv} />
              {#if kartePruefung?.fehler.cvv}<span class="kf-fehler">Ungültig</span>{/if}
            </div>
          </div>
          <p class="kf-hinweis">🔒 Testkarte: 4242 4242 4242 4242, beliebiges zukünftiges Datum + 3 Ziffern.</p>
        </div>
      {/if}
    </section>

    <!-- 💶 Trinkgeld -->
    <section class="block">
      <h2>Trinkgeld für den Fahrer 🚴</h2>
      <div class="trinkgeld-optionen">
        {#each [0, 5, 10, 15] as prozent}
          <button
            type="button"
            class="tg-btn"
            class:aktiv={trinkgeldProzent === prozent}
            onclick={() => (trinkgeldProzent = prozent)}
          >
            {prozent === 0 ? 'Kein' : prozent + '%'}
          </button>
        {/each}
      </div>
      {#if trinkgeld > 0}
        <p class="tg-info">+ {trinkgeld.toFixed(2)}€ Trinkgeld – danke! 💜</p>
      {/if}
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

    <!-- ⭐ Treuepunkte -->
    <section class="block">
      <h2>Treuepunkte ⭐</h2>
      <p class="punkte-stand">Du hast <strong>{$treuepunkte}</strong> Punkte.</p>
      {#if $treuepunkte >= EINLOESE_SCHRITT}
        <label class="punkte-einloesen">
          <input type="checkbox" bind:checked={punkteEinloesen} />
          {EINLOESE_SCHRITT} Punkte einlösen (−{EINLOESE_WERT.toFixed(2)}€)
        </label>
      {:else}
        <p class="gutschein-tipp">Ab {EINLOESE_SCHRITT} Punkten kannst du {EINLOESE_WERT}€ Rabatt einlösen. (1 Punkt je 1€ Bestellwert)</p>
      {/if}
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
      {#if punkteRabatt > 0}
        <div class="uebersicht-zeile rabatt">
          <span>Treuepunkte ({EINLOESE_SCHRITT})</span>
          <span>−{punkteRabatt.toFixed(2)}€</span>
        </div>
      {/if}
      <div class="uebersicht-zeile">
        <span>Liefergebühr</span>
        <span>{lieferkosten === 0 ? 'Gratis 🎉' : lieferkosten.toFixed(2) + '€'}</span>
      </div>
      {#if trinkgeld > 0}
        <div class="uebersicht-zeile">
          <span>Trinkgeld ({trinkgeldProzent}%)</span>
          <span>{trinkgeld.toFixed(2)}€</span>
        </div>
      {/if}
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
  .karte-vorschau { width: 100%; height: 200px; border: 0; border-radius: 12px; margin-top: 10px; }

  .btn { display: inline-block; background: #673ab7; color: white; text-decoration: none; padding: 14px 22px; border-radius: 12px; font-weight: bold; border: none; cursor: pointer; font-size: 1rem; }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .voll { width: 100%; }

  /* 💳 Kreditkarten-Formular */
  .karte-formular { margin-top: 14px; display: flex; flex-direction: column; gap: 12px; }
  .kf-row { display: flex; gap: 12px; }
  .kf-feld { display: flex; flex-direction: column; gap: 4px; flex: 1; }
  .kf-feld label { font-size: 0.82rem; font-weight: 600; color: #444; }
  .kf-feld input { padding: 11px; border: 1px solid #ddd; border-radius: 10px; font-size: 0.95rem; }
  .kartentyp { background: #673ab7; color: white; padding: 1px 8px; border-radius: 10px; font-size: 0.72rem; margin-left: 6px; }
  .kf-fehler { color: #dc3545; font-size: 0.8rem; font-weight: 600; }
  .kf-hinweis { font-size: 0.78rem; color: #999; margin: 0; }

  /* 💶 Trinkgeld */
  .trinkgeld-optionen { display: flex; gap: 10px; }
  .tg-btn { flex: 1; padding: 12px; border: 1px solid #ddd; border-radius: 12px; background: white; cursor: pointer; font-weight: 600; }
  .tg-btn.aktiv { border-color: #673ab7; background: #faf7ff; color: #673ab7; }
  .tg-info { color: #34c759; font-weight: 600; margin: 10px 0 0; font-size: 0.9rem; }

  .bestellnr { font-size: 1rem; color: #555; margin: 4px 0; }
  .bestellnr strong { color: #673ab7; letter-spacing: 1px; }

  /* ⏰ Lieferzeit + Adresswahl */
  .zeit-optionen { display: flex; flex-direction: column; gap: 10px; }
  .zeit-optionen label { border: 1px solid #ddd; border-radius: 12px; padding: 12px; cursor: pointer; }
  .zeit-optionen label.aktiv { border-color: #673ab7; background: #faf7ff; }
  .zeit-feld { margin-top: 12px; padding: 11px; border: 1px solid #ddd; border-radius: 10px; font-size: 0.95rem; }
  .adress-wahl { width: 100%; padding: 11px; border: 1px solid #ddd; border-radius: 10px; margin-bottom: 12px; }

  /* ⭐ Treuepunkte */
  .punkte-stand { color: #555; }
  .punkte-einloesen { display: flex; align-items: center; gap: 8px; font-weight: 600; color: #673ab7; cursor: pointer; }

  /* 🧾 Druckbare Rechnung: am Bildschirm versteckt, nur beim Drucken sichtbar */
  .rechnung-druck { display: none; }
  @media print {
    /* Beim Drucken alles ausblenden außer der Rechnung */
    :global(body *) { visibility: hidden !important; }
    .rechnung-druck, .rechnung-druck * { visibility: visible !important; }
    .rechnung-druck { display: block; position: absolute; top: 0; left: 0; width: 100%; padding: 20px; font-family: sans-serif; color: #000; }
    .rechnung-druck table { width: 100%; border-collapse: collapse; margin: 16px 0; }
    .rechnung-druck td { padding: 4px 0; border-bottom: 1px solid #eee; }
    .rechnung-druck .r-gesamt td { font-weight: bold; border-top: 2px solid #000; border-bottom: none; }
  }

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
