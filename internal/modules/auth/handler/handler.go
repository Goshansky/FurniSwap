package handler

import (
	"FurniSwap/internal/modules/auth/model"
	"FurniSwap/internal/modules/auth/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler provides authentication handlers
type Handler struct {
	service *service.Service
}

// NewHandler creates a new auth handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers auth routes to router
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/verify", h.Verify)
		auth.POST("/login", h.Login)
		auth.POST("/verify-2fa", h.Verify2FA)
	}
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	code, err := h.service.Register(req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
			return
		}
		log.Printf("Error during registration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification code sent",
		"code":    code,
	})
}

// Verify handles email verification
func (h *Handler) Verify(c *gin.Context) {
	var req model.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	err := h.service.VerifyUser(req)
	if err != nil {
		if err.Error() == "invalid or expired verification code" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification code"})
			return
		}
		log.Printf("Error during verification: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account verified successfully"})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	user, code, err := h.service.Login(req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		if err.Error() == "account not verified; new verification code sent" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account not verified; new verification code sent"})
			return
		}
		log.Printf("Error during login: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA code sent to your email",
		"code":    code,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"name":      user.Name,
			"last_name": user.LastName,
		},
	})
}

// Verify2FA handles two-factor authentication
func (h *Handler) Verify2FA(c *gin.Context) {
	var req model.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	user, err := h.service.Verify2FA(req)
	if err != nil {
		if err.Error() == "invalid or expired 2FA code" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired 2FA code"})
			return
		}
		log.Printf("Error during 2FA verification: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "2FA verification error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
