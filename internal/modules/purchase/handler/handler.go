package handler

import (
	"FurniSwap/internal/modules/purchase/model"
	"FurniSwap/internal/modules/purchase/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler provides purchase handlers
type Handler struct {
	service *service.Service
}

// NewHandler creates a new purchase handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers purchase routes to router
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/listings/:id/buy", h.BuyListing)
	router.GET("/purchases", h.GetUserPurchases)
	router.GET("/sales", h.GetUserSales)
}

// BuyListing handles buying a listing
func (h *Handler) BuyListing(c *gin.Context) {
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

	// Create buy request
	buyReq := model.BuyRequest{
		ListingID: listingID,
	}

	// Buy listing
	purchaseID, err := h.service.BuyListing(userID.(int), buyReq)
	if err != nil {
		if err.Error() == "listing not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Listing not found"})
			return
		}
		if err.Error() == "listing is not available for purchase" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Listing is not available for purchase"})
			return
		}
		if err.Error() == "you cannot buy your own listing" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot buy your own listing"})
			return
		}
		log.Printf("Error buying listing: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error buying listing"})
		return
	}

	// Get the purchase
	purchase, err := h.service.GetPurchaseByID(purchaseID)
	if err != nil {
		log.Printf("Error getting purchase: %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "Purchase successful", "purchase_id": purchaseID})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

// GetUserPurchases handles getting the user's purchases
func (h *Handler) GetUserPurchases(c *gin.Context) {
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

	// Get purchases
	purchases, err := h.service.GetUserPurchases(userID.(int), page, limit)
	if err != nil {
		log.Printf("Error getting purchases: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting purchases"})
		return
	}

	c.JSON(http.StatusOK, purchases)
}

// GetUserSales handles getting the user's sales
func (h *Handler) GetUserSales(c *gin.Context) {
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

	// Get sales
	sales, err := h.service.GetUserSales(userID.(int), page, limit)
	if err != nil {
		log.Printf("Error getting sales: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting sales"})
		return
	}

	c.JSON(http.StatusOK, sales)
}
