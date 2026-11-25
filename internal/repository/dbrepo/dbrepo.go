package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sanket9162/rehand/internal/models"
	"github.com/sanket9162/rehand/internal/repository"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

func NewPostgresDBRepo(conn *sql.DB) repository.DatabaseRepo {
	return &PostgresDBRepo{
		DB: conn,
	}
}

const dbTimeout = 3 * time.Second

func (p *PostgresDBRepo) Listings(ctx context.Context) ([]*models.Listing, error) {
	query := `
		SELECT id, title, description, price, city, created_at, updated_at
		FROM listings
		ORDER BY created_at DESC
		LIMIT 100
	`

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query listings: %w", err)
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
			return nil, fmt.Errorf("scan listing: %w", err)
		}
		list = append(list, &l)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listings: %w", err)
	}

	return list, nil
}

func (p *PostgresDBRepo) DeleteListing(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM listings WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	result, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete listing: %w", err)
	}

	r, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete listing rows affected: %w", err)
	}

	if r == 0 {
		return repository.ErrNotFound
	}

	return nil
}
