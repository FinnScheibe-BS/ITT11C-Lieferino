<script>
  import { emojiCursor, saison } from '$lib/stores/funmodus.js';
  import { onMount } from 'svelte';

  const CURSOR_EMOJIS = ['🍕', '🍔', '🌮', '🍣', '🍟', '🥗'];
  const SAISON_EMOJIS = ['❄️', '🍂', '⭐', '✨'];

  let letzteSpur = 0;

  // 🖱️ Emoji-Spur: bei Mausbewegung kleine Emojis hinterlassen.
  function maus(e) {
    if (!$emojiCursor) return;
    const jetzt = performance.now();
    if (jetzt - letzteSpur < 70) return; // throttlen
    letzteSpur = jetzt;

    const el = document.createElement('div');
    el.textContent = CURSOR_EMOJIS[Math.floor(Math.random() * CURSOR_EMOJIS.length)];
    el.style.cssText = `position:fixed;left:${e.clientX}px;top:${e.clientY}px;font-size:20px;pointer-events:none;z-index:999990;transition:transform .8s ease,opacity .8s ease;`;
    document.body.appendChild(el);
    requestAnimationFrame(() => {
      el.style.transform = 'translateY(24px) scale(0.4)';
      el.style.opacity = '0';
    });
    setTimeout(() => el.remove(), 850);
  }

 onMount(() => {
  window.addEventListener('mousemove', maus);
  return () => window.removeEventListener('mousemove', maus);
});

  // Feste Anzahl fallender Teilchen für den Saison-Effekt.
  const teilchen = Array.from({ length: 28 }, (_, i) => ({
    emoji: SAISON_EMOJIS[i % SAISON_EMOJIS.length],
    left: Math.random() * 100,
    dauer: 5 + Math.random() * 6,
    delay: Math.random() * 6,
    size: 14 + Math.random() * 16
  }));
</script>

{#if $saison}
  <div class="saison-overlay" aria-hidden="true">
    {#each teilchen as p}
      <span
        class="saison-teil"
        style="left:{p.left}%; animation-duration:{p.dauer}s; animation-delay:{p.delay}s; font-size:{p.size}px;"
      >{p.emoji}</span>
    {/each}
  </div>
{/if}

<style>
  .saison-overlay { position: fixed; inset: 0; pointer-events: none; z-index: 999989; overflow: hidden; }
  .saison-teil { position: absolute; top: -30px; animation-name: saison-fall; animation-timing-function: linear; animation-iteration-count: infinite; }
  @keyframes saison-fall {
    0% { transform: translateY(-30px) rotate(0deg); opacity: 0.9; }
    100% { transform: translateY(105vh) rotate(360deg); opacity: 0.9; }
  }
</style>
