<script>
  import { onMount } from 'svelte';
  // HIER DIE NEUE ZEILE: Füge sie einfach direkt unter 'onMount' ein
  import { warenkorb } from '$lib/stores/cart.js';

  // Neue Services für E-Mail-Verifizierung, Passwort-Sicherheit und Adress-Check.
  import { generiereCode, sendeVerifizierungsEmail } from '$lib/services/email.js';
  import { pruefePasswortStaerke } from '$lib/services/passwort.js';
  import { pruefeAdresse } from '$lib/services/adresse.js';

  // 🍽️ Aktive Restaurants (vom Admin gelöschte werden ausgeblendet).
  import { aktiveRestaurants } from '$lib/stores/lieferanten.js';

  // 🔑 Anmelde-Status kommt jetzt zentral aus dem Auth-Store.
  import { eingeloggt, login, hatKonto } from '$lib/stores/auth.js';
  // 🛠️ dev ist true bei "npm run dev" – damit zeigen wir den Dev-Überspringen-Button.
  import { dev } from '$app/environment';
  // 🌍 Übersetzungen
  import { t } from '$lib/i18n.js';
  // 🥚 Easter Eggs
  import { drachenlordAusloesen } from '$lib/stores/easteregg.js';
  import { konfetti, eierToast } from '$lib/confetti.js';

  // Merkt sich, ob schon ein Konto existiert (dann zeigen wir "Willkommen zurück").
  let kontoVorhanden = $state(false);
  let loginSchritt = $state(1);

  // Schritt 1
  let emailInput = $state("");
  let passwortInput = $state("");

  // 🔒 Passwort-Stärke: $derived berechnet sich automatisch neu, sobald sich
  // das Passwort ändert. So bleibt die Anzeige immer aktuell.
  let passwortStaerke = $derived(pruefePasswortStaerke(passwortInput));

  // 📧 E-Mail-Verifizierung (zwischen Schritt 1 und Schritt 2)
  let zeigeVerifizierung = $state(false); // Zeigt das Code-Eingabefeld an
  let korrekterCode = $state("");         // Der erzeugte Code (im Test sichtbar)
  let eingegebenerCode = $state("");      // Was der Nutzer eintippt
  let codeFehler = $state("");            // Fehlermeldung bei falschem Code
  let testCodeHinweis = $state("");       // Nur für den Test: zeigt den Code an

  // 🏠 Adress-Prüfung (Schritt 3)
  let prueftAdresse = $state(false); // true, während wir online prüfen
  let adressFehler = $state("");     // Fehlermeldung, falls Adresse nicht existiert

  // Schritt 2 (mit neuen Validierungsvariablen)
  let usernameInput = $state("");
  let vornameInput = $state("");
  let nachnameInput = $state("");
  let zweitnameInput = $state("");
  let zeigtZweitname = $state(false);
  let geburtsdatumInput = $state("");
  let altersFehler = $state("");
  let namensFehler = $state(""); // Fehlermeldung für Vor-/Nachname

  // Prüft einen Namen: mind. 2 Zeichen, nur Buchstaben (inkl. Umlaute),
  // Leerzeichen, Bindestrich und Apostroph erlaubt – KEINE Zahlen.
  function nameGueltig(name) {
    return /^[A-Za-zÄÖÜäöüß' -]{2,}$/.test(name.trim());
  }

  // Dynamisches maximales Jahr berechnen (heutiges Datum im YYYY-MM-DD Format)
  let heutigesDatumIso = new Date().toISOString().split('T')[0];

  // Schritt 3
  let strasseInput = $state("");
  let hausnummerInput = $state("");
  let plzInput = $state("");
  let ortInput = $state("");

  // 🥚 Hero-Klick-Zähler (10× klicken = Easter Egg)
  let heroKlicks = $state(0);
  function heroKlick() {
    heroKlicks += 1;
    if (heroKlicks === 10) {
      heroKlicks = 0;
      konfetti({ anzahl: 100, dauer: 2500, emojis: ['🍕', '🍔', '🌮', '🍣', '🎉'] });
      eierToast('🥚 Easter Egg gefunden! Du klickst gerne, was? 😄');
    }
  }

  // 🥚 Suche überwacht versteckte Codewörter.
  $effect(() => {
    const s = suche.toLowerCase().replace(/\s/g, '');
    if (s === 'drachenlord') {
      drachenlordAusloesen(); // 🐉 große Show
    } else if (s === 'pizzapizzapizza') {
      konfetti({ anzahl: 120, dauer: 3000, emojis: ['🍕'] });
      eierToast('🍕 Geheimcode entdeckt! Nutze PIZZAPARTY für 25% Rabatt 🎉');
    }
  });

  onMount(() => {
    // Anmelde-Status liefert der Auth-Store. Hier merken wir uns nur, ob bereits
    // ein Konto existiert – dann zeigen wir den Login-Hinweis statt der Registrierung.
    kontoVorhanden = hatKonto();

    // 🎂 Geburtstags-Überraschung: Wenn heute Geburtstag ist (laut Konto).
    const gespeichert = localStorage.getItem('lieferino_user');
    if (gespeichert) {
      const user = JSON.parse(gespeichert);
      if (user.geburtsdatum) {
        const heute = new Date();
        const geb = new Date(user.geburtsdatum);
        if (heute.getDate() === geb.getDate() && heute.getMonth() === geb.getMonth()) {
          konfetti({ anzahl: 150, dauer: 4000, emojis: ['🎂', '🎈', '🎉', '🥳'] });
          eierToast('🎂 Alles Gute zum Geburtstag! Code GEBURTSTAG = 20% Rabatt 🥳', 6000);
        }
      }
    }
  });

  async function geheZuSchritt2(e) {
    e.preventDefault();

    // 🔒 Sicherheits-Check: Passwort muss stark genug sein, bevor es weitergeht.
    if (!passwortStaerke.istSicher) {
      codeFehler = "";
      return; // Die Stärke-Anzeige im Formular sagt dem Nutzer, was fehlt.
    }

    if (emailInput.trim() !== "" && passwortInput.trim() !== "") {
      // 📧 Wir erzeugen einen Verifizierungscode und "senden" ihn an die E-Mail.
      // (Im Test landet der Code nur in der Konsole – siehe email.js.)
      korrekterCode = generiereCode();
      const ergebnis = await sendeVerifizierungsEmail(emailInput, korrekterCode);

      // Solange noch kein echtes Backend dranhängt, zeigen wir den Test-Code an.
      testCodeHinweis = ergebnis.testCode ?? "";

      // Statt direkt zu Schritt 2 zeigen wir erst die Code-Eingabe.
      zeigeVerifizierung = true;
      eingegebenerCode = "";
      codeFehler = "";
    }
  }

  // 📧 Prüft den eingegebenen Verifizierungscode.
  function bestaetigeCode(e) {
    e.preventDefault();
    if (eingegebenerCode.trim() === korrekterCode) {
      // Code korrekt -> E-Mail gilt als verifiziert, weiter zu Schritt 2.
      zeigeVerifizierung = false;
      codeFehler = "";
      loginSchritt = 2;
    } else {
      codeFehler = "Der Code ist leider falsch. Bitte versuche es erneut. 🔁";
    }
  }

  // 📧 Sendet auf Wunsch einen neuen Code (z.B. wenn die Mail nicht ankam).
  async function codeErneutSenden() {
    korrekterCode = generiereCode();
    const ergebnis = await sendeVerifizierungsEmail(emailInput, korrekterCode);
    testCodeHinweis = ergebnis.testCode ?? "";
    codeFehler = "";
  }

  function geheZuSchritt3(e) {
    e.preventDefault();

    // Namens-Check: Vor- und Nachname müssen gültig sein (min. 2 Buchstaben,
    // keine Zahlen). Ein optionaler Zweitname wird nur geprüft, wenn er ausgefüllt ist.
    if (!nameGueltig(vornameInput) || !nameGueltig(nachnameInput)) {
      namensFehler = "Vor- und Nachname brauchen mind. 2 Buchstaben und dürfen keine Zahlen enthalten. ✍️";
      return;
    }
    if (zeigtZweitname && zweitnameInput.trim() !== "" && !nameGueltig(zweitnameInput)) {
      namensFehler = "Der zweite Name darf keine Zahlen enthalten. ✍️";
      return;
    }
    namensFehler = "";

    if (!geburtsdatumInput) return;

    const geburtstag = new Date(geburtsdatumInput);
    const jahr = geburtstag.getFullYear();
    const heute = new Date();

    // 1. Check: Nicht vor 1900, nicht in der Zukunft
    if (jahr < 1900 || geburtstag > heute) {
      altersFehler = "Bitte gib ein realistisches Geburtsdatum an (ab 1900 bis heute)! 🤔";
      return;
    }
    
    // 2. Check: Mindestens 18 Jahre alt
    let alter = heute.getFullYear() - geburtstag.getFullYear();
    const monatsDiff = heute.getMonth() - geburtstag.getMonth();
    if (monatsDiff < 0 || (monatsDiff === 0 && heute.getDate() < geburtstag.getDate())) {
      alter--;
    }

    if (alter < 18) {
      altersFehler = "Du musst mindestens 18 Jahre alt sein, um Lieferino zu nutzen! 🔞";
      return;
    }

    altersFehler = "";
    loginSchritt = 3; 
  }

  async function registrationAbschliessen(e) {
    e.preventDefault();

    // PLZ-Format prüfen: in Deutschland genau 5 Ziffern.
    if (!/^\d{5}$/.test(plzInput.trim())) {
      adressFehler = "Die PLZ muss aus genau 5 Ziffern bestehen. 🔢";
      return;
    }

    // 🏠 Bevor wir die Registrierung abschließen, prüfen wir online, ob die
    // eingegebene Adresse überhaupt existiert.
    adressFehler = "";
    prueftAdresse = true;
    const adressErgebnis = await pruefeAdresse({
      strasse: strasseInput,
      hausnummer: hausnummerInput,
      plz: plzInput,
      ort: ortInput
    });
    prueftAdresse = false;

    if (!adressErgebnis.gefunden) {
      // Wenn der Online-Dienst gerade nicht erreichbar ist, lassen wir den
      // Nutzer trotzdem durch (sonst blockiert ein Server-Ausfall alles).
      if (!adressErgebnis.fehler) {
        adressFehler = "Diese Adresse konnten wir nicht finden. Bitte überprüfe deine Eingabe. 🗺️";
        return;
      }
    }

    const userDaten = {
      email: emailInput, passwort: passwortInput, username: usernameInput,
      vorname: vornameInput, nachname: nachnameInput, zweitname: zeigtZweitname ? zweitnameInput : "",
      strasse: strasseInput, hausnummer: hausnummerInput, plz: plzInput, ort: ortInput,
      geburtsdatum: geburtsdatumInput
    };
    // 🚨 BACKEND-HINWEIS: Diese Nutzerdaten landen aktuell nur im localStorage
    // des Browsers. Das Backend muss hier später einen Endpunkt anbieten
    // (z.B. POST /api/auth/register), der den Nutzer in der Datenbank anlegt.
    localStorage.setItem("lieferino_user", JSON.stringify(userDaten));
    // Nach der Registrierung gleich anmelden (Session setzen) – kein Reload nötig.
    login();
  }

  // 🍽️ Restaurants kommen jetzt aus der zentralen Quelle (siehe $lib/data)

  let gewaehlterTyp = $state("alle");
  let maxMinBestellwert = $state(30);
  let suche = $state(""); // 🔍 Suche auf der Startseite

  let gefilterteRestaurants = $derived(
    $aktiveRestaurants.filter(r => {
      let typPasst = gewaehlterTyp === "alle" || r.typ === gewaehlterTyp;
      let preisPasst = r.minBestell <= maxMinBestellwert;
      let suchePasst = r.name.toLowerCase().includes(suche.toLowerCase());
      return typPasst && preisPasst && suchePasst;
    })
  );

  let top10Restaurants = $derived(
    [...$aktiveRestaurants].sort((a, b) => b.bewertung - a.bewertung).slice(0, 10)
  );
</script>

{#if !$eingeloggt && kontoVorhanden}
  <!-- 🔑 Konto existiert, aber nicht eingeloggt: zum Login schicken -->
  <div class="login-wrapper">
    <div class="login-box">
      <div class="hero-box compact-hero">
        <h2>👋 Willkommen zurück!</h2>
      </div>
      <p>Du hast bereits ein Konto. Bitte melde dich an.</p>
      <a href="/login" class="login-btn" style="display:block; text-align:center; text-decoration:none;">Zum Login 🔑</a>
    </div>
  </div>
{:else if !$eingeloggt}
  <div class="login-wrapper">
    <div class="login-box">
      <!-- Hier ebenfalls das neue lila quadratische Banner für den Login -->
      <div class="hero-box compact-hero">
        <h2>🔑 Lieferino Account</h2>
      </div>
      
      {#if loginSchritt === 1 && !zeigeVerifizierung}
        <p>Schritt 1/3: Erstelle deine Login-Daten</p>
        <form onsubmit={geheZuSchritt2}>
          <div class="input-group">
            <label for="email">E-Mail Adresse</label>
            <input type="email" id="email" placeholder="name@beispiel.de" bind:value={emailInput} required />
          </div>
          <div class="input-group">
            <label for="password">Passwort</label>
            <input type="password" id="password" placeholder="••••••••" bind:value={passwortInput} required />

            <!-- 🔒 Passwort-Stärke-Anzeige: erscheint, sobald getippt wird -->
            {#if passwortInput.length > 0}
              <div class="pw-staerke">
                <!-- Der Balken füllt sich je nach Score (0-5). -->
                <div class="pw-balken-hintergrund">
                  <div
                    class="pw-balken"
                    style="width: {passwortStaerke.score * 20}%; background: {passwortStaerke.farbe};"
                  ></div>
                </div>
                <span class="pw-text" style="color: {passwortStaerke.farbe};">
                  {passwortStaerke.text}
                </span>

                <!-- Checkliste: was ist schon erfüllt, was fehlt noch? -->
                <ul class="pw-regeln">
                  <li class:erfuellt={passwortStaerke.regeln.laenge}>Mind. 8 Zeichen</li>
                  <li class:erfuellt={passwortStaerke.regeln.grossbuchstabe}>Großbuchstabe</li>
                  <li class:erfuellt={passwortStaerke.regeln.kleinbuchstabe}>Kleinbuchstabe</li>
                  <li class:erfuellt={passwortStaerke.regeln.zahl}>Zahl</li>
                  <li class:erfuellt={passwortStaerke.regeln.sonderzeichen}>Sonderzeichen</li>
                </ul>
              </div>
            {/if}
          </div>
          <!-- Button ist gesperrt, bis das Passwort sicher genug ist. -->
          <button type="submit" class="login-btn" disabled={!passwortStaerke.istSicher}>
            Weiter zu deinen Details ➡️
          </button>
        </form>
        <p class="login-hinweis">Schon registriert? <a href="/login">Hier einloggen 🔑</a></p>

      {:else if loginSchritt === 1 && zeigeVerifizierung}
        <!-- 📧 E-MAIL-VERIFIZIERUNG: Code eingeben, der per Mail kam -->
        <p class="step-title">Bestätige deine E-Mail 📧</p>
        <p class="verify-info">
          Wir haben einen 6-stelligen Code an <strong>{emailInput}</strong> geschickt.
        </p>

        <!-- ⚠️ NUR FÜR DEN TEST: Solange kein echtes Backend Mails verschickt,
             zeigen wir den Code hier direkt an. Das MUSS entfernt werden,
             sobald der echte E-Mail-Versand übers Backend läuft! -->
        {#if testCodeHinweis}
          <p class="test-code-hinweis">🧪 Test-Code (kommt später per Mail): <strong>{testCodeHinweis}</strong></p>
        {/if}

        <form onsubmit={bestaetigeCode}>
          <div class="input-group">
            <label for="code">Verifizierungscode</label>
            <input
              type="text"
              id="code"
              placeholder="z.B. 048213"
              maxlength="6"
              bind:value={eingegebenerCode}
              required
            />
            {#if codeFehler}
              <p class="error-msg">{codeFehler}</p>
            {/if}
          </div>
          <div class="button-row">
            <button type="button" class="back-btn" onclick={() => zeigeVerifizierung = false}>⬅️ Zurück</button>
            <button type="submit" class="login-btn compact-btn">Code bestätigen ✅</button>
          </div>
        </form>
        <button type="button" class="link-btn" onclick={codeErneutSenden}>Code erneut senden 🔁</button>

        <!-- 🛠️ Nur im Entwicklungsmodus: Verifizierung überspringen -->
        {#if dev}
          <button type="button" class="dev-skip" onclick={() => { zeigeVerifizierung = false; loginSchritt = 2; }}>
            🛠️ Dev: überspringen
          </button>
        {/if}

      {:else if loginSchritt === 2}
        <p class="step-title">Schritt 2/3: Wer bist du?</p>
        <form onsubmit={geheZuSchritt3}>
          <div class="input-group">
            <label for="username">Username</label>
            <input type="text" id="username" placeholder="z.B. max_power" bind:value={usernameInput} required />
          </div>

          <div class="name-row">
            <div class="input-group">
              <label for="vorname">Vorname</label>
              <div class="input-with-plus">
                <input type="text" id="vorname" placeholder="Max" bind:value={vornameInput} required />
                {#if !zeigtZweitname}
                  <button type="button" class="plus-btn" onclick={() => zeigtZweitname = true} title="Zweitname hinzufügen">+</button>
                {/if}
              </div>
            </div>

            {#if zeigtZweitname}
              <div class="input-group">
                <label for="zweitname">2. Name</label>
                <input type="text" id="zweitname" placeholder="Maria" bind:value={zweitnameInput} />
              </div>
            {/if}

            <div class="input-group">
              <label for="nachname">Nachname</label>
              <input type="text" id="nachname" placeholder="Mustermann" bind:value={nachnameInput} required />
            </div>
          </div>

          <!-- Fehlermeldung bei ungültigem Namen -->
          {#if namensFehler}
            <p class="error-msg">{namensFehler}</p>
          {/if}

          <div class="input-group">
            <label for="geburt">Geburtsdatum (Alters-Check)</label>
            <input type="date" id="geburt" min="1900-01-01" max={heutigesDatumIso} bind:value={geburtsdatumInput} required />
            {#if altersFehler}
              <p class="error-msg">{altersFehler}</p>
            {/if}
          </div>

          <div class="button-row">
            <button type="button" onclick={() => loginSchritt = 1} class="back-btn">⬅️ Zurück</button>
            <button type="submit" class="login-btn compact-btn">Weiter zur Adresse 🏠</button>
          </div>
        </form>

      {:else if loginSchritt === 3}
        <p class="step-title">Schritt 3/3: Wohin soll das Essen geliefert werden?</p>
        <form onsubmit={registrationAbschliessen}>
          <div class="address-row-1">
            <div class="input-group strasse">
              <label for="strasse">Straße</label>
              <input type="text" id="strasse" placeholder="Musterstraße" bind:value={strasseInput} required />
            </div>
            <div class="input-group nr">
              <label for="nr">Nr.</label>
              <input type="text" id="nr" placeholder="12a" bind:value={hausnummerInput} required />
            </div>
          </div>

          <div class="address-row-2">
            <div class="input-group plz">
              <label for="plz">PLZ</label>
              <input type="text" id="plz" placeholder="12345" bind:value={plzInput} required />
            </div>
            <div class="input-group ort">
              <label for="ort">Ort</label>
              <input type="text" id="ort" placeholder="Musterstadt" bind:value={ortInput} required />
            </div>
          </div>

          <!-- 🏠 Fehlermeldung, falls die Adresse nicht gefunden wurde -->
          {#if adressFehler}
            <p class="error-msg">{adressFehler}</p>
          {/if}

          <div class="button-row">
            <button type="button" onclick={() => loginSchritt = 2} class="back-btn">⬅️ Zurück</button>
            <!-- Während der Online-Prüfung ist der Button gesperrt + zeigt Ladetext. -->
            <button type="submit" class="login-btn compact-btn" disabled={prueftAdresse}>
              {prueftAdresse ? "Prüfe Adresse… ⏳" : "Registrierung abschließen 🎉"}
            </button>
          </div>
        </form>
      {/if}
    </div>
  </div>
{:else}
  <!-- 🍔 RESTLICHE STARTSEITE -->
  <div class="welcome-container">
    
    <!-- 🟪 LILA HERO-QUADRAT 🟪 (Titel 10× klicken = verstecktes Egg) -->
    <div class="hero-box">
      <h1 onclick={heroKlick} role="presentation">{$t('home.hero_title')}</h1>
      <p>{$t('home.hero_sub')}</p>
    </div>

    <!-- 🔍 Suchleiste -->
    <div class="filter-group" style="margin-bottom: 16px;">
      <input id="suche" type="search" placeholder={$t('common.search_placeholder')} bind:value={suche} class="such-input" />
    </div>

    <!-- Filter-Bar -->
    <div class="filter-bar">
      <div class="filter-group">
        <label for="kueche">{$t('home.cuisine_label')}</label>
        <select id="kueche" bind:value={gewaehlterTyp}>
          <option value="alle">{$t('common.all_cuisines')}</option>
          <option value="italienisch">🍕 Italienisch</option>
          <option value="spanisch">🥘 Spanisch</option>
          <option value="amerikanisch">🍔 Amerikanisch</option>
          <option value="asiatisch">🍣 Asiatisch</option>
          <option value="mexikanisch">🌮 Mexikanisch</option>
        </select>
      </div>

      <div class="filter-group">
        <label for="mindest">{$t('home.min_order_label')} <strong>{maxMinBestellwert}€</strong></label>
        <input type="range" id="mindest" min="5" max="30" step="5" bind:value={maxMinBestellwert} />
      </div>
    </div>

    <!-- ⭐ TOP 10 NETFLIX SCROLLER -->
    {#if gewaehlterTyp === "alle"}
      <div class="section-title">
        <h2>{$t('home.top10_title')}</h2>
        <p class="subtitle">{$t('home.top10_sub')}</p>
      </div>
      
      <div class="netflix-scroll-container">
        {#each top10Restaurants as restaurant, index}
          <a href="/restaurant/{restaurant.slug}" class="netflix-card">
            <div class="rank-number">#{index + 1}</div>
            <div class="image-placeholder emoji-bild">
              <span class="emoji-gross">{restaurant.emoji}</span>
              <span class="rating-badge">⭐ {restaurant.bewertung}</span>
            </div>
            <h3>{restaurant.name}</h3>
            <p class="desc">{restaurant.beschreibung}</p>
            <div class="card-footer">
              <span class="tag">{restaurant.typ}</span>
              <span class="min-order">{$t('common.min')}: {restaurant.minBestell}€</span>
            </div>
          </a>
        {/each}
      </div>
    {/if}

    <!-- Alle Restaurants -->
    <div class="section-title scroll-section">
      <h2>{$t('home.discover')}</h2>
    </div>

    <div class="restaurant-grid">
      {#each gefilterteRestaurants as restaurant}
        <a href="/restaurant/{restaurant.slug}" class="restaurant-card">
          <div class="image-placeholder emoji-bild">
            <span class="emoji-gross">{restaurant.emoji}</span>
            <span class="rating-badge">⭐ {restaurant.bewertung}</span>
          </div>
          <h3>{restaurant.name}</h3>
          <p class="desc">{restaurant.beschreibung}</p>
          <div class="card-footer">
            <span class="tag">{restaurant.typ}</span>
            <span class="min-order">Min: {restaurant.minBestell}€</span>
          </div>
        </a>
      {/each}
    </div>
  </div>
{/if}

<style>
  /* 🟪 LILA BANNER STYLES */
  .hero-box { background: #673ab7; color: white; padding: 40px 20px; border-radius: 24px; text-align: center; margin-bottom: 40px; box-shadow: 0 8px 25px rgba(103, 58, 183, 0.2); }
  .hero-box h1 { font-size: 2.8rem; color: white !important; margin: 0 0 10px 0; font-weight: 800; }
  .hero-box p { color: #e1d5f5 !important; font-size: 1.2rem; margin: 0; }
  .compact-hero { padding: 20px; margin-bottom: 25px; border-radius: 16px; }
  .compact-hero h2 { color: white !important; margin: 0; font-size: 1.5rem; }

  /* Login & Tab Styles */
  .login-wrapper { display: flex; justify-content: center; align-items: center; min-height: 90vh; font-family: sans-serif; }
  .login-box { background: white; padding: 35px; border-radius: 24px; box-shadow: 0 10px 30px rgba(0,0,0,0.1); max-width: 480px; width: 100%; border: 1px solid #eee; }
  .login-box p { color: #777; margin-bottom: 25px; font-size: 0.95rem; text-align: center; }
  .step-title { color: #673ab7 !important; font-weight: bold; }
  
  .input-group { text-align: left; margin-bottom: 15px; display: flex; flex-direction: column; gap: 4px; width: 100%; }
  .input-group label { font-weight: 600; color: #444; font-size: 0.85rem; }
  .input-group input, select { padding: 11px; border: 1px solid #ddd; border-radius: 10px; font-size: 0.95rem; box-sizing: border-box; width: 100%; }
  
  .name-row, .address-row-1, .address-row-2 { display: flex; gap: 10px; width: 100%; }
  .input-with-plus { display: flex; gap: 6px; align-items: center; width: 100%; }
  .plus-btn { padding: 10px 14px; background: #f3e5f5; color: #673ab7; border: none; border-radius: 10px; font-weight: bold; cursor: pointer; font-size: 1.1rem; }
  
  .nr { width: 30%; }
  .plz { width: 35%; }
  
  .error-msg { color: #dc3545; font-size: 0.85rem; font-weight: bold; margin-top: 5px; text-align: left !important; }

  /* 🔒 Passwort-Stärke-Anzeige */
  .pw-staerke { margin-top: 8px; text-align: left; }
  .pw-balken-hintergrund { background: #eee; border-radius: 6px; height: 8px; overflow: hidden; }
  .pw-balken { height: 100%; transition: width 0.3s ease, background 0.3s ease; }
  .pw-text { font-size: 0.8rem; font-weight: bold; display: inline-block; margin-top: 4px; }
  .pw-regeln { list-style: none; padding: 0; margin: 8px 0 0; display: flex; flex-wrap: wrap; gap: 6px 14px; }
  .pw-regeln li { font-size: 0.78rem; color: #999; }
  .pw-regeln li::before { content: "✗ "; color: #dc3545; }
  .pw-regeln li.erfuellt { color: #34c759; }
  .pw-regeln li.erfuellt::before { content: "✓ "; color: #34c759; }

  /* 📧 E-Mail-Verifizierung */
  .verify-info { color: #777; }
  .test-code-hinweis { background: #fff8e1; border: 1px dashed #ffc107; border-radius: 10px; padding: 8px; color: #8a6d00 !important; font-size: 0.85rem; }
  .link-btn { background: none; border: none; color: #673ab7; text-decoration: underline; cursor: pointer; font-size: 0.85rem; margin-top: 10px; }

  /* Gesperrte Buttons (z.B. solange Passwort zu schwach oder Adresse geprüft wird) */
  .login-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .login-hinweis { text-align: center; font-size: 0.9rem; color: #777; margin-top: 14px; }
  .login-hinweis a { color: #673ab7; font-weight: 600; }
  .such-input { padding: 11px; border: 1px solid #ddd; border-radius: 10px; font-size: 0.95rem; width: 100%; box-sizing: border-box; }
  .dev-skip { display: block; width: 100%; margin-top: 10px; padding: 10px; background: #fff3cd; color: #8a6d00; border: 1px dashed #ffc107; border-radius: 10px; cursor: pointer; font-weight: 600; font-size: 0.85rem; }

  .button-row { display: flex; gap: 12px; margin-top: 15px; }
  .back-btn { padding: 14px; background: #f1f1f1; color: #333; border: none; border-radius: 10px; font-weight: bold; cursor: pointer; }
  .login-btn { width: 100%; padding: 14px; background: #673ab7; color: white; border: none; border-radius: 10px; font-weight: bold; cursor: pointer; }
  .login-btn:hover { background: #542f95; }
  .compact-btn { flex: 1; }

  /* 🍔 RESTAURANT CSS */
  .welcome-container { max-width: 1200px; margin: 0 auto; padding-top: 40px; font-family: sans-serif; overflow-x: hidden; }
  .filter-bar { display: flex; justify-content: space-between; align-items: center; background: #f8f9fa; padding: 20px 30px; border-radius: 16px; margin-bottom: 50px; box-shadow: 0 4px 15px rgba(0,0,0,0.02); gap: 30px; }
  .filter-group { display: flex; flex-direction: column; gap: 8px; flex: 1; text-align: left; }
  .filter-group label { font-weight: 600; color: #444; font-size: 0.95rem; }
  input[type="range"] { accent-color: #673ab7; cursor: pointer; }
  .section-title { text-align: left; padding: 0 20px; margin-bottom: 25px; }
  .section-title h2 { font-size: 1.8rem; color: #222; margin-bottom: 5px; }
  .subtitle { color: #888; font-size: 1rem; margin: 0; }
  .scroll-section { margin-top: 60px; border-top: 1px solid #eee; padding-top: 40px; }
  
  .netflix-scroll-container { display: flex; gap: 25px; overflow-x: auto; padding: 20px 20px 30px 20px; scroll-behavior: smooth; }
  .netflix-scroll-container::-webkit-scrollbar { height: 8px; }
  .netflix-scroll-container::-webkit-scrollbar-thumb { background: #673ab7; border-radius: 10px; }
  
  .netflix-card, .restaurant-card {
    flex: 0 0 280px;
    background: white;
    border-radius: 16px;
    padding: 15px;
    box-shadow: 0 6px 18px rgba(0, 0, 0, 0.05);
    border: 1px solid #f0f0f0;
    display: flex;
    flex-direction: column;
    position: relative;
    text-decoration: none;
    color: inherit;
    transition: transform 0.2s;
  }
  .netflix-card:hover, .restaurant-card:hover { transform: translateY(-5px); }
  .restaurant-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 30px; padding: 0 20px; margin-bottom: 40px; }
  .restaurant-card { flex: none; }
  
  .rank-number { position: absolute; top: -15px; left: -10px; font-size: 3.5rem; font-weight: 900; color: #673ab7; text-shadow: 2px 2px 0px white, 4px 4px 10px rgba(0,0,0,0.15); z-index: 10; }
  .image-placeholder { width: 100%; height: 150px; border-radius: 12px; overflow: hidden; position: relative; }
  .image-placeholder img { width: 100%; height: 100%; object-fit: cover; }
  /* 🍽️ Emoji-Variante statt echtem Bild */
  .emoji-bild { display: flex; align-items: center; justify-content: center; background: linear-gradient(135deg, #f3e5f5, #ede7f6); }
  .emoji-gross { font-size: 4rem; }
  .rating-badge { position: absolute; top: 10px; right: 10px; background: white; padding: 4px 8px; border-radius: 8px; font-weight: bold; font-size: 0.85rem; }
  h3 { font-size: 1.25rem; color: #222; margin: 15px 0 5px 0; text-align: left; }
  .desc { font-size: 0.9rem; color: #777; margin: 0 0 15px 0; flex-grow: 1; text-align: left; }
  .card-footer { display: flex; justify-content: space-between; align-items: center; border-top: 1px solid #f5f5f5; padding-top: 12px; font-size: 0.85rem; }
  .tag { background: #f3e5f5; color: #673ab7; padding: 4px 10px; border-radius: 20px; font-weight: 600; }
  .min-order { color: #666; font-weight: 500; }
</style>