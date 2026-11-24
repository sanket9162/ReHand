package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sanket9162/rehand/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate < up | down>")
	}

	cfg := config.LoadConfig()

	m, err := migrate.New(
		"file://migrations",
		cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "up":
		log.Println("up called")
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	case "down":
		log.Println("down called")
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("invalid command %v", os.Args[1])
	}

	fmt.Println("running migrate")

}
