package config

import (
	"log"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	cfg.Domain = os.Getenv("DOMAIN")
	cfg.Port = os.Getenv("PORT")
	cfg.DSN = os.Getenv("DSN")

	return &cfg
}
