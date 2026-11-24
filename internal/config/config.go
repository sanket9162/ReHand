package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Domain string
	Port   string
	DSN    string
}

func LoadConfig() *Config {
	var cfg Config

	_ = godotenv.Load()

	cfg.Domain = os.Getenv("DOMAIN")
	cfg.Port = os.Getenv("PORT")
	cfg.DSN = os.Getenv("DSN")

	return &cfg
}
