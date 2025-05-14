package service

import (
	"FurniSwap/internal/modules/auth/model"
	"FurniSwap/internal/modules/auth/repository"
	"FurniSwap/pkg/utils"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Service provides authentication operations
type Service struct {
	repo *repository.Repository
}

// NewService creates a new auth service
func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Register handles user registration
func (s *Service) Register(req model.RegisterRequest) (string, error) {
	// Check if user already exists
	exists, err := s.repo.UserExists(req.Email)
	if err != nil {
		return "", fmt.Errorf("error checking user existence: %w", err)
	}
	if exists {
		return "", errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}

	// Default values for city and avatar
	defaultCity := "Москва"
	defaultAvatar := "https://cs12.pikabu.ru/post_img/2021/05/08/12/1620504640126273687.jpg"

	// Create user
	userID, err := s.repo.CreateUser(req.Email, string(hashedPassword), req.Name, req.LastName, defaultCity, defaultAvatar)
	if err != nil {
		return "", fmt.Errorf("error creating user: %w", err)
	}

	// Generate verification code
	code := utils.GenerateCode()

	// Save code in database with 10-minute expiration
	err = s.repo.SaveVerificationCode(userID, code, time.Now().Add(10*time.Minute))
	if err != nil {
		return "", fmt.Errorf("error saving verification code: %w", err)
	}

	// Send verification email
	err = utils.SendEmail(req.Email, "Verification Code", "Your code: "+code)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
	}

	return code, nil
}

// VerifyUser verifies a user's email
func (s *Service) VerifyUser(req model.VerifyRequest) error {
	// Verify the code
	userID, err := s.repo.VerifyCode(req.Email, req.Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("invalid or expired verification code")
		}
		return fmt.Errorf("error verifying code: %w", err)
	}

	// Mark user as verified
	err = s.repo.SetUserVerified(userID)
	if err != nil {
		return fmt.Errorf("error setting user as verified: %w", err)
	}

	// Delete used verification codes
	err = s.repo.DeleteVerificationCodes(userID)
	if err != nil {
		log.Printf("Error deleting verification codes: %v", err)
	}

	return nil
}

// Login authenticates a user and returns user info with token
func (s *Service) Login(req model.LoginRequest) (*model.UserResponse, string, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", errors.New("invalid email or password")
		}
		return nil, "", fmt.Errorf("error getting user: %w", err)
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Check if user is verified
	if !user.IsVerified {
		// Generate new verification code
		code := utils.GenerateCode()

		// Save code in database with 10-minute expiration
		err = s.repo.SaveVerificationCode(user.ID, code, time.Now().Add(10*time.Minute))
		if err != nil {
			return nil, "", fmt.Errorf("error saving verification code: %w", err)
		}

		// Send verification email
		err = utils.SendEmail(user.Email, "Verification Code", "Your code: "+code)
		if err != nil {
			log.Printf("Error sending verification email: %v", err)
		}

		return nil, "", errors.New("account not verified; new verification code sent")
	}

	// Generate 2FA code for login
	code := utils.GenerateCode()

	// Save code in database with 10-minute expiration
	err = s.repo.SaveVerificationCode(user.ID, code, time.Now().Add(10*time.Minute))
	if err != nil {
		return nil, "", fmt.Errorf("error saving 2FA code: %w", err)
	}

	// Send 2FA code via email
	err = utils.SendEmail(user.Email, "Two-Factor Authentication Code", "Your 2FA code: "+code)
	if err != nil {
		log.Printf("Error sending 2FA email: %v", err)
	}

	// Return partial user info without token (will be completed after 2FA)
	return &model.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		LastName: user.LastName,
	}, code, nil
}

// Verify2FA verifies the two-factor authentication code
func (s *Service) Verify2FA(req model.Verify2FARequest) (*model.UserResponse, error) {
	// Verify the code
	userID, err := s.repo.VerifyCode(req.Email, req.Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid or expired 2FA code")
		}
		return nil, fmt.Errorf("error verifying 2FA code: %w", err)
	}

	// Get user by ID
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	// Delete used verification codes
	err = s.repo.DeleteVerificationCodes(user.ID)
	if err != nil {
		log.Printf("Error deleting verification codes: %v", err)
	}

	// Return user info with token
	return &model.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Name:     user.Name,
		LastName: user.LastName,
		Token:    token,
	}, nil
}
