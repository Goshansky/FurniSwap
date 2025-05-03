package handlers

import (
	"FurniSwap/utils"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler обрабатывает вход пользователя
func LoginHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
			return
		}

		// Получаем пользователя из БД
		var userID int
		var passwordHash string
		var isVerified bool
		err := db.QueryRow(`
			SELECT id, password_hash, is_verified FROM users WHERE email = $1
		`, req.Email).Scan(&userID, &passwordHash, &isVerified)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка авторизации"})
			return
		}

		// Проверяем пароль
		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		}

		// Проверяем подтвержден ли email
		if !isVerified {
			c.JSON(http.StatusForbidden, gin.H{"error": "Email не подтвержден"})
			return
		}

		// Генерируем новый 2FA-код
		code := utils.GenerateCode()

		// Сначала удаляем старые коды, если они есть
		_, err = db.Exec(`DELETE FROM two_factor_codes WHERE user_id = $1`, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке аутентификации"})
			return
		}

		// Сохраняем код в БД с 10-минутным сроком действия
		_, err = db.Exec(`
			INSERT INTO two_factor_codes (user_id, code, expires_at) VALUES ($1, $2, $3)
		`, userID, code, time.Now().Add(10*time.Minute))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации кода"})
			return
		}

		// Отправляем код пользователю
		err = utils.SendEmail(req.Email, "Код для входа", "Ваш код: "+code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправки email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Код подтверждения отправлен"})
	}
}
