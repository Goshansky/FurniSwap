package repository

import (
	"FurniSwap/internal/modules/auth/model"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for the auth module
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetUserByEmail retrieves a user by their email
func (r *Repository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, nil
}

// GetUserByID retrieves a user by their ID
func (r *Repository) GetUserByID(id int) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, nil
}

// CreateUser creates a new user in the database
func (r *Repository) CreateUser(email, passwordHash, name, lastName, city, avatar string) (int, error) {
	var userID int
	err := r.db.QueryRow(`
		INSERT INTO users (email, password_hash, name, last_name, city, avatar) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`, email, passwordHash, name, lastName, city, avatar).Scan(&userID)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return 0, fmt.Errorf("error creating user: %w", err)
	}
	return userID, nil
}

// UserExists checks if a user with the given email exists
func (r *Repository) UserExists(email string) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email)
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		return false, fmt.Errorf("error checking if user exists: %w", err)
	}
	return exists, nil
}

// SaveVerificationCode saves a verification code for a user
func (r *Repository) SaveVerificationCode(userID int, code string, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		INSERT INTO two_factor_codes (user_id, code, expires_at) VALUES ($1, $2, $3)
	`, userID, code, expiresAt)
	if err != nil {
		log.Printf("Error saving verification code: %v", err)
		return fmt.Errorf("error saving verification code: %w", err)
	}
	return nil
}

// VerifyCode checks if a verification code is valid
func (r *Repository) VerifyCode(email, code string) (int, error) {
	var userID int
	err := r.db.Get(&userID, `
		SELECT u.id FROM users u
		JOIN two_factor_codes tfc ON u.id = tfc.user_id
		WHERE u.email = $1 AND tfc.code = $2 AND tfc.expires_at > NOW()
	`, email, code)
	if err != nil {
		log.Printf("Error verifying code: %v", err)
		return 0, fmt.Errorf("error verifying code: %w", err)
	}
	return userID, nil
}

// SetUserVerified sets a user as verified
func (r *Repository) SetUserVerified(userID int) error {
	_, err := r.db.Exec("UPDATE users SET is_verified = true WHERE id = $1", userID)
	if err != nil {
		log.Printf("Error setting user as verified: %v", err)
		return fmt.Errorf("error setting user as verified: %w", err)
	}
	return nil
}

// DeleteVerificationCodes deletes all verification codes for a user
func (r *Repository) DeleteVerificationCodes(userID int) error {
	_, err := r.db.Exec("DELETE FROM two_factor_codes WHERE user_id = $1", userID)
	if err != nil {
		log.Printf("Error deleting verification codes: %v", err)
		return fmt.Errorf("error deleting verification codes: %w", err)
	}
	return nil
}
