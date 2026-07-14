<script>
  import { eingeloggt, logout } from '$lib/stores/auth.js';
  import { treuepunkte } from '$lib/stores/treue.js';
  import { api, getToken } from '$lib/api/api.js';
  import { t } from '$lib/utils/i18n.js';

  let user = $state({ username: "", vorname: "", zweitname: "", nachname: "", strasse: "", hausnummer: "", plz: "", ort: "", email: "", passwort: "" });
  let geladen = $state(false);

  // Edit-Modus gilt jetzt pro BEREICH (Kategorie) statt pro Feld
  let editBereich = $state({
    login: false,
    persoenlich: false,
    adresse: false
  });

  // Temporäre Inputs für die Bearbeitung
  let inputs = $state({ username: "", vorname: "", zweitname: "", nachname: "", strasse: "", hausnummer: "", plz: "", ort: "", email: "", passwort: "" });

  // Hinweis: onMount feuert in diesem Setup nicht zuverlässig -> $effect (läuft sicher).
  let initGeladen = false;
  $effect(() => {
    if (initGeladen) return;
    initGeladen = true;
    const gespeicherterUser = localStorage.getItem("lieferino_user");
    if (gespeicherterUser) {
      user = JSON.parse(gespeicherterUser);
      // Inputs initialisieren
      Object.keys(user).forEach(key => {
        inputs[key] = user[key] || "";
      });
      inputs.passwort = "";

      // 📍 Adressen initialisieren: falls noch keine Liste existiert, die
      // Registrierungs-Adresse als erste Adresse ("Zuhause") übernehmen.
      if (!user.adressen || user.adressen.length === 0) {
        user.adressen = [
          { label: 'Zuhause', strasse: user.strasse, hausnummer: user.hausnummer, plz: user.plz, ort: user.ort }
        ];
        localStorage.setItem('lieferino_user', JSON.stringify(user));
      }

      geladen = true;
    }

    // 🗄️ Falls eingeloggt: aktuelle Daten aus dem Backend (DB) nachladen.
    (async () => {
      if (!getToken()) return;
      const res = await api('/api/me');
      if (res.ok && res.daten) {
        // Backend ist die Quelle der Wahrheit für Profil + Adressen.
        user = { ...user, ...res.daten, passwort: user.passwort };
        if (res.daten.adressen?.length) user.adressen = res.daten.adressen;
        Object.keys(inputs).forEach((key) => { if (user[key] != null) inputs[key] = user[key]; });
        inputs.passwort = "";
        localStorage.setItem('lieferino_user', JSON.stringify(user));
        geladen = true;
      }
    })();
  });

  // 🗄️ Speichert Profil + Adressen ins Backend (DB), wenn eingeloggt.
  async function backendSync() {
    if (!getToken()) return;
    await api('/api/me', {
      method: 'PUT',
      body: {
        username: user.username, vorname: user.vorname,
        nachname: user.nachname, geburtsdatum: user.geburtsdatum || '',
        adressen: user.adressen || []
      }
    });
  }

  // 📍 Eingabefelder für eine neue Adresse
  let neueAdresse = $state({ label: '', strasse: '', hausnummer: '', plz: '', ort: '' });

  function adresseHinzufuegen(e) {
    e.preventDefault();
    // Minimal-Validierung: PLZ = 5 Ziffern, Felder nicht leer.
    if (!/^\d{5}$/.test(neueAdresse.plz.trim())) return;
    if (!neueAdresse.strasse.trim() || !neueAdresse.ort.trim()) return;
    user.adressen = [...user.adressen, { ...neueAdresse, label: neueAdresse.label || 'Adresse' }];
    localStorage.setItem('lieferino_user', JSON.stringify(user));
    neueAdresse = { label: '', strasse: '', hausnummer: '', plz: '', ort: '' };
    backendSync();
  }

  function adresseLoeschen(index) {
    user.adressen = user.adressen.filter((_, i) => i !== index);
    localStorage.setItem('lieferino_user', JSON.stringify(user));
    backendSync();
  }

  // Speichert alle Felder einer bestimmten Gruppe
  function bereichSpeichern(bereich, felder) {
    felder.forEach(feld => {
      if (feld === 'passwort') {
        if (inputs.passwort.trim() !== "") {
          user.passwort = inputs.passwort;
        }
      } else {
        user[feld] = inputs[feld];
      }
    });

    localStorage.setItem("lieferino_user", JSON.stringify(user));
    editBereich[bereich] = false; // Bearbeitungsmodus schließen
    if (felder.includes('passwort')) inputs.passwort = ""; // Passwort-Feld leeren
    backendSync(); // 🗄️ Änderungen auch ins Backend (DB)
  }

  // 🔓 Ausloggen beendet nur die Session – das Konto bleibt erhalten!
  function ausloggen() {
    logout();
    window.location.href = "/";
  }

  // 🗑️ Konto löschen (zweistufig, damit es nicht aus Versehen passiert).
  let loeschBestaetigung = $state(false);
  function kontoLoeschen() {
    localStorage.removeItem('lieferino_user');
    localStorage.removeItem('lieferino_session');
    localStorage.removeItem('lieferino_bestellungen');
    localStorage.removeItem('lieferino_favoriten');
    logout();
    window.location.href = '/';
  }

</script>

<div class="account-container">
  
  <div class="hero-box">
    <h2>{$t('acc.title')}</h2>
    <p>{$t('acc.subtitle')}</p>
  </div>

  {#if $eingeloggt && geladen}
    <div class="profile-quadrat">
      
      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">{$t('acc.login_data')}</h3>
          {#if !editBereich.login}
            <button onclick={() => editBereich.login = true} class="edit-icon-btn" title="Gruppe bearbeiten">{$t('acc.edit')}</button>
          {/if}
        </div>
        
        <div class="info-block">
          <span class="label">{$t('acc.username')}</span>
          <div class="block-content">
            {#if !editBereich.login} <span class="value">@{user.username}</span> {:else} <input type="text" bind:value={inputs.username} class="inline-input" /> {/if}
          </div>
        </div>

        <div class="info-block">
          <span class="label">{$t('auth.email')}</span>
          <div class="block-content">
            {#if !editBereich.login} <span class="value email-value">{user.email}</span> {:else} <input type="email" bind:value={inputs.email} class="inline-input" /> {/if}
          </div>
        </div>

        <div class="info-block">
          <span class="label">{$t('auth.password')}</span>
          <div class="block-content">
            <a href="/passwort-vergessen" class="value pw-link">{$t('acc.change_pw')}</a>
          </div>
        </div>

        {#if editBereich.login}
          <button onclick={() => bereichSpeichern('login', ['username', 'email'])} class="save-group-btn">{$t('acc.save_login')}</button>
        {/if}
      </div>

      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">{$t('acc.personal')}</h3>
          {#if !editBereich.persoenlich}
            <button onclick={() => editBereich.persoenlich = true} class="edit-icon-btn" title="Gruppe bearbeiten">{$t('acc.edit')}</button>
          {/if}
        </div>

        <div class="info-block">
          <span class="label">{$t('acc.firstname')}</span>
          <div class="block-content">
            {#if !editBereich.persoenlich} <span class="value">{user.vorname}</span> {:else} <input type="text" bind:value={inputs.vorname} class="inline-input" /> {/if}
          </div>
        </div>

        <div class="info-block">
          <span class="label">{$t('acc.middlename')}</span>
          <div class="block-content">
            {#if !editBereich.persoenlich}
              <span class="value {user.zweitname ? '' : 'placeholder-text'}">{user.zweitname ? user.zweitname : $t('acc.no_middle')}</span>
            {:else} 
              <input type="text" bind:value={inputs.zweitname} placeholder="Optional" class="inline-input" /> 
            {/if}
          </div>
        </div>

        <div class="info-block">
          <span class="label">{$t('acc.lastname')}</span>
          <div class="block-content">
            {#if !editBereich.persoenlich} <span class="value">{user.nachname}</span> {:else} <input type="text" bind:value={inputs.nachname} class="inline-input" /> {/if}
          </div>
        </div>

        {#if editBereich.persoenlich}
          <button onclick={() => bereichSpeichern('persoenlich', ['vorname', 'zweitname', 'nachname'])} class="save-group-btn">{$t('acc.save_names')}</button>
        {/if}
      </div>

      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">{$t('acc.address')}</h3>
          {#if !editBereich.adresse}
            <button onclick={() => editBereich.adresse = true} class="edit-icon-btn" title="Gruppe bearbeiten">{$t('acc.edit')}</button>
          {/if}
        </div>

        <div class="split-row">
          <div class="info-block grow">
            <span class="label">{$t('acc.street')}</span>
            <div class="block-content">
              {#if !editBereich.adresse} <span class="value">{user.strasse}</span> {:else} <input type="text" bind:value={inputs.strasse} class="inline-input" /> {/if}
            </div>
          </div>

          <div class="info-block narrow">
            <span class="label">{$t('acc.number')}</span>
            <div class="block-content">
              {#if !editBereich.adresse} <span class="value">{user.hausnummer}</span> {:else} <input type="text" bind:value={inputs.hausnummer} class="inline-input" /> {/if}
            </div>
          </div>
        </div>

        <div class="split-row">
          <div class="info-block narrow">
            <span class="label">{$t('acc.zip')}</span>
            <div class="block-content">
              {#if !editBereich.adresse} <span class="value">{user.plz}</span> {:else} <input type="text" bind:value={inputs.plz} class="inline-input" /> {/if}
            </div>
          </div>

          <div class="info-block grow">
            <span class="label">{$t('acc.city')}</span>
            <div class="block-content">
              {#if !editBereich.adresse} <span class="value">{user.ort}</span> {:else} <input type="text" bind:value={inputs.ort} class="inline-input" /> {/if}
            </div>
          </div>
        </div>

        {#if editBereich.adresse}
          <button onclick={() => bereichSpeichern('adresse', ['strasse', 'hausnummer', 'plz', 'ort'])} class="save-group-btn">{$t('acc.save_address')}</button>
        {/if}
      </div>

      <!-- ⭐ TREUEPUNKTE -->
      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">{$t('acc.points_title')}</h3>
        </div>
        <p class="punkte-zahl">{$t('acc.points_n').replace('{n}', $treuepunkte)}</p>
        <p class="punkte-hint">{$t('acc.points_hint')}</p>
      </div>

      <!-- 📍 WEITERE ADRESSEN -->
      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">{$t('acc.my_addresses')}</h3>
        </div>

        {#each user.adressen || [] as adr, i}
          <div class="adress-eintrag">
            <span><strong>{adr.label}</strong>: {adr.strasse} {adr.hausnummer}, {adr.plz} {adr.ort}</span>
            {#if (user.adressen || []).length > 1}
              <button class="adr-loeschen" onclick={() => adresseLoeschen(i)} aria-label="Adresse löschen">🗑️</button>
            {/if}
          </div>
        {/each}

        <!-- Neue Adresse hinzufügen -->
        <form class="adress-form" onsubmit={adresseHinzufuegen}>
          <input type="text" placeholder={$t('acc.addr_label_ph')} bind:value={neueAdresse.label} />
          <div class="adr-row">
            <input type="text" placeholder={$t('acc.street')} bind:value={neueAdresse.strasse} required />
            <input type="text" placeholder={$t('acc.number')} bind:value={neueAdresse.hausnummer} required />
          </div>
          <div class="adr-row">
            <input type="text" placeholder={$t('acc.zip')} bind:value={neueAdresse.plz} required />
            <input type="text" placeholder={$t('acc.city')} bind:value={neueAdresse.ort} required />
          </div>
          <button type="submit" class="save-group-btn">{$t('acc.add_address')}</button>
        </form>
      </div>

      <!-- 🔐 ZWEI-FAKTOR-AUTHENTIFIZIERUNG -->
      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">{$t('acc.mfa_title')}</h3>
        </div>

        {#if user.mfaAktiv}
          <p class="mfa-status an">{$t('acc.mfa_active')}</p>
        {:else}
          <p class="mfa-status aus">{$t('acc.mfa_inactive')}</p>
        {/if}
        <p class="mfa-info">{$t('acc.mfa_info')}</p>
      </div>

      <button onclick={ausloggen} class="logout-btn">{$t('acc.logout')}</button>

      <!-- 🗑️ Konto löschen (Gefahrenzone) -->
      {#if !loeschBestaetigung}
        <button onclick={() => (loeschBestaetigung = true)} class="delete-btn">{$t('acc.delete')}</button>
      {:else}
        <div class="delete-confirm">
          <p>{$t('acc.delete_confirm')}</p>
          <div class="delete-row">
            <button onclick={() => (loeschBestaetigung = false)} class="mfa-btn grau">{$t('acc.cancel')}</button>
            <button onclick={kontoLoeschen} class="delete-btn endgueltig">{$t('acc.delete_final')}</button>
          </div>
        </div>
      {/if}
    </div>
  {:else}
    <div class="no-user">
      <p>{$t('acc.not_logged')}</p>
      <a href="/login" class="login-redirect-btn">{$t('acc.to_login')}</a>
    </div>
  {/if}
</div>

<style>
  .account-container {
    max-width: 760px;
    margin: 0 auto;
    padding: 28px 20px 56px;
  }

  .hero-box {
    position: relative;
    overflow: hidden;
    padding: 34px 24px;
    border-radius: 24px;
    text-align: center;
    margin-bottom: 26px;
    background:
      radial-gradient(circle at 20% 20%, rgba(249, 201, 50, 0.28), transparent 38%),
      linear-gradient(135deg, rgba(230, 168, 0, 0.26), rgba(20, 12, 0, 0.82));
    border: 1px solid rgba(230, 168, 0, 0.28);
    box-shadow: 0 12px 42px rgba(0, 0, 0, 0.38), 0 0 28px rgba(230, 168, 0, 0.10);
    backdrop-filter: blur(18px) saturate(1.4);
    -webkit-backdrop-filter: blur(18px) saturate(1.4);
  }

  .hero-box::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(120deg, transparent 0%, rgba(255,255,255,0.08) 45%, transparent 70%);
    pointer-events: none;
  }

  .hero-box h2 {
    position: relative;
    margin: 0 0 8px;
    color: #fff !important;
    font-size: clamp(1.7rem, 4vw, 2.35rem);
    font-weight: 800;
    letter-spacing: -0.035em;
    z-index: 1;
  }

  .hero-box p {
    position: relative;
    margin: 0;
    color: rgba(255, 232, 160, 0.86) !important;
    font-size: 1rem;
    z-index: 1;
  }

  .profile-quadrat,
  .no-user {
    background: rgba(255, 248, 220, 0.06);
    border: 1px solid rgba(230, 168, 0, 0.18);
    border-radius: 24px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.34);
    backdrop-filter: blur(18px) saturate(1.4);
    -webkit-backdrop-filter: blur(18px) saturate(1.4);
  }

  .profile-quadrat {
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 18px;
  }

  .category-section {
    display: flex;
    flex-direction: column;
    gap: 14px;
    padding: 20px;
    border-radius: 18px;
    background: linear-gradient(180deg, rgba(255, 248, 220, 0.075), rgba(255, 248, 220, 0.035));
    border: 1px solid rgba(230, 168, 0, 0.18);
    box-shadow: inset 0 1px 0 rgba(255,255,255,0.04), 0 4px 20px rgba(0,0,0,0.18);
  }

  .category-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
    padding-bottom: 10px;
    border-bottom: 1px solid rgba(230, 168, 0, 0.16);
  }

  .category-title {
    margin: 0;
    color: #f5d87c !important;
    font-size: 1.05rem;
    font-weight: 750;
    letter-spacing: -0.02em;
  }

  .info-block {
    display: flex;
    flex-direction: column;
    gap: 5px;
    width: 100%;
  }

  .split-row {
    display: flex;
    gap: 14px;
    width: 100%;
  }

  .grow { flex: 1; }
  .narrow { width: 32%; min-width: 90px; }

  .label {
    font-size: 0.72rem;
    color: rgba(245, 216, 124, 0.72);
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }

  .block-content {
    min-height: 36px;
    display: flex;
    align-items: center;
  }

  .value {
    display: inline-block;
    padding: 4px 0;
    color: #f5f0e8;
    font-size: 1.02rem;
    font-weight: 700;
    word-break: break-word;
  }

  .email-value,
  .password-dots,
  .punkte-zahl {
    color: #f9c932 !important;
  }

  .placeholder-text {
    color: rgba(245, 240, 232, 0.42) !important;
    font-weight: 400;
    font-style: italic;
    font-size: 0.95rem;
  }

  .password-dots { letter-spacing: 2px; }

  .inline-input {
    width: 100%;
    font-weight: 650;
    box-sizing: border-box;
  }

  .edit-icon-btn,
  .save-group-btn,
  .mfa-btn,
  .login-redirect-btn {
    border-radius: 12px !important;
  }

  .edit-icon-btn {
    padding: 8px 12px !important;
    font-size: 0.82rem !important;
    white-space: nowrap;
    background: rgba(230, 168, 0, 0.12) !important;
    color: #f9c932 !important;
    border: 1px solid rgba(230, 168, 0, 0.28) !important;
    box-shadow: none !important;
  }

  .edit-icon-btn:hover {
    background: rgba(230, 168, 0, 0.20) !important;
    transform: translateY(-1px) !important;
    opacity: 1 !important;
  }

  .save-group-btn {
    margin-top: 8px;
    width: 100%;
    animation: fadeIn 0.2s ease;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(-5px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .punkte-zahl {
    margin: 0;
    font-size: 1.75rem;
    font-weight: 850;
    letter-spacing: -0.03em;
  }

  .punkte-hint,
  .mfa-info,
  .mfa-mini,
  .keine-reviews {
    color: rgba(245, 240, 232, 0.68);
  }

  .punkte-hint {
    font-size: 0.88rem;
    margin: 0;
    line-height: 1.45;
  }

  .adress-eintrag {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
    background: rgba(255, 248, 220, 0.06);
    border: 1px solid rgba(230, 168, 0, 0.16);
    border-radius: 14px;
    padding: 12px 14px;
    color: #f5f0e8;
    font-size: 0.92rem;
  }

  .adr-loeschen {
    width: 34px;
    height: 34px;
    padding: 0 !important;
    border-radius: 10px !important;
    background: rgba(220, 53, 69, 0.10) !important;
    color: #ff7b86 !important;
    border: 1px solid rgba(220, 53, 69, 0.30) !important;
    box-shadow: none !important;
    flex-shrink: 0;
  }

  .adress-form {
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin-top: 4px;
  }

  .adr-row {
    display: flex;
    gap: 10px;
  }

  .adr-row input { flex: 1; min-width: 0; }

  .mfa-status {
    font-weight: 750;
    margin: 0;
    line-height: 1.4;
  }

  .mfa-status.an { color: #63d88a; }
  .mfa-status.aus { color: #f5d87c; }

  .mfa-info {
    font-size: 0.92rem;
    margin: 4px 0;
    line-height: 1.45;
  }

  .mfa-mini {
    font-size: 0.82rem;
    margin: 8px 0 4px;
    line-height: 1.45;
  }

  .mfa-fehler {
    color: #ff7b86;
    font-weight: 700;
    font-size: 0.86rem;
    margin: 4px 0 0;
  }

  .mfa-btn-row,
  .delete-row {
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
    margin-top: 8px;
  }

  .mfa-btn {
    flex: 1;
    min-width: 150px;
  }

  .mfa-btn.grau {
    background: rgba(255, 248, 220, 0.10) !important;
    color: #f5f0e8 !important;
    border: 1px solid rgba(230, 168, 0, 0.18) !important;
    box-shadow: none !important;
  }

  .mfa-btn.aus,
  .logout-btn,
  .delete-btn {
    background: rgba(220, 53, 69, 0.08) !important;
    color: #ff7b86 !important;
    border: 1px solid rgba(220, 53, 69, 0.34) !important;
    box-shadow: none !important;
  }

  .mfa-btn.aus:hover,
  .logout-btn:hover,
  .delete-btn:hover {
    background: rgba(220, 53, 69, 0.18) !important;
    border-color: rgba(220, 53, 69, 0.65) !important;
    color: #fff !important;
    opacity: 1 !important;
  }

  .totp-secret,
  .backup-box,
  .delete-confirm,
  .no-user {
    border-radius: 16px;
    background: rgba(255, 248, 220, 0.06);
    border: 1px solid rgba(230, 168, 0, 0.18);
    padding: 14px;
  }

  .totp-secret {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin: 4px 0;
    border-style: dashed;
  }

  code,
  .totp-secret code,
  .totp-uri,
  .backup-codes code {
    font-family: 'Geist Mono', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  }

  .totp-secret code {
    font-size: 1.15rem;
    letter-spacing: 2px;
    color: #f9c932;
    font-weight: 800;
    word-break: break-all;
  }

  .totp-uri {
    display: block;
    word-break: break-all;
    background: rgba(0, 0, 0, 0.22);
    border: 1px solid rgba(230, 168, 0, 0.12);
    padding: 10px;
    border-radius: 10px;
    font-size: 0.75rem;
    color: rgba(245, 240, 232, 0.72);
  }

  .backup-box {
    background: rgba(230, 168, 0, 0.08);
    border-style: dashed;
  }

  .backup-titel {
    font-weight: 800;
    color: #f9c932;
    margin: 0 0 4px;
  }

  .backup-codes {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 8px;
    margin-top: 10px;
  }

  .backup-codes code {
    background: rgba(0, 0, 0, 0.22);
    border: 1px solid rgba(230, 168, 0, 0.18);
    border-radius: 9px;
    padding: 8px;
    text-align: center;
    font-weight: 800;
    letter-spacing: 1px;
    color: #f9c932;
  }

  .logout-btn,
  .delete-btn {
    width: 100%;
    padding: 13px !important;
    border-radius: 14px !important;
    font-size: 0.96rem !important;
    font-weight: 800 !important;
  }

  .delete-confirm {
    background: rgba(220, 53, 69, 0.08);
    border-color: rgba(220, 53, 69, 0.28);
  }

  .delete-confirm p {
    color: #ffb3ba;
    font-size: 0.9rem;
    margin: 0 0 12px;
    line-height: 1.45;
  }

  .delete-btn.endgueltig {
    flex: 1;
    background: linear-gradient(135deg, #dc3545, #9d1f2d) !important;
    color: #fff !important;
    border: none !important;
  }

  .no-user {
    padding: 36px 24px;
    text-align: center;
  }

  .no-user p {
    margin: 0;
    color: rgba(245, 240, 232, 0.78);
  }

  .login-redirect-btn {
    display: inline-block;
    margin-top: 16px;
    padding: 11px 20px;
    font-weight: 800;
  }

  :global(html[data-theme='light']) .hero-box {
    background:
      radial-gradient(circle at 20% 20%, rgba(230, 168, 0, 0.22), transparent 40%),
      linear-gradient(135deg, rgba(255, 252, 235, 0.92), rgba(253, 232, 160, 0.52));
    box-shadow: 0 12px 36px rgba(184, 124, 0, 0.15);
  }

  :global(html[data-theme='light']) .hero-box h2 { color: #1a0f00 !important; }
  :global(html[data-theme='light']) .hero-box p { color: #7a5000 !important; }

  :global(html[data-theme='light']) .profile-quadrat,
  :global(html[data-theme='light']) .category-section,
  :global(html[data-theme='light']) .adress-eintrag,
  :global(html[data-theme='light']) .totp-secret,
  :global(html[data-theme='light']) .backup-box,
  :global(html[data-theme='light']) .no-user {
    background: rgba(255, 252, 235, 0.78);
    border-color: rgba(230, 168, 0, 0.30);
    box-shadow: 0 6px 24px rgba(184, 124, 0, 0.10);
  }

  :global(html[data-theme='light']) .value,
  :global(html[data-theme='light']) .adress-eintrag,
  :global(html[data-theme='light']) .no-user p {
    color: #1a1200;
  }

  :global(html[data-theme='light']) .punkte-hint,
  :global(html[data-theme='light']) .mfa-info,
  :global(html[data-theme='light']) .mfa-mini {
    color: rgba(26, 18, 0, 0.65);
  }

  :global(html[data-theme='light']) .totp-uri,
  :global(html[data-theme='light']) .backup-codes code {
    background: rgba(255, 248, 220, 0.70);
    color: #7a5000;
  }

  @media (max-width: 640px) {
    .account-container { padding: 20px 14px 44px; }
    .profile-quadrat { padding: 16px; border-radius: 20px; }
    .category-section { padding: 16px; }
    .category-header { align-items: flex-start; flex-direction: column; }
    .edit-icon-btn { width: 100%; }
    .split-row, .adr-row { flex-direction: column; gap: 10px; }
    .narrow { width: 100%; }
    .backup-codes { grid-template-columns: 1fr; }
  }
</style>