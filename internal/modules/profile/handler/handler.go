package handler

import (
	"FurniSwap/internal/modules/profile/model"
	"FurniSwap/internal/modules/profile/service"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler provides profile handlers
type Handler struct {
	service *service.Service
}

// NewHandler creates a new profile handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers profile routes to router
func (h *Handler) RegisterRoutes(apiRouter *gin.RouterGroup) {
	profile := apiRouter.Group("/profile")
	{
		profile.GET("", h.GetProfile)
		profile.PUT("", h.UpdateProfile)
		profile.POST("/avatar", h.UploadAvatar)
		profile.POST("/avatar/url", h.SetAvatarURL)
	}

	// Add a public endpoint for getting user avatars
	publicProfile := apiRouter.Group("/users")
	{
		publicProfile.GET("/:id/avatar", h.GetUserAvatar)
	}
}

// GetProfile handles getting the user's profile
func (h *Handler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get profile
	profile, err := h.service.GetProfile(userID.(int))
	if err != nil {
		log.Printf("Error getting profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile handles updating the user's profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Update profile
	err := h.service.UpdateProfile(userID.(int), req)
	if err != nil {
		log.Printf("Error updating profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating profile"})
		return
	}

	// Get updated profile
	profile, err := h.service.GetProfile(userID.(int))
	if err != nil {
		log.Printf("Error getting updated profile: %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UploadAvatar handles uploading a user avatar
func (h *Handler) UploadAvatar(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Check file type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File must be an image (JPEG, PNG, or GIF)"})
		return
	}

	// Upload avatar
	err = h.service.UploadAvatar(userID.(int), file)
	if err != nil {
		log.Printf("Error uploading avatar: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading avatar"})
		return
	}

	// Get updated profile
	profile, err := h.service.GetProfile(userID.(int))
	if err != nil {
		log.Printf("Error getting updated profile: %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded successfully"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetUserAvatar handles getting a user's avatar image
func (h *Handler) GetUserAvatar(c *gin.Context) {
	// Parse user ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get the user's avatar path
	avatarPath, err := h.service.GetAvatarPath(userID)
	if err != nil {
		log.Printf("Error getting avatar path: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}

	// Если аватар не задан, возвращаем 404
	if avatarPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "User has no avatar"})
		return
	}

	// Добавляем заголовок для кеширования на 1 день (86400 секунд)
	c.Header("Cache-Control", "public, max-age=86400")

	// If the avatar is a URL, redirect to it
	if h.service.IsAvatarURL(avatarPath) {
		c.Redirect(http.StatusTemporaryRedirect, avatarPath)
		return
	}

	// Otherwise, serve the file
	fullPath := filepath.Join("uploads", avatarPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar file not found"})
		return
	}

	c.File(fullPath)
}

// SetAvatarURL sets a URL as the user's avatar
func (h *Handler) SetAvatarURL(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Set avatar URL
	err := h.service.SetAvatarURL(userID.(int), req.URL)
	if err != nil {
		log.Printf("Error setting avatar URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error setting avatar URL"})
		return
	}

	// Get updated profile
	profile, err := h.service.GetProfile(userID.(int))
	if err != nil {
		log.Printf("Error getting updated profile: %v", err)
		c.JSON(http.StatusOK, gin.H{"message": "Avatar URL set successfully"})
		return
	}

	c.JSON(http.StatusOK, profile)
}
