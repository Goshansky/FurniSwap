package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// VerifyHandler проверяет код и подтверждает email
func VerifyHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
			return
		}

		var userID int
		var expiresAt time.Time
		err := db.QueryRow(`
			SELECT user_id, expires_at FROM two_factor_codes 
			WHERE code = $1 AND user_id = (SELECT id FROM users WHERE email = $2)
		`, req.Code, req.Email).Scan(&userID, &expiresAt)

		if err == sql.ErrNoRows || time.Now().After(expiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный или просроченный код"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки кода"})
			return
		}

		// Обновляем статус пользователя
		_, err = db.Exec(`UPDATE users SET is_verified = true WHERE id = $1`, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления пользователя"})
			return
		}

		// Удаляем использованный код
		_, _ = db.Exec(`DELETE FROM two_factor_codes WHERE user_id = $1`, userID)

		c.JSON(http.StatusOK, gin.H{"message": "Email подтвержден!"})
	}
}

type Verify2FARequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// Verify2FAHandler проверяет код 2FA и выдает JWT
func Verify2FAHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Verify2FARequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
			return
		}

		var userID int
		var expiresAt time.Time
		err := db.QueryRow(`
			SELECT user_id, expires_at FROM two_factor_codes 
			WHERE code = $1 AND user_id = (SELECT id FROM users WHERE email = $2)
		`, req.Code, req.Email).Scan(&userID, &expiresAt)

		if err == sql.ErrNoRows || time.Now().After(expiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный или просроченный код"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки кода"})
			return
		}

		// Удаляем использованный код
		_, _ = db.Exec(`DELETE FROM two_factor_codes WHERE user_id = $1`, userID)

		// Генерируем JWT
		token := generateJWT(userID)

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

// generateJWT создает JWT-токен
func generateJWT(userID int) string {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
