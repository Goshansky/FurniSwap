package model

import (
	"FurniSwap/internal/modules/listing/model"
	"time"
)

// Purchase represents a purchase transaction
type Purchase struct {
	ID         int            `db:"id" json:"id"`
	UserID     int            `db:"user_id" json:"user_id"`
	ListingID  int            `db:"listing_id" json:"listing_id"`
	SellerID   int            `db:"seller_id" json:"seller_id"`
	Price      float64        `db:"price" json:"price"`
	Status     string         `db:"status" json:"status"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at" json:"updated_at"`
	Listing    *model.Listing `json:"listing,omitempty"`
	SellerName string         `db:"seller_name" json:"seller_name,omitempty"`
	BuyerName  string         `db:"buyer_name" json:"buyer_name,omitempty"`
}

// PurchaseResponse represents a list of purchases with pagination
type PurchaseResponse struct {
	Purchases   []Purchase `json:"purchases"`
	TotalCount  int        `json:"total_count"`
	CurrentPage int        `json:"current_page"`
	TotalPages  int        `json:"total_pages"`
}

// BuyRequest represents the data needed to buy a listing
type BuyRequest struct {
	ListingID int `json:"listing_id" binding:"required"`
}
