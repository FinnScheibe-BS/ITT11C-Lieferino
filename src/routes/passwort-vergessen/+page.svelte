<script>
  import { pruefePasswortStaerke } from '$lib/services/passwort.js';
  import { api } from '$lib/api.js';
  import { t } from '$lib/i18n.js';

  // Ablauf: 'email' -> 'code' -> 'neu' -> 'fertig'
  let schritt = $state('email');

  let email = $state('');
  let fehler = $state('');
  let laedt = $state(false);
  let resetToken = $state('');

  let codeEingabe = $state('');

  let neuesPasswort = $state('');
  let passwortBestaetigung = $state('');
  let staerke = $derived(pruefePasswortStaerke(neuesPasswort));
  let anforderungen = $state([]);

  // Schritt 1: Code anfordern. Antwort ist IMMER gleich (kein Verraten der E-Mail).
  async function emailAbsenden(e) {
    e.preventDefault();
    fehler = '';
    laedt = true;
    await api('/api/auth/reset-anfordern', { method: 'POST', body: { email: email.trim() } });
    laedt = false;
    schritt = 'code';
  }

  // Schritt 2: Code prüfen -> Reset-Token.
  async function codeAbsenden(e) {
    e.preventDefault();
    fehler = '';
    laedt = true;
    const res = await api('/api/auth/reset-code', { method: 'POST', body: { email: email.trim(), code: codeEingabe.trim() } });
    laedt = false;
    if (res.ok && res.daten?.resetToken) {
      resetToken = res.daten.resetToken;
      schritt = 'neu';
    } else {
      fehler = res.daten?.fehler || 'Der Code ist nicht korrekt.';
    }
  }

  async function codeErneut() {
    fehler = '';
    await api('/api/auth/reset-anfordern', { method: 'POST', body: { email: email.trim() } });
  }

  // Schritt 3: neues Passwort setzen (mit Bestätigung).
  async function passwortSpeichern(e) {
    e.preventDefault();
    fehler = '';
    anforderungen = [];
    if (!staerke.istSicher) return;
    if (neuesPasswort !== passwortBestaetigung) {
      fehler = $t('reset.mismatch');
      return;
    }
    laedt = true;
    const res = await api('/api/auth/reset-neu', { method: 'POST', token: resetToken, body: { passwort: neuesPasswort } });
    laedt = false;
    if (res.ok) {
      schritt = 'fertig';
    } else {
      fehler = res.daten?.fehler || 'Passwort konnte nicht geändert werden.';
      anforderungen = res.daten?.anforderungen || [];
    }
  }
</script>

<div class="wrap">
  <div class="box">
    <div class="hero"><h2>{$t('reset.title')}</h2></div>

    {#if schritt === 'email'}
      <form onsubmit={emailAbsenden}>
        <p>{$t('reset.email_sub')}</p>
        <input type="email" placeholder="name@beispiel.de" bind:value={email} required />
        <button type="submit" disabled={laedt}>{laedt ? $t('auth.sending') : $t('reset.request_btn')}</button>
      </form>
      <p class="zurueck"><a href="/login">{$t('auth.back_login')}</a></p>

    {:else if schritt === 'code'}
      <form onsubmit={codeAbsenden}>
        <p>{$t('reset.code_sub').replace('{email}', email)}</p>
        <input type="text" inputmode="numeric" maxlength="6" placeholder="123456" bind:value={codeEingabe} required />
        {#if fehler}<p class="fehler">{fehler}</p>{/if}
        <button type="submit" disabled={laedt}>{laedt ? $t('auth.checking') : $t('reset.code_btn')}</button>
      </form>
      <button type="button" class="link" onclick={codeErneut}>{$t('auth.resend')}</button>

    {:else if schritt === 'neu'}
      <form onsubmit={passwortSpeichern}>
        <p>{$t('reset.new_sub')}</p>
        <input type="password" placeholder={$t('reset.new_ph')} bind:value={neuesPasswort} required />
        {#if neuesPasswort.length > 0}
          <ul class="pw-rules">
            <li class:ok={staerke.regeln.laenge}>{$t('pw.rule_length')}</li>
            <li class:ok={staerke.regeln.grossbuchstabe}>{$t('pw.rule_upper')}</li>
            <li class:ok={staerke.regeln.kleinbuchstabe}>{$t('pw.rule_lower')}</li>
            <li class:ok={staerke.regeln.zahl}>{$t('pw.rule_digit')}</li>
            <li class:ok={staerke.regeln.sonderzeichen}>{$t('pw.rule_special')}</li>
          </ul>
        {/if}
        <input type="password" placeholder={$t('reset.confirm_ph')} bind:value={passwortBestaetigung} required />
        {#if passwortBestaetigung.length > 0 && neuesPasswort !== passwortBestaetigung}
          <p class="fehler">{$t('reset.mismatch')}</p>
        {/if}
        {#if fehler}<p class="fehler">{fehler}</p>{/if}
        {#if anforderungen.length > 0}
          <ul class="fehler-liste">{#each anforderungen as a}<li>{a}</li>{/each}</ul>
        {/if}
        <button type="submit" disabled={laedt || !staerke.istSicher || neuesPasswort !== passwortBestaetigung}>
          {laedt ? $t('reset.saving') : $t('reset.save_btn')}
        </button>
      </form>

    {:else}
      <div class="fertig">
        <p>{$t('reset.done')}</p>
        <a href="/login" class="btn-link">{$t('reset.login_now')}</a>
      </div>
    {/if}
  </div>
</div>

<style>
  .wrap { display: flex; justify-content: center; align-items: center; min-height: 90vh; font-family: sans-serif; }
  .box { background: rgba(255, 248, 220, 0.06); backdrop-filter: blur(16px) saturate(1.4); padding: 35px; border-radius: 24px; box-shadow: 0 10px 30px rgba(0,0,0,0.3); max-width: 400px; width: 100%; border: 1px solid rgba(230,168,0,0.25); color: #f5f0e8; }
  :global(html[data-theme='light']) .box { background: rgba(255, 252, 235, 0.85); color: #1a1200; }
  .hero { background: linear-gradient(135deg, #e6a800, #b87c00); color: #1a0f00; padding: 18px; border-radius: 16px; text-align: center; margin-bottom: 20px; }
  .hero h2 { margin: 0; }
  p { font-size: 0.92rem; opacity: 0.9; }
  form { display: flex; flex-direction: column; gap: 12px; }
  button[type='submit'] { width: 100%; }
  button:disabled { opacity: 0.5; cursor: not-allowed; }
  .fehler { color: #ff453a; font-weight: 600; font-size: 0.85rem; margin: 0; }
  .fehler-liste { margin: 0; padding-left: 18px; }
  .fehler-liste li { color: #ff6961; font-size: 0.8rem; }
  .link { background: none; border: none; color: #f9c932; cursor: pointer; font-size: 0.85rem; text-decoration: underline; margin-top: 10px; width: 100%; }
  .pw-rules { list-style: none; padding: 0; margin: 0; display: grid; grid-template-columns: 1fr 1fr; gap: 2px 10px; font-size: 0.78rem; opacity: 0.7; }
  .pw-rules li::before { content: '○ '; }
  .pw-rules li.ok { opacity: 1; color: #34c759; }
  .pw-rules li.ok::before { content: '✓ '; }
  .zurueck { text-align: center; margin-top: 14px; }
  .zurueck a { color: #f9c932; }
  .fertig { text-align: center; }
  .btn-link { display: inline-block; margin-top: 14px; }
</style>
