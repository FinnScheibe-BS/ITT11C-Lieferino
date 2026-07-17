<script>
  import { onMount } from 'svelte';
  import { warenkorb } from '$lib/stores/cart.js';
  import { pruefePasswortStaerke } from '$lib/services/passwort.js';
  import { pruefeAdresse } from '$lib/services/adresse.js';
  import { api } from '$lib/api/api.js';
  import AuthVollenden from '$lib/components/auth/AuthVollenden.svelte';
  import { aktiveRestaurants } from '$lib/stores/lieferanten.js';
  import { bewertungsSchnitt } from '$lib/stores/bewertungsSchnitt.js';
  import { eingeloggt, login } from '$lib/stores/auth.js';
  import { t } from '$lib/utils/i18n.js';

  // ════════════════ AUTH STATE ════════════════
  let loginSchritt = $state(1);
  let emailInput = $state("");
  let passwortInput = $state("");
  let passwortStaerke = $derived(pruefePasswortStaerke(passwortInput));
  let zeigeAbschluss = $state(false);
  let prueftAdresse = $state(false);
  let adressFehler = $state("");
  let registrierFehler = $state("");
  let registrierAnforderungen = $state([]);
  let emailFehler = $state("");
  let prueftEmail = $state(false);
  let agbAkzeptiert = $state(false);
  let vornameInput = $state("");
  let nachnameInput = $state("");
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

  // ════════════════ AUTH FUNCTIONS ════════════════
  async function geheZuSchritt2(e) {
    e.preventDefault();
    if (!passwortStaerke.istSicher) return;
    if (emailInput.trim() === "" || passwortInput.trim() === "") return;
    emailFehler = "";
    prueftEmail = true;
    try {
      const res = await api('/api/auth/email-check', { method: 'POST', body: { email: emailInput.trim() } });
      prueftEmail = false;
      if (res.ok && res.daten && res.daten.frei === false) {
        emailFehler = "Diese E-Mail ist bereits registriert. Bitte melde dich an. 📧";
        return;
      }
    } catch (err) {
      console.error('E-Mail-Prüfung fehlgeschlagen:', err);
      prueftEmail = false;
    }
    registrierFehler = "";
    loginSchritt = 2;
  }

  function geheZuSchritt3(e) {
    e.preventDefault();
    if (!nameGueltig(vornameInput) || !nameGueltig(nachnameInput)) {
      namensFehler = "Vor- und Nachname brauchen mind. 2 Buchstaben und dürfen keine Zahlen enthalten. ✍️";
      return;
    }
    namensFehler = "";
    if (!geburtsdatumInput) {
      altersFehler = "Bitte gib dein Geburtsdatum an. 📅";
      return;
    }
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
    if (!agbAkzeptiert) {
      adressFehler = "Bitte akzeptiere die AGB und die Datenschutzerklärung. ✅";
      return;
    }
    if (!strasseInput.trim()) {
      adressFehler = "Bitte gib deine Straße an. 📍";
      return;
    }
    if (!hausnummerInput.trim()) {
      adressFehler = "Bitte gib deine Hausnummer an. 🏠";
      return;
    }
    if (!plzInput.trim()) {
      adressFehler = "Bitte gib deine PLZ an. 🔢";
      return;
    }
    if (!ortInput.trim()) {
      adressFehler = "Bitte gib deinen Ort an. 🌆";
      return;
    }
    if (!/^\d{5}$/.test(plzInput.trim())) {
      adressFehler = "Die PLZ muss aus genau 5 Ziffern bestehen. 🔢";
      return;
    }
    adressFehler = "";
    prueftAdresse = true;
    try {
      const adressErgebnis = await pruefeAdresse({
        strasse: strasseInput,
        hausnummer: hausnummerInput,
        plz: plzInput,
        ort: ortInput
      });
      if (adressErgebnis.fehler) {
        adressFehler = adressErgebnis.fehler;
        prueftAdresse = false;
        return;
      }
      prueftAdresse = false;
    } catch (err) {
      console.error('Adressprüfung fehlgeschlagen:', err);
      prueftAdresse = false;
    }
    const userDaten = {
      email: emailInput,
      passwort: passwortInput,
      vorname: vornameInput,
      nachname: nachnameInput,
      strasse: strasseInput,
      hausnummer: hausnummerInput,
      plz: plzInput,
      ort: ortInput,
      geburtsdatum: geburtsdatumInput
    };
    localStorage.setItem("lieferino_user", JSON.stringify(userDaten));
    registrierFehler = "";
    registrierAnforderungen = [];
    try {
      const reg = await api('/api/auth/register', {
        method: 'POST',
        body: {
          email: emailInput,
          passwort: passwortInput,
          vorname: vornameInput,
          nachname: nachnameInput,
          geburtsdatum: geburtsdatumInput
        }
      });
      if (reg.ok && reg.daten?.needsVerification) {
        zeigeAbschluss = true;
        return;
      }
      if (reg.offline) {
        login();
        return;
      }
      registrierFehler = reg.daten?.fehler || "Registrierung fehlgeschlagen.";
      registrierAnforderungen = reg.daten?.anforderungen || [];
    } catch (err) {
      console.error('Registrierung fehlgeschlagen:', err);
      registrierFehler = "Netzwerkfehler.";
    }
  }

  async function abschlussFertig() {
    try {
      await api('/api/me', {
        method: 'PUT',
        body: {
          vorname: vornameInput, nachname: nachnameInput, geburtsdatum: geburtsdatumInput,
          adressen: [{ label: 'Zuhause', strasse: strasseInput, hausnummer: hausnummerInput, plz: plzInput, ort: ortInput }]
        }
      });
      login();
    } catch (err) {
      console.error('Abschluss fehlgeschlagen:', err);
      login();
    }
  }

  // ════════════════ RESTAURANT FILTER ════════════════
  let gewaehlterTyp = $state("alle");
  let sortierung = $state("standard");
  let maxMinBestellwert = $state(30);
  let suche = $state("");

  let anzahl = $derived($warenkorb?.length ?? 0);
  let warenkorbSumme = $derived($warenkorb?.reduce((sum, item) => sum + (item.preis ?? 0), 0) ?? 0);

  function schnittVon(r) {
    const s = $bewertungsSchnitt[r.slug];
    return s && s.anzahl > 0 ? s.schnitt : r.bewertung;
  }

  let gefilterteRestaurants = $derived(
    $aktiveRestaurants.filter(r => {
      let typPasst = gewaehlterTyp === "alle" || r.typ === gewaehlterTyp;
      let preisPasst = r.minBestell <= maxMinBestellwert;
      let suchePasst = r.name.toLowerCase().includes(suche.toLowerCase());
      return typPasst && preisPasst && suchePasst;
    })
  );

  // ⭐ TOP 5 statt TOP 10
  let top5Restaurants = $derived(
    [...$aktiveRestaurants].sort((a, b) => schnittVon(b) - schnittVon(a)).slice(0, 5)
  );

  let typen = $derived([...new Set($aktiveRestaurants.map((r) => r.typ))]);
</script>

<!-- ░░░ NICHT EINGELOGGT: Registrierung ░░░ -->
{#if !$eingeloggt}
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
      {#if !zeigeAbschluss}
        <div class="step-dots">
          <span class="dot" class:aktiv={loginSchritt >= 1}></span>
          <span class="dot-line"></span>
          <span class="dot" class:aktiv={loginSchritt >= 2}></span>
          <span class="dot-line"></span>
          <span class="dot" class:aktiv={loginSchritt >= 3}></span>
        </div>
      {/if}
      {#if zeigeAbschluss}
        <AuthVollenden start="verify" email={emailInput} onFertig={abschlussFertig} />
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
          {#if emailFehler}<p class="err">{emailFehler}</p>{/if}
          <button type="submit" class="gold-btn block-btn" disabled={!passwortStaerke.istSicher || prueftEmail}>
            {prueftEmail ? 'Prüfe E-Mail…' : 'Weiter zu deinen Details →'}
          </button>
          <p class="auth-hint">Schon registriert? <a href="/login">Hier einloggen 🔑</a></p>
        </form>
      {:else if loginSchritt === 2}
        <form onsubmit={geheZuSchritt3} class="auth-form">
          <div class="name-row">
            <div class="field grow">
              <label for="vorname">Vorname</label>
              <input type="text" id="vorname" placeholder="Max" bind:value={vornameInput} required />
            </div>
            <div class="field grow">
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
            <button type="button" class="ghost-btn flex1" onclick={() => loginSchritt = 1}>← Zurück</button>
            <button type="submit" class="gold-btn flex1">Weiter zur Adresse 🏠</button>
          </div>
        </form>
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
          <label class="agb-check">
            <input type="checkbox" bind:checked={agbAkzeptiert} />
            <span>Ich akzeptiere die <a href="/agb" target="_blank" rel="noopener">AGB</a> und die <a href="/datenschutz" target="_blank" rel="noopener">Datenschutzerklärung</a>.</span>
          </label>
          <div class="btn-row">
            <button type="button" class="ghost-btn flex1" onclick={() => loginSchritt = 2}>← Zurück</button>
            <button type="submit" class="gold-btn flex1" disabled={prueftAdresse || !agbAkzeptiert}>
              {prueftAdresse ? "Prüfe Adresse… ⏳" : "Registrierung abschließen 🎉"}
            </button>
          </div>
        </form>
      {/if}
    </div>
  </div>
{:else}
  <!-- ░░░ HERO BANNER ░░░ -->
  <section class="hero-banner">
    <div class="hero-bg-layer"></div>
    <div class="hero-gradient-orb orb-1"></div>
    <div class="hero-gradient-orb orb-2"></div>
    <div class="hero-gradient-orb orb-3"></div>
    <div class="hero-pattern"></div>
    <div class="hero-container">
      <div class="hero-text-col">
        <div class="hero-badge-wrapper">
          <span class="hero-badge">
            <span class="badge-dot"></span>
            <span>Das edelste Essen, direkt zu Ihnen nach hause</span>
          </span>
        </div>
        <h1 class="hero-headline">
          <span class="headline-line-1">Essen, das</span>
          <span class="headline-line-2">
            <span class="text-gold">Begeisterung</span>
            <span class="text-white">liefert.</span>
          </span>
        </h1>
        <p class="hero-description">
          Von lokalen Favoriten bis zu gastronomischen Highlights –
          <strong>frisch zubereitet</strong> und in unter 30 Minuten an deiner Tür.
        </p>
        <div class="hero-cta-group">
          <a href="#restaurants" class="cta-primary">
            <span>Jetzt bestellen</span>
            <svg class="cta-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </a>
          <a href="/bestellungen" class="cta-secondary">
            <span>Meine Bestellungen</span>
          </a>
        </div>
        <div class="trust-indicators">
          <div class="trust-item">
            <span class="trust-icon">⚡</span>
            <span class="trust-text"><strong>30 min</strong> Ø Lieferzeit</span>
          </div>
          <div class="trust-divider"></div>
          <div class="trust-item">
            <span class="trust-icon">⭐</span>
            <span class="trust-text"><strong>4.8/5</strong> Kundenbewertung</span>
          </div>
          <div class="trust-divider"></div>
          <div class="trust-item">
            <span class="trust-icon">🍽</span>
            <span class="trust-text"><strong>500+</strong> Restaurants</span>
          </div>
        </div>
      </div>
      <div class="hero-visual-col">
        <div class="cart-preview-card">
          <div class="card-header">
            <span class="cart-icon-wrapper">🛒</span>
            <span class="card-title">Dein Warenkorb</span>
          </div>
          <div class="card-body">
            {#if anzahl > 0}
              <div class="cart-items-summary">
                <span class="item-count">{anzahl} Artikel</span>
                <span class="item-total">{warenkorbSumme.toFixed(2)} €</span>
              </div>
              <div class="cart-progress">
                <div class="progress-bar">
                  <div class="progress-fill" style="width: {Math.min(anzahl / 5 * 100, 100)}%"></div>
                </div>
                <span class="progress-text">Noch {Math.max(0, 5 - anzahl)} bis zur kostenlosen Lieferung</span>
              </div>
            {:else}
              <div class="cart-empty-state">
                <span class="empty-icon">🍕</span>
                <p>Noch keine Artikel im Warenkorb</p>
              </div>
            {/if}
          </div>
          <a href="/cart" class="card-footer">
            <span>Zum Warenkorb</span>
            <span class="footer-arrow">→</span>
          </a>
        </div>
      </div>
    </div>
    <div class="floating-foods" aria-hidden="true">
      <span class="float-food food-1">🍕</span>
      <span class="float-food food-2">🍔</span>
      <span class="float-food food-3">🍣</span>
      <span class="float-food food-4">🌮</span>
      <span class="float-food food-5">🍰</span>
      <span class="float-food food-6">🥗</span>
    </div>
    <div class="scroll-indicator">
      <span class="scroll-text">Scrollen</span>
      <div class="scroll-mouse">
        <div class="scroll-wheel"></div>
      </div>
    </div>
  </section>

  <!-- ░░░ CONTROLS & FILTER ░░░ -->
  <div class="controls">
    <div class="search-wrap">
      <span class="search-icon">🔍</span>
      <input type="search" placeholder={$t('common.search_placeholder')} bind:value={suche} class="search-input" />
    </div>
    <div class="filter-wrap">
      <div class="filter-field">
        <label for="kueche">{$t('home.cuisine_label')}</label>
        <select id="kueche" bind:value={gewaehlterTyp}>
          <option value="alle">{$t('common.all_cuisines')}</option>
          {#each typen as typ}
            <option value={typ}>{typ}</option>
          {/each}
        </select>
      </div>
      <div class="filter-field">
        <label for="mindest">{$t('home.min_order_label')} <strong class="gold-val">{maxMinBestellwert}€</strong></label>
        <input type="range" id="mindest" min="5" max="30" step="5" bind:value={maxMinBestellwert} />
      </div>
    </div>
  </div>

  <!-- ░░░ TOP 5 RESTAURANTS ░░░ -->
  {#if gewaehlterTyp === "alle"}
    <div class="section-header">
      <h2>{$t('home.top5_title') || 'Top 5 Restaurants'}</h2>
      <span class="section-sub">{$t('home.top5_sub') || 'Die besten Bewertungen dieser Woche'}</span>
    </div>
    <div class="top5-scroll-track" id="restaurants">
      {#each top5Restaurants as restaurant, index}
        <a href="/restaurant/{restaurant.slug}" class="top5-card restaurant-card">
          <span class="rank">{index + 1}</span>
          <div class="card-bild">
            <div class="bewertung-badge">
              <span class="stern">⭐️</span>
              <span class="bewertung-wert">{schnittVon(restaurant).toFixed(1)}</span>
            </div>
            <span class="emoji-bild">{restaurant.emoji}</span>
          </div>
          <div class="card-info">
            <h3>{restaurant.name}</h3>
            <div class="card-meta">
              <span class="tag">{restaurant.typ}</span>
            </div>
          </div>
        </a>
      {/each}
    </div>
  {/if}

  <!-- ░░░ ALL RESTAURANTS GRID ░░░ -->
  <div class="section-header grid-section-header">
    <h2>{$t('home.discover')}</h2>
  </div>
  <div class="restaurant-grid">
    {#each gefilterteRestaurants as restaurant}
      <a href="/restaurant/{restaurant.slug}" class="restaurant-card">
        <div class="card-bild">
          <div class="bewertung-badge">
            <span class="stern">⭐️</span>
            <span class="bewertung-wert">{schnittVon(restaurant).toFixed(1)}</span>
          </div>
          <span class="emoji-bild">{restaurant.emoji}</span>
        </div>
        <div class="card-info">
          <h3>{restaurant.name}</h3>
          <div class="card-meta">
            <span class="tag">{restaurant.typ}</span>
          </div>
        </div>
      </a>
    {/each}
  </div>
{/if}

<style>
/* ════════════════════════════════════════════════════════════════
     DESIGN TOKENS
     ════════════════════════════════════════════════════════════════ */
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
:global(html[data-theme='light']) {
  --gold-text: #9a6600;
  --surface: rgba(255, 252, 235, 0.75);
  --surface-hover: rgba(255, 250, 230, 0.92);
  --border: rgba(230, 168, 0, 0.35);
  --border-hover: rgba(230, 168, 0, 0.65);
  --text: #1a1200;
  --text-muted: rgba(60, 45, 10, 0.72);
}

/* ════════════════════════════════════════════════════════════════
     HERO BANNER
     ════════════════════════════════════════════════════════════════ */
.hero-banner {
  position: relative;
  min-height: 95vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 140px 24px 80px;
  overflow: hidden;
  isolation: isolate;
  width: 100vw;
  margin-left: calc(-50vw + 50%);
  margin-right: calc(-50vw + 50%);
}
.hero-bg-layer {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(ellipse at 20% 10%, rgba(230,168,0,0.15) 0%, transparent 50%),
    radial-gradient(ellipse at 80% 90%, rgba(184,124,0,0.12) 0%, transparent 50%),
    linear-gradient(180deg, #0f0a05 0%, #1a1205 100%);
  z-index: -3;
}
.hero-gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: orb-float 12s ease-in-out infinite;
  z-index: -2;
}
.orb-1 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(230,168,0,0.4) 0%, transparent 70%);
  top: -150px;
  left: -100px;
  animation-delay: 0s;
}
.orb-2 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(249,201,50,0.3) 0%, transparent 70%);
  bottom: -100px;
  right: -50px;
  animation-delay: 4s;
}
.orb-3 {
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(184,124,0,0.35) 0%, transparent 70%);
  top: 40%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: 8s;
}
@keyframes orb-float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(20px, -30px) scale(1.05); }
  50% { transform: translate(-15px, 20px) scale(0.95); }
  75% { transform: translate(25px, 15px) scale(1.02); }
}
.hero-pattern {
  position: absolute;
  inset: 0;
  background-image:
    repeating-linear-gradient(0deg, transparent, transparent 40px, rgba(230,168,0,0.03) 40px, rgba(230,168,0,0.03) 41px),
    repeating-linear-gradient(90deg, transparent, transparent 40px, rgba(230,168,0,0.03) 40px, rgba(230,168,0,0.03) 41px);
  mask-image: radial-gradient(ellipse at 50% 50%, black 0%, transparent 70%);
  -webkit-mask-image: radial-gradient(ellipse at 50% 50%, black 0%, transparent 70%);
  z-index: -1;
  pointer-events: none;
}
.hero-container {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 60px;
  max-width: 1200px;
  width: 100%;
  align-items: center;
  position: relative;
  z-index: 10;
}
.hero-text-col {
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.hero-badge-wrapper {
  display: inline-flex;
  width: fit-content;
}
.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: rgba(230,168,0,0.12);
  border: 1px solid rgba(230,168,0,0.35);
  border-radius: 50px;
  font-size: 0.8rem;
  font-weight: 600;
  color: #f9c932;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  letter-spacing: 0.02em;
  text-transform: uppercase;
}
.badge-dot {
  width: 8px;
  height: 8px;
  background: #f9c932;
  border-radius: 50%;
  animation: pulse-dot 2s ease-in-out infinite;
}
@keyframes pulse-dot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(1.2); }
}
.hero-headline {
  font-size: clamp(2.5rem, 6vw, 4.5rem);
  font-weight: 800;
  line-height: 1.05;
  letter-spacing: -0.04em;
  margin: 0;
}
.headline-line-1 {
  display: block;
  color: #fff;
  margin-bottom: 4px;
}
.headline-line-2 {
  display: block;
}
.text-gold {
  background: linear-gradient(135deg, #f9c932 0%, #e6a800 50%, #fde8a0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.text-white {
  color: #fff;
  margin-left: 12px;
}
.hero-description {
  font-size: 1.1rem;
  line-height: 1.7;
  color: rgba(245,240,232,0.75);
  max-width: 540px;
  margin: 0;
}
.hero-description strong {
  color: #f9c932;
  font-weight: 700;
}
.hero-cta-group {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  margin-top: 8px;
}
.cta-primary {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 16px 32px;
  background: linear-gradient(135deg, #e6a800 0%, #b87c00 100%);
  color: #1a0f00;
  border-radius: 14px;
  font-size: 1rem;
  font-weight: 700;
  text-decoration: none;
  box-shadow: 0 8px 32px rgba(230,168,0,0.35);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.cta-primary:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 40px rgba(230,168,0,0.5);
}
.cta-arrow {
  width: 20px;
  height: 20px;
  transition: transform 0.3s ease;
}
.cta-primary:hover .cta-arrow {
  transform: translateX(4px);
}
.cta-secondary {
  display: inline-flex;
  align-items: center;
  padding: 16px 28px;
  background: rgba(255,248,220,0.06);
  color: #f9c932;
  border: 1px solid rgba(230,168,0,0.35);
  border-radius: 14px;
  font-size: 1rem;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.3s ease;
}
.cta-secondary:hover {
  background: rgba(230,168,0,0.12);
  border-color: rgba(230,168,0,0.6);
  transform: translateY(-2px);
}
.trust-indicators {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-top: 16px;
  padding: 20px 24px;
  background: rgba(255,248,220,0.04);
  border: 1px solid rgba(230,168,0,0.15);
  border-radius: 16px;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}
.trust-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.85rem;
  color: rgba(245,240,232,0.7);
}
.trust-icon {
  font-size: 1.2rem;
}
.trust-text strong {
  color: #f9c932;
}
.trust-divider {
  width: 1px;
  height: 32px;
  background: rgba(230,168,0,0.25);
}
.hero-visual-col {
  display: flex;
  justify-content: center;
  align-items: center;
}
.cart-preview-card {
  width: 100%;
  max-width: 380px;
  background: rgba(20,14,5,0.85);
  backdrop-filter: blur(40px) saturate(1.8);
  -webkit-backdrop-filter: blur(40px) saturate(1.8);
  border: 1px solid rgba(230,168,0,0.25);
  border-radius: 24px;
  box-shadow: 0 24px 80px rgba(0,0,0,0.5), 0 0 0 1px rgba(230,168,0,0.1), inset 0 1px 0 rgba(255,255,255,0.05);
  overflow: hidden;
  animation: card-float 6s ease-in-out infinite;
}
@keyframes card-float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-12px); }
}
.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 24px;
  border-bottom: 1px solid rgba(230,168,0,0.15);
}
.cart-icon-wrapper {
  font-size: 1.4rem;
}
.card-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: #f9c932;
}
.card-body {
  padding: 24px;
}
.cart-items-summary {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
}
.item-count {
  color: rgba(245,240,232,0.7);
  font-size: 0.9rem;
}
.item-total {
  color: #fff;
  font-weight: 700;
  font-size: 1.1rem;
}
.cart-progress {
  margin-top: 12px;
}
.progress-bar {
  height: 6px;
  background: rgba(230,168,0,0.15);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 8px;
}
.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #e6a800, #f9c932);
  border-radius: 3px;
  transition: width 0.3s ease;
}
.progress-text {
  font-size: 0.8rem;
  color: rgba(245,240,232,0.6);
}
.cart-empty-state {
  text-align: center;
  padding: 20px 0;
}
.empty-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 12px;
  opacity: 0.7;
}
.cart-empty-state p {
  color: rgba(245,240,232,0.6);
  font-size: 0.95rem;
  margin: 0;
}
.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: rgba(230,168,0,0.08);
  border-top: 1px solid rgba(230,168,0,0.15);
  text-decoration: none;
  color: #f9c932;
  font-weight: 600;
  font-size: 0.95rem;
  transition: all 0.3s ease;
}
.card-footer:hover {
  background: rgba(230,168,0,0.15);
}
.footer-arrow {
  transition: transform 0.3s ease;
}
.card-footer:hover .footer-arrow {
  transform: translateX(4px);
}
.floating-foods {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 1;
  overflow: hidden;
}
.float-food {
  position: absolute;
  font-size: 2.5rem;
  opacity: 0.15;
  animation: float-food 8s ease-in-out infinite;
  filter: drop-shadow(0 4px 8px rgba(0,0,0,0.3));
}
.food-1 { top: 15%; left: 10%; animation-delay: 0s; }
.food-2 { top: 25%; right: 15%; animation-delay: 1.5s; }
.food-3 { bottom: 30%; left: 15%; animation-delay: 3s; }
.food-4 { bottom: 20%; right: 20%; animation-delay: 4.5s; }
.food-5 { top: 50%; left: 5%; animation-delay: 6s; }
.food-6 { top: 60%; right: 10%; animation-delay: 2s; }
@keyframes float-food {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(5deg); }
}
.scroll-indicator {
  position: absolute;
  bottom: 32px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  z-index: 10;
}
.scroll-text {
  font-size: 0.75rem;
  color: rgba(245,240,232,0.5);
  text-transform: uppercase;
  letter-spacing: 0.1em;
}
.scroll-mouse {
  width: 26px;
  height: 40px;
  border: 2px solid rgba(230,168,0,0.4);
  border-radius: 13px;
  position: relative;
}
.scroll-wheel {
  width: 4px;
  height: 8px;
  background: #f9c932;
  border-radius: 2px;
  position: absolute;
  top: 8px;
  left: 50%;
  transform: translateX(-50%);
  animation: scroll-wheel 1.5s ease-in-out infinite;
}
@keyframes scroll-wheel {
  0%, 100% { opacity: 1; transform: translateX(-50%) translateY(0); }
  50% { opacity: 0.3; transform: translateX(-50%) translateY(8px); }
}

/* ════════════════════════════════════════════════════════════════
     CONTROLS SECTION
     ════════════════════════════════════════════════════════════════ */
.controls {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  justify-content: space-between;
  align-items: flex-end;
  max-width: 1200px;
  margin: 0 auto 32px;
  padding: 0 24px;
}
.search-wrap {
  position: relative;
  flex: 1;
  min-width: 280px;
  max-width: 400px;
}
.search-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 1.1rem;
  opacity: 0.6;
}
.search-input {
  width: 100%;
  padding: 14px 16px 14px 48px;
  background: rgba(255,248,220,0.05);
  border: 1px solid rgba(230,168,0,0.2);
  border-radius: 12px;
  color: #f5f0e8;
  font-size: 1rem;
  transition: all 0.3s ease;
}
.search-input:focus {
  outline: none;
  border-color: rgba(230,168,0,0.5);
  background: rgba(255,248,220,0.08);
  box-shadow: 0 0 0 3px rgba(230,168,0,0.1);
}
.search-input::placeholder {
  color: rgba(245,240,232,0.4);
}
.filter-wrap {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}
.filter-field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.filter-field label {
  font-size: 0.85rem;
  color: rgba(245,240,232,0.7);
  font-weight: 500;
}
.gold-val {
  color: #f9c932;
}
.filter-field select,
.filter-field input[type="range"] {
  padding: 12px 16px;
  background: rgba(255,248,220,0.05);
  border: 1px solid rgba(230,168,0,0.2);
  border-radius: 10px;
  color: #f5f0e8;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s ease;
}
.filter-field select {
  min-width: 180px;
}
.filter-field select:focus {
  outline: none;
  border-color: rgba(230,168,0,0.5);
  box-shadow: 0 0 0 3px rgba(230,168,0,0.1);
}
.filter-field input[type="range"] {
  -webkit-appearance: none;
  appearance: none;
  width: 180px;
  height: 8px;
  padding: 0;
  background: rgba(230,168,0,0.2);
  border-radius: 4px;
  cursor: pointer;
}
.filter-field input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 20px;
  height: 20px;
  background: linear-gradient(135deg, #e6a800, #f9c932);
  border-radius: 50%;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(230,168,0,0.4);
  transition: transform 0.2s ease;
}
.filter-field input[type="range"]::-webkit-slider-thumb:hover {
  transform: scale(1.1);
}

/* ════════════════════════════════════════════════════════════════
     SECTION HEADERS
     ════════════════════════════════════════════════════════════════ */
.section-header {
  text-align: center;
  margin: 48px auto 32px;
  max-width: 1200px;
  padding: 0 24px;
}
.section-header h2 {
  font-size: clamp(1.8rem, 4vw, 2.5rem);
  font-weight: 800;
  color: #fff;
  margin: 0 0 8px;
  letter-spacing: -0.02em;
}
.section-sub {
  display: block;
  color: rgba(245,240,232,0.6);
  font-size: 1rem;
}
.grid-section-header {
  text-align: left;
}

/* ════════════════════════════════════════════════════════════════
     TOP 5 SCROLL TRACK - Gleiche Karten wie Grid
     ════════════════════════════════════════════════════════════════ */
.top5-scroll-track {
  display: flex;
  gap: 24px;
  overflow-x: auto;
  padding: 0 24px 24px;
  max-width: 1200px;
  margin: 0 auto;
  scrollbar-width: thin;
  scrollbar-color: rgba(230,168,0,0.4) transparent;
}
.top5-scroll-track::-webkit-scrollbar {
  height: 8px;
}
.top5-scroll-track::-webkit-scrollbar-track {
  background: transparent;
}
.top5-scroll-track::-webkit-scrollbar-thumb {
  background: rgba(230,168,0,0.4);
  border-radius: 4px;
}
.top5-scroll-track::-webkit-scrollbar-thumb:hover {
  background: rgba(230,168,0,0.6);
}

/* TOP 5 Karte - nutzt jetzt restaurant-card Basis */
.top5-card {
  flex: 0 0 300px;
  /* Rank Badge für Top 5 */
}
.top5-card .rank {
  position: absolute;
  top: 16px;
  left: 16px;
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #e6a800, #b87c00);
  color: #1a0f00;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 900;
  font-size: 1.4rem;
  z-index: 25;
  box-shadow: 0 4px 16px rgba(230, 168, 0, 0.5);
  border: 2px solid rgba(255, 255, 255, 0.2);
}

/* ════════════════════════════════════════════════════════════════
     RESTAURANT CARD - EINHEITLICHES DESIGN für TOP 5 & GRID
     ════════════════════════════════════════════════════════════════ */
.restaurant-card {
  display: flex !important;
  flex-direction: column !important;
  align-items: stretch !important;
  justify-content: flex-start !important;
  text-decoration: none;
  color: inherit;
  overflow: hidden;
  border-radius: 16px;
  background: #1c1710;
  border: 1px solid rgba(230, 168, 0, 0.15);
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
  position: relative;
}

.restaurant-card:hover {
  border-color: rgba(230, 168, 0, 0.4);
  box-shadow: 0 8px 24px rgba(230, 168, 0, 0.15);
}

.card-bild {
  position: relative;
  width: 100%;
  height: 250px;
  flex: 0 0 250px !important;
  display: flex !important;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: radial-gradient(circle at 50% 40%, rgba(230, 168, 0, 0.12), rgba(0, 0, 0, 0) 70%);
  z-index: 1;
}

.bewertung-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 10px;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  border-radius: 20px;
  border: 1px solid rgba(230, 168, 0, 0.3);
  z-index: 20;
  pointer-events: none;
}

.bewertung-badge .stern {
  font-size: 0.9rem;
  line-height: 1;
}

.bewertung-badge .bewertung-wert {
  font-size: 0.85rem;
  font-weight: 700;
  color: #f9c932;
  line-height: 1;
}

.card-bild::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 30%;
  backdrop-filter: blur(6px);
  -webkit-backdrop-filter: blur(6px);
  background: linear-gradient(to bottom, rgba(28, 23, 16, 0) 0%, rgba(28, 23, 16, 0.5) 55%, rgba(28, 23, 16, 0.92) 100%);
  z-index: 1;
  pointer-events: none;
}

.card-bild::before {
  content: '';
  position: absolute;
  left: 16px;
  right: 16px;
  bottom: 30%;
  height: 1px;
  background: linear-gradient(to right, rgba(255, 255, 255, 0) 0%, rgba(255, 255, 255, 0.4) 50%, rgba(255, 255, 255, 0) 100%);
  box-shadow: 0 1px 4px rgba(255, 255, 255, 0.15);
  z-index: 2;
  pointer-events: none;
}

.emoji-bild {
  font-size: 6.2rem;
  line-height: 1;
  filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.4));
  transition: transform 0.3s ease;
  position: relative;
  z-index: 10;
}

.restaurant-card:hover .emoji-bild {
  transform: scale(1.08);
}

.card-info {
  width: 100%;
  flex: 1 1 auto !important;
  min-height: 64px;
  box-sizing: border-box;
  padding: 6px 16px 18px;
  margin-top: -6px;
  display: flex !important;
  flex-direction: column !important;
  align-items: stretch !important;
  justify-content: flex-start !important;
  gap: 3px;
  position: relative;
  z-index: 3;
}

.card-info h3 {
  width: 100%;
  margin: 0 !important;
  padding: 0 !important;
  font-size: 1.1rem;
  font-weight: 600;
  color: #fff !important;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;
  min-height: 1.3em;
  visibility: visible !important;
  opacity: 1 !important;
  display: -webkit-box !important;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.25);
}

.card-meta {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.tag {
  display: inline-flex;
  align-items: center;
  box-sizing: border-box;
  background: rgba(230, 168, 0, 0.2);
  color: #f9c932 !important;
  padding: 3px 10px;
  line-height: 1.3;
  border-radius: 12px;
  font-size: 0.72rem;
  font-weight: 600;
  border: 1px solid rgba(230, 168, 0, 0.25);
  white-space: nowrap;
  flex-shrink: 0;
}

/* ════════════════════════════════════════════════════════════════
     RESTAURANT GRID
     ════════════════════════════════════════════════════════════════ */
.restaurant-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
  max-width: 1200px;
  margin: 0 auto 60px;
  padding: 0 24px;
  align-items: stretch;
}

/* ════════════════════════════════════════════════════════════════
     AUTH / LOGIN STYLES
     ════════════════════════════════════════════════════════════════ */
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
.pw-block { margin-top: 8px; }
.pw-track { background: rgba(255,248,220,0.10); border-radius: 6px; height: 6px; overflow: hidden; }
.pw-fill { height: 100%; border-radius: 6px; transition: width 0.3s, background 0.3s; }
.pw-label { font-size: 0.78rem; font-weight: 700; display: block; margin-top: 5px; }
.pw-rules { list-style: none; padding: 0; margin: 8px 0 0; display: flex; flex-wrap: wrap; gap: 5px 14px; }
.pw-rules li { font-size: 0.76rem; color: rgba(245,240,232,0.78); }
.pw-rules li::before { content: "✗ "; color: #ff453a; font-weight: 700; }
.pw-rules li.ok { color: #30d158; }
.pw-rules li.ok::before { content: "✓ "; color: #30d158; }
:global(html[data-theme='light']) .pw-rules li { color: rgba(40, 28, 0, 0.75); }
.err { color: #ff453a; font-size: 0.82rem; font-weight: 600; margin: 2px 0 0; }
.err-box { background: rgba(255, 69, 58, 0.10); border: 1px solid rgba(255, 69, 58, 0.35); border-radius: 10px; padding: 10px 12px; margin: 4px 0; }
.err-list { margin: 6px 0 0; padding-left: 18px; }
.err-list li { color: #ff6961; font-size: 0.8rem; line-height: 1.5; }
.auth-hint { text-align: center; font-size: 0.84rem; color: var(--text-muted); margin: 4px 0 0; }
.auth-hint a { color: var(--gold-text); font-weight: 600; }
.agb-check { display: flex; align-items: flex-start; gap: 8px; font-size: 0.84rem; color: var(--text-muted); text-align: left; margin: 4px 0; }
.agb-check input { margin-top: 3px; flex-shrink: 0; width: auto; }
.agb-check a { color: var(--gold-text); font-weight: 600; text-decoration: underline; }

/* ════════════════════════════════════════════════════════════════
     RESPONSIVE
     ════════════════════════════════════════════════════════════════ */
@media (max-width: 900px) {
  .hero-container {
    grid-template-columns: 1fr;
    gap: 40px;
    text-align: center;
  }
  .hero-text-col {
    align-items: center;
  }
  .hero-description {
    max-width: 100%;
  }
  .hero-cta-group {
    justify-content: center;
  }
  .trust-indicators {
    flex-wrap: wrap;
    justify-content: center;
  }
  .trust-divider {
    display: none;
  }
  .controls {
    flex-direction: column;
    align-items: stretch;
  }
  .search-wrap {
    max-width: 100%;
  }
  .filter-wrap {
    flex-direction: column;
  }
  .filter-field select,
  .filter-field input[type="range"] {
    width: 100%;
  }
  .grid-section-header {
    text-align: center;
  }
}
@media (max-width: 600px) {
  .hero-banner {
    padding: 100px 16px 60px;
  }
  .hero-headline {
    font-size: 2.2rem;
  }
  .cart-preview-card {
    max-width: 100%;
  }
  .top5-card {
    flex: 0 0 280px;
  }
  .top5-card .rank {
    width: 42px;
    height: 42px;
    font-size: 1.2rem;
  }
  .restaurant-grid {
    grid-template-columns: 1fr;
  }
  .emoji-bild {
    font-size: 5rem;
  }
}
</style>