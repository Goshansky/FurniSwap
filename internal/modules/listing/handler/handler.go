package handler

import (
	"FurniSwap/internal/modules/listing/model"
	"FurniSwap/internal/modules/listing/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler provides listing handlers
type Handler struct {
	service *service.Service
}

// NewHandler creates a new listing handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterPublicRoutes registers public listing routes (no auth required)
func (h *Handler) RegisterPublicRoutes(router *gin.RouterGroup) {
	router.GET("", h.GetListings)
	router.GET("/:id", h.GetListing)
}

// RegisterProtectedRoutes registers protected listing routes (auth required)
func (h *Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.POST("/listings", h.CreateListing)
	router.PUT("/listings/:id", h.UpdateListing)
	router.DELETE("/listings/:id", h.DeleteListing)
	router.POST("/listings/:id/images", h.UploadListingImage)
	router.DELETE("/listings/:id/images/:imageId", h.DeleteListingImage)
	router.PUT("/listings/:id/images/:imageId/main", h.SetMainImage)
	router.GET("/listings/my", h.GetUserListings)
}

// RegisterRoutes registers listing routes to router
// Deprecated: Use RegisterPublicRoutes and RegisterProtectedRoutes instead
func (h *Handler) RegisterRoutes(apiRouter *gin.RouterGroup, publicRouter *gin.RouterGroup) {
	// Public listing routes (no auth required)
	if publicRouter != nil {
		publicRouter.GET("/listings", h.GetListings)
		publicRouter.GET("/listings/:id", h.GetListing)
	}

	// Protected listing routes (auth required)
	if apiRouter != nil {
		apiRouter.POST("/listings", h.CreateListing)
		apiRouter.PUT("/listings/:id", h.UpdateListing)
		apiRouter.DELETE("/listings/:id", h.DeleteListing)
		apiRouter.POST("/listings/:id/images", h.UploadListingImage)
		apiRouter.DELETE("/listings/:id/images/:imageId", h.DeleteListingImage)
		apiRouter.PUT("/listings/:id/images/:imageId/main", h.SetMainImage)
	}
}

// GetListings handles getting listings with filtering and pagination
func (h *Handler) GetListings(c *gin.Context) {
	var filter model.ListingFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters"})
		return
	}

	// Check if there's a search query
	search := c.Query("search")
	var response *model.ListingResponse
	var err error

	if search != "" {
		// Search listings
		response, err = h.service.SearchListings(search, filter)
	} else {
		// Get listings without search
		response, err = h.service.GetListings(filter)
	}

	if err != nil {
		log.Printf("Error getting listings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting listings"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetListing handles getting a single listing
func (h *Handler) GetListing(c *gin.Context) {
	// Parse listing ID
	listingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	// Get listing
	listing, err := h.service.GetListing(listingID)
	if err != nil {
		log.Printf("Error getting listing: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
		return
	}

	c.JSON(http.StatusOK, listing)
}

// CreateListing handles creating a new listing
func (h *Handler) CreateListing(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var req model.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Create listing
	listingID, err := h.service.CreateListing(userID.(int), req)
	if err != nil {
		log.Printf("Error creating listing: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating listing"})
		return
	}

	// Get the created listing
	listing, err := h.service.GetListing(listingID)
	if err != nil {
		log.Printf("Error getting created listing: %v", err)
		c.JSON(http.StatusOK, gin.H{"id": listingID, "message": "Listing created successfully"})
		return
	}

	c.JSON(http.StatusOK, listing)
}

// UpdateListing handles updating an existing listing
func (h *Handler) UpdateListing(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse listing ID
	listingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	// Parse request body
	var req model.UpdateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Update listing
	err = h.service.UpdateListing(listingID, userID.(int), req)
	if err != nil {
		if err.Error() == "listing not found or does not belong to the user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Listing not found or you don't have permission to update it"})
			return
		}
		log.Printf("Error updating listing: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating listing"})
		return
	}

	// Get the updated listing
	listing, err := h.service.GetListing(listingID)
	if err != nil {
		log.Printf("Error getting updated listing: %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "Listing updated successfully"})
		return
	}

	c.JSON(http.StatusOK, listing)
}

// DeleteListing handles deleting a listing
func (h *Handler) DeleteListing(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse listing ID
	listingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	// Delete listing
	err = h.service.DeleteListing(listingID, userID.(int))
	if err != nil {
		if err.Error() == "listing not found or does not belong to the user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Listing not found or you don't have permission to delete it"})
			return
		}
		log.Printf("Error deleting listing: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting listing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Listing deleted successfully"})
}

// UploadListingImage handles uploading an image for a listing
func (h *Handler) UploadListingImage(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse listing ID
	listingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	// Check if the request has a file or an image URL
	var imageID int
	var imageErr error

	if url := c.PostForm("image_url"); url != "" {
		// Use the provided image URL
		imageID, imageErr = h.service.AddListingImageURL(listingID, userID.(int), url)
	} else {
		// Get file from form for backwards compatibility
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No image URL or file provided"})
			return
		}

		// Check file type
		contentType := file.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image (JPEG, PNG, or GIF)"})
			return
		}

		// Upload image
		imageID, imageErr = h.service.UploadListingImage(listingID, userID.(int), file)
	}

	if imageErr != nil {
		if imageErr.Error() == "listing does not belong to the user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to upload images to this listing"})
			return
		}
		log.Printf("Error adding image: %v", imageErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding image"})
		return
	}

	// Get the updated listing
	listing, err := h.service.GetListing(listingID)
	if err != nil {
		log.Printf("Error getting updated listing: %v", err)
		c.JSON(http.StatusOK, gin.H{"id": imageID, "message": "Image added successfully"})
		return
	}

	c.JSON(http.StatusOK, listing)
}

// DeleteListingImage handles deleting an image from a listing
func (h *Handler) DeleteListingImage(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse listing ID and image ID
	listingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	imageID, err := strconv.Atoi(c.Param("imageId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	// Delete image
	err = h.service.DeleteListingImage(imageID, listingID, userID.(int))
	if err != nil {
		if err.Error() == "image not found or does not belong to the user's listing" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Image not found or you don't have permission to delete it"})
			return
		}
		log.Printf("Error deleting image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting image"})
		return
	}

	// Get the updated listing
	listing, err := h.service.GetListing(listingID)
	if err != nil {
		log.Printf("Error getting updated listing: %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
		return
	}

	c.JSON(http.StatusOK, listing)
}

// SetMainImage handles setting an image as the main image for a listing
func (h *Handler) SetMainImage(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse listing ID and image ID
	listingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing ID"})
		return
	}

	imageID, err := strconv.Atoi(c.Param("imageId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	// Set main image
	err = h.service.SetMainImage(imageID, listingID, userID.(int))
	if err != nil {
		if err.Error() == "image not found or does not belong to the user's listing" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Image not found or you don't have permission to modify it"})
			return
		}
		log.Printf("Error setting main image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error setting main image"})
		return
	}

	// Get the updated listing
	listing, err := h.service.GetListing(listingID)
	if err != nil {
		log.Printf("Error getting updated listing: %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "Main image set successfully"})
		return
	}

	c.JSON(http.StatusOK, listing)
}

// GetUserListings handles getting listings belonging to the authenticated user
func (h *Handler) GetUserListings(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user's listings
	listings, err := h.service.GetUserListings(userID.(int))
	if err != nil {
		log.Printf("Error getting user listings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user listings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"listings": listings,
		"count":    len(listings),
	})
}
