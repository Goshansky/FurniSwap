package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWT claims struct
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken генерирует новый JWT токен для пользователя
func GenerateToken(userID int) (string, error) {
	// Получение секретного ключа из переменных окружения
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Println("ВНИМАНИЕ: JWT_SECRET_KEY не настроен, используется значение по умолчанию")
		secretKey = "default_secret_key_change_in_production" // Запасной вариант для разработки
	}

	// Создание новых claims с ID пользователя
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен действителен 24 часа
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Создание токена с указанными claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("ошибка подписи токена: %w", err)
	}

	return tokenString, nil
}

// ValidateToken проверяет валидность токена и возвращает ID пользователя
func ValidateToken(tokenString string) (int, error) {
	// Получение секретного ключа из переменных окружения
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Println("ВНИМАНИЕ: JWT_SECRET_KEY не настроен, используется значение по умолчанию")
		secretKey = "default_secret_key_change_in_production" // Запасной вариант для разработки
	}

	// Парсинг и валидация токена
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, используется ли правильный алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи токена: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга токена: %w", err)
	}

	// Проверяем, прошел ли токен валидацию
	if !token.Valid {
		return 0, errors.New("недействительный токен")
	}

	// Получаем claims из токена
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return 0, errors.New("невозможно получить claims из токена")
	}

	return claims.UserID, nil
}
