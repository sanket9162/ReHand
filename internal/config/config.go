package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
}

type Config struct {
	Domain string
	Port   string
	DSN    string
	Auth   Auth
}

func LoadConfig() *Config {
	var cfg Config

	_ = godotenv.Load()

	cfg.Domain = getEnv("DOMAIN", "ReHand")
	cfg.Port = getEnv("PORT", "8080")
	cfg.DSN = os.Getenv("DSN")

	cfg.Auth = Auth{
		Issuer:        getEnv("JWT_ISSUER", "rehand"),
		Audience:      getEnv("JWT_AUDIENCE", "rehand"),
		Secret:        os.Getenv("JWT_SECRET"),
		TokenExpiry:   1 * time.Hour,
		RefreshExpiry: 7 * 24 * time.Hour,
		CookieDomain:  getEnv("COOKIE_DOMAIN", "localhost"),
		CookiePath:    getEnv("COOKIE_PATH", "/"),
	}

	return &cfg
}

func (c *Config) ValidateAPI() error {
	if c.DSN == "" {
		return errors.New("DSN environment variable is required")
	}
	if c.Auth.Secret == "" {
		return errors.New("JWT_SECRET environment variable is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
