package middlewares

import (
	"FurniSwap/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// AuthRequired проверяет JWT токен и авторизует пользователя
func AuthRequired(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Отсутствует токен авторизации в запросе")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "отсутствует токен авторизации"})
			c.Abort()
			return
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Неверный формат токена:", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный формат токена авторизации"})
			c.Abort()
			return
		}

		token := parts[1]
		log.Printf("Длина JWT токена: %d символов\n", len(token))

		// Проверка, не является ли токен пустым
		if token == "" {
			log.Println("Ошибка: пустой JWT токен")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "пустой токен авторизации"})
			c.Abort()
			return
		}

		// Валидируем токен и получаем ID пользователя
		userID, err := utils.ValidateToken(token)
		if err != nil {
			log.Printf("Ошибка валидации токена: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
			c.Abort()
			return
		}

		// Проверяем, что userID положительный
		if userID <= 0 {
			log.Printf("Недопустимый userID из токена: %d\n", userID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
			c.Abort()
			return
		}

		log.Printf("Токен успешно валидирован для пользователя ID: %d\n", userID)

		// Проверяем, существует ли пользователь в БД
		var exists bool
		err = db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND is_verified = true)", userID)
		if err != nil {
			log.Printf("Ошибка проверки пользователя в БД: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при проверке пользователя"})
			c.Abort()
			return
		}

		if !exists {
			log.Printf("Пользователь ID: %d не найден или не подтвержден\n", userID)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не найден или не подтвержден"})
			c.Abort()
			return
		}

		// Устанавливаем ID пользователя в контекст для дальнейшего использования
		c.Set("userID", userID)
		log.Printf("Пользователь ID: %d успешно аутентифицирован\n", userID)
		c.Next()
	}
}
