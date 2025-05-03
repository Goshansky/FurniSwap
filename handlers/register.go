package handlers

import (
	"FurniSwap/utils"
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
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
}

// RegisterHandler обрабатывает регистрацию пользователя
func RegisterHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
			return
		}

		// Проверяем, существует ли пользователь с таким email
		var exists bool
		err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки email"})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email уже существует"})
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
			INSERT INTO users (email, password_hash, name, last_name) 
			VALUES ($1, $2, $3, $4) RETURNING id
		`, req.Email, string(hashedPassword), req.Name, req.LastName).Scan(&userID)
		if err != nil {
			log.Printf("Ошибка сохранения пользователя: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения пользователя"})
			return
		}

		// Генерируем 6-значный код
		code := utils.GenerateCode()

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
