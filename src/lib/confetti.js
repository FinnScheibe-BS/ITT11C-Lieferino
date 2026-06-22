// =====================================================================
// 🎉 KONFETTI + MINI-TOAST (für Easter Eggs)
// Reine DOM-Helfer ohne Fremd-Bibliothek. Können von überall aufgerufen werden.
// =====================================================================

const FARBEN = ['#673ab7', '#ff3b30', '#34c759', '#ff9500', '#0a84ff', '#ffd60a', '#ff2d55'];

// Sorgt dafür, dass die Animations-Styles nur EINMAL in die Seite kommen.
function styleEinmalEinfuegen() {
  if (typeof document === 'undefined') return;
  if (document.getElementById('konfetti-style')) return;
  const style = document.createElement('style');
  style.id = 'konfetti-style';
  style.textContent = `
    .konfetti-container { position: fixed; inset: 0; pointer-events: none; z-index: 999998; overflow: hidden; }
    .konfetti-teil { position: absolute; top: -20px; width: 10px; height: 14px; border-radius: 2px; animation: konfetti-fall linear forwards; }
    .konfetti-teil.emoji { width: auto; height: auto; font-size: 22px; background: none !important; }
    @keyframes konfetti-fall {
      0% { transform: translateY(-20px) rotate(0deg); opacity: 1; }
      100% { transform: translateY(105vh) rotate(720deg); opacity: 1; }
    }
    .eier-toast { position: fixed; top: 24px; left: 50%; transform: translateX(-50%); background: #673ab7; color: #fff;
      padding: 14px 22px; border-radius: 30px; font-weight: bold; box-shadow: 0 8px 24px rgba(0,0,0,0.25);
      z-index: 999999; font-family: sans-serif; text-align: center; max-width: 90%; animation: eier-pop 0.3s ease; }
    @keyframes eier-pop { from { transform: translateX(-50%) scale(0.7); opacity: 0; } to { transform: translateX(-50%) scale(1); opacity: 1; } }
  `;
  document.head.appendChild(style);
}

// 🎉 Lässt Konfetti regnen. Optional mit Emojis statt bunten Schnipseln.
export function konfetti({ anzahl = 80, dauer = 2500, emojis = null } = {}) {
  if (typeof document === 'undefined') return;
  styleEinmalEinfuegen();

  const container = document.createElement('div');
  container.className = 'konfetti-container';

  for (let i = 0; i < anzahl; i++) {
    const teil = document.createElement('div');
    teil.className = 'konfetti-teil';
    teil.style.left = Math.random() * 100 + '%';
    teil.style.animationDelay = Math.random() * 0.6 + 's';
    teil.style.animationDuration = (dauer / 1000) * (0.6 + Math.random() * 0.8) + 's';
    if (emojis) {
      teil.classList.add('emoji');
      teil.textContent = emojis[Math.floor(Math.random() * emojis.length)];
    } else {
      teil.style.background = FARBEN[i % FARBEN.length];
    }
    container.appendChild(teil);
  }

  document.body.appendChild(container);
  // Nach der Animation wieder aufräumen.
  setTimeout(() => container.remove(), dauer + 1500);
}

// 💬 Kurze, mittige Hinweis-Blase (z.B. für Geheimcodes).
export function eierToast(text, dauer = 4000) {
  if (typeof document === 'undefined') return;
  styleEinmalEinfuegen();
  const toast = document.createElement('div');
  toast.className = 'eier-toast';
  toast.textContent = text;
  document.body.appendChild(toast);
  setTimeout(() => toast.remove(), dauer);
}
