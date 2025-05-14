package service

import (
	listingModel "FurniSwap/internal/modules/listing/model"
	"FurniSwap/internal/modules/listing/repository"
	listingRepo "FurniSwap/internal/modules/listing/repository"
	"FurniSwap/internal/modules/purchase/model"
	purchaseRepo "FurniSwap/internal/modules/purchase/repository"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// Service provides purchase operations
type Service struct {
	repo        *purchaseRepo.Repository
	listingRepo *listingRepo.Repository
}

// NewService creates a new purchase service
func NewService(repo *purchaseRepo.Repository, listingRepo *repository.Repository) *Service {
	return &Service{
		repo:        repo,
		listingRepo: listingRepo,
	}
}

// BuyListing handles purchasing a listing
func (s *Service) BuyListing(userID int, req model.BuyRequest) (int, error) {
	// Get the listing
	listing, err := s.listingRepo.GetListing(req.ListingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Listing not found: %d", req.ListingID)
			return 0, errors.New("listing not found")
		}
		log.Printf("Error getting listing %d: %v", req.ListingID, err)
		return 0, fmt.Errorf("error getting listing: %w", err)
	}

	// Check if listing is available
	if listing.Status != "active" {
		log.Printf("Listing %d is not available, status: %s", req.ListingID, listing.Status)
		return 0, errors.New("listing is not available for purchase")
	}

	// Check if user is trying to buy own listing
	if listing.UserID == userID {
		return 0, errors.New("you cannot buy your own listing")
	}

	// Create purchase record
	purchaseID, err := s.repo.CreatePurchase(userID, listing.ID, listing.UserID, listing.Price)
	if err != nil {
		log.Printf("Error creating purchase for listing %d: %v", req.ListingID, err)
		return 0, fmt.Errorf("error creating purchase: %w", err)
	}

	// Update listing status to "sold"
	updateReq := listingModel.UpdateListingRequest{
		Status: "sold",
	}
	err = s.listingRepo.UpdateListing(listing.ID, listing.UserID, updateReq)
	if err != nil {
		log.Printf("Error updating listing %d status to sold: %v", req.ListingID, err)
		return 0, fmt.Errorf("error updating listing status: %w", err)
	}

	return purchaseID, nil
}

// GetPurchaseByID gets a purchase by ID
func (s *Service) GetPurchaseByID(purchaseID int) (*model.Purchase, error) {
	return s.repo.GetPurchaseByID(purchaseID)
}

// GetUserPurchases gets purchases made by a user
func (s *Service) GetUserPurchases(userID, page, limit int) (*model.PurchaseResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.repo.GetUserPurchases(userID, page, limit)
}

// GetUserSales gets sales made by a user
func (s *Service) GetUserSales(userID, page, limit int) (*model.PurchaseResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.repo.GetUserSales(userID, page, limit)
}
