<script>
  import { onMount } from 'svelte';
  import { generiereCode, sendeVerifizierungsEmail } from '$lib/services/email.js';
  import { pruefePasswortStaerke } from '$lib/services/passwort.js';
  import { dev } from '$app/environment';

  // Ablauf: 'email' -> 'code' -> 'neu' -> 'fertig'
  let schritt = $state('email');

  let email = $state('');
  let user = $state(null);
  let fehler = $state('');

  let korrekterCode = $state('');
  let testCodeHinweis = $state('');
  let codeEingabe = $state('');

  let neuesPasswort = $state('');
  let staerke = $derived(pruefePasswortStaerke(neuesPasswort));

  onMount(() => {
    const gespeichert = localStorage.getItem('lieferino_user');
    if (gespeichert) user = JSON.parse(gespeichert);
  });

  // Schritt 1: E-Mail prüfen + Code "senden".
  // 🚨 BACKEND-HINWEIS: Den Reset-Code MUSS später das Backend erzeugen und per
  // echter E-Mail verschicken (POST /api/auth/passwort-reset). Aus Sicherheits-
  // gründen sollte die Meldung immer gleich sein (egal ob die E-Mail existiert).
  async function emailAbsenden(e) {
    e.preventDefault();
    fehler = '';
    if (!user || user.email !== email.trim()) {
      fehler = 'Kein Konto mit dieser E-Mail gefunden. 🤔';
      return;
    }
    korrekterCode = generiereCode();
    const ergebnis = await sendeVerifizierungsEmail(user.email, korrekterCode);
    testCodeHinweis = ergebnis.testCode ?? '';
    schritt = 'code';
  }

  function codeAbsenden(e) {
    e.preventDefault();
    if (codeEingabe.trim() === korrekterCode) {
      fehler = '';
      schritt = 'neu';
    } else {
      fehler = 'Der Code ist leider falsch. 🔁';
    }
  }

  function passwortSpeichern(e) {
    e.preventDefault();
    if (!staerke.istSicher) return;
    user.passwort = neuesPasswort;
    localStorage.setItem('lieferino_user', JSON.stringify(user));
    schritt = 'fertig';
  }
</script>

<div class="wrap">
  <div class="box">
    <div class="hero"><h2>🔑 Passwort zurücksetzen</h2></div>

    {#if schritt === 'email'}
      <form onsubmit={emailAbsenden}>
        <p>Gib deine E-Mail ein – wir schicken dir einen Code.</p>
        <input type="email" placeholder="name@beispiel.de" bind:value={email} required />
        {#if fehler}<p class="fehler">{fehler}</p>{/if}
        <button type="submit">Code anfordern</button>
      </form>
      <p class="zurueck"><a href="/login">← Zurück zum Login</a></p>

    {:else if schritt === 'code'}
      <form onsubmit={codeAbsenden}>
        <p>Wir haben einen Code an <strong>{email}</strong> geschickt.</p>
        {#if testCodeHinweis}
          <p class="test-code">🧪 Test-Code: <strong>{testCodeHinweis}</strong></p>
        {/if}
        <input type="text" maxlength="6" placeholder="123456" bind:value={codeEingabe} required />
        {#if fehler}<p class="fehler">{fehler}</p>{/if}
        <button type="submit">Code bestätigen</button>
      </form>
      <!-- 🛠️ Nur im Entwicklungsmodus -->
      {#if dev}
        <button type="button" class="dev-skip" onclick={() => (schritt = 'neu')}>🛠️ Dev: überspringen</button>
      {/if}

    {:else if schritt === 'neu'}
      <form onsubmit={passwortSpeichern}>
        <p>Wähle ein neues, sicheres Passwort.</p>
        <input type="password" placeholder="Neues Passwort" bind:value={neuesPasswort} required />
        {#if neuesPasswort.length > 0}
          <div class="staerke">
            <div class="balken-bg"><div class="balken" style="width:{staerke.score * 20}%; background:{staerke.farbe};"></div></div>
            <span style="color:{staerke.farbe};">{staerke.text}</span>
          </div>
        {/if}
        <button type="submit" disabled={!staerke.istSicher}>Passwort speichern</button>
      </form>

    {:else}
      <div class="fertig">
        <p>✅ Dein Passwort wurde geändert!</p>
        <a href="/login" class="btn-link">Jetzt einloggen 🔑</a>
      </div>
    {/if}
  </div>
</div>

<style>
  .wrap { display: flex; justify-content: center; align-items: center; min-height: 90vh; font-family: sans-serif; }
  .box { background: white; padding: 35px; border-radius: 24px; box-shadow: 0 10px 30px rgba(0,0,0,0.1); max-width: 400px; width: 100%; border: 1px solid #eee; }
  .hero { background: #673ab7; color: white; padding: 18px; border-radius: 16px; text-align: center; margin-bottom: 20px; }
  .hero h2 { margin: 0; }
  p { color: #666; font-size: 0.92rem; }
  form { display: flex; flex-direction: column; gap: 12px; }
  input { padding: 11px; border: 1px solid #ddd; border-radius: 10px; font-size: 0.95rem; }
  button { padding: 13px; background: #673ab7; color: white; border: none; border-radius: 12px; font-weight: bold; cursor: pointer; }
  button:disabled { opacity: 0.5; cursor: not-allowed; }
  .fehler { color: #dc3545; font-weight: 600; font-size: 0.85rem; margin: 0; }
  .test-code { background: #fff8e1; border: 1px dashed #ffc107; border-radius: 10px; padding: 8px; color: #8a6d00; text-align: center; }
  .staerke { display: flex; flex-direction: column; gap: 4px; }
  .balken-bg { background: #eee; border-radius: 6px; height: 8px; overflow: hidden; }
  .balken { height: 100%; transition: width 0.3s ease; }
  .staerke span { font-size: 0.8rem; font-weight: bold; }
  .zurueck { text-align: center; margin-top: 14px; }
  .zurueck a { color: #673ab7; }
  .fertig { text-align: center; }
  .btn-link { display: inline-block; margin-top: 14px; background: #673ab7; color: white; text-decoration: none; padding: 12px 20px; border-radius: 12px; font-weight: bold; }
  .dev-skip { width: 100%; margin-top: 10px; padding: 10px; background: #fff3cd; color: #8a6d00; border: 1px dashed #ffc107; border-radius: 10px; cursor: pointer; font-weight: 600; font-size: 0.85rem; }
</style>
