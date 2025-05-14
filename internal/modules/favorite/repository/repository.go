package repository

import (
	"FurniSwap/internal/modules/favorite/model"
	listingModel "FurniSwap/internal/modules/listing/model"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for the favorite module
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new favorite repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// AddFavorite adds a listing to a user's favorites
func (r *Repository) AddFavorite(userID, listingID int) error {
	// Check if the listing exists
	var exists bool
	err := r.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM listings WHERE id = $1)", listingID)
	if err != nil {
		log.Printf("Error checking if listing exists: %v", err)
		return fmt.Errorf("error checking if listing exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("listing does not exist")
	}

	// Check if already in favorites
	err = r.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND listing_id = $2)", userID, listingID)
	if err != nil {
		log.Printf("Error checking if already in favorites: %v", err)
		return fmt.Errorf("error checking if already in favorites: %w", err)
	}

	if exists {
		return nil // Already in favorites, just return success
	}

	// Add to favorites
	_, err = r.db.Exec(`
		INSERT INTO favorites (user_id, listing_id, created_at)
		VALUES ($1, $2, $3)
	`, userID, listingID, time.Now())
	if err != nil {
		log.Printf("Error adding to favorites: %v", err)
		return fmt.Errorf("error adding to favorites: %w", err)
	}

	return nil
}

// RemoveFavorite removes a listing from a user's favorites
func (r *Repository) RemoveFavorite(userID, listingID int) error {
	_, err := r.db.Exec("DELETE FROM favorites WHERE user_id = $1 AND listing_id = $2", userID, listingID)
	if err != nil {
		log.Printf("Error removing from favorites: %v", err)
		return fmt.Errorf("error removing from favorites: %w", err)
	}

	return nil
}

// IsFavorite checks if a listing is in a user's favorites
func (r *Repository) IsFavorite(userID, listingID int) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND listing_id = $2)", userID, listingID)
	if err != nil {
		log.Printf("Error checking if listing is in favorites: %v", err)
		return false, fmt.Errorf("error checking if listing is in favorites: %w", err)
	}

	return exists, nil
}

// GetFavorites gets a user's favorite listings with pagination
func (r *Repository) GetFavorites(userID, page, limit int) (*model.FavoriteResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get total count
	var totalCount int
	err := r.db.Get(&totalCount, "SELECT COUNT(*) FROM favorites WHERE user_id = $1", userID)
	if err != nil {
		log.Printf("Error getting favorites count: %v", err)
		return nil, fmt.Errorf("error getting favorites count: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Get favorites
	var favorites []model.Favorite
	err = r.db.Select(&favorites, `
		SELECT f.id, f.user_id, f.listing_id, f.created_at
		FROM favorites f
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		log.Printf("Error getting favorites: %v", err)
		return nil, fmt.Errorf("error getting favorites: %w", err)
	}

	// Get listings for each favorite
	for i := range favorites {
		var listing listingModel.Listing
		err = r.db.Get(&listing, `
			SELECT l.*, u.name as user_name
			FROM listings l
			JOIN users u ON l.user_id = u.id
			WHERE l.id = $1
		`, favorites[i].ListingID)
		if err != nil {
			log.Printf("Error getting listing for favorite %d: %v", favorites[i].ID, err)
			continue // Skip this listing if error
		}

		// Get images for the listing
		err = r.db.Select(&listing.Images, "SELECT * FROM listing_images WHERE listing_id = $1 ORDER BY is_main DESC, created_at ASC", listing.ID)
		if err != nil {
			log.Printf("Error getting images for listing %d: %v", listing.ID, err)
			// Continue without images if there's an error
		}

		favorites[i].Listing = &listing
	}

	return &model.FavoriteResponse{
		Favorites:   favorites,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}
