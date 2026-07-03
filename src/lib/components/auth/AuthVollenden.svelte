<script>
  // 🔐 Schließt die Anmeldung ab – je nach Phase:
  //   'verify'    -> E-Mail-Code eingeben  (-> danach MFA-Setup)
  //   'mfaSetup'  -> Authenticator einrichten (QR + erster Code)
  //   'mfaVerify' -> beim Login den MFA-Code eingeben
  // Bei Erfolg ruft die Komponente onFertig({ token, user }) auf.
  import { api, setzeToken } from '$lib/api/api.js';
  import { t } from '$lib/utils/i18n.js';

  let { start = 'verify', email = '', setupToken = '', mfaToken = '', onFertig } = $props();

  let phase = $state(start);
  let stoken = $state(setupToken); // Setup-Token (für MFA-Einrichtung)
  let mtoken = $state(mfaToken);   // MFA-Token (für Login-Schritt 2)

  let code = $state('');
  let fehler = $state('');
  let laedt = $state(false);
  let infoText = $state('');

  let qr = $state('');
  let secret = $state('');
  let setupGeladen = $state(false);

  // Sobald wir in der MFA-Einrichtung sind: QR-Code vom Backend holen.
  $effect(() => {
    if (phase === 'mfaSetup' && !setupGeladen && stoken) {
      setupGeladen = true;
      ladeMfaSetup();
    }
  });

  async function ladeMfaSetup() {
    laedt = true; fehler = '';
    const res = await api('/api/auth/mfa/setup', { method: 'POST', token: stoken });
    laedt = false;
    if (res.ok && res.daten?.qr) {
      qr = res.daten.qr;
      secret = res.daten.secret;
    } else {
      fehler = res.daten?.fehler || 'Die MFA-Einrichtung konnte nicht gestartet werden.';
    }
  }

  // E-Mail-Code prüfen -> danach MFA einrichten.
  async function verifyAbsenden(e) {
    e.preventDefault(); fehler = ''; laedt = true;
    const res = await api('/api/auth/verify-email', { method: 'POST', body: { email, code: code.trim() } });
    laedt = false;
    if (res.ok && res.daten?.setupToken) {
      stoken = res.daten.setupToken;
      code = '';
      phase = 'mfaSetup';
    } else {
      fehler = res.daten?.fehler || 'Der Code ist nicht korrekt.';
    }
  }

  // Ersten Authenticator-Code prüfen -> MFA aktivieren -> fertig.
  async function enableAbsenden(e) {
    e.preventDefault(); fehler = ''; laedt = true;
    const res = await api('/api/auth/mfa/enable', { method: 'POST', token: stoken, body: { code: code.trim() } });
    laedt = false;
    if (res.ok && res.daten?.token) fertig(res.daten);
    else fehler = res.daten?.fehler || 'Der Code stimmt nicht.';
  }

  // Login-Schritt 2: MFA-Code prüfen.
  async function verifyLoginAbsenden(e) {
    e.preventDefault(); fehler = ''; laedt = true;
    const res = await api('/api/auth/mfa/verify', { method: 'POST', token: mtoken, body: { code: code.trim() } });
    laedt = false;
    if (res.ok && res.daten?.token) fertig(res.daten);
    else fehler = res.daten?.fehler || 'Der Code stimmt nicht.';
  }

  async function codeErneutSenden() {
    fehler = ''; infoText = '';
    await api('/api/auth/resend-code', { method: 'POST', body: { email } });
    infoText = $t('auth.resent');
  }

  function fertig(daten) {
    setzeToken(daten.token);
    onFertig?.(daten);
  }
</script>

<div class="auth-vollenden">
  {#if phase === 'verify'}
    <h3>{$t('auth.verify_title')}</h3>
    <p class="hint">{$t('auth.verify_sub').replace('{email}', email)}</p>
    <form onsubmit={verifyAbsenden}>
      <input type="text" inputmode="numeric" maxlength="6" placeholder="123456" bind:value={code} required />
      {#if fehler}<p class="err">{fehler}</p>{/if}
      {#if infoText}<p class="info">{infoText}</p>{/if}
      <button type="submit" class="btn" disabled={laedt}>{laedt ? $t('auth.checking') : $t('auth.verify_btn')}</button>
    </form>
    <button type="button" class="link" onclick={codeErneutSenden}>{$t('auth.resend')}</button>

  {:else if phase === 'mfaSetup'}
    <h3>{$t('auth.mfa_setup_title')}</h3>
    <p class="hint">{$t('auth.mfa_setup_sub')}</p>
    {#if qr}
      <div class="qr-box"><img src={qr} alt="QR-Code" /></div>
      <p class="secret">{$t('auth.mfa_secret_hint')}<br /><code>{secret}</code></p>
    {:else if laedt}
      <p class="hint">{$t('auth.qr_loading')}</p>
    {/if}
    <form onsubmit={enableAbsenden}>
      <input type="text" inputmode="numeric" maxlength="6" placeholder="123456" bind:value={code} required />
      {#if fehler}<p class="err">{fehler}</p>{/if}
      <button type="submit" class="btn" disabled={laedt || !qr}>{laedt ? $t('auth.checking') : $t('auth.mfa_enable_btn')}</button>
    </form>

  {:else if phase === 'mfaVerify'}
    <h3>{$t('auth.mfa_verify_title')}</h3>
    <p class="hint">{$t('auth.mfa_verify_sub')}</p>
    <form onsubmit={verifyLoginAbsenden}>
      <input type="text" inputmode="numeric" maxlength="6" placeholder="123456" bind:value={code} required />
      {#if fehler}<p class="err">{fehler}</p>{/if}
      <button type="submit" class="btn" disabled={laedt}>{laedt ? $t('auth.checking') : $t('auth.confirm')}</button>
    </form>
  {/if}
</div>

<style>
  .auth-vollenden { display: flex; flex-direction: column; gap: 12px; text-align: center; }
  .auth-vollenden h3 { margin: 0; }
  .hint { font-size: 0.9rem; opacity: 0.85; margin: 0; }
  .auth-vollenden form { display: flex; flex-direction: column; gap: 10px; margin: 0; }
  .auth-vollenden input { text-align: center; letter-spacing: 0.3em; font-size: 1.2rem; }
  .btn { width: 100%; }
  .err { color: #ff453a; font-weight: 600; font-size: 0.85rem; margin: 0; }
  .info { color: #34c759; font-weight: 600; font-size: 0.85rem; margin: 0; }
  .link { background: none; border: none; color: #f9c932; cursor: pointer; font-size: 0.85rem; text-decoration: underline; padding: 0; }
  /* QR immer auf weißem Grund (damit er in Dark & Light scanbar bleibt) */
  .qr-box { background: #fff; padding: 12px; border-radius: 12px; display: inline-block; margin: 0 auto; }
  .qr-box img { display: block; width: 200px; height: 200px; }
  .secret { font-size: 0.8rem; opacity: 0.85; margin: 0; }
  .secret code { font-size: 0.9rem; letter-spacing: 0.1em; word-break: break-all; }
</style>
