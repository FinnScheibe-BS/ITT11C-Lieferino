<script>
  import { onMount } from 'svelte';

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

  onMount(() => {
    const gespeicherterUser = localStorage.getItem("lieferino_user");
    if (gespeicherterUser) {
      user = JSON.parse(gespeicherterUser);
      // Inputs initialisieren
      Object.keys(user).forEach(key => {
        inputs[key] = user[key] || "";
      });
      inputs.passwort = ""; 
      geladen = true;
    }
  });

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
  }

  function ausloggen() {
    localStorage.removeItem("lieferino_user");
    window.location.href = "/";
  }
</script>

<div class="account-container">
  
  <div class="hero-box">
    <h2>👤 Mein Lieferino Profil</h2>
    <p>Verwalte deine Daten und Lieferadressen</p>
  </div>

  {#if geladen}
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
          <span class="label">Passwort ändern</span>
          <div class="block-content">
            {#if !editBereich.login} <span class="value password-dots">••••••••</span> {:else} <input type="password" bind:value={inputs.passwort} placeholder="Neues Passwort eingeben..." class="inline-input" /> {/if}
          </div>
        </div>

        {#if editBereich.login}
          <button onclick={() => bereichSpeichern('login', ['username', 'email', 'passwort'])} class="save-group-btn">💾 Logindaten speichern</button>
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

      <button onclick={ausloggen} class="logout-btn">🔴 Ausloggen</button>
    </div>
  {:else}
    <div class="no-user">
      <p>Du bist aktuell nicht eingeloggt.</p>
      <a href="/" class="login-redirect-btn">Zur Startseite / Login</a>
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
  .password-dots { letter-spacing: 2px; color: #673ab7; }
  
  .inline-input { width: 100%; padding: 10px; border: 2px solid #e1bee7; border-radius: 8px; font-size: 0.95rem; font-family: sans-serif; font-weight: 600; box-sizing: border-box; outline: none; background: #fff; }
  .inline-input:focus { border-color: #673ab7; }
  
  .logout-btn { margin-top: 5px; padding: 14px; background: #fff; color: #dc3545; border: 2px solid #dc3545; border-radius: 12px; font-size: 1rem; font-weight: bold; cursor: pointer; }
  .logout-btn:hover { background: #dc3545; color: white; }
  .no-user { padding: 40px; background: #f9f9f9; border-radius: 12px; text-align: center; }
  .login-redirect-btn { display: inline-block; margin-top: 15px; padding: 10px 20px; background: #673ab7; color: white; text-decoration: none; border-radius: 8px; font-weight: bold; }
</style>