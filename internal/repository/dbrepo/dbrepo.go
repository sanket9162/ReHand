package dbrepo

import (
	"context"
	"database/sql"

	"github.com/sanket9162/rehand/internal/config"
	"github.com/sanket9162/rehand/internal/models"
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

func (p *PostgresDBRepo) Listings(ctx context.Context) ([]*models.Listing, error) {
	query := `
		SELECT id, title, description, price, city, created_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT 100
	`

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*models.Listing, 0)
	for rows.Next() {
		var l models.Listing
		err := rows.Scan(
			&l.ID,
			&l.Title,
			&l.Description,
			&l.Price,
			&l.City,
			&l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &l)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
