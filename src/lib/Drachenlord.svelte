<script>
  import { drachenAktiv } from '$lib/stores/easteregg.js';
  import { konfetti } from '$lib/confetti.js';
  // Wir zeigen hier ein verstecktes Bild aus dem Projekt-Ordner (als URL geladen).
  import drachenBild from '../../rw-26-cd.gif';

  // Phasen: 'idle' (nichts) -> 'fly' (Bild fliegt + Konfetti) -> 'bye' (Tschüss)
  let phase = $state('idle');

  // Sobald der Store auslöst, starten wir die Show (nur einmal).
  $effect(() => {
    if ($drachenAktiv && phase === 'idle') {
      starte();
    }
  });

  function starte() {
    phase = 'fly';
    // Konfetti-Regen aus Drachen & Party-Emojis
    konfetti({ anzahl: 140, dauer: 4000, emojis: ['🐉', '🎉', '🎊', '✨', '🔥'] });

    // Nach dem Überflug: Tschüss-Bildschirm
    setTimeout(() => (phase = 'bye'), 3000);

    // Danach versuchen wir, die Seite zu schließen.
    // Hinweis: Browser erlauben window.close() oft nur bei selbst geöffneten Tabs;
    // klappt das nicht, bleibt der Tschüss-Bildschirm einfach stehen.
    setTimeout(() => {
      try {
        window.close();
      } catch {
        /* ignorieren */
      }
    }, 4200);
  }
</script>

{#if phase === 'fly'}
  <div class="drachen-overlay">
    <img src={drachenBild} alt="" class="drachen-bild" />
  </div>
{:else if phase === 'bye'}
  <div class="drachen-overlay bye">
    <div class="drachen-bye">👋 Tschüss!</div>
  </div>
{/if}

<style>
  .drachen-overlay {
    position: fixed;
    inset: 0;
    z-index: 999997;
    pointer-events: none;
    overflow: hidden;
  }
  /* Das Bild fliegt von links unten nach rechts oben über den Bildschirm. */
  .drachen-bild {
    position: absolute;
    bottom: 10%;
    left: -300px;
    width: 240px;
    height: auto;
    animation: drachen-flug 3s ease-in-out forwards;
  }
  @keyframes drachen-flug {
    0% { left: -300px; bottom: 5%; transform: rotate(-10deg) scale(0.8); }
    50% { bottom: 55%; transform: rotate(10deg) scale(1.2); }
    100% { left: 110%; bottom: 20%; transform: rotate(-5deg) scale(0.9); }
  }
  /* Tschüss-Bildschirm: alles schwarz */
  .drachen-overlay.bye {
    background: #000;
    display: flex;
    align-items: center;
    justify-content: center;
    pointer-events: all;
  }
  .drachen-bye {
    color: #fff;
    font-size: 3rem;
    font-weight: bold;
    font-family: sans-serif;
  }
</style>
