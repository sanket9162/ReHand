package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sanket9162/rehand/internal/models"
)

var (
	ErrNotFound       = errors.New("record not found")
	ErrDuplicateEmail = errors.New("email already exists")
	ErrForbidden      = errors.New("forbidden")
)

type DatabaseRepo interface {
	Ping(ctx context.Context) error

	Listings(ctx context.Context) ([]*models.Listing, error)
	GetListingByID(ctx context.Context, id uuid.UUID) (*models.Listing, error)
	CreateListing(ctx context.Context, l *models.Listing) (*models.Listing, error)
	UpdateListing(ctx context.Context, l *models.Listing) (*models.Listing, error)
	DeleteListing(ctx context.Context, id, userID uuid.UUID) error

	CreateUser(ctx context.Context, u *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}
