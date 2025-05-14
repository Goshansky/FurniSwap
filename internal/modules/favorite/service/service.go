package service

import (
	"FurniSwap/internal/modules/favorite/model"
	"FurniSwap/internal/modules/favorite/repository"
)

// Service provides favorite operations
type Service struct {
	repo *repository.Repository
}

// NewService creates a new favorite service
func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// AddFavorite adds a listing to a user's favorites
func (s *Service) AddFavorite(userID, listingID int) error {
	return s.repo.AddFavorite(userID, listingID)
}

// RemoveFavorite removes a listing from a user's favorites
func (s *Service) RemoveFavorite(userID, listingID int) error {
	return s.repo.RemoveFavorite(userID, listingID)
}

// IsFavorite checks if a listing is in a user's favorites
func (s *Service) IsFavorite(userID, listingID int) (bool, error) {
	return s.repo.IsFavorite(userID, listingID)
}

// GetFavorites gets a user's favorite listings with pagination
func (s *Service) GetFavorites(userID, page, limit int) (*model.FavoriteResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.repo.GetFavorites(userID, page, limit)
}
