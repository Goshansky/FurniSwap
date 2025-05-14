package model

import "time"

// Profile represents a user profile
type Profile struct {
	ID        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	LastName  string    `db:"last_name" json:"last_name"`
	City      string    `db:"city" json:"city"`
	Avatar    string    `db:"avatar" json:"avatar"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// UpdateProfileRequest represents the data needed to update a profile
type UpdateProfileRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	City     string `json:"city"`
	Avatar   string `json:"avatar"`
}

// PublicProfile is a subset of profile information for public viewing
type PublicProfile struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	City      string    `json:"city"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}
