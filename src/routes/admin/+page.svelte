<script>
  import { onMount } from 'svelte';
  import { restaurants } from '$lib/data';
  import { geloeschteLieferanten, loescheLieferant, stelleLieferantWiederHer } from '$lib/stores/lieferanten.js';
  import { bewertungen, bewertungLoeschen, bewertungBearbeiten } from '$lib/stores/bewertungen.js';

  // 🛠️ ADMIN-DASHBOARD mit Tabs + Editier-Modus.
  // 🚨 BACKEND-HINWEIS: Alles hier läuft lokal im Browser. Echte Verwaltung
  // (Restaurants, Bewertungen, Nutzer, Sperren, MFA-Reset) muss das Backend
  // über geschützte Endpunkte (mit Admin-Rechten) bereitstellen.

  let bestellungen = $state([]);

  // 🔐 Sicherheitsschlüssel (im Klartext gewünscht – nur Schein-Schutz, siehe Hinweis).
  const ADMIN_KEY = 'Lieferino#2026_DevOps';
  let entsperrt = $state(false);
  let schluesselEingabe = $state('');
  let schluesselFehler = $state('');

  // Welcher Tab ist aktiv + ob der Editier-Modus an ist.
  let aktiverTab = $state('uebersicht');
  let editierModus = $state(false);

  // Nutzer (in diesem Frontend-Demo gibt es das eine lokale Konto).
  let nutzer = $state([]);

  // Review-Bearbeitung (welche Bewertung wird gerade bearbeitet?)
  let editSlug = $state(null);
  let editIndex = $state(-1);
  let editText = $state('');
  let editSterne = $state(5);

  onMount(() => {
    bestellungen = JSON.parse(localStorage.getItem('lieferino_bestellungen') || '[]');
    entsperrt = sessionStorage.getItem('lieferino_admin') === 'true';
    ladeNutzer();
  });

  function entsperren(e) {
    e.preventDefault();
    if (schluesselEingabe === ADMIN_KEY) {
      entsperrt = true;
      sessionStorage.setItem('lieferino_admin', 'true');
      schluesselFehler = '';
    } else {
      schluesselFehler = 'Falscher Sicherheitsschlüssel. 🔒';
    }
  }

  function adminSperren() {
    sessionStorage.removeItem('lieferino_admin');
    entsperrt = false;
    editierModus = false;
    schluesselEingabe = '';
  }

  // --- Kennzahlen ---
  let umsatz = $derived(bestellungen.reduce((s, b) => s + b.summe, 0));
  let verkaufteArtikel = $derived(
    bestellungen.reduce((s, b) => s + b.artikel.reduce((a, art) => a + art.menge, 0), 0)
  );

  let umsatzProTag = $derived.by(() => {
    const tage = [];
    const heute = new Date();
    for (let i = 6; i >= 0; i--) {
      const d = new Date(heute);
      d.setDate(heute.getDate() - i);
      const schluessel = d.toISOString().slice(0, 10);
      tage.push({ schluessel, label: d.toLocaleDateString('de-DE', { weekday: 'short' }), summe: 0 });
    }
    for (const b of bestellungen) {
      const schluessel = (b.datum || '').slice(0, 10);
      const tag = tage.find((t) => t.schluessel === schluessel);
      if (tag) tag.summe += b.summe;
    }
    return tage;
  });
  let maxTag = $derived(Math.max(1, ...umsatzProTag.map((t) => t.summe)));

  // --- Nutzerverwaltung (lokales Konto) ---
  function ladeNutzer() {
    const u = JSON.parse(localStorage.getItem('lieferino_user') || 'null');
    nutzer = u ? [u] : [];
  }
  function speichereNutzer(u) {
    localStorage.setItem('lieferino_user', JSON.stringify(u));
    ladeNutzer();
  }
  function bannen(u) {
    speichereNutzer({ ...u, gesperrt: true });
  }
  function entbannen(u) {
    speichereNutzer({ ...u, gesperrt: false });
  }
  function mfaZuruecksetzen(u) {
    speichereNutzer({ ...u, mfa: { aktiv: false } });
  }

  // --- Review-Bearbeitung ---
  function starteReviewEdit(slug, index, review) {
    editSlug = slug;
    editIndex = index;
    editText = review.text;
    editSterne = review.sterne;
  }
  function speichereReviewEdit() {
    bewertungBearbeiten(editSlug, editIndex, { text: editText, sterne: Number(editSterne) });
    editSlug = null;
    editIndex = -1;
  }
</script>

<div class="seite">
  {#if !entsperrt}
    <!-- 🔐 Sperrbildschirm -->
    <div class="sperre">
      <h1>🔐 Admin-Bereich</h1>
      <p>Dieser Bereich ist geschützt. Bitte gib den Sicherheitsschlüssel ein.</p>
      <form onsubmit={entsperren}>
        <input type="password" placeholder="Sicherheitsschlüssel" bind:value={schluesselEingabe} />
        {#if schluesselFehler}<p class="sperre-fehler">{schluesselFehler}</p>{/if}
        <button type="submit">Entsperren 🔓</button>
      </form>
    </div>
  {:else}
    <div class="kopf">
      <h1>🛠️ Admin-Dashboard</h1>
      <button class="sperren-btn" onclick={adminSperren}>🔒 Sperren</button>
    </div>

    <!-- ✏️ Editier-Modus-Schalter -->
    <label class="edit-switch" class:an={editierModus}>
      <input type="checkbox" bind:checked={editierModus} />
      ✏️ Editier-Modus {editierModus ? 'AN' : 'aus'}
      <span class="edit-hinweis">{editierModus ? '(Änderungen/Löschen möglich)' : '(nur Ansicht)'}</span>
    </label>

    <!-- Tabs -->
    <div class="tabs">
      <button class:aktiv={aktiverTab === 'uebersicht'} onclick={() => (aktiverTab = 'uebersicht')}>📊 Übersicht</button>
      <button class:aktiv={aktiverTab === 'lieferanten'} onclick={() => (aktiverTab = 'lieferanten')}>🍽️ Lieferanten</button>
      <button class:aktiv={aktiverTab === 'bewertungen'} onclick={() => (aktiverTab = 'bewertungen')}>⭐ Bewertungen</button>
      <button class:aktiv={aktiverTab === 'nutzer'} onclick={() => (aktiverTab = 'nutzer')}>👥 Nutzer</button>
    </div>

    <!-- TAB: Übersicht -->
    {#if aktiverTab === 'uebersicht'}
      <div class="kacheln">
        <div class="kachel"><span class="zahl">{$geloeschteLieferanten.length ? restaurants.length - $geloeschteLieferanten.length : restaurants.length}</span><span class="label">Aktive Restaurants</span></div>
        <div class="kachel"><span class="zahl">{bestellungen.length}</span><span class="label">Bestellungen</span></div>
        <div class="kachel"><span class="zahl">{verkaufteArtikel}</span><span class="label">Verkaufte Artikel</span></div>
        <div class="kachel"><span class="zahl">{umsatz.toFixed(2)}€</span><span class="label">Umsatz</span></div>
      </div>

      <h2>📊 Umsatz (letzte 7 Tage)</h2>
      <div class="chart">
        {#each umsatzProTag as tag}
          <div class="balken-spalte">
            <span class="balken-wert">{tag.summe > 0 ? tag.summe.toFixed(0) + '€' : ''}</span>
            <div class="balken" style="height: {(tag.summe / maxTag) * 100}%"></div>
            <span class="balken-label">{tag.label}</span>
          </div>
        {/each}
      </div>
    {/if}

    <!-- TAB: Lieferanten -->
    {#if aktiverTab === 'lieferanten'}
      <h2>🍽️ Lieferanten verwalten</h2>
      {#if !editierModus}<p class="hint">ℹ️ Aktiviere den Editier-Modus, um Lieferanten zu löschen.</p>{/if}
      <div class="tabelle">
        <div class="zeile kopf-zeile">
          <span>Restaurant</span><span>Typ</span><span>Bewertung</span><span>Status / Aktion</span>
        </div>
        {#each restaurants as r}
          {@const geloescht = $geloeschteLieferanten.includes(r.slug)}
          <div class="zeile" class:inaktiv={geloescht}>
            <span>{r.emoji} {r.name}</span>
            <span>{r.typ}</span>
            <span>⭐ {r.bewertung}</span>
            <span>
              {#if geloescht}
                <span class="badge rot">gelöscht</span>
                {#if editierModus}
                  <button class="mini-btn" onclick={() => stelleLieferantWiederHer(r.slug)}>↩️ Wiederherstellen</button>
                {/if}
              {:else if editierModus}
                <button class="mini-btn gefahr" onclick={() => loescheLieferant(r.slug)}>🗑️ Löschen</button>
              {:else}
                <span class="badge gruen">aktiv</span>
              {/if}
            </span>
          </div>
        {/each}
      </div>
    {/if}

    <!-- TAB: Bewertungen -->
    {#if aktiverTab === 'bewertungen'}
      <h2>⭐ Bewertungen verwalten</h2>
      {#if !editierModus}<p class="hint">ℹ️ Aktiviere den Editier-Modus, um Bewertungen zu bearbeiten oder zu löschen.</p>{/if}
      {#each restaurants as r}
        {#if ($bewertungen[r.slug] || []).length > 0}
          <h3 class="rev-titel">{r.emoji} {r.name}</h3>
          {#each $bewertungen[r.slug] as review, index}
            <div class="review-zeile">
              {#if editierModus && editSlug === r.slug && editIndex === index}
                <!-- Bearbeiten-Ansicht -->
                <select bind:value={editSterne}>
                  {#each [5, 4, 3, 2, 1] as s}<option value={s}>{s} ⭐</option>{/each}
                </select>
                <input type="text" bind:value={editText} class="edit-input" />
                <button class="mini-btn" onclick={speichereReviewEdit}>💾</button>
                <button class="mini-btn" onclick={() => (editSlug = null)}>✖️</button>
              {:else}
                <div class="rev-inhalt">
                  <strong>{review.name}</strong> {'⭐'.repeat(review.sterne)}
                  <p>{review.text}</p>
                </div>
                {#if editierModus}
                  <div class="rev-aktionen">
                    <button class="mini-btn" onclick={() => starteReviewEdit(r.slug, index, review)}>✏️</button>
                    <button class="mini-btn gefahr" onclick={() => bewertungLoeschen(r.slug, index)}>🗑️</button>
                  </div>
                {/if}
              {/if}
            </div>
          {/each}
        {/if}
      {/each}
      {#if Object.values($bewertungen).every((l) => l.length === 0)}
        <p class="leer">Noch keine Bewertungen vorhanden.</p>
      {/if}
    {/if}

    <!-- TAB: Nutzer -->
    {#if aktiverTab === 'nutzer'}
      <h2>👥 Nutzer verwalten</h2>
      <p class="hint">🚨 In diesem Frontend-Demo gibt es nur das lokal angelegte Konto.
        Eine echte Nutzerliste aller Kunden muss das Backend liefern (GET /api/admin/users).</p>
      {#if !editierModus}<p class="hint">ℹ️ Aktiviere den Editier-Modus, um Nutzer zu bannen oder MFA zurückzusetzen.</p>{/if}

      {#if nutzer.length === 0}
        <p class="leer">Kein aktives Nutzerkonto gefunden.</p>
      {:else}
        <div class="tabelle">
          <div class="zeile kopf-zeile"><span>Nutzer</span><span>E-Mail</span><span>2FA</span><span>Status / Aktion</span></div>
          {#each nutzer as u}
            <div class="zeile" class:inaktiv={u.gesperrt}>
              <span>{u.vorname} {u.nachname} <small>(@{u.username})</small></span>
              <span>{u.email}</span>
              <span>{u.mfa?.aktiv ? '🔐 ' + (u.mfa.methode === 'totp' ? 'App' : 'E-Mail') : '—'}</span>
              <span class="nutzer-aktion">
                {#if u.gesperrt}<span class="badge rot">gebannt</span>{:else}<span class="badge gruen">aktiv</span>{/if}
                {#if editierModus}
                  {#if u.gesperrt}
                    <button class="mini-btn" onclick={() => entbannen(u)}>✅ Entbannen</button>
                  {:else}
                    <button class="mini-btn gefahr" onclick={() => bannen(u)}>🚫 Bannen</button>
                  {/if}
                  {#if u.mfa?.aktiv}
                    <button class="mini-btn" onclick={() => mfaZuruecksetzen(u)}>🔓 MFA reset</button>
                  {/if}
                {/if}
              </span>
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  {/if}
</div>

<style>
  .seite { max-width: 940px; margin: 0 auto; padding: 20px; font-family: sans-serif; }

  /* 🔐 Sperrbildschirm */
  .sperre { max-width: 380px; margin: 60px auto; background: white; border: 1px solid #eee; border-radius: 20px; padding: 30px; text-align: center; box-shadow: 0 10px 30px rgba(0,0,0,0.06); }
  .sperre p { color: #777; font-size: 0.9rem; }
  .sperre form { display: flex; flex-direction: column; gap: 12px; margin-top: 16px; }
  .sperre input { padding: 12px; border: 1px solid #ddd; border-radius: 10px; font-size: 0.95rem; }
  .sperre button { padding: 13px; background: #673ab7; color: white; border: none; border-radius: 12px; font-weight: bold; cursor: pointer; }
  .sperre-fehler { color: #dc3545; font-weight: 600; font-size: 0.85rem; margin: 0; }

  .kopf { display: flex; justify-content: space-between; align-items: center; }
  .sperren-btn { background: #f1f1f1; border: none; border-radius: 10px; padding: 8px 14px; cursor: pointer; font-weight: 600; color: #555; }

  /* ✏️ Editier-Modus-Schalter */
  .edit-switch { display: inline-flex; align-items: center; gap: 8px; background: #f7f7f7; border: 1px solid #e5e5e5; border-radius: 12px; padding: 10px 14px; cursor: pointer; font-weight: 600; margin: 10px 0 18px; }
  .edit-switch.an { background: #fff8e1; border-color: #ffc107; color: #8a6d00; }
  .edit-hinweis { font-weight: 400; color: #888; font-size: 0.82rem; }

  /* Tabs */
  .tabs { display: flex; gap: 8px; flex-wrap: wrap; margin-bottom: 18px; border-bottom: 2px solid #eee; }
  .tabs button { background: none; border: none; padding: 10px 14px; cursor: pointer; font-weight: 600; color: #777; border-bottom: 3px solid transparent; }
  .tabs button.aktiv { color: #673ab7; border-bottom-color: #673ab7; }

  .hint { background: #faf7ff; border: 1px solid #ede7f6; border-radius: 10px; padding: 10px 14px; color: #666; font-size: 0.85rem; }
  .leer { color: #777; }

  .kacheln { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 16px; }
  .kachel { background: #673ab7; color: white; border-radius: 16px; padding: 24px; text-align: center; }
  .kachel .zahl { display: block; font-size: 2rem; font-weight: bold; }
  .kachel .label { opacity: 0.85; font-size: 0.9rem; }

  .chart { display: flex; align-items: flex-end; gap: 10px; height: 200px; background: white; border: 1px solid #eee; border-radius: 14px; padding: 20px; }
  .balken-spalte { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: flex-end; height: 100%; gap: 6px; }
  .balken-wert { font-size: 0.72rem; color: #673ab7; font-weight: bold; }
  .balken { width: 100%; max-width: 46px; background: linear-gradient(180deg, #7e57c2, #673ab7); border-radius: 8px 8px 0 0; min-height: 3px; transition: height 0.3s ease; }
  .balken-label { font-size: 0.78rem; color: #888; }

  .tabelle { background: white; border: 1px solid #eee; border-radius: 14px; overflow: hidden; }
  .zeile { display: grid; grid-template-columns: 2fr 1.5fr 1fr 1.6fr; padding: 12px 16px; border-bottom: 1px solid #f0f0f0; align-items: center; gap: 8px; }
  .zeile:last-child { border-bottom: none; }
  .zeile.kopf-zeile { background: #faf7ff; font-weight: bold; }
  .zeile.inaktiv { opacity: 0.55; }

  .badge { padding: 2px 10px; border-radius: 12px; font-size: 0.75rem; font-weight: bold; }
  .badge.gruen { background: #e8f5e9; color: #2e7d32; }
  .badge.rot { background: #ffebee; color: #c62828; }

  .mini-btn { background: #f1f1f1; border: none; border-radius: 8px; padding: 5px 10px; cursor: pointer; font-size: 0.8rem; margin-left: 4px; }
  .mini-btn.gefahr { background: #ffebee; color: #c62828; }
  .nutzer-aktion { display: flex; flex-wrap: wrap; align-items: center; gap: 4px; }

  .rev-titel { margin: 18px 0 6px; }
  .review-zeile { display: flex; justify-content: space-between; align-items: center; gap: 10px; background: white; border: 1px solid #eee; border-radius: 12px; padding: 12px; margin-bottom: 8px; }
  .rev-inhalt p { margin: 4px 0 0; color: #555; font-size: 0.9rem; }
  .rev-aktionen { white-space: nowrap; }
  .edit-input { flex: 1; padding: 8px; border: 1px solid #ddd; border-radius: 8px; }
</style>
