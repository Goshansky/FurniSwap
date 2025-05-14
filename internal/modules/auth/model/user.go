package model

import "time"

// User represents a user in the system
type User struct {
	ID           int       `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Name         string    `db:"name" json:"name"`
	LastName     string    `db:"last_name" json:"last_name"`
	City         string    `db:"city" json:"city"`
	Avatar       string    `db:"avatar" json:"avatar"`
	IsVerified   bool      `db:"is_verified" json:"is_verified"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// TwoFactorCode represents a 2FA verification code
type TwoFactorCode struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Code      string    `db:"code"`
	ExpiresAt time.Time `db:"expires_at"`
}

// UserProfile represents a user profile for public access
type UserProfile struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	City      string    `json:"city"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

// UserResponse is returned on successful authentication
type UserResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Token    string `json:"token"`
}

// RegisterRequest represents the data needed for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
}

// LoginRequest represents the data needed for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// VerifyRequest represents the data needed for verification code validation
type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// Verify2FARequest represents the data needed for two-factor authentication
type Verify2FARequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}
