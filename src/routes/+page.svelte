<script>
  import { onMount } from 'svelte';
  import { warenkorb } from '$lib/stores/cart.js';
  import { pruefePasswortStaerke } from '$lib/services/passwort.js';
  import { pruefeAdresse } from '$lib/services/adresse.js';
  import { api } from '$lib/api.js';
  import AuthVollenden from '$lib/AuthVollenden.svelte';
  import { aktiveRestaurants } from '$lib/stores/lieferanten.js';
  import { bewertungsSchnitt } from '$lib/stores/bewertungsSchnitt.js';
  import { eingeloggt, login, hatKonto } from '$lib/stores/auth.js';
  import { t } from '$lib/i18n.js';
  import { drachenlordAusloesen } from '$lib/stores/easteregg.js';
  import { konfetti, eierToast } from '$lib/confetti.js';
  import { toggleEmojiCursor, toggleSaison } from '$lib/stores/funmodus.js';
  import { geheimFreischalten } from '$lib/stores/lieferanten.js';

  let kontoVorhanden = $state(false);
  let loginSchritt = $state(1);
  let emailInput = $state("");
  let passwortInput = $state("");
  let passwortStaerke = $derived(pruefePasswortStaerke(passwortInput));
  let zeigeAbschluss = $state(false); // true -> E-Mail-Code + MFA-Einrichtung (AuthVollenden)
  let prueftAdresse = $state(false);
  let adressFehler = $state("");
  let registrierFehler = $state("");
  let registrierAnforderungen = $state([]);
  let usernameInput = $state("");
  let vornameInput = $state("");
  let nachnameInput = $state("");
  let zweitnameInput = $state("");
  let zeigtZweitname = $state(false);
  let geburtsdatumInput = $state("");
  let altersFehler = $state("");
  let namensFehler = $state("");

  function nameGueltig(name) {
    return /^[A-Za-zÄÖÜäöüß' -]{2,}$/.test(name.trim());
  }

  let heutigesDatumIso = new Date().toISOString().split('T')[0];
  let strasseInput = $state("");
  let hausnummerInput = $state("");
  let plzInput = $state("");
  let ortInput = $state("");
  let heroKlicks = $state(0);

  function heroKlick() {
    heroKlicks += 1;
    if (heroKlicks === 10) {
      heroKlicks = 0;
      konfetti({ anzahl: 100, dauer: 2500, emojis: ['🍕','🍔','🌮','🍣','🎉'] });
      eierToast('🥚 Easter Egg gefunden! Du klickst gerne, was? 😄');
    }
  }

  $effect(() => {
    const s = suche.toLowerCase().replace(/\s/g, '');
    if (s === 'drachenlord') {
      drachenlordAusloesen();
    } else if (s === 'pizzapizzapizza') {
      konfetti({ anzahl: 120, dauer: 3000, emojis: ['🍕'] });
      eierToast('🍕 Geheimcode entdeckt! Nutze PIZZAPARTY für 25% Rabatt 🎉');
    } else if (s === 'foodcursor') {
      toggleEmojiCursor();
      eierToast('🖱️ Emoji-Cursor umgeschaltet!');
    } else if (s === 'schnee' || s === 'winter') {
      toggleSaison();
      eierToast('❄️ Saison-Effekt umgeschaltet!');
    } else if (s === 'dragonpizza') {
      geheimFreischalten();
      eierToast('🐲 Geheimes Restaurant freigeschaltet! Schau in die Restaurant-Liste 🔥');
    }
  });

  onMount(() => {
    kontoVorhanden = hatKonto();
    const gespeichert = localStorage.getItem('lieferino_user');
    if (gespeichert) {
      const user = JSON.parse(gespeichert);
      if (user.geburtsdatum) {
        const heute = new Date();
        const geb = new Date(user.geburtsdatum);
        if (heute.getDate() === geb.getDate() && heute.getMonth() === geb.getMonth()) {
          konfetti({ anzahl: 150, dauer: 4000, emojis: ['🎂','🎈','🎉','🥳'] });
          eierToast('🎂 Alles Gute zum Geburtstag! Code GEBURTSTAG = 20% Rabatt 🥳', 6000);
        }
      }
    }
  });

  // Schritt 1 -> 2: Das Passwort ist (per Button-disabled) schon stark genug.
  function geheZuSchritt2(e) {
    e.preventDefault();
    if (!passwortStaerke.istSicher) return;
    if (emailInput.trim() === "" || passwortInput.trim() === "") return;
    registrierFehler = "";
    loginSchritt = 2;
  }

  function geheZuSchritt3(e) {
    e.preventDefault();
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
    if (jahr < 1900 || geburtstag > heute) {
      altersFehler = "Bitte gib ein realistisches Geburtsdatum an (ab 1900 bis heute)! 🤔";
      return;
    }
    let alter = heute.getFullYear() - geburtstag.getFullYear();
    const monatsDiff = heute.getMonth() - geburtstag.getMonth();
    if (monatsDiff < 0 || (monatsDiff === 0 && heute.getDate() < geburtstag.getDate())) alter--;
    if (alter < 18) {
      altersFehler = "Du musst mindestens 18 Jahre alt sein, um Lieferino zu nutzen! 🔞";
      return;
    }
    altersFehler = "";
    loginSchritt = 3;
  }

  async function registrationAbschliessen(e) {
    e.preventDefault();
    if (!/^\d{5}$/.test(plzInput.trim())) {
      adressFehler = "Die PLZ muss aus genau 5 Ziffern bestehen. 🔢";
      return;
    }
    adressFehler = "";
    prueftAdresse = true;
    const adressErgebnis = await pruefeAdresse({ strasse: strasseInput, hausnummer: hausnummerInput, plz: plzInput, ort: ortInput });
    prueftAdresse = false;
    if (!adressErgebnis.gefunden && !adressErgebnis.fehler) {
      adressFehler = "Diese Adresse konnten wir nicht finden. Bitte überprüfe deine Eingabe. 🗺️";
      return;
    }
    const userDaten = {
      email: emailInput, passwort: passwortInput, username: usernameInput,
      vorname: vornameInput, nachname: nachnameInput, zweitname: zeigtZweitname ? zweitnameInput : "",
      strasse: strasseInput, hausnummer: hausnummerInput, plz: plzInput, ort: ortInput,
      geburtsdatum: geburtsdatumInput
    };
    localStorage.setItem("lieferino_user", JSON.stringify(userDaten));

    // 🗄️ Nutzer im Backend (Datenbank) anlegen. Es kommt KEIN Token zurück –
    // das Konto ist erst nutzbar nach E-Mail-Bestätigung + MFA-Einrichtung.
    registrierFehler = "";
    registrierAnforderungen = [];
    const reg = await api('/api/auth/register', {
      method: 'POST',
      body: {
        email: emailInput, passwort: passwortInput, username: usernameInput,
        vorname: vornameInput, nachname: nachnameInput, geburtsdatum: geburtsdatumInput
      }
    });

    if (reg.ok && reg.daten?.needsVerification) {
      // Backend hat einen Code per E-Mail geschickt -> Bestätigung + MFA einrichten.
      zeigeAbschluss = true;
      return;
    }
    if (reg.offline) {
      // Backend nicht erreichbar -> App lokal weiternutzen (Daten bleiben lokal).
      login();
      return;
    }
    // 🛡️ Server hat abgelehnt (E-Mail vergeben, Passwort, zu viele Anfragen) -> klar zeigen.
    registrierFehler = reg.daten?.fehler || "Registrierung fehlgeschlagen. Bitte später erneut versuchen.";
    registrierAnforderungen = reg.daten?.anforderungen || [];
  }

  // E-Mail bestätigt + MFA eingerichtet -> volles Token liegt vor. Jetzt Adresse/Profil
  // im Backend speichern und einloggen.
  async function abschlussFertig() {
    await api('/api/me', {
      method: 'PUT',
      body: {
        username: usernameInput, vorname: vornameInput, nachname: nachnameInput, geburtsdatum: geburtsdatumInput,
        adressen: [{ label: 'Zuhause', strasse: strasseInput, hausnummer: hausnummerInput, plz: plzInput, ort: ortInput }]
      }
    });
    login();
  }

  let gewaehlterTyp = $state("alle");
  let maxMinBestellwert = $state(30);
  let suche = $state("");

  let gefilterteRestaurants = $derived(
    $aktiveRestaurants.filter(r => {
      let typPasst = gewaehlterTyp === "alle" || r.typ === gewaehlterTyp;
      let preisPasst = r.minBestell <= maxMinBestellwert;
      let suchePasst = r.name.toLowerCase().includes(suche.toLowerCase());
      return typPasst && preisPasst && suchePasst;
    })
  );

  // Live-Sterne aus der DB (Durchschnitt aller Bewertungen), sonst statischer Wert.
  function schnittVon(r) {
    const s = $bewertungsSchnitt[r.slug];
    return s && s.anzahl > 0 ? s.schnitt : r.bewertung;
  }

  let top10Restaurants = $derived(
    [...$aktiveRestaurants].sort((a, b) => schnittVon(b) - schnittVon(a)).slice(0, 10)
  );
</script>

<!-- ░░░ NICHT EINGELOGGT: Konto vorhanden ░░░ -->
{#if !$eingeloggt && kontoVorhanden}
  <div class="auth-wrapper">
    <div class="auth-card">
      <div class="auth-hero">
        <span class="auth-hero-icon">👋</span>
        <h2>Willkommen zurück!</h2>
        <p>Du hast bereits ein Konto. Melde dich an.</p>
      </div>
      <a href="/login" class="gold-btn block-btn">Zum Login 🔑</a>
    </div>
  </div>

<!-- ░░░ NICHT EINGELOGGT: Registrierung ░░░ -->
{:else if !$eingeloggt}
  <div class="auth-wrapper">
    <div class="auth-card">

      <div class="auth-hero">
        <span class="auth-hero-icon">🍕</span>
        <h2>Lieferino Account</h2>
        <p>
          {#if zeigeAbschluss}Fast geschafft – bestätigen & absichern 🔐
          {:else if loginSchritt === 1}Schritt 1 / 3 – Login-Daten
          {:else if loginSchritt === 2}Schritt 2 / 3 – Deine Angaben
          {:else}Schritt 3 / 3 – Lieferadresse{/if}
        </p>
      </div>

      <!-- Schrittanzeige (nur während des Formulars) -->
      {#if !zeigeAbschluss}
        <div class="step-dots">
          <span class="dot" class:aktiv={loginSchritt >= 1}></span>
          <span class="dot-line"></span>
          <span class="dot" class:aktiv={loginSchritt >= 2}></span>
          <span class="dot-line"></span>
          <span class="dot" class:aktiv={loginSchritt >= 3}></span>
        </div>
      {/if}

      <!-- ── Abschluss: E-Mail-Code bestätigen + MFA einrichten ── -->
      {#if zeigeAbschluss}
        <AuthVollenden start="verify" email={emailInput} onFertig={abschlussFertig} />

      <!-- ── Schritt 1: E-Mail + Passwort ── -->
      {:else if loginSchritt === 1}
        <form onsubmit={geheZuSchritt2} class="auth-form">
          <div class="field">
            <label for="email">E-Mail</label>
            <input type="email" id="email" placeholder="name@beispiel.de" bind:value={emailInput} required />
          </div>
          <div class="field">
            <label for="password">Passwort</label>
            <input type="password" id="password" placeholder="••••••••" bind:value={passwortInput} required />
            {#if passwortInput.length > 0}
              <div class="pw-block">
                <div class="pw-track">
                  <div class="pw-fill" style="width:{passwortStaerke.score * 20}%; background:{passwortStaerke.farbe};"></div>
                </div>
                <span class="pw-label" style="color:{passwortStaerke.farbe};">{$t(passwortStaerke.stufeKey)}</span>
                <ul class="pw-rules">
                  <li class:ok={passwortStaerke.regeln.laenge}>{$t('pw.rule_length')}</li>
                  <li class:ok={passwortStaerke.regeln.grossbuchstabe}>{$t('pw.rule_upper')}</li>
                  <li class:ok={passwortStaerke.regeln.kleinbuchstabe}>{$t('pw.rule_lower')}</li>
                  <li class:ok={passwortStaerke.regeln.zahl}>{$t('pw.rule_digit')}</li>
                  <li class:ok={passwortStaerke.regeln.sonderzeichen}>{$t('pw.rule_special')}</li>
                </ul>
              </div>
            {/if}
          </div>
          <button type="submit" class="gold-btn" disabled={!passwortStaerke.istSicher}>
            Weiter zu deinen Details →
          </button>
          <p class="auth-hint">Schon registriert? <a href="/login">Hier einloggen 🔑</a></p>
        </form>

      <!-- ── Schritt 2: Persönliche Daten ── -->
      {:else if loginSchritt === 2}
        <form onsubmit={geheZuSchritt3} class="auth-form">
          <div class="field">
            <label for="username">Username</label>
            <input type="text" id="username" placeholder="z.B. max_power" bind:value={usernameInput} required />
          </div>
          <div class="name-row">
            <div class="field">
              <label for="vorname">Vorname</label>
              <div class="input-plus">
                <input type="text" id="vorname" placeholder="Max" bind:value={vornameInput} required />
                {#if !zeigtZweitname}
                  <button type="button" class="plus-btn" onclick={() => zeigtZweitname = true} title="Zweitname">+</button>
                {/if}
              </div>
            </div>
            {#if zeigtZweitname}
              <div class="field">
                <label for="zweitname">2. Name</label>
                <input type="text" id="zweitname" placeholder="Maria" bind:value={zweitnameInput} />
              </div>
            {/if}
            <div class="field">
              <label for="nachname">Nachname</label>
              <input type="text" id="nachname" placeholder="Mustermann" bind:value={nachnameInput} required />
            </div>
          </div>
          {#if namensFehler}<p class="err">{namensFehler}</p>{/if}
          <div class="field">
            <label for="geburt">Geburtsdatum</label>
            <input type="date" id="geburt" min="1900-01-01" max={heutigesDatumIso} bind:value={geburtsdatumInput} required />
            {#if altersFehler}<p class="err">{altersFehler}</p>{/if}
          </div>
          <div class="btn-row">
            <button type="button" class="ghost-btn" onclick={() => loginSchritt = 1}>← Zurück</button>
            <button type="submit" class="gold-btn flex1">Weiter zur Adresse 🏠</button>
          </div>
        </form>

      <!-- ── Schritt 3: Adresse ── -->
      {:else if loginSchritt === 3}
        <form onsubmit={registrationAbschliessen} class="auth-form">
          <div class="addr-row">
            <div class="field grow">
              <label for="strasse">Straße</label>
              <input type="text" id="strasse" placeholder="Musterstraße" bind:value={strasseInput} required />
            </div>
            <div class="field narrow">
              <label for="nr">Nr.</label>
              <input type="text" id="nr" placeholder="12a" bind:value={hausnummerInput} required />
            </div>
          </div>
          <div class="addr-row">
            <div class="field narrow">
              <label for="plz">PLZ</label>
              <input type="text" id="plz" placeholder="12345" bind:value={plzInput} required />
            </div>
            <div class="field grow">
              <label for="ort">Ort</label>
              <input type="text" id="ort" placeholder="Musterstadt" bind:value={ortInput} required />
            </div>
          </div>
          {#if adressFehler}<p class="err">{adressFehler}</p>{/if}
          {#if registrierFehler}
            <div class="err-box">
              <p class="err">{registrierFehler}</p>
              {#if registrierAnforderungen.length > 0}
                <ul class="err-list">
                  {#each registrierAnforderungen as a}<li>{a}</li>{/each}
                </ul>
              {/if}
            </div>
          {/if}
          <div class="btn-row">
            <button type="button" class="ghost-btn" onclick={() => loginSchritt = 2}>← Zurück</button>
            <button type="submit" class="gold-btn flex1" disabled={prueftAdresse}>
              {prueftAdresse ? "Prüfe Adresse… ⏳" : "Registrierung abschließen 🎉"}
            </button>
          </div>
        </form>
      {/if}

    </div>
  </div>

<!-- ░░░ EINGELOGGT: Startseite ░░░ -->
{:else}
  <div class="home">

    <!-- ── Hero ── -->
    <div class="hero" onclick={heroKlick} role="presentation">
      <div class="hero-glow"></div>
      <div class="hero-content">
        <p class="hero-eyebrow">🍕 Schnell. Frisch. Lecker.</p>
        <h1>{$t('home.hero_title')}</h1>
        <p class="hero-sub">{$t('home.hero_sub')}</p>
      </div>
    </div>

    <!-- ── Suche + Filter ── -->
    <div class="controls">
      <div class="search-wrap">
        <span class="search-icon">🔍</span>
        <input
          type="search"
          placeholder={$t('common.search_placeholder')}
          bind:value={suche}
          class="search-input"
        />
      </div>

      <div class="filter-wrap">
        <div class="filter-field">
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

        <div class="filter-field">
          <label for="mindest">{$t('home.min_order_label')} <strong class="gold-val">{maxMinBestellwert}€</strong></label>
          <input type="range" id="mindest" min="5" max="30" step="5" bind:value={maxMinBestellwert} />
        </div>
      </div>
    </div>

    <!-- ── Top 10 Scroll ── -->
    {#if gewaehlterTyp === "alle"}
      <div class="section-header">
        <h2>{$t('home.top10_title')}</h2>
        <span class="section-sub">{$t('home.top10_sub')}</span>
      </div>

      <div class="scroll-track">
        {#each top10Restaurants as restaurant, index}
          <a href="/restaurant/{restaurant.slug}" class="scroll-card">
            <span class="rank">{index + 1}</span>
            <div class="card-img emoji-bild">
              <span class="emoji-big">{restaurant.emoji}</span>
              <span class="star-badge">⭐ {schnittVon(restaurant).toFixed(1)}</span>
            </div>
            <!-- Blur-Label: das Apple-Musik-Muster -->
            <div class="card-blur-label">
              <span class="card-name">{restaurant.name}</span>
              <span class="card-meta">{restaurant.beschreibung}</span>
            </div>
            <div class="card-footer">
              <span class="chip">{restaurant.typ}</span>
              <span class="min-order">{$t('common.min')}: {restaurant.minBestell}€</span>
            </div>
          </a>
        {/each}
      </div>
    {/if}

    <!-- ── Alle Restaurants Grid ── -->
    <div class="section-header grid-section-header">
      <h2>{$t('home.discover')}</h2>
    </div>

    <div class="restaurant-grid">
      {#each gefilterteRestaurants as restaurant}
        <a href="/restaurant/{restaurant.slug}" class="grid-card">
          <div class="card-img emoji-bild">
            <span class="emoji-big">{restaurant.emoji}</span>
            <span class="star-badge">⭐ {schnittVon(restaurant).toFixed(1)}</span>
          </div>
          <!-- Blur-Label -->
          <div class="card-blur-label">
            <span class="card-name">{restaurant.name}</span>
            <span class="card-meta">{restaurant.beschreibung}</span>
          </div>
          <div class="card-footer">
            <span class="chip">{restaurant.typ}</span>
            <span class="min-order">Min: {restaurant.minBestell}€</span>
          </div>
        </a>
      {/each}
    </div>

  </div>
{/if}

<style>
  /* ══════════════════════════════════════════════════════════════
     DESIGN TOKENS
  ══════════════════════════════════════════════════════════════ */
  :root {
    --g1: #e6a800;
    --g2: #b87c00;
    --g3: #7a5000;
    --gold-text: #f9c932;
    --surface: rgba(255, 248, 220, 0.06);
    --surface-hover: rgba(255, 248, 220, 0.10);
    --border: rgba(230, 168, 0, 0.20);
    --border-hover: rgba(230, 168, 0, 0.50);
    --text: #f5f0e8;
    --text-muted: rgba(245, 240, 232, 0.55);
    --card-r: 18px;
  }

  /* ══════════════════════════════════════════════════════════════
     AUTH / LOGIN
  ══════════════════════════════════════════════════════════════ */
  .auth-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 90vh;
    padding: 24px;
  }

  .auth-card {
    background: rgba(20, 12, 0, 0.75);
    backdrop-filter: blur(32px) saturate(1.6);
    -webkit-backdrop-filter: blur(32px) saturate(1.6);
    border: 1px solid var(--border);
    border-radius: 24px;
    box-shadow: 0 24px 80px rgba(0,0,0,0.55), 0 0 0 0.5px rgba(230,168,0,0.12);
    padding: 36px 32px 32px;
    max-width: 460px;
    width: 100%;
  }

  /* Auth Hero Banner — ersetzt den lila Block */
  .auth-hero {
    text-align: center;
    margin-bottom: 28px;
  }

  .auth-hero-icon {
    font-size: 2.8rem;
    display: block;
    margin-bottom: 10px;
    filter: drop-shadow(0 4px 12px rgba(230,168,0,0.40));
  }

  .auth-hero h2 {
    font-size: 1.6rem;
    font-weight: 700;
    letter-spacing: -0.02em;
    color: #fff;
    margin: 0 0 6px;
  }

  .auth-hero p {
    color: var(--text-muted);
    font-size: 0.88rem;
    margin: 0;
  }

  /* Schritt-Punkte */
  .step-dots {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0;
    margin-bottom: 28px;
  }

  .dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: rgba(230,168,0,0.20);
    border: 1px solid rgba(230,168,0,0.35);
    transition: background 0.3s, transform 0.3s;
  }

  .dot.aktiv {
    background: linear-gradient(135deg, var(--g1), var(--g2));
    border-color: var(--g1);
    transform: scale(1.25);
    box-shadow: 0 0 8px rgba(230,168,0,0.50);
  }

  .dot-line {
    flex: 1;
    max-width: 60px;
    height: 1px;
    background: rgba(230,168,0,0.20);
    margin: 0 6px;
  }

  /* Auth Form */
  .auth-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .field label {
    font-size: 0.80rem;
    font-weight: 600;
    color: rgba(245,216,124,0.80);
    letter-spacing: 0.02em;
  }

  .field input,
  .field select {
    background: rgba(255,248,220,0.07) !important;
    border: 1px solid var(--border) !important;
    border-radius: 11px !important;
    color: var(--text) !important;
    padding: 11px 14px !important;
    font-size: 0.92rem !important;
    width: 100%;
    transition: border-color 0.18s, box-shadow 0.18s;
  }

  .field input:focus,
  .field select:focus {
    border-color: rgba(230,168,0,0.70) !important;
    box-shadow: 0 0 0 3px rgba(230,168,0,0.14) !important;
    outline: none !important;
  }

  /* Buttons */
  .gold-btn {
    background: linear-gradient(135deg, var(--g1), var(--g2)) !important;
    color: #1a0f00 !important;
    border: none !important;
    border-radius: 12px !important;
    padding: 13px 20px !important;
    font-size: 0.92rem !important;
    font-weight: 700 !important;
    cursor: pointer !important;
    box-shadow: 0 4px 16px rgba(230,168,0,0.30) !important;
    transition: opacity 0.15s, transform 0.15s !important;
    text-align: center;
    text-decoration: none;
    display: inline-block;
  }

  .gold-btn:hover { opacity: 0.88 !important; transform: scale(0.98) !important; }
  .gold-btn:disabled { opacity: 0.40 !important; cursor: not-allowed !important; }
  .block-btn { display: block !important; width: 100%; }
  .flex1 { flex: 1; }

  .ghost-btn {
    background: rgba(255,248,220,0.06) !important;
    border: 1px solid var(--border) !important;
    border-radius: 12px !important;
    color: var(--text-muted) !important;
    padding: 13px 18px !important;
    font-size: 0.88rem !important;
    font-weight: 600 !important;
    cursor: pointer !important;
    transition: background 0.15s !important;
  }
  .ghost-btn:hover { background: rgba(255,248,220,0.12) !important; }

  .plus-btn {
    width: 40px !important;
    height: 40px !important;
    min-width: 40px !important;
    padding: 0 !important;
    border-radius: 10px !important;
    background: rgba(230,168,0,0.15) !important;
    border: 1px solid var(--border) !important;
    color: var(--gold-text) !important;
    font-size: 1.3rem !important;
    font-weight: 700 !important;
    cursor: pointer !important;
  }

  .btn-row { display: flex; gap: 10px; }
  .addr-row { display: flex; gap: 10px; }
  .name-row { display: flex; gap: 10px; }
  .grow { flex: 1; }
  .narrow { width: 90px; min-width: 70px; }
  .input-plus { display: flex; gap: 6px; align-items: center; }

  /* Passwort-Stärke */
  .pw-block { margin-top: 8px; }
  .pw-track { background: rgba(255,248,220,0.10); border-radius: 6px; height: 6px; overflow: hidden; }
  .pw-fill { height: 100%; border-radius: 6px; transition: width 0.3s, background 0.3s; }
  .pw-label { font-size: 0.78rem; font-weight: 700; display: block; margin-top: 5px; }
  .pw-rules { list-style: none; padding: 0; margin: 8px 0 0; display: flex; flex-wrap: wrap; gap: 5px 14px; }
  .pw-rules li { font-size: 0.76rem; color: rgba(245,240,232,0.40); }
  .pw-rules li::before { content: "✗ "; color: #ff453a; }
  .pw-rules li.ok { color: #30d158; }
  .pw-rules li.ok::before { content: "✓ "; color: #30d158; }

  /* Feedback */
  .err { color: #ff453a; font-size: 0.82rem; font-weight: 600; margin: 2px 0 0; }
  .err-box { background: rgba(255, 69, 58, 0.10); border: 1px solid rgba(255, 69, 58, 0.35); border-radius: 10px; padding: 10px 12px; margin: 4px 0; }
  .err-list { margin: 6px 0 0; padding-left: 18px; }
  .err-list li { color: #ff6961; font-size: 0.8rem; line-height: 1.5; }
  .auth-hint { text-align: center; font-size: 0.84rem; color: var(--text-muted); margin: 4px 0 0; }
  .auth-hint a { color: var(--gold-text); font-weight: 600; }

  /* ══════════════════════════════════════════════════════════════
     STARTSEITE
  ══════════════════════════════════════════════════════════════ */
  .home {
    max-width: 1200px;
    margin: 0 auto;
    padding-bottom: 60px;
  }

  /* Hero */
  .hero {
    position: relative;
    overflow: hidden;
    border-radius: 24px;
    margin-bottom: 32px;
    padding: 56px 36px;
    background: linear-gradient(135deg, rgba(30,18,0,0.95) 0%, rgba(14,9,0,0.98) 100%);
    border: 1px solid rgba(230,168,0,0.22);
    box-shadow: 0 12px 48px rgba(0,0,0,0.50), inset 0 1px 0 rgba(230,168,0,0.15);
    cursor: pointer;
    user-select: none;
  }

  .hero-glow {
    position: absolute;
    top: -60px;
    left: 10%;
    width: 60%;
    height: 200px;
    background: radial-gradient(ellipse, rgba(230,168,0,0.22) 0%, transparent 70%);
    pointer-events: none;
  }

  .hero-content { position: relative; z-index: 1; }

  .hero-eyebrow {
    font-size: 0.80rem;
    font-weight: 600;
    letter-spacing: 0.08em;
    color: var(--gold-text);
    text-transform: uppercase;
    margin: 0 0 12px;
  }

  .hero h1 {
    font-size: clamp(2rem, 5vw, 3.2rem);
    font-weight: 800;
    letter-spacing: -0.04em;
    color: #fff;
    margin: 0 0 12px;
    line-height: 1.1;
  }

  .hero-sub {
    color: rgba(245,240,232,0.60);
    font-size: 1.05rem;
    margin: 0;
    line-height: 1.5;
  }

  /* Suche + Filter */
  .controls {
    display: flex;
    flex-direction: column;
    gap: 14px;
    margin-bottom: 40px;
  }

  .search-wrap {
    position: relative;
  }

  .search-icon {
    position: absolute;
    left: 14px;
    top: 50%;
    transform: translateY(-50%);
    font-size: 1rem;
    pointer-events: none;
  }

  .search-input {
    width: 100%;
    padding: 13px 14px 13px 42px !important;
    background: rgba(255,248,220,0.06) !important;
    border: 1px solid var(--border) !important;
    border-radius: 14px !important;
    color: var(--text) !important;
    font-size: 0.95rem !important;
    box-sizing: border-box;
    transition: border-color 0.18s !important;
  }

  .search-input:focus {
    border-color: rgba(230,168,0,0.65) !important;
    box-shadow: 0 0 0 3px rgba(230,168,0,0.12) !important;
    outline: none !important;
  }

  .filter-wrap {
    display: flex;
    gap: 20px;
    background: rgba(255,248,220,0.04);
    border: 1px solid var(--border);
    border-radius: 16px;
    padding: 18px 22px;
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
  }

  .filter-field {
    display: flex;
    flex-direction: column;
    gap: 8px;
    flex: 1;
  }

  .filter-field label {
    font-size: 0.80rem;
    font-weight: 600;
    color: rgba(245,216,124,0.70);
    letter-spacing: 0.02em;
  }

  .gold-val { color: var(--gold-text); }

  .filter-field select {
    background: rgba(255,248,220,0.07) !important;
    border: 1px solid var(--border) !important;
    border-radius: 10px !important;
    color: var(--text) !important;
    padding: 9px 12px !important;
    font-size: 0.88rem !important;
    cursor: pointer;
  }

  .filter-field input[type="range"] {
    accent-color: var(--g1);
    cursor: pointer;
    border: none !important;
    background: none !important;
    padding: 4px 0 !important;
    border-radius: 0 !important;
  }

  /* Abschnitt-Überschriften */
  .section-header {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-bottom: 20px;
  }

  .section-header h2 {
    font-size: 1.5rem;
    font-weight: 700;
    letter-spacing: -0.02em;
    color: #fff;
    margin: 0;
  }

  .section-sub {
    font-size: 0.84rem;
    color: var(--text-muted);
  }

  .grid-section-header {
    margin-top: 56px;
    padding-top: 36px;
    border-top: 1px solid rgba(230,168,0,0.10);
  }

  /* ── Horizontaler Scroll-Track (Top 10) ── */
  .scroll-track {
    display: flex;
    gap: 18px;
    overflow-x: auto;
    padding: 8px 4px 28px;
    scroll-behavior: smooth;
    -webkit-overflow-scrolling: touch;
    margin-bottom: 8px;
  }

  .scroll-track::-webkit-scrollbar { height: 5px; }
  .scroll-track::-webkit-scrollbar-track { background: transparent; }
  .scroll-track::-webkit-scrollbar-thumb { background: rgba(230,168,0,0.35); border-radius: 3px; }

  /* ── Karten (Scroll + Grid) ── */
  .scroll-card,
  .grid-card {
    position: relative;
    overflow: hidden;
    border-radius: var(--card-r);
    border: 1px solid rgba(230,168,0,0.20);
    box-shadow: 0 8px 28px rgba(0,0,0,0.45);
    background: #100a00;
    text-decoration: none;
    color: inherit;
    display: flex;
    flex-direction: column;
    transition: transform 0.24s cubic-bezier(0.4,0,0.2,1), box-shadow 0.24s;
  }

  .scroll-card { flex: 0 0 240px; }

  .scroll-card:hover,
  .grid-card:hover {
    transform: scale(1.025) translateY(-3px);
    box-shadow: 0 16px 48px rgba(0,0,0,0.55), 0 0 0 1px rgba(230,168,0,0.45);
  }

  /* Rang-Nummer (groß links oben) */
  .rank {
    position: absolute;
    top: -16px;
    left: -10px;
    font-size: 5.5rem;
    font-weight: 900;
    line-height: 1;
    letter-spacing: -4px;
    user-select: none;
    z-index: 4;
    background: linear-gradient(160deg, rgba(249,201,50,0.85) 0%, rgba(184,124,0,0.60) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    filter: drop-shadow(0 4px 10px rgba(230,168,0,0.30));
  }

  /* Bild-Platzhalter (Emoji) */
  .card-img {
    width: 100%;
    aspect-ratio: 4/3;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    background: linear-gradient(135deg, #2a1f00, #0d0800);
  }

  .emoji-big { font-size: 3.6rem; }

  .star-badge {
    position: absolute;
    top: 10px;
    right: 10px;
    background: rgba(14,9,0,0.72);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    border: 1px solid rgba(230,168,0,0.30);
    color: #f9c932;
    font-size: 0.78rem;
    font-weight: 700;
    padding: 4px 9px;
    border-radius: 8px;
  }

  /* 🍎 Apple-Musik Blur-Label */
  .card-blur-label {
    position: relative;
    padding: 10px 13px 8px;
    background: rgba(14,9,0,0.65);
    backdrop-filter: blur(18px) saturate(1.5);
    -webkit-backdrop-filter: blur(18px) saturate(1.5);
    border-top: 0.5px solid rgba(230,168,0,0.16);
    display: flex;
    flex-direction: column;
    gap: 3px;
  }

  .card-name {
    font-size: 0.90rem;
    font-weight: 700;
    color: #fff;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    letter-spacing: -0.01em;
  }

  .card-meta {
    font-size: 0.74rem;
    color: rgba(249,201,50,0.72);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  /* Card Footer */
  .card-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 9px 13px 12px;
    background: rgba(14,9,0,0.50);
  }

  .chip {
    background: rgba(230,168,0,0.15);
    border: 1px solid rgba(230,168,0,0.28);
    color: var(--gold-text);
    font-size: 0.72rem;
    font-weight: 700;
    padding: 3px 10px;
    border-radius: 20px;
    letter-spacing: 0.02em;
  }

  .min-order {
    font-size: 0.76rem;
    color: var(--text-muted);
    font-weight: 500;
  }

  /* Restaurant Grid */
  .restaurant-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 22px;
  }

  /* ══════════════════════════════════════════════════════════════
     HELL-MODUS OVERRIDES
  ══════════════════════════════════════════════════════════════ */
  :global(html[data-theme='light']) .auth-card {
    background: rgba(255, 252, 240, 0.90);
    border-color: rgba(184,124,0,0.28);
  }

  :global(html[data-theme='light']) .auth-hero h2 { color: #1a0f00; }
  :global(html[data-theme='light']) .auth-hero p { color: #7a5000; }
  :global(html[data-theme='light']) .field label { color: #7a5000; }
  :global(html[data-theme='light']) .field input,
  :global(html[data-theme='light']) .field select {
    background: rgba(255,248,220,0.70) !important;
    color: #1a0f00 !important;
  }

  :global(html[data-theme='light']) .hero { background: linear-gradient(135deg, #fff8e0, #fef3c0); }
  :global(html[data-theme='light']) .hero h1 { color: #1a0f00; }
  :global(html[data-theme='light']) .hero-sub { color: rgba(26,15,0,0.60); }
  :global(html[data-theme='light']) .hero-eyebrow { color: #b87c00; }

  :global(html[data-theme='light']) .filter-wrap { background: rgba(255,248,220,0.65); }
  :global(html[data-theme='light']) .filter-field label { color: #7a5000; }
  :global(html[data-theme='light']) .search-input { background: rgba(255,248,220,0.80) !important; color: #1a0f00 !important; }

  :global(html[data-theme='light']) .card-blur-label {
    background: rgba(255,248,220,0.80);
    border-top-color: rgba(184,124,0,0.18);
  }
  :global(html[data-theme='light']) .card-name { color: #1a0f00; }
  :global(html[data-theme='light']) .card-footer { background: rgba(255,248,220,0.60); }

  :global(html[data-theme='light']) .scroll-card,
  :global(html[data-theme='light']) .grid-card {
    background: #fff8e7;
    border-color: rgba(184,124,0,0.22);
  }

  :global(html[data-theme='light']) .section-header h2 { color: #1a0f00; }
</style>