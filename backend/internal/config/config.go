package config

import (
	"fmt"
	"os"
	"strconv"
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
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:             getenv("APP_ENV", "development"),
		HTTPPort:           getenv("HTTP_PORT", "8080"),
		DatabaseURL:        getenv("DATABASE_URL", "postgres://kareuser:karepass@localhost:5432/karerehber?sslmode=disable"),
		JWTSecret:          getenv("JWT_SECRET", "change-me-in-production"),
		BcryptCost:         getenvInt("BCRYPT_COST", 12),
		CORSAllowedOrigins: getenv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
	}

	jwtHours := getenvInt("JWT_EXPIRES_HOURS", 24)
	cfg.JWTExpiresIn = time.Duration(jwtHours) * time.Hour

	if cfg.JWTSecret == "change-me-in-production" && cfg.AppEnv == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return cfg, nil
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
