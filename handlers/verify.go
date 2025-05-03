package handlers

import (
	"FurniSwap/utils"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

		// Проверяем существование пользователя
		var exists bool
		err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки пользователя"})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}

		var userID int
		var expiresAt time.Time
		err = db.QueryRow(`
			SELECT user_id, expires_at FROM two_factor_codes 
			WHERE code = $1 AND user_id = (SELECT id FROM users WHERE email = $2)
		`, req.Code, req.Email).Scan(&userID, &expiresAt)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный код"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки кода"})
			return
		}

		// Проверяем срок действия кода
		if time.Now().After(expiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Код просрочен"})
			return
		}

		// Обновляем статус пользователя
		_, err = db.Exec(`UPDATE users SET is_verified = true WHERE id = $1`, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления пользователя"})
			return
		}

		// Удаляем использованный код
		_, err = db.Exec(`DELETE FROM two_factor_codes WHERE user_id = $1`, userID)
		if err != nil {
			// Логируем ошибку, но продолжаем выполнение
			// Это некритичная ошибка, так как пользователь уже верифицирован
		}

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

		// Проверяем существование и верификацию пользователя
		var userID int
		var isVerified bool
		err := db.QueryRow(`
			SELECT id, is_verified FROM users WHERE email = $1
		`, req.Email).Scan(&userID, &isVerified)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки пользователя"})
			return
		}

		if !isVerified {
			c.JSON(http.StatusForbidden, gin.H{"error": "Email не подтвержден"})
			return
		}

		// Проверяем код
		var expiresAt time.Time
		err = db.QueryRow(`
			SELECT expires_at FROM two_factor_codes 
			WHERE code = $1 AND user_id = $2
		`, req.Code, userID).Scan(&expiresAt)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный код"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки кода"})
			return
		}

		// Проверяем срок действия кода
		if time.Now().After(expiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Код просрочен"})
			return
		}

		// Удаляем использованный код
		_, err = db.Exec(`DELETE FROM two_factor_codes WHERE user_id = $1`, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке аутентификации"})
			return
		}

		// Генерируем JWT
		token, err := utils.GenerateToken(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания токена"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
