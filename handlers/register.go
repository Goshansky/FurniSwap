package handlers

import (
	"FurniSwap/utils"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterHandler обрабатывает регистрацию пользователя
func RegisterHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
			return
		}

		// Хэшируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки пароля"})
			return
		}

		// Сохраняем пользователя в БД
		var userID int
		err = db.QueryRow(`
			INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id
		`, req.Email, string(hashedPassword)).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения пользователя"})
			return
		}

		// Генерируем 6-значный код
		code := generateCode()

		// Сохраняем код в БД с 10-минутным сроком действия
		_, err = db.Exec(`
			INSERT INTO two_factor_codes (user_id, code, expires_at) VALUES ($1, $2, $3)
		`, userID, code, time.Now().Add(10*time.Minute))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации кода"})
			return
		}

		// Отправляем код на email
		err = utils.SendEmail(req.Email, "Код подтверждения", "Ваш код: "+code)
		if err != nil {
			log.Println("Ошибка отправки email:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправки email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Код подтверждения отправлен"})
	}
}

// generateCode генерирует 6-значный код
func generateCode() string {
	b := make([]byte, 6) // 6 байтов для кода
	_, err := rand.Read(b)
	if err != nil {
		return "000000" // На случай ошибки, вернуть фиксированный код
	}

	code := base64.StdEncoding.EncodeToString(b)

	// Убеждаемся, что код >= 6 символов
	if len(code) < 6 {
		return code
	}

	return code[:6]
}
