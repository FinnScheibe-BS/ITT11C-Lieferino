<script>
  import { onMount, onDestroy } from 'svelte';
  import { api, getToken } from '$lib/api.js';

  // 🚚 LIVE-TRACKING der zuletzt aufgegebenen Bestellung.
  // Eingeloggt: der Status kommt SERVERSEITIG (manipulationssicher) vom Backend.
  // Sonst: lokale Berechnung aus der vergangenen Zeit (Fallback / offline).

  let PHASEN = $state(['Bestellung erhalten', 'Wird zubereitet', 'Unterwegs', 'Geliefert']);
  let bestellung = $state(null);
  let phase = $state(0);
  let timer;

  // 🗄️ Status vom Backend holen.
  async function ladeStatusVomBackend() {
    if (!bestellung?.nummer || !getToken()) return;
    const res = await api('/api/orders/' + bestellung.nummer + '/status');
    if (res.ok && res.daten) {
      phase = res.daten.phase;
      if (Array.isArray(res.daten.phasen)) PHASEN = res.daten.phasen;
    }
  }

  // Lokale Berechnung (Fallback, wenn nicht eingeloggt / Backend aus).
  function berechnePhaseLokal() {
    if (!bestellung) return;
    const start = new Date(bestellung.datum).getTime();
    const gesamtMs = 30 * 60 * 1000; // 30 Min, 4 Phasen
    const anteil = Math.min(1, Math.max(0, (Date.now() - start) / gesamtMs));
    phase = Math.min(PHASEN.length - 1, Math.floor(anteil * PHASEN.length));
  }

  onMount(() => {
    (async () => {
      // Eingeloggt: neueste Bestellung + echten Status vom Backend.
      if (getToken()) {
        const res = await api('/api/orders');
        if (res.ok && Array.isArray(res.daten) && res.daten.length) {
          bestellung = res.daten[0];
          await ladeStatusVomBackend();
          timer = setInterval(ladeStatusVomBackend, 5000);
          return;
        }
      }
      // Fallback: lokale Kopie.
      const verlauf = JSON.parse(localStorage.getItem('lieferino_bestellungen') || '[]');
      bestellung = verlauf[0] || null;
      berechnePhaseLokal();
      timer = setInterval(berechnePhaseLokal, 5000);
    })();
  });

  onDestroy(() => clearInterval(timer));
</script>

<div class="seite">
  <h1>🚚 Live-Tracking</h1>

  {#if !bestellung}
    <div class="leer">
      <p>Du hast aktuell keine Bestellung zum Verfolgen.</p>
      <a href="/restaurants" class="btn">Jetzt bestellen</a>
    </div>
  {:else}
    <div class="karte">
      <div class="kopf">
        {#if bestellung.nummer}<span class="nummer">{bestellung.nummer}</span>{/if}
        <span class="termin">🕒 Lieferung bis {bestellung.liefertermin} Uhr</span>
      </div>

      <!-- Großer Status-Indikator -->
      <div class="status-gross">
        {#if phase === 0}📝{:else if phase === 1}👨‍🍳{:else if phase === 2}🛵{:else}✅{/if}
        <span>{PHASEN[phase]}</span>
      </div>

      <!-- Fortschritts-Tracker -->
      <div class="tracker">
        {#each PHASEN as p, i}
          <div class="schritt" class:erreicht={i <= phase} class:aktuell={i === phase}>
            <div class="punkt">{i < phase ? '✅' : i === phase ? '🔄' : '⬜'}</div>
            <span>{p}</span>
          </div>
        {/each}
      </div>

      <!-- Fortschrittsbalken -->
      <div class="balken-bg">
        <div class="balken" style="width: {((phase + 1) / PHASEN.length) * 100}%"></div>
      </div>

      <p class="hinweis">Diese Seite aktualisiert sich automatisch. 🔄</p>
      <a href="/bestellungen" class="btn sekundaer">Zu meinen Bestellungen</a>
    </div>
  {/if}
</div>

<style>
  .seite { max-width: 640px; margin: 0 auto; padding: 20px; font-family: sans-serif; }
  h1 { margin-bottom: 20px; }
  .leer { text-align: center; padding: 40px; color: #777; }

  .karte { background: white; border: 1px solid #eee; border-radius: 18px; padding: 24px; }
  .kopf { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; flex-wrap: wrap; gap: 8px; }
  .nummer { background: #f3e5f5; color: #673ab7; font-weight: bold; padding: 3px 12px; border-radius: 12px; letter-spacing: 1px; }
  .termin { color: #555; }

  .status-gross { text-align: center; font-size: 3.5rem; margin: 10px 0 24px; }
  .status-gross span { display: block; font-size: 1.2rem; font-weight: bold; color: #673ab7; margin-top: 8px; }

  .tracker { display: flex; justify-content: space-between; gap: 6px; margin-bottom: 18px; }
  .schritt { flex: 1; text-align: center; opacity: 0.4; transition: opacity 0.3s; }
  .schritt.erreicht { opacity: 1; }
  .schritt .punkt { font-size: 1.4rem; margin-bottom: 6px; }
  .schritt span { font-size: 0.75rem; display: block; }
  .schritt.aktuell span { font-weight: bold; color: #673ab7; }

  .balken-bg { background: #eee; border-radius: 6px; height: 10px; overflow: hidden; }
  .balken { height: 100%; background: linear-gradient(90deg, #7e57c2, #673ab7); transition: width 0.5s ease; }

  .hinweis { color: #999; font-size: 0.85rem; text-align: center; margin: 16px 0; }
  .btn { display: inline-block; background: #673ab7; color: white; text-decoration: none; padding: 12px 20px; border-radius: 12px; font-weight: bold; }
  .btn.sekundaer { background: #f1f1f1; color: #333; display: block; text-align: center; }
</style>
