package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims struct contains custom claims for JWT
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID int) (string, error) {
	// Get secret key from environment variables
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Println("WARNING: JWT_SECRET_KEY not configured, using default value")
		secretKey = "default_secret_key_change_in_production" // Fallback for development
	}

	// Create new claims with user ID
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with specified claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates token and returns user ID
func ValidateToken(tokenString string) (int, error) {
	// Get secret key from environment variables
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Println("WARNING: JWT_SECRET_KEY not configured, using default value")
		secretKey = "default_secret_key_change_in_production" // Fallback for development
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check if the correct signing method is used
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("error parsing token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Get claims from token
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return 0, errors.New("unable to get claims from token")
	}

	return claims.UserID, nil
}
