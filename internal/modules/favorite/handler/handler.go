package handler

import (
	"FurniSwap/internal/modules/favorite/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler provides favorite handlers
type Handler struct {
	service *service.Service
}

// NewHandler creates a new favorite handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers favorite routes to router
func (h *Handler) RegisterRoutes(apiRouter *gin.RouterGroup) {
	apiRouter.POST("/listings/:id/favorite", h.AddFavorite)
	apiRouter.DELETE("/listings/:id/favorite", h.RemoveFavorite)
	apiRouter.GET("/listings/:id/favorite", h.IsFavorite)
	apiRouter.GET("/favorites", h.GetFavorites)
}

// AddFavorite handles adding a listing to favorites
func (h *Handler) AddFavorite(c *gin.Context) {
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

	// Add to favorites
	err = h.service.AddFavorite(userID.(int), listingID)
	if err != nil {
		if err.Error() == "listing does not exist" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
			return
		}
		log.Printf("Error adding to favorites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding to favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to favorites"})
}

// RemoveFavorite handles removing a listing from favorites
func (h *Handler) RemoveFavorite(c *gin.Context) {
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

	// Remove from favorites
	err = h.service.RemoveFavorite(userID.(int), listingID)
	if err != nil {
		log.Printf("Error removing from favorites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing from favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from favorites"})
}

// IsFavorite handles checking if a listing is in favorites
func (h *Handler) IsFavorite(c *gin.Context) {
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

	// Check if in favorites
	isFavorite, err := h.service.IsFavorite(userID.(int), listingID)
	if err != nil {
		log.Printf("Error checking if in favorites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if in favorites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorite": isFavorite})
}

// GetFavorites handles getting the user's favorites
func (h *Handler) GetFavorites(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// Get favorites
	favorites, err := h.service.GetFavorites(userID.(int), page, limit)
	if err != nil {
		log.Printf("Error getting favorites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting favorites"})
		return
	}

	c.JSON(http.StatusOK, favorites)
}
