package model

import (
	"FurniSwap/internal/modules/listing/model"
	"time"
)

// Favorite represents a favorite listing
type Favorite struct {
	ID        int            `db:"id" json:"id"`
	UserID    int            `db:"user_id" json:"user_id"`
	ListingID int            `db:"listing_id" json:"listing_id"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	Listing   *model.Listing `json:"listing,omitempty"`
}

// FavoriteResponse represents a list of favorites with pagination
type FavoriteResponse struct {
	Favorites   []Favorite `json:"favorites"`
	TotalCount  int        `json:"total_count"`
	CurrentPage int        `json:"current_page"`
	TotalPages  int        `json:"total_pages"`
}
