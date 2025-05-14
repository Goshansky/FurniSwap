package repository

import (
	"FurniSwap/internal/modules/purchase/model"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for the purchase module
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new purchase repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// CreatePurchase creates a new purchase
func (r *Repository) CreatePurchase(userID, listingID, sellerID int, price float64) (int, error) {
	var purchaseID int
	err := r.db.QueryRow(`
		INSERT INTO purchases (buyer_id, listing_id, seller_id, price, purchased_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, userID, listingID, sellerID, price, time.Now()).Scan(&purchaseID)

	if err != nil {
		log.Printf("Error creating purchase: %v", err)
		return 0, fmt.Errorf("error creating purchase: %w", err)
	}

	return purchaseID, nil
}

// GetPurchaseByID gets a purchase by ID
func (r *Repository) GetPurchaseByID(purchaseID int) (*model.Purchase, error) {
	var purchase model.Purchase
	err := r.db.Get(&purchase, `
		SELECT p.id, p.buyer_id as user_id, p.listing_id, p.seller_id, p.price, 
		       p.purchased_at as created_at, p.purchased_at as updated_at, 
			   u1.name || ' ' || COALESCE(u1.last_name, '') as buyer_name,
			   u2.name || ' ' || COALESCE(u2.last_name, '') as seller_name
		FROM purchases p
		JOIN users u1 ON p.buyer_id = u1.id
		JOIN users u2 ON p.seller_id = u2.id
		WHERE p.id = $1
	`, purchaseID)
	if err != nil {
		log.Printf("Error getting purchase by ID: %v", err)
		return nil, fmt.Errorf("error getting purchase: %w", err)
	}
	return &purchase, nil
}

// GetUserPurchases gets purchases made by a user with pagination
func (r *Repository) GetUserPurchases(userID, page, limit int) (*model.PurchaseResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get total count
	var totalCount int
	err := r.db.Get(&totalCount, "SELECT COUNT(*) FROM purchases WHERE buyer_id = $1", userID)
	if err != nil {
		log.Printf("Error getting purchases count: %v", err)
		return nil, fmt.Errorf("error getting purchases count: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Get purchases
	var purchases []model.Purchase
	err = r.db.Select(&purchases, `
		SELECT p.id, p.buyer_id as user_id, p.listing_id, p.seller_id, p.price, 
		       p.purchased_at as created_at, p.purchased_at as updated_at, 
			   u.name || ' ' || COALESCE(u.last_name, '') as seller_name
		FROM purchases p
		JOIN users u ON p.seller_id = u.id
		WHERE p.buyer_id = $1
		ORDER BY p.purchased_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		log.Printf("Error getting purchases: %v", err)
		return nil, fmt.Errorf("error getting purchases: %w", err)
	}

	return &model.PurchaseResponse{
		Purchases:   purchases,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}

// GetUserSales gets sales made by a user with pagination
func (r *Repository) GetUserSales(userID, page, limit int) (*model.PurchaseResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get total count
	var totalCount int
	err := r.db.Get(&totalCount, "SELECT COUNT(*) FROM purchases WHERE seller_id = $1", userID)
	if err != nil {
		log.Printf("Error getting sales count: %v", err)
		return nil, fmt.Errorf("error getting sales count: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Get sales
	var purchases []model.Purchase
	err = r.db.Select(&purchases, `
		SELECT p.id, p.buyer_id as user_id, p.listing_id, p.seller_id, p.price, 
		       p.purchased_at as created_at, p.purchased_at as updated_at, 
			   u.name || ' ' || COALESCE(u.last_name, '') as buyer_name
		FROM purchases p
		JOIN users u ON p.buyer_id = u.id
		WHERE p.seller_id = $1
		ORDER BY p.purchased_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		log.Printf("Error getting sales: %v", err)
		return nil, fmt.Errorf("error getting sales: %w", err)
	}

	return &model.PurchaseResponse{
		Purchases:   purchases,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}
