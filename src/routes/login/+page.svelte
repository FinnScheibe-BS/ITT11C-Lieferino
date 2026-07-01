<script>
  import { goto } from '$app/navigation';
  import { login } from '$lib/stores/auth.js';
  import { api } from '$lib/api.js';
  import AuthVollenden from '$lib/AuthVollenden.svelte';

  // 'login' = E-Mail+Passwort. 'weiter' = Code/MFA über die AuthVollenden-Komponente.
  let schritt = $state('login');

  let email = $state('');
  let passwort = $state('');
  let fehler = $state('');
  let laedt = $state(false);
  let zeigePasswort = $state(false);
  let angemeldetBleiben = $state(false);

  // Übergabe an die Abschluss-Komponente:
  let authStart = $state('mfaVerify');
  let authEmail = $state('');
  let authSetupToken = $state('');
  let authMfaToken = $state('');

  // Schritt 1: Anmeldedaten ans Backend. Passwortprüfung passiert SERVERSEITIG.
  async function loginAbsenden(e) {
    e.preventDefault();
    fehler = '';
    laedt = true;
    const res = await api('/api/auth/login', { method: 'POST', body: { email: email.trim(), passwort } });
    laedt = false;

    if (res.offline) {
      fehler = 'Server nicht erreichbar. Bitte später erneut versuchen. 📡';
      return;
    }
    const d = res.daten || {};

    // Passwort ok -> nur noch MFA-Code (Schritt 2).
    if (res.ok && d.mfaRequired) {
      authStart = 'mfaVerify';
      authMfaToken = d.mfaToken;
      schritt = 'weiter';
      return;
    }
    // E-Mail noch nicht bestätigt -> Code eingeben (Backend hat einen geschickt).
    if (res.status === 403 && d.needsVerification) {
      authStart = 'verify';
      authEmail = email.trim();
      schritt = 'weiter';
      return;
    }
    // MFA noch nicht eingerichtet -> jetzt einrichten.
    if (res.status === 403 && d.needsMfaSetup) {
      authStart = 'mfaSetup';
      authSetupToken = d.setupToken;
      schritt = 'weiter';
      return;
    }
    // Sonst: Fehler (falsches Passwort, gesperrt, zu viele Versuche …) klar anzeigen.
    fehler = d.fehler || 'Anmeldung fehlgeschlagen.';
  }

  // Wird aufgerufen, sobald die Anmeldung komplett ist (volles Token liegt vor).
  function authFertig(daten) {
    if (daten?.user) localStorage.setItem('lieferino_user', JSON.stringify(daten.user));
    login(angemeldetBleiben);
    goto('/');
  }
</script>

<div class="login-wrapper">
  <div class="login-box">
    <div class="hero-box">
      <h2>🔑 Anmelden</h2>
    </div>

    {#if schritt === 'login'}
      <form onsubmit={loginAbsenden}>
        <div class="input-group">
          <label for="email">E-Mail Adresse</label>
          <input type="email" id="email" placeholder="name@beispiel.de" bind:value={email} required />
        </div>
        <div class="input-group">
          <label for="pw">Passwort</label>
          <div class="pw-feld">
            {#if zeigePasswort}
              <input type="text" id="pw" placeholder="••••••••" bind:value={passwort} required />
            {:else}
              <input type="password" id="pw" placeholder="••••••••" bind:value={passwort} required />
            {/if}
            <button type="button" class="auge" onclick={() => (zeigePasswort = !zeigePasswort)}>
              {zeigePasswort ? '🙈' : '👁️'}
            </button>
          </div>
        </div>

        <label class="bleiben">
          <input type="checkbox" bind:checked={angemeldetBleiben} /> Angemeldet bleiben
        </label>

        {#if fehler}<p class="error-msg">{fehler}</p>{/if}

        <button type="submit" class="login-btn" disabled={laedt}>{laedt ? 'Prüfe…' : 'Einloggen ➡️'}</button>
      </form>

      <p class="hinweis">Noch kein Konto? <a href="/">Jetzt registrieren</a></p>
      <p class="hinweis"><a href="/passwort-vergessen">Passwort vergessen? 🔑</a></p>

    {:else if schritt === 'weiter'}
      <AuthVollenden
        start={authStart}
        email={authEmail}
        setupToken={authSetupToken}
        mfaToken={authMfaToken}
        onFertig={authFertig}
      />
      <button type="button" class="zurueck" onclick={() => { schritt = 'login'; fehler = ''; }}>← Zurück zum Login</button>
    {/if}
  </div>
</div>

<style>
  .login-wrapper { display: flex; justify-content: center; align-items: center; min-height: 90vh; font-family: sans-serif; }
  .login-box { background: rgba(255, 248, 220, 0.06); backdrop-filter: blur(16px) saturate(1.4); padding: 35px; border-radius: 24px; box-shadow: 0 10px 30px rgba(0,0,0,0.3); max-width: 420px; width: 100%; border: 1px solid rgba(230,168,0,0.25); color: #f5f0e8; }
  :global(html[data-theme='light']) .login-box { background: rgba(255, 252, 235, 0.85); color: #1a1200; }
  .hero-box { background: linear-gradient(135deg, #e6a800, #b87c00); color: #1a0f00; padding: 20px; border-radius: 16px; text-align: center; margin-bottom: 24px; }
  .hero-box h2 { margin: 0; }

  .input-group { text-align: left; margin-bottom: 15px; display: flex; flex-direction: column; gap: 4px; }
  .input-group label { font-weight: 600; font-size: 0.85rem; }
  .input-group input { width: 100%; box-sizing: border-box; }

  .pw-feld { display: flex; gap: 8px; }
  .pw-feld input { flex: 1; }
  .auge { background: rgba(230,168,0,0.15); border-radius: 10px; padding: 0 12px; cursor: pointer; font-size: 1.1rem; }

  .error-msg { color: #ff453a; font-size: 0.85rem; font-weight: bold; margin: 5px 0; }
  .login-btn { width: 100%; }
  .hinweis { text-align: center; font-size: 0.9rem; opacity: 0.8; margin-top: 16px; }
  .bleiben { display: flex; align-items: center; gap: 8px; font-size: 0.88rem; }
  .bleiben input { width: auto; }
  .zurueck { background: none; border: none; color: #f9c932; cursor: pointer; font-size: 0.85rem; margin-top: 14px; width: 100%; }
</style>
