<script>
  import { onMount } from 'svelte';
  import { api, getToken } from '$lib/api.js';

  // 🛠️ ADMIN-DASHBOARD – läuft jetzt komplett über das Backend (mit Admin-Rechten).
  // Zugang nur für den Admin-Nutzer (Konto mit der konfigurierten Admin-E-Mail).

  let status = $state('pruefe'); // 'pruefe' | 'keinZugang' | 'ok'
  let aktiverTab = $state('uebersicht');
  let editierModus = $state(false);

  let stats = $state(null);
  let nutzer = $state([]);
  let restaurants = $state([]);
  let reviews = $state([]);

  onMount(() => {
    (async () => {
      if (!getToken()) { status = 'keinZugang'; return; }
      const me = await api('/api/me');
      if (!me.ok || !me.daten?.istAdmin) { status = 'keinZugang'; return; }
      status = 'ok';
      await alleLaden();
    })();
  });

  async function alleLaden() {
    const [s, u, r, b] = await Promise.all([
      api('/api/admin/stats'),
      api('/api/admin/users'),
      api('/api/restaurants'),
      api('/api/admin/reviews')
    ]);
    if (s.ok) stats = s.daten;
    if (u.ok) nutzer = u.daten;
    if (r.ok) restaurants = r.daten;
    if (b.ok) reviews = b.daten;
  }

  // --- Aktionen (Backend) ---
  async function sperren(u, gesperrt) {
    await api(`/api/admin/users/${u.id}/gesperrt`, { method: 'PATCH', body: { gesperrt } });
    await alleLaden();
  }
  async function mfaReset(u) {
    await api(`/api/admin/users/${u.id}/mfa-reset`, { method: 'POST' });
    await alleLaden();
  }
  async function restaurantAktiv(r, aktiv) {
    await api(`/api/admin/restaurants/${r.slug}/aktiv`, { method: 'PATCH', body: { aktiv } });
    await alleLaden();
  }
  async function reviewLoeschen(id) {
    await api(`/api/admin/reviews/${id}`, { method: 'DELETE' });
    await alleLaden();
  }

  function formatiere(iso) {
    return new Date(iso).toLocaleDateString('de-DE', { day: '2-digit', month: '2-digit', year: 'numeric' });
  }
</script>

<div class="seite">
  {#if status === 'pruefe'}
    <p class="lade">Prüfe Admin-Zugang…</p>

  {:else if status === 'keinZugang'}
    <div class="sperre">
      <h1>🔐 Admin-Bereich</h1>
      <p>Dieser Bereich ist nur für Administratoren. Bitte melde dich mit dem Admin-Konto an.</p>
      <a href="/login" class="btn">Zum Login 🔑</a>
    </div>

  {:else}
    <div class="kopf">
      <h1>🛠️ Admin-Dashboard</h1>
    </div>

    <label class="edit-switch" class:an={editierModus}>
      <input type="checkbox" bind:checked={editierModus} />
      ✏️ Editier-Modus {editierModus ? 'AN' : 'aus'}
      <span class="edit-hinweis">{editierModus ? '(Änderungen möglich)' : '(nur Ansicht)'}</span>
    </label>

    <div class="tabs">
      <button class:aktiv={aktiverTab === 'uebersicht'} onclick={() => (aktiverTab = 'uebersicht')}>📊 Übersicht</button>
      <button class:aktiv={aktiverTab === 'lieferanten'} onclick={() => (aktiverTab = 'lieferanten')}>🍽️ Lieferanten</button>
      <button class:aktiv={aktiverTab === 'bewertungen'} onclick={() => (aktiverTab = 'bewertungen')}>⭐ Bewertungen</button>
      <button class:aktiv={aktiverTab === 'nutzer'} onclick={() => (aktiverTab = 'nutzer')}>👥 Nutzer</button>
    </div>

    <!-- TAB: Übersicht -->
    {#if aktiverTab === 'uebersicht' && stats}
      <div class="kacheln">
        <div class="kachel"><span class="zahl">{stats.aktiveRestaurants}</span><span class="label">Aktive Restaurants</span></div>
        <div class="kachel"><span class="zahl">{stats.bestellungen}</span><span class="label">Bestellungen</span></div>
        <div class="kachel"><span class="zahl">{stats.verkaufteArtikel}</span><span class="label">Verkaufte Artikel</span></div>
        <div class="kachel"><span class="zahl">{stats.umsatz.toFixed(2)}€</span><span class="label">Umsatz</span></div>
        <div class="kachel"><span class="zahl">{stats.nutzer}</span><span class="label">Nutzer</span></div>
      </div>
    {/if}

    <!-- TAB: Lieferanten -->
    {#if aktiverTab === 'lieferanten'}
      <h2>🍽️ Lieferanten verwalten</h2>
      {#if !editierModus}<p class="hint">ℹ️ Editier-Modus aktivieren, um zu (de)aktivieren.</p>{/if}
      <div class="tabelle">
        <div class="zeile kopf-zeile"><span>Restaurant</span><span>Typ</span><span>Status</span><span>Aktion</span></div>
        {#each restaurants as r}
          <div class="zeile" class:inaktiv={!r.aktiv}>
            <span>{r.emoji} {r.name}</span>
            <span>{r.typ}</span>
            <span>{#if r.aktiv}<span class="badge gruen">✅ Aktiv</span>{:else}<span class="badge rot">🚫 Deaktiviert</span>{/if}</span>
            <span>
              {#if editierModus}
                <button class="mini-btn {r.aktiv ? 'gefahr' : 'gut'}" onclick={() => restaurantAktiv(r, !r.aktiv)}>
                  {r.aktiv ? '🚫 Deaktivieren' : '✅ Aktivieren'}
                </button>
              {:else}<span class="dim">—</span>{/if}
            </span>
          </div>
        {/each}
      </div>
    {/if}

    <!-- TAB: Bewertungen -->
    {#if aktiverTab === 'bewertungen'}
      <h2>⭐ Bewertungen moderieren</h2>
      {#if !editierModus}<p class="hint">ℹ️ Editier-Modus aktivieren, um zu löschen.</p>{/if}
      {#if reviews.length === 0}
        <p class="leer">Noch keine Bewertungen vorhanden.</p>
      {:else}
        {#each reviews as review}
          <div class="review-zeile">
            <div class="rev-inhalt">
              <strong>{review.name}</strong> {'⭐'.repeat(review.sterne)} · <small>{review.restaurantName} · {formatiere(review.datum)}</small>
              <p>{review.text}</p>
            </div>
            {#if editierModus}
              <button class="mini-btn gefahr" onclick={() => reviewLoeschen(review.id)}>🗑️</button>
            {/if}
          </div>
        {/each}
      {/if}
    {/if}

    <!-- TAB: Nutzer -->
    {#if aktiverTab === 'nutzer'}
      <h2>👥 Nutzer verwalten</h2>
      {#if !editierModus}<p class="hint">ℹ️ Editier-Modus aktivieren, um zu sperren oder MFA zurückzusetzen.</p>{/if}
      <div class="tabelle">
        <div class="zeile kopf-zeile"><span>Nutzer</span><span>E-Mail</span><span>Status</span><span>Aktion</span></div>
        {#each nutzer as u}
          <div class="zeile" class:inaktiv={u.gesperrt}>
            <span>{u.vorname} {u.nachname} {#if u.istAdmin}<span class="badge gold">Admin</span>{/if}</span>
            <span>{u.email}</span>
            <span>
              {#if u.gesperrt}<span class="badge rot">gesperrt</span>{:else}<span class="badge gruen">aktiv</span>{/if}
              {#if u.mfaAktiv}<span class="badge">🔐</span>{/if}
            </span>
            <span class="nutzer-aktion">
              {#if editierModus && !u.istAdmin}
                {#if u.gesperrt}
                  <button class="mini-btn gut" onclick={() => sperren(u, false)}>✅ Entsperren</button>
                {:else}
                  <button class="mini-btn gefahr" onclick={() => sperren(u, true)}>🚫 Sperren</button>
                {/if}
                {#if u.mfaAktiv}<button class="mini-btn" onclick={() => mfaReset(u)}>🔓 MFA reset</button>{/if}
              {:else}<span class="dim">—</span>{/if}
            </span>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>

<style>
  .seite { max-width: 940px; margin: 0 auto; padding: 20px; font-family: sans-serif; }
  .lade { text-align: center; color: #888; padding: 40px; }
  .sperre { max-width: 400px; margin: 60px auto; background: rgba(255,248,220,0.06); backdrop-filter: blur(16px); border: 1px solid rgba(230,168,0,0.25); border-radius: 20px; padding: 30px; text-align: center; color: #f5f0e8; }
  :global(html[data-theme='light']) .sperre { background: rgba(255,252,235,0.85); color: #1a1200; }
  .sperre .btn { display: inline-block; margin-top: 16px; }

  .kopf { display: flex; justify-content: space-between; align-items: center; }
  .edit-switch { display: inline-flex; align-items: center; gap: 8px; background: rgba(230,168,0,0.10); border: 1px solid rgba(230,168,0,0.30); border-radius: 12px; padding: 10px 14px; cursor: pointer; font-weight: 600; margin: 10px 0 18px; }
  .edit-switch.an { background: rgba(230,168,0,0.22); }
  .edit-switch input { width: auto; }
  .edit-hinweis { font-weight: 400; opacity: 0.7; font-size: 0.82rem; }

  .tabs { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 18px; border-bottom: 2px solid rgba(230,168,0,0.2); }
  .tabs button { background: none; border: none; padding: 10px 14px; cursor: pointer; font-weight: 600; opacity: 0.7; border-bottom: 3px solid transparent; color: inherit; }
  .tabs button.aktiv { opacity: 1; color: #f9c932; border-bottom-color: #f9c932; }

  .hint { background: rgba(230,168,0,0.08); border: 1px solid rgba(230,168,0,0.2); border-radius: 10px; padding: 10px 14px; opacity: 0.85; font-size: 0.85rem; }
  .leer { opacity: 0.7; }

  .kacheln { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 16px; }
  .kachel { background: linear-gradient(135deg, #e6a800, #b87c00); color: #1a0f00; border-radius: 16px; padding: 22px; text-align: center; }
  .kachel .zahl { display: block; font-size: 1.9rem; font-weight: bold; }
  .kachel .label { opacity: 0.85; font-size: 0.88rem; }

  .tabelle { border: 1px solid rgba(230,168,0,0.2); border-radius: 14px; overflow: hidden; }
  .zeile { display: grid; grid-template-columns: 2fr 2fr 1.4fr 1.6fr; padding: 12px 16px; border-bottom: 1px solid rgba(230,168,0,0.12); align-items: center; gap: 8px; }
  .zeile:last-child { border-bottom: none; }
  .zeile.kopf-zeile { background: rgba(230,168,0,0.10); font-weight: bold; }
  .zeile.inaktiv { opacity: 0.55; }

  .badge { padding: 2px 10px; border-radius: 12px; font-size: 0.75rem; font-weight: bold; background: rgba(255,255,255,0.1); }
  .badge.gruen { background: #e8f5e9; color: #2e7d32; }
  .badge.rot { background: #ffebee; color: #c62828; }
  .badge.gold { background: #f9c932; color: #1a0f00; }

  .mini-btn { background: rgba(230,168,0,0.15); border: none; border-radius: 8px; padding: 5px 10px; cursor: pointer; font-size: 0.8rem; margin-left: 4px; color: inherit; }
  .mini-btn.gefahr { background: #ffebee; color: #c62828; }
  .mini-btn.gut { background: #e8f5e9; color: #2e7d32; }
  .dim { opacity: 0.4; }
  .nutzer-aktion { display: flex; flex-wrap: wrap; align-items: center; gap: 4px; }

  .review-zeile { display: flex; justify-content: space-between; align-items: center; gap: 10px; border: 1px solid rgba(230,168,0,0.2); border-radius: 12px; padding: 12px; margin-bottom: 8px; }
  .rev-inhalt p { margin: 4px 0 0; opacity: 0.85; font-size: 0.9rem; }
</style>
