package repository

import (
	"FurniSwap/internal/modules/listing/model"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for the listing module
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new listing repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// CreateListing creates a new listing
func (r *Repository) CreateListing(userID int, req model.CreateListingRequest) (int, error) {
	var listingID int
	err := r.db.QueryRow(`
		INSERT INTO listings (user_id, title, description, price, condition, city, category_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, userID, req.Title, req.Description, req.Price, req.Condition, req.City, req.CategoryID, "active", time.Now(), time.Now()).Scan(&listingID)

	if err != nil {
		log.Printf("Error creating listing: %v", err)
		return 0, fmt.Errorf("error creating listing: %w", err)
	}

	return listingID, nil
}

// UpdateListing updates an existing listing
func (r *Repository) UpdateListing(listingID, userID int, req model.UpdateListingRequest) error {
	// First check if the listing belongs to the user
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM listings WHERE id = $1 AND user_id = $2", listingID, userID)
	if err != nil {
		log.Printf("Error checking listing ownership: %v", err)
		return fmt.Errorf("error checking listing ownership: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("listing not found or does not belong to the user")
	}

	// Get current listing data
	var current model.Listing
	err = r.db.Get(&current, "SELECT * FROM listings WHERE id = $1", listingID)
	if err != nil {
		log.Printf("Error getting current listing data: %v", err)
		return fmt.Errorf("error getting current listing data: %w", err)
	}

	// Update only provided fields
	title := current.Title
	if req.Title != "" {
		title = req.Title
	}

	description := current.Description
	if req.Description != "" {
		description = req.Description
	}

	price := current.Price
	if req.Price != 0 {
		price = req.Price
	}

	condition := current.Condition
	if req.Condition != "" {
		condition = req.Condition
	}

	city := current.City
	if req.City != "" {
		city = req.City
	}

	categoryID := current.CategoryID
	if req.CategoryID != 0 {
		categoryID = req.CategoryID
	}

	status := current.Status
	if req.Status != "" {
		status = req.Status
	}

	// Update the listing
	_, err = r.db.Exec(`
		UPDATE listings
		SET title = $1, description = $2, price = $3, condition = $4, city = $5, category_id = $6, status = $7, updated_at = $8
		WHERE id = $9
	`, title, description, price, condition, city, categoryID, status, time.Now(), listingID)

	if err != nil {
		log.Printf("Error updating listing: %v", err)
		return fmt.Errorf("error updating listing: %w", err)
	}

	return nil
}

// DeleteListing deletes a listing
func (r *Repository) DeleteListing(listingID, userID int) error {
	// First check if the listing belongs to the user
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM listings WHERE id = $1 AND user_id = $2", listingID, userID)
	if err != nil {
		log.Printf("Error checking listing ownership: %v", err)
		return fmt.Errorf("error checking listing ownership: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("listing not found or does not belong to the user")
	}

	// Delete the listing
	_, err = r.db.Exec("DELETE FROM listings WHERE id = $1", listingID)
	if err != nil {
		log.Printf("Error deleting listing: %v", err)
		return fmt.Errorf("error deleting listing: %w", err)
	}

	return nil
}

// GetListing gets a single listing by ID
func (r *Repository) GetListing(listingID int) (*model.Listing, error) {
	var listing model.Listing
	err := r.db.Get(&listing, `
		SELECT l.*, COALESCE(u.name, '') as user_name
		FROM listings l
		LEFT JOIN users u ON l.user_id = u.id
		WHERE l.id = $1
	`, listingID)

	if err != nil {
		log.Printf("Error getting listing: %v", err)
		return nil, fmt.Errorf("error getting listing: %w", err)
	}

	// Get images for the listing
	listing.Images = []model.Image{} // Initialize with empty slice to avoid null in JSON
	err = r.db.Select(&listing.Images, "SELECT * FROM listing_images WHERE listing_id = $1 ORDER BY is_main DESC, created_at ASC", listingID)
	if err != nil {
		log.Printf("Error getting listing images: %v", err)
		// Continue without images if there's an error
	}

	return &listing, nil
}

// GetListings gets listings with filtering and pagination
func (r *Repository) GetListings(filter model.ListingFilter) (*model.ListingResponse, error) {
	// Build the query with filters
	query := "SELECT l.*, COALESCE(u.name, '') as user_name FROM listings l LEFT JOIN users u ON l.user_id = u.id WHERE l.status = 'active'"
	countQuery := "SELECT COUNT(*) FROM listings l WHERE l.status = 'active'"
	var args []interface{}
	var countArgs []interface{}
	argIndex := 1

	// Apply category filter
	if filter.CategoryID != nil && *filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND l.category_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.category_id = $%d", argIndex)
		args = append(args, *filter.CategoryID)
		countArgs = append(countArgs, *filter.CategoryID)
		argIndex++
	}

	// Apply city filter
	if filter.City != "" {
		query += fmt.Sprintf(" AND l.city ILIKE $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.city ILIKE $%d", argIndex)
		args = append(args, "%"+filter.City+"%")
		countArgs = append(countArgs, "%"+filter.City+"%")
		argIndex++
	}

	// Apply condition filter
	if filter.Condition != "" {
		query += fmt.Sprintf(" AND l.condition = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.condition = $%d", argIndex)
		args = append(args, filter.Condition)
		countArgs = append(countArgs, filter.Condition)
		argIndex++
	}

	// Apply min price filter
	if filter.MinPrice != nil && *filter.MinPrice >= 0 {
		query += fmt.Sprintf(" AND l.price >= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.price >= $%d", argIndex)
		args = append(args, *filter.MinPrice)
		countArgs = append(countArgs, *filter.MinPrice)
		argIndex++
	}

	// Apply max price filter
	if filter.MaxPrice != nil && *filter.MaxPrice > 0 {
		query += fmt.Sprintf(" AND l.price <= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.price <= $%d", argIndex)
		args = append(args, *filter.MaxPrice)
		countArgs = append(countArgs, *filter.MaxPrice)
		argIndex++
	}

	// Apply sorting
	switch filter.SortBy {
	case "date":
		query += " ORDER BY l.created_at ASC"
	case "-date":
		query += " ORDER BY l.created_at DESC"
	case "price":
		query += " ORDER BY l.price ASC"
	case "-price":
		query += " ORDER BY l.price DESC"
	default:
		query += " ORDER BY l.created_at DESC" // Default sort by newest
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, offset)
	argIndex += 2

	// Get total count
	var totalCount int
	err := r.db.Get(&totalCount, countQuery, countArgs...)
	if err != nil {
		log.Printf("Error getting listings count: %v", err)
		return nil, fmt.Errorf("error getting listings count: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(filter.Limit)))

	// Get listings
	var listings []model.Listing
	err = r.db.Select(&listings, query, args...)
	if err != nil {
		log.Printf("Error getting listings: %v", err)
		return nil, fmt.Errorf("error getting listings: %w", err)
	}

	// Get images for each listing
	for i := range listings {
		listings[i].Images = []model.Image{} // Initialize with empty slice to avoid null in JSON
		err = r.db.Select(&listings[i].Images, "SELECT * FROM listing_images WHERE listing_id = $1 ORDER BY is_main DESC, created_at ASC", listings[i].ID)
		if err != nil {
			log.Printf("Error getting images for listing %d: %v", listings[i].ID, err)
			// Continue without images if there's an error
		}
	}

	return &model.ListingResponse{
		Listings:    listings,
		TotalCount:  totalCount,
		CurrentPage: filter.Page,
		TotalPages:  totalPages,
	}, nil
}

// AddImage adds an image to a listing
func (r *Repository) AddImage(listingID int, imagePath string) (int, error) {
	var imageID int
	err := r.db.QueryRow(`
		INSERT INTO listing_images (listing_id, image_path, is_main, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, listingID, imagePath, false, time.Now()).Scan(&imageID)

	if err != nil {
		log.Printf("Error adding image: %v", err)
		return 0, fmt.Errorf("error adding image: %w", err)
	}

	// If this is the first image, make it the main image
	var count int
	err = r.db.Get(&count, "SELECT COUNT(*) FROM listing_images WHERE listing_id = $1", listingID)
	if err != nil {
		log.Printf("Error counting images: %v", err)
		return imageID, nil // Return imageID anyway
	}

	if count == 1 {
		_, err = r.db.Exec("UPDATE listing_images SET is_main = true WHERE id = $1", imageID)
		if err != nil {
			log.Printf("Error setting main image: %v", err)
			return imageID, nil // Return imageID anyway
		}
	}

	return imageID, nil
}

// DeleteImage deletes an image from a listing
func (r *Repository) DeleteImage(imageID, listingID, userID int) error {
	// First check if the listing belongs to the user
	var count int
	err := r.db.Get(&count, `
		SELECT COUNT(*) FROM listing_images li
		JOIN listings l ON li.listing_id = l.id
		WHERE li.id = $1 AND li.listing_id = $2 AND l.user_id = $3
	`, imageID, listingID, userID)
	if err != nil {
		log.Printf("Error checking image ownership: %v", err)
		return fmt.Errorf("error checking image ownership: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("image not found or does not belong to the user's listing")
	}

	// Get image path for deletion
	var imagePath string
	var isMain bool
	err = r.db.QueryRow("SELECT image_path, is_main FROM listing_images WHERE id = $1", imageID).Scan(&imagePath, &isMain)
	if err != nil {
		log.Printf("Error getting image path: %v", err)
		return fmt.Errorf("error getting image path: %w", err)
	}

	// Delete the image
	_, err = r.db.Exec("DELETE FROM listing_images WHERE id = $1", imageID)
	if err != nil {
		log.Printf("Error deleting image: %v", err)
		return fmt.Errorf("error deleting image: %w", err)
	}

	// If the deleted image was the main image, set another image as main
	if isMain {
		var newMainImageID int
		err = r.db.QueryRow("SELECT id FROM listing_images WHERE listing_id = $1 ORDER BY created_at ASC LIMIT 1", listingID).Scan(&newMainImageID)
		if err == nil {
			_, err = r.db.Exec("UPDATE listing_images SET is_main = true WHERE id = $1", newMainImageID)
			if err != nil {
				log.Printf("Error setting new main image: %v", err)
				// Not a critical error, continue
			}
		}
	}

	return nil
}

// SetMainImage sets an image as the main image for a listing
func (r *Repository) SetMainImage(imageID, listingID, userID int) error {
	// First check if the listing belongs to the user
	var count int
	err := r.db.Get(&count, `
		SELECT COUNT(*) FROM listing_images li
		JOIN listings l ON li.listing_id = l.id
		WHERE li.id = $1 AND li.listing_id = $2 AND l.user_id = $3
	`, imageID, listingID, userID)
	if err != nil {
		log.Printf("Error checking image ownership: %v", err)
		return fmt.Errorf("error checking image ownership: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("image not found or does not belong to the user's listing")
	}

	// Begin transaction
	tx, err := r.db.Beginx()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	// Reset all images for this listing to not be main
	_, err = tx.Exec("UPDATE listing_images SET is_main = false WHERE listing_id = $1", listingID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error resetting main images: %v", err)
		return fmt.Errorf("error resetting main images: %w", err)
	}

	// Set the selected image as main
	_, err = tx.Exec("UPDATE listing_images SET is_main = true WHERE id = $1", imageID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error setting main image: %v", err)
		return fmt.Errorf("error setting main image: %w", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// GetUserListings gets all listings for a user
func (r *Repository) GetUserListings(userID int) ([]model.Listing, error) {
	var listings []model.Listing
	err := r.db.Select(&listings, `
		SELECT l.*, COALESCE(u.name, '') as user_name 
		FROM listings l 
		LEFT JOIN users u ON l.user_id = u.id 
		WHERE l.user_id = $1 
		ORDER BY l.created_at DESC
	`, userID)
	if err != nil {
		log.Printf("Error getting user listings: %v", err)
		return nil, fmt.Errorf("error getting user listings: %w", err)
	}

	// Get images for each listing
	for i := range listings {
		listings[i].Images = []model.Image{} // Initialize with empty slice to avoid null in JSON
		err = r.db.Select(&listings[i].Images, "SELECT * FROM listing_images WHERE listing_id = $1 ORDER BY is_main DESC, created_at ASC", listings[i].ID)
		if err != nil {
			log.Printf("Error getting images for listing %d: %v", listings[i].ID, err)
			// Continue without images if there's an error
		}
	}

	return listings, nil
}

// SearchListings searches for listings by keyword in title or description
func (r *Repository) SearchListings(keyword string, filter model.ListingFilter) (*model.ListingResponse, error) {
	// Split the keyword into words for better search
	words := strings.Fields(keyword)
	if len(words) == 0 {
		return r.GetListings(filter) // If no search term, return regular listings
	}

	// Build the query with filters
	query := "SELECT l.*, COALESCE(u.name, '') as user_name FROM listings l LEFT JOIN users u ON l.user_id = u.id WHERE l.status = 'active' OR l.status = 'active'"
	countQuery := "SELECT COUNT(*) FROM listings l WHERE l.status = 'active'"
	var args []interface{}
	var countArgs []interface{}
	argIndex := 1

	// Add title and description search conditions
	titleConditions := make([]string, len(words))
	descConditions := make([]string, len(words))

	for i, word := range words {
		titleConditions[i] = fmt.Sprintf("l.title ILIKE $%d", argIndex)
		args = append(args, "%"+word+"%")
		countArgs = append(countArgs, "%"+word+"%")
		argIndex++

		descConditions[i] = fmt.Sprintf("l.description ILIKE $%d", argIndex)
		args = append(args, "%"+word+"%")
		countArgs = append(countArgs, "%"+word+"%")
		argIndex++
	}

	query += " AND (" + strings.Join(titleConditions, " OR ") + " OR " + strings.Join(descConditions, " OR ") + ")"
	countQuery += " AND (" + strings.Join(titleConditions, " OR ") + " OR " + strings.Join(descConditions, " OR ") + ")"

	// Apply category filter
	if filter.CategoryID != nil && *filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND l.category_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.category_id = $%d", argIndex)
		args = append(args, *filter.CategoryID)
		countArgs = append(countArgs, *filter.CategoryID)
		argIndex++
	}

	// Apply city filter
	if filter.City != "" {
		query += fmt.Sprintf(" AND l.city ILIKE $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.city ILIKE $%d", argIndex)
		args = append(args, "%"+filter.City+"%")
		countArgs = append(countArgs, "%"+filter.City+"%")
		argIndex++
	}

	// Apply condition filter
	if filter.Condition != "" {
		query += fmt.Sprintf(" AND l.condition = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.condition = $%d", argIndex)
		args = append(args, filter.Condition)
		countArgs = append(countArgs, filter.Condition)
		argIndex++
	}

	// Apply min price filter
	if filter.MinPrice != nil && *filter.MinPrice >= 0 {
		query += fmt.Sprintf(" AND l.price >= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.price >= $%d", argIndex)
		args = append(args, *filter.MinPrice)
		countArgs = append(countArgs, *filter.MinPrice)
		argIndex++
	}

	// Apply max price filter
	if filter.MaxPrice != nil && *filter.MaxPrice > 0 {
		query += fmt.Sprintf(" AND l.price <= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND l.price <= $%d", argIndex)
		args = append(args, *filter.MaxPrice)
		countArgs = append(countArgs, *filter.MaxPrice)
		argIndex++
	}

	// Apply sorting
	switch filter.SortBy {
	case "date":
		query += " ORDER BY l.created_at ASC"
	case "-date":
		query += " ORDER BY l.created_at DESC"
	case "price":
		query += " ORDER BY l.price ASC"
	case "-price":
		query += " ORDER BY l.price DESC"
	default:
		query += " ORDER BY l.created_at DESC" // Default sort by newest
	}

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, offset)
	argIndex += 2

	// Get total count
	var totalCount int
	err := r.db.Get(&totalCount, countQuery, countArgs...)
	if err != nil {
		log.Printf("Error getting listings count: %v", err)
		return nil, fmt.Errorf("error getting listings count: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(filter.Limit)))

	// Get listings
	var listings []model.Listing
	err = r.db.Select(&listings, query, args...)
	if err != nil {
		log.Printf("Error getting listings: %v", err)
		return nil, fmt.Errorf("error getting listings: %w", err)
	}

	// Get images for each listing
	for i := range listings {
		listings[i].Images = []model.Image{} // Initialize with empty slice to avoid null in JSON
		err = r.db.Select(&listings[i].Images, "SELECT * FROM listing_images WHERE listing_id = $1 ORDER BY is_main DESC, created_at ASC", listings[i].ID)
		if err != nil {
			log.Printf("Error getting images for listing %d: %v", listings[i].ID, err)
			// Continue without images if there's an error
		}
	}

	return &model.ListingResponse{
		Listings:    listings,
		TotalCount:  totalCount,
		CurrentPage: filter.Page,
		TotalPages:  totalPages,
	}, nil
}
