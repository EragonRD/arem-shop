package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// AppConfig contient toute la configuration runtime de l'API.
type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string

	CORSAllowedOrigins []string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimezone string

	JWTSecret   string
	JWTTTLHours int

	BcryptCost int

	LowStockThreshold int
}

// Load charge la configuration depuis l'environnement (.env facultatif).
func Load() (AppConfig, error) {
	_ = godotenv.Load()

	cfg := AppConfig{
		AppName: getenvDefault("APP_NAME", "arem-shop"),
		AppEnv:  getenvDefault("APP_ENV", "development"),
		AppPort: getenvDefault("APP_PORT", "8080"),
		CORSAllowedOrigins: splitCSV(
			getenvDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		),

		DBHost:     getenvDefault("DB_HOST", "localhost"),
		DBPort:     getenvDefault("DB_PORT", "5432"),
		DBUser:     getenvDefault("DB_USER", "postgres"),
		DBPassword: getenvDefault("DB_PASSWORD", "postgres"),
		DBName:     getenvDefault("DB_NAME", "arem_shop"),
		DBSSLMode:  getenvDefault("DB_SSLMODE", "disable"),
		DBTimezone: getenvDefault("DB_TIMEZONE", "UTC"),

		JWTSecret: getenvDefault("JWT_SECRET", "change-me-super-secret"),
	}

	jwtTTL, err := getenvIntDefault("JWT_TTL_HOURS", 24)
	if err != nil {
		return AppConfig{}, fmt.Errorf("invalid JWT_TTL_HOURS: %w", err)
	}
	cfg.JWTTTLHours = jwtTTL

	bcryptCost, err := getenvIntDefault("BCRYPT_COST", 12)
	if err != nil {
		return AppConfig{}, fmt.Errorf("invalid BCRYPT_COST: %w", err)
	}
	cfg.BcryptCost = bcryptCost

	lowStockThreshold, err := getenvIntDefault("LOW_STOCK_THRESHOLD", 5)
	if err != nil {
		return AppConfig{}, fmt.Errorf("invalid LOW_STOCK_THRESHOLD: %w", err)
	}
	if lowStockThreshold < 0 {
		return AppConfig{}, fmt.Errorf("LOW_STOCK_THRESHOLD must be >= 0")
	}
	cfg.LowStockThreshold = lowStockThreshold

	if cfg.JWTSecret == "" {
		return AppConfig{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getenvDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getenvIntDefault(key string, fallback int) (int, error) {
	value := getenvDefault(key, strconv.Itoa(fallback))
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	entries := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			entries = append(entries, trimmed)
		}
	}

	if len(entries) == 0 {
		return []string{"http://localhost:3000"}
	}

	return entries
}
