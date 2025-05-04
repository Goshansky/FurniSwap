package main

import (
	"FurniSwap/scripts/database"
	"fmt"
	"log"
	"strings"
)

// Карта соответствия кириллических и латинских символов
var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh",
	'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o",
	'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts",
	'ч': "ch", 'ш': "sh", 'щ': "sch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu",
	'я': "ya",
	'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo", 'Ж': "Zh",
	'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M", 'Н': "N", 'О': "O",
	'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U", 'Ф': "F", 'Х': "Kh", 'Ц': "Ts",
	'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch", 'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu",
	'Я': "Ya",
}

// Функция транслитерации строки из кириллицы в латиницу
func transliterate(input string) string {
	var result strings.Builder

	for _, r := range input {
		if replacement, ok := translitMap[r]; ok {
			result.WriteString(replacement)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

func main() {
	// Подключение к базе данных через общую утилиту
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Получаем всех пользователей
	type User struct {
		ID       int    `db:"id"`
		Email    string `db:"email"`
		Name     string `db:"name"`
		LastName string `db:"last_name"`
	}

	var users []User
	err = db.Select(&users, "SELECT id, email, name, last_name FROM users")
	if err != nil {
		log.Fatalf("Ошибка получения пользователей: %v", err)
	}

	// Обновляем email-адреса, содержащие кириллицу
	updateCount := 0
	for _, user := range users {
		// Проверяем, содержит ли email кириллические символы
		needsUpdate := false
		for _, r := range user.Email {
			if _, ok := translitMap[r]; ok {
				needsUpdate = true
				break
			}
		}

		if !needsUpdate {
			continue
		}

		// Создаем новый email с транслитерацией
		emailParts := strings.Split(user.Email, "@")
		if len(emailParts) != 2 {
			log.Printf("Некорректный формат email для пользователя %d: %s", user.ID, user.Email)
			continue
		}

		localPart := transliterate(emailParts[0])
		domain := emailParts[1]
		newEmail := fmt.Sprintf("%s@%s", localPart, domain)
		newEmail = strings.ToLower(newEmail)

		// Обновляем пользователя
		_, err := db.Exec("UPDATE users SET email = $1 WHERE id = $2", newEmail, user.ID)
		if err != nil {
			log.Printf("Ошибка обновления email для пользователя %d: %v", user.ID, err)
			continue
		}

		log.Printf("Обновлен email для пользователя %d: %s -> %s", user.ID, user.Email, newEmail)
		updateCount++
	}

	log.Printf("Обновлено %d email-адресов", updateCount)
}
