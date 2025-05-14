package middleware

import (
	"FurniSwap/pkg/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// AuthRequired middleware checks JWT token and authorizes the user
func AuthRequired(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Missing authorization token in request")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
			c.Abort()
			return
		}

		// Check "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Invalid token format:", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization token format"})
			c.Abort()
			return
		}

		token := parts[1]
		log.Printf("JWT token length: %d characters\n", len(token))

		// Check if token is empty
		if token == "" {
			log.Println("Error: empty JWT token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "empty authorization token"})
			c.Abort()
			return
		}

		// Validate token and get user ID
		userID, err := utils.ValidateToken(token)
		if err != nil {
			log.Printf("Token validation error: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Check if userID is positive
		if userID <= 0 {
			log.Printf("Invalid userID from token: %d\n", userID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		log.Printf("Token successfully validated for user ID: %d\n", userID)

		// Check if user exists in the database
		var exists bool
		err = db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND is_verified = true)", userID)
		if err != nil {
			log.Printf("Error checking user in database: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking user"})
			c.Abort()
			return
		}

		if !exists {
			log.Printf("User ID: %d not found or not verified\n", userID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found or not verified"})
			c.Abort()
			return
		}

		// Set user ID in context for further use
		c.Set("userID", userID)
		log.Printf("User ID: %d successfully authenticated\n", userID)
		c.Next()
	}
}
