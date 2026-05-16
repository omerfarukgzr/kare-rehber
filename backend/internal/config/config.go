package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	HTTPPort    string
	DatabaseURL string

	JWTSecret    string
	JWTExpiresIn time.Duration

	BcryptCost int

	CORSAllowedOrigins string

	// Otomatik migrate / seed (Railway gibi container ortamları için)
	AutoMigrate    bool
	SeedAdmin      bool // varsa atlanır
	SeedTestUsers  bool // canlı'da kapat (false default)
}

const defaultDatabaseURL = "postgres://kareuser:karepass@localhost:5432/karerehber?sslmode=disable"

func Load() (*Config, error) {
	_ = godotenv.Load()

	// Railway PORT env'i (kendi atadığı dinamik port). Yoksa HTTP_PORT,
	// onun da yoksa 8080 fallback.
	httpPort := getenv("PORT", "")
	if httpPort == "" {
		httpPort = getenv("HTTP_PORT", "8080")
	}

	rawDB := getenv("DATABASE_URL", defaultDatabaseURL)
	usingDefaultDB := rawDB == defaultDatabaseURL

	cfg := &Config{
		AppEnv:             getenv("APP_ENV", "development"),
		HTTPPort:           httpPort,
		DatabaseURL:        normalizeDatabaseURL(rawDB),
		JWTSecret:          getenv("JWT_SECRET", "change-me-in-production"),
		BcryptCost:         getenvInt("BCRYPT_COST", 12),
		CORSAllowedOrigins: getenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
		AutoMigrate:        getenvBool("AUTO_MIGRATE", true),
		SeedAdmin:          getenvBool("SEED_ADMIN", true),
		SeedTestUsers:      getenvBool("SEED_TEST_USERS", false),
	}

	jwtHours := getenvInt("JWT_EXPIRES_HOURS", 24)
	cfg.JWTExpiresIn = time.Duration(jwtHours) * time.Hour

	if cfg.AppEnv == "production" {
		if cfg.JWTSecret == "change-me-in-production" {
			return nil, fmt.Errorf("JWT_SECRET must be set in production")
		}
		if usingDefaultDB {
			return nil, fmt.Errorf("DATABASE_URL is not set in production " +
				"(binary fell back to the localhost default). " +
				"Add DATABASE_URL as a service variable on Railway and reference your Postgres service.")
		}
	}

	return cfg, nil
}

// RedactedDatabaseURL: log'lara basmak için güvenli (şifresiz) versiyon döner.
func (c *Config) RedactedDatabaseURL() string {
	u, err := url.Parse(c.DatabaseURL)
	if err != nil {
		return "(unparseable)"
	}
	if u.User != nil {
		u.User = url.UserPassword(u.User.Username(), "***")
	}
	return u.Redacted()
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

// normalizeDatabaseURL: Railway internal URL'leri sslmode olmadan geliyor.
// golang-migrate'in lib/pq sürücüsü default sslmode=require ile başarısız olur.
// sslmode parametresi yoksa disable ekleriz; kullanıcı kendi de set edebilir.
func normalizeDatabaseURL(u string) string {
	if u == "" {
		return u
	}
	if strings.Contains(u, "sslmode=") {
		return u
	}
	sep := "?"
	if strings.Contains(u, "?") {
		sep = "&"
	}
	return u + sep + "sslmode=disable"
}

func getenvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "1", "true", "yes", "on":
			return true
		case "0", "false", "no", "off":
			return false
		}
	}
	return fallback
}
