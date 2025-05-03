package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateCode генерирует 6-значный код для проверки email и двухфакторной аутентификации
func GenerateCode() string {
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
