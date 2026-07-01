<script>
  import { eingeloggt, logout } from '$lib/stores/auth.js';
  import { treuepunkte } from '$lib/stores/treue.js';
  import { api, getToken } from '$lib/api.js';

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
    <h2>👤 Mein Lieferino Profil</h2>
    <p>Verwalte deine Daten und Lieferadressen</p>
  </div>

  {#if $eingeloggt && geladen}
    <div class="profile-quadrat">
      
      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">🔐 Logindaten</h3>
          {#if !editBereich.login}
            <button onclick={() => editBereich.login = true} class="edit-icon-btn" title="Gruppe bearbeiten">✏️ Bearbeiten</button>
          {/if}
        </div>
        
        <div class="info-block">
          <span class="label">Username</span>
          <div class="block-content">
            {#if !editBereich.login} <span class="value">@{user.username}</span> {:else} <input type="text" bind:value={inputs.username} class="inline-input" /> {/if}
          </div>
        </div>

        <div class="info-block">
          <span class="label">E-Mail Adresse</span>
          <div class="block-content">
            {#if !editBereich.login} <span class="value email-value">{user.email}</span> {:else} <input type="email" bind:value={inputs.email} class="inline-input" /> {/if}
          </div>
        </div>

        <div class="info-block">
          <span class="label">Passwort</span>
          <div class="block-content">
            <a href="/passwort-vergessen" class="value pw-link">Passwort ändern 🔑</a>
          </div>
        </div>

        {#if editBereich.login}
          <button onclick={() => bereichSpeichern('login', ['username', 'email'])} class="save-group-btn">💾 Logindaten speichern</button>
        {/if}
      </div>

      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">👤 Persönliche Daten</h3>
          {#if !editBereich.persoenlich}
            <button onclick={() => editBereich.persoenlich = true} class="edit-icon-btn" title="Gruppe bearbeiten">✏️ Bearbeiten</button>
          {/if}
        </div>

        <div class="info-block">
          <span class="label">Vorname</span>
          <div class="block-content">
            {#if !editBereich.persoenlich} <span class="value">{user.vorname}</span> {:else} <input type="text" bind:value={inputs.vorname} class="inline-input" /> {/if}
          </div>
        </div>

        <div class="info-block">
          <span class="label">Zweiter Vorname</span>
          <div class="block-content">
            {#if !editBereich.persoenlich} 
              <span class="value {user.zweitname ? '' : 'placeholder-text'}">{user.zweitname ? user.zweitname : "Kein Zweitname"}</span> 
            {:else} 
              <input type="text" bind:value={inputs.zweitname} placeholder="Optional" class="inline-input" /> 
            {/if}
          </div>
        </div>

        <div class="info-block">
          <span class="label">Nachname</span>
          <div class="block-content">
            {#if !editBereich.persoenlich} <span class="value">{user.nachname}</span> {:else} <input type="text" bind:value={inputs.nachname} class="inline-input" /> {/if}
          </div>
        </div>

        {#if editBereich.persoenlich}
          <button onclick={() => bereichSpeichern('persoenlich', ['vorname', 'zweitname', 'nachname'])} class="save-group-btn">💾 Namen speichern</button>
        {/if}
      </div>

      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">🏠 Lieferadresse</h3>
          {#if !editBereich.adresse}
            <button onclick={() => editBereich.adresse = true} class="edit-icon-btn" title="Gruppe bearbeiten">✏️ Bearbeiten</button>
          {/if}
        </div>

        <div class="split-row">
          <div class="info-block grow">
            <span class="label">Straße</span>
            <div class="block-content">
              {#if !editBereich.adresse} <span class="value">{user.strasse}</span> {:else} <input type="text" bind:value={inputs.strasse} class="inline-input" /> {/if}
            </div>
          </div>

          <div class="info-block narrow">
            <span class="label">Nr.</span>
            <div class="block-content">
              {#if !editBereich.adresse} <span class="value">{user.hausnummer}</span> {:else} <input type="text" bind:value={inputs.hausnummer} class="inline-input" /> {/if}
            </div>
          </div>
        </div>

        <div class="split-row">
          <div class="info-block narrow">
            <span class="label">PLZ</span>
            <div class="block-content">
              {#if !editBereich.adresse} <span class="value">{user.plz}</span> {:else} <input type="text" bind:value={inputs.plz} class="inline-input" /> {/if}
            </div>
          </div>

          <div class="info-block grow">
            <span class="label">Ort</span>
            <div class="block-content">
              {#if !editBereich.adresse} <span class="value">{user.ort}</span> {:else} <input type="text" bind:value={inputs.ort} class="inline-input" /> {/if}
            </div>
          </div>
        </div>

        {#if editBereich.adresse}
          <button onclick={() => bereichSpeichern('adresse', ['strasse', 'hausnummer', 'plz', 'ort'])} class="save-group-btn">💾 Adresse speichern</button>
        {/if}
      </div>

      <!-- ⭐ TREUEPUNKTE -->
      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">⭐ Treuepunkte</h3>
        </div>
        <p class="punkte-zahl">{$treuepunkte} Punkte</p>
        <p class="punkte-hint">Sammle 1 Punkt je 1€ Bestellwert. 100 Punkte = 5€ Rabatt im Checkout.</p>
      </div>

      <!-- 📍 WEITERE ADRESSEN -->
      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">📍 Meine Adressen</h3>
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
          <input type="text" placeholder="Bezeichnung (z.B. Arbeit)" bind:value={neueAdresse.label} />
          <div class="adr-row">
            <input type="text" placeholder="Straße" bind:value={neueAdresse.strasse} required />
            <input type="text" placeholder="Nr." bind:value={neueAdresse.hausnummer} required />
          </div>
          <div class="adr-row">
            <input type="text" placeholder="PLZ" bind:value={neueAdresse.plz} required />
            <input type="text" placeholder="Ort" bind:value={neueAdresse.ort} required />
          </div>
          <button type="submit" class="save-group-btn">➕ Adresse hinzufügen</button>
        </form>
      </div>

      <!-- 🔐 ZWEI-FAKTOR-AUTHENTIFIZIERUNG -->
      <div class="category-section">
        <div class="category-header">
          <h3 class="category-title">🔐 Zwei-Faktor-Authentifizierung (2FA)</h3>
        </div>

        {#if user.mfaAktiv}
          <p class="mfa-status an">✅ Aktiv – Authenticator-App</p>
        {:else}
          <p class="mfa-status aus">⚠️ Nicht aktiv</p>
        {/if}
        <p class="mfa-info">
          Die Zwei-Faktor-Authentifizierung ist <strong>Pflicht</strong> und wird bei jeder Anmeldung
          abgefragt. Sie wurde bei der Registrierung mit deiner Authenticator-App eingerichtet.
        </p>
      </div>

      <button onclick={ausloggen} class="logout-btn">🔴 Ausloggen</button>

      <!-- 🗑️ Konto löschen (Gefahrenzone) -->
      {#if !loeschBestaetigung}
        <button onclick={() => (loeschBestaetigung = true)} class="delete-btn">🗑️ Konto löschen</button>
      {:else}
        <div class="delete-confirm">
          <p>Wirklich löschen? Alle deine Daten (Konto, Bestellungen, Favoriten) gehen verloren.</p>
          <div class="delete-row">
            <button onclick={() => (loeschBestaetigung = false)} class="mfa-btn grau">Abbrechen</button>
            <button onclick={kontoLoeschen} class="delete-btn endgueltig">Endgültig löschen</button>
          </div>
        </div>
      {/if}
    </div>
  {:else}
    <div class="no-user">
      <p>Du bist aktuell nicht eingeloggt.</p>
      <a href="/login" class="login-redirect-btn">Zum Login 🔑</a>
    </div>
  {/if}
</div>

<style>
  .account-container { max-width: 540px; margin: 40px auto; font-family: sans-serif; padding: 0 20px; }
  
  /* 🟪 Lila Hero-Box */
  .hero-box { background: #673ab7; color: white; padding: 30px 20px; border-radius: 24px; text-align: center; margin-bottom: 30px; box-shadow: 0 8px 25px rgba(103, 58, 183, 0.2); }
  .hero-box h2 { font-size: 2rem; color: white !important; margin: 0 0 8px 0; font-weight: 800; }
  .hero-box p { color: #e1d5f5 !important; font-size: 1rem; margin: 0; }

  .profile-quadrat { background: white; border-radius: 24px; padding: 35px; border: 1px solid #eee; box-shadow: 0 10px 35px rgba(0,0,0,0.05); display: flex; flex-direction: column; gap: 24px; }
  
  /* Kategorien */
  .category-section { display: flex; flex-direction: column; gap: 14px; background: #fafafa; padding: 20px; border-radius: 16px; border: 1px solid #f0f0f0; transition: all 0.2s ease; }
  .category-header { display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid #eee; padding-bottom: 8px; margin-bottom: 4px; }
  .category-title { font-size: 1.05rem; color: #333; font-weight: bold; margin: 0; }
  
  /* Neuer lila Gruppen-Stift */
  .edit-icon-btn { background: #f3e5f5; color: #673ab7; border: none; cursor: pointer; font-size: 0.85rem; font-weight: bold; padding: 6px 12px; border-radius: 20px; transition: background 0.2s; }
  .edit-icon-btn:hover { background: #e1bee7; }

  /* 🟪 Dynamischer lila Speicher-Button */
  .save-group-btn { margin-top: 10px; padding: 12px; background: #673ab7; color: white; border: none; border-radius: 10px; font-weight: bold; cursor: pointer; font-size: 0.95rem; box-shadow: 0 4px 12px rgba(103, 58, 183, 0.15); animation: fadeIn 0.2s ease; }
  .save-group-btn:hover { background: #542f95; }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(-5px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .info-block { display: flex; flex-direction: column; gap: 4px; width: 100%; }
  .split-row { display: flex; gap: 15px; width: 100%; }
  .grow { flex: 1; }
  .narrow { width: 30%; }
  
  .label { font-size: 0.75rem; color: #888; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px; }
  .value { font-size: 1.05rem; color: #222; font-weight: 700; padding: 4px 0; display: inline-block; }
  .email-value { color: #673ab7; }
  .placeholder-text { color: #ccc; font-weight: 400; font-style: italic; font-size: 0.95rem; }
  .pw-link { color: #673ab7; }

  .inline-input { width: 100%; padding: 10px; border: 2px solid #e1bee7; border-radius: 8px; font-size: 0.95rem; font-family: sans-serif; font-weight: 600; box-sizing: border-box; outline: none; background: #fff; }
  .inline-input:focus { border-color: #673ab7; }
  
  .logout-btn { margin-top: 5px; padding: 14px; background: #fff; color: #dc3545; border: 2px solid #dc3545; border-radius: 12px; font-size: 1rem; font-weight: bold; cursor: pointer; }
  .logout-btn:hover { background: #dc3545; color: white; }
  /* 🔐 MFA / 2FA */
  .mfa-status { font-weight: 700; margin: 0; }
  .mfa-status.an { color: #34c759; }
  .mfa-status.aus { color: #d97706; }
  .mfa-info { font-size: 0.9rem; color: #555; margin: 6px 0; }
  .mfa-btn { flex: 1; min-width: 140px; padding: 12px; background: #673ab7; color: white; border: none; border-radius: 10px; font-weight: bold; cursor: pointer; }
  .mfa-btn:hover { background: #542f95; }
  .mfa-btn.grau { background: #f1f1f1; color: #333; }

  /* ⭐ Treuepunkte */
  .punkte-zahl { font-size: 1.6rem; font-weight: 800; color: #673ab7; margin: 0; }
  .punkte-hint { font-size: 0.85rem; color: #888; margin: 4px 0 0; }

  /* 📍 Adressen */
  .adress-eintrag { display: flex; justify-content: space-between; align-items: center; gap: 10px; background: #fff; border: 1px solid #eee; border-radius: 10px; padding: 10px 12px; font-size: 0.9rem; }
  .adr-loeschen { background: none; border: none; cursor: pointer; font-size: 1rem; }
  .adress-form { display: flex; flex-direction: column; gap: 8px; margin-top: 8px; }
  .adress-form input { padding: 10px; border: 1px solid #ddd; border-radius: 8px; font-size: 0.9rem; }
  .adr-row { display: flex; gap: 8px; }
  .adr-row input { flex: 1; min-width: 0; }

  /* 🗑️ Konto löschen */
  .delete-btn { padding: 12px; background: none; color: #999; border: 1px solid #ddd; border-radius: 12px; font-weight: 600; cursor: pointer; }
  .delete-btn:hover { color: #dc3545; border-color: #dc3545; }
  .delete-btn.endgueltig { flex: 1; background: #dc3545; color: white; border: none; }
  .delete-confirm { background: #fff5f5; border: 1px solid #f5c6cb; border-radius: 12px; padding: 16px; }
  .delete-confirm p { color: #842029; font-size: 0.9rem; margin: 0 0 12px; }
  .delete-row { display: flex; gap: 10px; }

  .no-user { padding: 40px; background: #f9f9f9; border-radius: 12px; text-align: center; }
  .login-redirect-btn { display: inline-block; margin-top: 15px; padding: 10px 20px; background: #673ab7; color: white; text-decoration: none; border-radius: 8px; font-weight: bold; }
</style>