package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/sanket9162/rehand/internal/models"
)

type DatabaseRepo interface {
	Listings(ctx context.Context) ([]*models.Listing, error)
	DeleteListing(ctx context.Context, id uuid.UUID) error
}
