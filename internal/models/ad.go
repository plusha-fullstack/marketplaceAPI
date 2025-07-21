package models

import (
	"time"
)

type Ad struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	AuthorID    int       `json:"author_id"`
	AuthorLogin string    `json:"author_login"`
	IsOwner     bool      `json:"is_owner,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateAdRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
}

type AdFilters struct {
	Page      int     `json:"page"`
	Limit     int     `json:"limit"`
	SortBy    string  `json:"sort_by"`    // "date" or "price"
	SortOrder string  `json:"sort_order"` // "asc" or "desc"
	MinPrice  float64 `json:"min_price"`
	MaxPrice  float64 `json:"max_price"`
}
