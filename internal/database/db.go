package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB holds the database connection pools
type DB struct {
	SQL *sql.DB
}

// ConnectPostgres opens a connection to the database
func ConnectPostgres(dsn string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{
		SQL: db,
	}, nil
}
