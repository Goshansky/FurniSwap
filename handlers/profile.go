package handlers

import (
	"FurniSwap/models"
	"FurniSwap/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetProfileHandler возвращает профиль пользователя
func GetProfileHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			log.Println("Ошибка: userID не найден в контексте")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
			return
		}

		log.Printf("Получение профиля для пользователя ID: %v (тип: %T)\n", userID, userID)

		// Конвертируем userID в int, если он не такого типа
		var userIDInt int
		switch v := userID.(type) {
		case int:
			userIDInt = v
		case float64:
			userIDInt = int(v)
		case float32:
			userIDInt = int(v)
		case int64:
			userIDInt = int(v)
		case int32:
			userIDInt = int(v)
		case string:
			// Try to parse the string as int if it's a string
			var err error
			_, err = fmt.Sscanf(v, "%d", &userIDInt)
			if err != nil {
				log.Printf("Не удалось преобразовать строку userID в int: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка сервера"})
				return
			}
		default:
			log.Printf("Неподдерживаемый тип userID: %T\n", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка сервера"})
			return
		}

		// Проверяем что userID > 0
		if userIDInt <= 0 {
			log.Printf("Недопустимый userID: %d\n", userIDInt)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "недопустимый ID пользователя"})
			return
		}

		var user models.User
		query := `
			SELECT id, email, name, last_name, city, avatar, is_verified, created_at
			FROM users
			WHERE id = $1
		`
		err := db.Get(&user, query, userIDInt)
		if err != nil {
			log.Printf("Ошибка при получении профиля пользователя %d: %v\n", userIDInt, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения профиля"})
			return
		}

		log.Printf("Профиль для пользователя ID: %d успешно получен\n", userIDInt)
		c.JSON(http.StatusOK, user)
	}
}

type UpdateProfileRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	City     string `json:"city"`
	Avatar   string `json:"avatar,omitempty"`
}

// UpdateProfileHandler обновляет профиль пользователя
func UpdateProfileHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			log.Println("Ошибка: userID не найден в контексте")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
			return
		}

		// Конвертируем userID в int, если он не такого типа
		var userIDInt int
		switch v := userID.(type) {
		case int:
			userIDInt = v
		case float64:
			userIDInt = int(v)
		case float32:
			userIDInt = int(v)
		case int64:
			userIDInt = int(v)
		case int32:
			userIDInt = int(v)
		case string:
			// Try to parse the string as int if it's a string
			var err error
			_, err = fmt.Sscanf(v, "%d", &userIDInt)
			if err != nil {
				log.Printf("Не удалось преобразовать строку userID в int: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка сервера"})
				return
			}
		default:
			log.Printf("Неподдерживаемый тип userID: %T\n", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка сервера"})
			return
		}

		var req UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные"})
			return
		}

		_, err := db.Exec(`
			UPDATE users 
			SET name = COALESCE(NULLIF($1, ''), name),
				last_name = COALESCE(NULLIF($2, ''), last_name),
				city = COALESCE(NULLIF($3, ''), city)
			WHERE id = $4
		`, req.Name, req.LastName, req.City, userIDInt)
		if err != nil {
			log.Printf("Ошибка при обновлении профиля пользователя %d: %v\n", userIDInt, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления профиля"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "профиль успешно обновлен"})
	}
}

// UploadAvatarHandler загружает аватар пользователя
func UploadAvatarHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			log.Println("Ошибка: userID не найден в контексте")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
			return
		}

		// Конвертируем userID в int, если он не такого типа
		var userIDInt int
		switch v := userID.(type) {
		case int:
			userIDInt = v
		case float64:
			userIDInt = int(v)
		case float32:
			userIDInt = int(v)
		case int64:
			userIDInt = int(v)
		case int32:
			userIDInt = int(v)
		case string:
			// Try to parse the string as int if it's a string
			var err error
			_, err = fmt.Sscanf(v, "%d", &userIDInt)
			if err != nil {
				log.Printf("Не удалось преобразовать строку userID в int: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка сервера"})
				return
			}
		default:
			log.Printf("Неподдерживаемый тип userID: %T\n", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка сервера"})
			return
		}

		// Получаем загружаемый файл
		file, err := c.FormFile("avatar")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "файл не найден"})
			return
		}

		// Сохраняем файл
		filePath, err := utils.UploadImage(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Обновляем аватар пользователя
		_, err = db.Exec("UPDATE users SET avatar = $1 WHERE id = $2", filePath, userIDInt)
		if err != nil {
			// Удаляем файл в случае ошибки
			utils.DeleteImage(filePath)
			log.Printf("Ошибка при обновлении аватара пользователя %d: %v\n", userIDInt, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления аватара"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "аватар успешно обновлен",
			"avatar":  filePath,
		})
	}
}
