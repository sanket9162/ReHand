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
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `josn:"updated_at"`
}
