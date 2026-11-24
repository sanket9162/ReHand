package repository

import (
	"context"

	"github.com/sanket9162/rehand/internal/models"
)

type DatabaseRepo interface {
	Listings(ctx context.Context) ([]*models.Listing, error)
}
