<script>
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { generiereCode, sendeVerifizierungsEmail } from '$lib/services/email.js';
  import { pruefeTotp } from '$lib/services/mfa.js';
  import { login, loginStatus, registriereFehlversuch } from '$lib/stores/auth.js';
  import { dev } from '$app/environment';

  // Ablauf: 'login' (E-Mail+Passwort) -> ggf. 'mfa' (zweiter Faktor) -> fertig.
  let schritt = $state('login');

  let email = $state('');
  let passwort = $state('');
  let fehler = $state('');
  let zeigePasswort = $state(false); // 👁️ Passwort anzeigen/verstecken
  let angemeldetBleiben = $state(false); // ✅ Session nicht ablaufen lassen

  // Der gespeicherte Nutzer (zum Abgleich).
  let user = $state(null);

  // MFA-Zustand
  let mfaCode = $state('');          // erwarteter E-Mail-Code (nur Demo)
  let mfaEingabe = $state('');       // was der Nutzer eintippt
  let mfaFehler = $state('');
  let testCodeHinweis = $state('');  // zeigt den E-Mail-Code im Test an

  onMount(() => {
    const gespeichert = localStorage.getItem('lieferino_user');
    if (gespeichert) user = JSON.parse(gespeichert);
  });

  // Schritt 1: E-Mail + Passwort prüfen.
  // 🚨 BACKEND-HINWEIS: Diese Prüfung MUSS später serverseitig passieren
  // (POST /api/auth/login). Das Vergleichen von Klartext-Passwörtern im Browser
  // ist nur für dieses Schulprojekt so erlaubt.
  async function loginAbsenden(e) {
    e.preventDefault();
    fehler = '';

    // 🔒 Erst prüfen, ob der Login gerade wegen zu vieler Fehlversuche gesperrt ist.
    const status = loginStatus();
    if (status.gesperrt) {
      fehler = `Zu viele Fehlversuche. Bitte ${status.restSek} Sekunden warten. ⏳`;
      return;
    }

    if (!user || user.email !== email.trim()) {
      registriereFehlversuch();
      fehler = 'Kein Konto mit dieser E-Mail gefunden. 🤔';
      return;
    }
    // 🚫 Gebannte Nutzer dürfen sich nicht einloggen.
    if (user.gesperrt) {
      fehler = 'Dieses Konto wurde gesperrt. Bitte wende dich an den Support. 🚫';
      return;
    }
    if (user.passwort !== passwort) {
      registriereFehlversuch();
      const s = loginStatus();
      fehler = s.gesperrt
        ? `Zu viele Fehlversuche. Bitte ${s.restSek} Sekunden warten. ⏳`
        : 'Falsches Passwort. 🔒';
      return;
    }

    // Passwort stimmt – ist 2FA aktiv?
    if (user.mfa?.aktiv) {
      if (user.mfa.methode === 'email') {
        // E-Mail-Code erzeugen und "senden".
        mfaCode = generiereCode();
        const ergebnis = await sendeVerifizierungsEmail(user.email, mfaCode);
        testCodeHinweis = ergebnis.testCode ?? '';
      }
      mfaEingabe = '';
      mfaFehler = '';
      schritt = 'mfa';
    } else {
      einloggenFertig();
    }
  }

  // Schritt 2: zweiten Faktor prüfen.
  async function mfaAbsenden(e) {
    e.preventDefault();
    mfaFehler = '';

    let okay = false;
    if (user.mfa.methode === 'totp') {
      okay = await pruefeTotp(user.mfa.secret, mfaEingabe);
    } else {
      okay = mfaEingabe.trim() === mfaCode;
    }

    // 🆘 Alternativ: ein gültiger Backup-Code (wird danach verbraucht).
    if (!okay && user.mfa.backupCodes?.length) {
      const eingabe = mfaEingabe.trim().toUpperCase();
      const idx = user.mfa.backupCodes.indexOf(eingabe);
      if (idx !== -1) {
        user.mfa.backupCodes.splice(idx, 1); // Code verbrauchen
        localStorage.setItem('lieferino_user', JSON.stringify(user));
        okay = true;
      }
    }

    if (okay) {
      einloggenFertig();
    } else {
      mfaFehler = 'Code/Backup-Code ist leider falsch. 🔁';
    }
  }

  // Login abschließen: Session über den Auth-Store setzen und zur Startseite.
  function einloggenFertig() {
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
            <!-- 👁️ Typ wechselt je nach "anzeigen" -->
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

        <button type="submit" class="login-btn">Einloggen ➡️</button>
      </form>

      <p class="hinweis">Noch kein Konto? <a href="/">Jetzt registrieren</a></p>
      <p class="hinweis"><a href="/passwort-vergessen">Passwort vergessen? 🔑</a></p>

    {:else if schritt === 'mfa'}
      <!-- 🔐 Zweiter Faktor -->
      <p class="step-title">Zwei-Faktor-Bestätigung 🔐</p>
      {#if user.mfa.methode === 'email'}
        <p class="mfa-text">Wir haben einen Code an <strong>{user.email}</strong> geschickt.</p>
        {#if testCodeHinweis}
          <p class="test-code">🧪 Test-Code (kommt später per Mail): <strong>{testCodeHinweis}</strong></p>
        {/if}
      {:else}
        <p class="mfa-text">Gib den 6-stelligen Code aus deiner Authenticator-App ein.</p>
      {/if}

      <form onsubmit={mfaAbsenden}>
        <div class="input-group">
          <label for="code">Code</label>
          <input type="text" id="code" maxlength="6" placeholder="123456" bind:value={mfaEingabe} required />
        </div>
        {#if mfaFehler}<p class="error-msg">{mfaFehler}</p>{/if}
        <button type="submit" class="login-btn">Bestätigen ✅</button>
      </form>

      <!-- 🛠️ Nur im Entwicklungsmodus: 2FA überspringen -->
      {#if dev}
        <button type="button" class="dev-skip" onclick={einloggenFertig}>🛠️ Dev: überspringen</button>
      {/if}
    {/if}
  </div>
</div>

<style>
  .login-wrapper { display: flex; justify-content: center; align-items: center; min-height: 90vh; font-family: sans-serif; }
  .login-box { background: white; padding: 35px; border-radius: 24px; box-shadow: 0 10px 30px rgba(0,0,0,0.1); max-width: 420px; width: 100%; border: 1px solid #eee; }
  .hero-box { background: #673ab7; color: white; padding: 20px; border-radius: 16px; text-align: center; margin-bottom: 24px; }
  .hero-box h2 { margin: 0; }

  .input-group { text-align: left; margin-bottom: 15px; display: flex; flex-direction: column; gap: 4px; }
  .input-group label { font-weight: 600; color: #444; font-size: 0.85rem; }
  .input-group input { padding: 11px; border: 1px solid #ddd; border-radius: 10px; font-size: 0.95rem; box-sizing: border-box; width: 100%; }

  .pw-feld { display: flex; gap: 8px; }
  .pw-feld input { flex: 1; }
  .auge { border: 1px solid #ddd; background: #f7f7f7; border-radius: 10px; padding: 0 12px; cursor: pointer; font-size: 1.1rem; }

  .step-title { color: #673ab7; font-weight: bold; text-align: center; }
  .mfa-text { color: #666; text-align: center; font-size: 0.9rem; }
  .test-code { background: #fff8e1; border: 1px dashed #ffc107; border-radius: 10px; padding: 8px; color: #8a6d00; font-size: 0.85rem; text-align: center; }

  .error-msg { color: #dc3545; font-size: 0.85rem; font-weight: bold; margin: 5px 0; }
  .login-btn { width: 100%; padding: 14px; background: #673ab7; color: white; border: none; border-radius: 12px; font-weight: bold; font-size: 1rem; cursor: pointer; }
  .login-btn:hover { background: #542f95; }
  .hinweis { text-align: center; font-size: 0.9rem; color: #777; margin-top: 16px; }
  .hinweis a { color: #673ab7; font-weight: 600; }
  .bleiben { display: flex; align-items: center; gap: 8px; font-size: 0.88rem; color: #555; }
  .bleiben input { width: auto; }
  .dev-skip { display: block; width: 100%; margin-top: 10px; padding: 10px; background: #fff3cd; color: #8a6d00; border: 1px dashed #ffc107; border-radius: 10px; cursor: pointer; font-weight: 600; font-size: 0.85rem; }
</style>
