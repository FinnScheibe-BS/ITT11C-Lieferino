package config

import "os"

// ⚙️ Config liest alle Einstellungen aus Umgebungsvariablen (mit Standardwerten).
// So kann man im Docker-Container alles über die docker-compose.yml setzen.
type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPasswort string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
	AdminEmail string
}

func Laden() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "lieferino"),
		DBPasswort: getEnv("DB_PASSWORD", "lieferino"),
		DBName:     getEnv("DB_NAME", "lieferino"),
		// docker-compose: "disable". Bei Patroni/Spilo (K8s) z.B. "require".
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-bitte-aendern"),
		// Wer sich mit dieser E-Mail registriert/anmeldet, wird Admin.
		AdminEmail: getEnv("ADMIN_EMAIL", "admin@lieferino.de"),
	}
}

// Liest eine Umgebungsvariable, oder gibt den Standardwert zurück.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
