package service

import (
	"FurniSwap/internal/modules/listing/model"
	"FurniSwap/internal/modules/listing/repository"
	"FurniSwap/pkg/utils"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

// Service provides listing operations
type Service struct {
	repo *repository.Repository
}

// NewService creates a new listing service
func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateListing creates a new listing
func (s *Service) CreateListing(userID int, req model.CreateListingRequest) (int, error) {
	return s.repo.CreateListing(userID, req)
}

// UpdateListing updates an existing listing
func (s *Service) UpdateListing(listingID, userID int, req model.UpdateListingRequest) error {
	return s.repo.UpdateListing(listingID, userID, req)
}

// DeleteListing deletes a listing
func (s *Service) DeleteListing(listingID, userID int) error {
	// Get listing to retrieve images
	listing, err := s.repo.GetListing(listingID)
	if err != nil {
		log.Printf("Error getting listing for deletion: %v", err)
		// Continue with deletion attempt anyway
	} else {
		// Delete all local image files (skip URLs)
		for _, image := range listing.Images {
			// Only delete the file if it's a local image (not a URL)
			if !strings.HasPrefix(image.ImagePath, "http://") && !strings.HasPrefix(image.ImagePath, "https://") {
				err = utils.DeleteFile(image.ImagePath)
				if err != nil {
					log.Printf("Error deleting image file: %v", err)
					// Continue deletion process anyway
				}
			}
		}
	}

	// Delete listing from database
	return s.repo.DeleteListing(listingID, userID)
}

// GetListing gets a single listing by ID
func (s *Service) GetListing(listingID int) (*model.Listing, error) {
	return s.repo.GetListing(listingID)
}

// GetListings gets listings with filtering and pagination
func (s *Service) GetListings(filter model.ListingFilter) (*model.ListingResponse, error) {
	return s.repo.GetListings(filter)
}

// GetUserListings gets all listings for a user
func (s *Service) GetUserListings(userID int) ([]model.Listing, error) {
	return s.repo.GetUserListings(userID)
}

// SearchListings searches for listings by keyword
func (s *Service) SearchListings(keyword string, filter model.ListingFilter) (*model.ListingResponse, error) {
	return s.repo.SearchListings(keyword, filter)
}

// UploadListingImage uploads an image for a listing
func (s *Service) UploadListingImage(listingID, userID int, file *multipart.FileHeader) (int, error) {
	// First check if the listing exists and belongs to the user
	listing, err := s.repo.GetListing(listingID)
	if err != nil {
		return 0, fmt.Errorf("error getting listing: %w", err)
	}

	if listing.UserID != userID {
		return 0, fmt.Errorf("listing does not belong to the user")
	}

	// Upload the image
	imagePath, err := utils.UploadFile(file, fmt.Sprintf("listings/%d", listingID))
	if err != nil {
		return 0, fmt.Errorf("error uploading image: %w", err)
	}

	// Add the image to the database
	imageID, err := s.repo.AddImage(listingID, imagePath)
	if err != nil {
		// If there's an error adding to the database, delete the uploaded file
		_ = utils.DeleteFile(imagePath)
		return 0, fmt.Errorf("error adding image to database: %w", err)
	}

	return imageID, nil
}

// DeleteListingImage deletes an image from a listing
func (s *Service) DeleteListingImage(imageID, listingID, userID int) error {
	// Get listing to retrieve image path
	listing, err := s.repo.GetListing(listingID)
	if err != nil {
		return fmt.Errorf("error getting listing: %w", err)
	}

	// Find the image
	var imagePath string
	for _, image := range listing.Images {
		if image.ID == imageID {
			imagePath = image.ImagePath
			break
		}
	}

	if imagePath == "" {
		return fmt.Errorf("image not found")
	}

	// Delete image from database
	err = s.repo.DeleteImage(imageID, listingID, userID)
	if err != nil {
		return err
	}

	// Only delete the file if it's a local image (not a URL)
	if !strings.HasPrefix(imagePath, "http://") && !strings.HasPrefix(imagePath, "https://") {
		// Delete image file
		err = utils.DeleteFile(imagePath)
		if err != nil {
			log.Printf("Error deleting image file: %v", err)
			// Continue anyway as the database entry is already deleted
		}
	}

	return nil
}

// SetMainImage sets an image as the main image for a listing
func (s *Service) SetMainImage(imageID, listingID, userID int) error {
	return s.repo.SetMainImage(imageID, listingID, userID)
}

// AddListingImageURL adds an image URL to a listing
func (s *Service) AddListingImageURL(listingID, userID int, imageURL string) (int, error) {
	// First check if the listing exists and belongs to the user
	listing, err := s.repo.GetListing(listingID)
	if err != nil {
		return 0, fmt.Errorf("error getting listing: %w", err)
	}

	if listing.UserID != userID {
		return 0, fmt.Errorf("listing does not belong to the user")
	}

	// Basic validation of the URL
	if !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
		return 0, fmt.Errorf("invalid image URL, must start with http:// or https://")
	}

	// Optional: Verify the URL points to an image
	// This is a basic check and could be improved
	resp, err := http.Head(imageURL)
	if err != nil {
		return 0, fmt.Errorf("error validating image URL: %w", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return 0, fmt.Errorf("URL does not point to a valid image: %s", contentType)
	}

	// Add the image URL to the database
	imageID, err := s.repo.AddImage(listingID, imageURL)
	if err != nil {
		return 0, fmt.Errorf("error adding image to database: %w", err)
	}

	return imageID, nil
}
