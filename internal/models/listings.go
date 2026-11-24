package models

import (
	"time"
)

type listings struct {
	id          int       `josn:"id"`
	title       string    `josn:"title"`
	description string    `josn:"description"`
	price       float64   `josn:"price"`
	city        string    `json:"city"`
	updatedAt   time.Time `josn:"updated_at"`
	deletedAt   time.Time `josn:"deleted_at"`
}
