package dbrepo

import (
	"database/sql"

	"github.com/sanket9162/rehand/internal/config"
	"github.com/sanket9162/rehand/internal/repository"
)

type PostgresDBRepo struct {
	App *config.Config
	DB  *sql.DB
}

func NewPostgresDBRepo(conn *sql.DB, app *config.Config) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App: app,
		DB:  conn,
	}
}
