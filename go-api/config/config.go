package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config agrupa la configuración del servicio leída del entorno.
type Config struct {
	ServerPort         string
	JWTSecret          string
	AuthUser           string
	AuthPass           string
	NodeAPIURL         string
	TokenExpiryMinutes int
	CORSOrigins        string
}

// Load construye la Config desde .env (si existe) y las variables del entorno.
// Si falta una variable requerida, aborta el arranque.
func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No se cargó .env, usando variables del entorno...")
	}

	var missing []string

	cfg := &Config{
		// Requerida.
		JWTSecret: mustEnv("JWT_SECRET", &missing),

		// Opcionales, con valor por defecto.
		NodeAPIURL:         getEnv("NODE_API_URL", "http://localhost:4000"),
		ServerPort:         getEnv("SERVER_PORT", "3000"),
		AuthUser:           getEnv("AUTH_USER", "admin"),
		AuthPass:           getEnv("AUTH_PASS", "admin123"),
		TokenExpiryMinutes: getEnvInt("TOKEN_EXPIRY_MINUTES", 60),
		CORSOrigins:        getEnv("CORS_ORIGINS", "*"),
	}

	if len(missing) > 0 {
		log.Fatalf("Faltan variables de entorno requeridas: %s", strings.Join(missing, ", "))
	}

	return cfg
}

// TokenExpiry devuelve la expiración del token como time.Duration.
func (c *Config) TokenExpiry() time.Duration {
	return time.Duration(c.TokenExpiryMinutes) * time.Minute
}

// mustEnv lee una variable requerida; si está vacía, la registra en missing.
func mustEnv(key string, missing *[]string) string {
	v := os.Getenv(key)
	if v == "" {
		*missing = append(*missing, key)
	}
	return v
}

// getEnv devuelve la variable de entorno o un valor por defecto.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvInt devuelve la variable de entorno como entero o un valor por defecto.
func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Printf("La variable %s='%s' no es un entero válido; usando %d", key, v, fallback)
	}
	return fallback
}
