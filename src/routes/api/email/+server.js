import { json } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

// =====================================================================
// 📧 E-MAIL-VERSAND (Server-Route)
// Diese Route läuft auf dem SERVER (nicht im Browser) und verschickt echte
// E-Mails über das Gmail-Konto. Sie ist bereits fertig – sie wird automatisch
// aktiv, SOBALD die Zugangsdaten in der .env-Datei stehen (siehe EMAIL-SETUP.md).
//
// Ohne Zugangsdaten läuft alles im "Test-Modus": es wird nichts verschickt,
// die App zeigt den Code stattdessen direkt an.
// =====================================================================

export async function POST({ request }) {
  const { an, betreff, text } = await request.json();

  // Keine Zugangsdaten gesetzt? -> Test-Modus (es wird NICHTS verschickt).
  if (!env.GMAIL_USER || !env.GMAIL_APP_PASSWORD) {
    return json({ gesendet: false, testModus: true });
  }

  try {
    // Nodemailer erst hier laden (nur auf dem Server, nur wenn wirklich gesendet wird).
    const nodemailer = (await import('nodemailer')).default;
    const transporter = nodemailer.createTransport({
      service: 'gmail',
      auth: { user: env.GMAIL_USER, pass: env.GMAIL_APP_PASSWORD }
    });

    await transporter.sendMail({
      from: `Lieferino <${env.GMAIL_USER}>`,
      to: an,
      subject: betreff,
      text
    });

    return json({ gesendet: true });
  } catch (fehler) {
    console.error('E-Mail-Versand fehlgeschlagen:', fehler);
    return json({ gesendet: false, fehler: true }, { status: 500 });
  }
}
