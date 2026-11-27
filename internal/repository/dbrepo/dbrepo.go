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

func (p *PostgresDBRepo) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	return p.DB.PingContext(ctx)
}

func (p *PostgresDBRepo) Listings(ctx context.Context) ([]*models.Listing, error) {
	query := `
		SELECT id, user_id, title, description, price, city, created_at, updated_at
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
			&l.UserID,
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

func (p *PostgresDBRepo) GetListingByID(ctx context.Context, id uuid.UUID) (*models.Listing, error) {
	query := `
		SELECT id, user_id, title, description, price, city, created_at, updated_at
		FROM listings
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var l models.Listing
	err := p.DB.QueryRowContext(ctx, query, id).Scan(
		&l.ID,
		&l.UserID,
		&l.Title,
		&l.Description,
		&l.Price,
		&l.City,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get listing by id: %w", err)
	}

	return &l, nil
}

func (p *PostgresDBRepo) CreateListing(ctx context.Context, l *models.Listing) (*models.Listing, error) {
	query := `
		INSERT INTO listings (user_id, title, description, price, city)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, title, description, price, city, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var created models.Listing
	err := p.DB.QueryRowContext(
		ctx,
		query,
		l.UserID,
		l.Title,
		l.Description,
		l.Price,
		l.City,
	).Scan(
		&created.ID,
		&created.UserID,
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

func (p *PostgresDBRepo) UpdateListing(ctx context.Context, l *models.Listing) (*models.Listing, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var ownerID uuid.UUID
	if err := p.DB.QueryRowContext(ctx, `SELECT user_id FROM listings WHERE id = $1`, l.ID).Scan(&ownerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("check listing owner: %w", err)
	}
	if ownerID != l.UserID {
		return nil, repository.ErrForbidden
	}

	query := `
		UPDATE listings
		SET title = $1, description = $2, price = $3, city = $4, updated_at = NOW()
		WHERE id = $5 AND user_id = $6
		RETURNING id, user_id, title, description, price, city, created_at, updated_at
	`

	var updated models.Listing
	err := p.DB.QueryRowContext(ctx, query, l.Title, l.Description, l.Price, l.City, l.ID, l.UserID).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Title,
		&updated.Description,
		&updated.Price,
		&updated.City,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("update listing: %w", err)
	}

	return &updated, nil
}

func (p *PostgresDBRepo) DeleteListing(ctx context.Context, id, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var ownerID uuid.UUID
	if err := p.DB.QueryRowContext(ctx, `SELECT user_id FROM listings WHERE id = $1`, id).Scan(&ownerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("check listing owner: %w", err)
	}
	if ownerID != userID {
		return repository.ErrForbidden
	}

	result, err := p.DB.ExecContext(ctx, `DELETE FROM listings WHERE id = $1 AND user_id = $2`, id, userID)
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

func (p *PostgresDBRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var u models.User
	err := p.DB.QueryRowContext(ctx, query, id).Scan(
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
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &u, nil
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
