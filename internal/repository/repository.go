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
)

type DatabaseRepo interface {
	Listings(ctx context.Context) ([]*models.Listing, error)
	DeleteListing(ctx context.Context, id uuid.UUID) error
	CreateListing(ctx context.Context, l *models.Listing) (*models.Listing, error)
	CreateUser(ctx context.Context, u *models.User) (*models.User, error)
}
