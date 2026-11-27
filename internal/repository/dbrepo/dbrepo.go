package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

func (p *PostgresDBRepo) CreateListing(ctx context.Context, l *models.Listing) (*models.Listing, error) {
	query := `
		INSERT INTO listings (title, description, price, city)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, description, price, city, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var created models.Listing
	err := p.DB.QueryRowContext(
		ctx,
		query,
		l.Title,
		l.Description,
		l.Price,
		l.City,
	).Scan(
		&created.ID,
		&created.Title,
		&created.Description,
		&created.Price,
		&created.City,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create listing: %w", err)
	}

	return &created, nil
}

func (p *PostgresDBRepo) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var created models.User
	err := p.DB.QueryRowContext(
		ctx,
		query,
		u.Name,
		u.Email,
		u.Password,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Email,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, repository.ErrDuplicateEmail
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &created, nil
}

func (p *PostgresDBRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var u models.User
	err := p.DB.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil

}
