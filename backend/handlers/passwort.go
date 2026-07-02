package handlers

import "unicode"

// 🔐 Passwort-Anforderungen (zentral, damit Backend & Frontend dieselben Regeln
// kennen). Wird bei Registrierung UND Passwort-Reset geprüft.

// passwortAnforderungen prüft ein Passwort und gibt eine Liste der NICHT
// erfüllten Anforderungen zurück (leer = alles ok). Die Texte sind so
// formuliert, dass man dem Nutzer direkt anzeigen kann, was noch fehlt.
func passwortAnforderungen(pw string) []string {
	var gross, klein, zahl, sonder bool
	for _, r := range pw {
		switch {
		case unicode.IsUpper(r):
			gross = true
		case unicode.IsLower(r):
			klein = true
		case unicode.IsDigit(r):
			zahl = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			sonder = true
		}
	}

	var fehlt []string
	if len([]rune(pw)) < 10 {
		fehlt = append(fehlt, "Mindestens 10 Zeichen")
	}
	if !gross {
		fehlt = append(fehlt, "Mindestens ein Großbuchstabe (A–Z)")
	}
	if !klein {
		fehlt = append(fehlt, "Mindestens ein Kleinbuchstabe (a–z)")
	}
	if !zahl {
		fehlt = append(fehlt, "Mindestens eine Zahl (0–9)")
	}
	if !sonder {
		fehlt = append(fehlt, "Mindestens ein Sonderzeichen (z. B. ! ? @ #)")
	}
	return fehlt
}
