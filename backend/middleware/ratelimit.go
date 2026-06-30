package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 🛡️ Einfacher Rate-Limiter pro IP (Schutz gegen Brute-Force / Spam).
// Merkt sich pro IP die Zeitpunkte der letzten Anfragen und erlaubt höchstens
// "max" Anfragen innerhalb von "fenster". Alles im Arbeitsspeicher – reicht für
// ein einzelnes Backend (kein externes Redis nötig).

type ratenZaehler struct {
	mu       sync.Mutex
	zugriffe map[string][]time.Time
	max      int
	fenster  time.Duration
}

func neuerZaehler(max int, fenster time.Duration) *ratenZaehler {
	return &ratenZaehler{
		zugriffe: make(map[string][]time.Time),
		max:      max,
		fenster:  fenster,
	}
}

// erlaubt prüft, ob die IP gerade noch eine Anfrage machen darf.
func (z *ratenZaehler) erlaubt(ip string) bool {
	z.mu.Lock()
	defer z.mu.Unlock()

	jetzt := time.Now()
	grenze := jetzt.Add(-z.fenster)

	// Nur die Zeitpunkte behalten, die noch im Zeitfenster liegen.
	var aktuell []time.Time
	for _, t := range z.zugriffe[ip] {
		if t.After(grenze) {
			aktuell = append(aktuell, t)
		}
	}

	if len(aktuell) >= z.max {
		z.zugriffe[ip] = aktuell
		return false
	}

	z.zugriffe[ip] = append(aktuell, jetzt)
	return true
}

// 🛡️ RateLimit gibt eine Middleware zurück, die pro IP nur "max" Anfragen je
// "fenster" erlaubt. Danach kommt 429 (Too Many Requests).
func RateLimit(max int, fenster time.Duration) gin.HandlerFunc {
	z := neuerZaehler(max, fenster)
	return func(c *gin.Context) {
		if !z.erlaubt(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"fehler": "Zu viele Anfragen. Bitte kurz warten und erneut versuchen.",
			})
			return
		}
		c.Next()
	}
}
