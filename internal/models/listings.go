package models

import (
	"time"
)

type Listings struct {
	ID          int       `josn:"id"`
	Title       string    `josn:"title"`
	Description string    `josn:"description"`
	Price       float64   `josn:"price"`
	City        string    `json:"city"`
	UpdatedAt   time.Time `josn:"updated_at"`
	DeletedAt   time.Time `josn:"deleted_at"`
}
