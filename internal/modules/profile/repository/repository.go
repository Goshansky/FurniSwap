package repository

import (
	"FurniSwap/internal/modules/profile/model"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for the profile module
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new profile repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetProfileByID retrieves a user profile by ID
func (r *Repository) GetProfileByID(userID int) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.Get(&profile, `
		SELECT id, email, name, last_name, city, avatar, created_at
		FROM users 
		WHERE id = $1
	`, userID)
	if err != nil {
		log.Printf("Error getting profile by ID: %v", err)
		return nil, fmt.Errorf("error getting profile: %w", err)
	}
	return &profile, nil
}

// UpdateProfile updates a user profile
func (r *Repository) UpdateProfile(userID int, req model.UpdateProfileRequest) error {
	// Если задано поле аватара, добавляем его в обновление
	if req.Avatar != "" {
		_, err := r.db.Exec(`
			UPDATE users 
			SET name = $1, last_name = $2, city = $3, avatar = $4
			WHERE id = $5
		`, req.Name, req.LastName, req.City, req.Avatar, userID)
		if err != nil {
			log.Printf("Error updating profile with avatar: %v", err)
			return fmt.Errorf("error updating profile: %w", err)
		}
	} else {
		// Обновление без изменения аватара
		_, err := r.db.Exec(`
			UPDATE users 
			SET name = $1, last_name = $2, city = $3
			WHERE id = $4
		`, req.Name, req.LastName, req.City, userID)
		if err != nil {
			log.Printf("Error updating profile: %v", err)
			return fmt.Errorf("error updating profile: %w", err)
		}
	}
	return nil
}

// UpdateAvatar updates a user's avatar
func (r *Repository) UpdateAvatar(userID int, avatarPath string) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET avatar = $1
		WHERE id = $2
	`, avatarPath, userID)
	if err != nil {
		log.Printf("Error updating avatar: %v", err)
		return fmt.Errorf("error updating avatar: %w", err)
	}
	return nil
}

// GetAvatarPath gets the current avatar path for a user
func (r *Repository) GetAvatarPath(userID int) (string, error) {
	var avatarPath string
	err := r.db.Get(&avatarPath, "SELECT avatar FROM users WHERE id = $1", userID)
	if err != nil {
		log.Printf("Error getting avatar path: %v", err)
		return "", fmt.Errorf("error getting avatar path: %w", err)
	}
	return avatarPath, nil
}
