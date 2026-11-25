package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
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
		SELECT id, title, description, price, city, created_at, updated_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT 100
	`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

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
			&l.UpdatedAt,
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

func (p *PostgresDBRepo) DeleteListing(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM listings WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	r, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if r == 0 {
		return errors.New("no listing found with this id")
	}

	return nil
}
