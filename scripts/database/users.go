package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// CheckLastNameColumn проверяет наличие колонки last_name в таблице users и добавляет ее при необходимости
func CheckLastNameColumn(db *sqlx.DB) error {
	// Проверяем структуру таблицы users
	var columns []string
	err := db.Select(&columns, `
		SELECT column_name FROM information_schema.columns 
		WHERE table_name = 'users'
		ORDER BY ordinal_position
	`)
	if err != nil {
		return err
	}

	log.Printf("Текущие столбцы таблицы users: %v", columns)

	// Проверяем, существует ли столбец last_name
	var lastNameExists bool
	for _, column := range columns {
		if column == "last_name" {
			lastNameExists = true
			break
		}
	}

	if !lastNameExists {
		log.Println("Столбец last_name не найден в таблице users. Добавляю...")

		// Добавляем столбец last_name, если его нет
		_, err = db.Exec("ALTER TABLE users ADD COLUMN last_name TEXT")
		if err != nil {
			return err
		}

		// Обновляем существующие записи
		_, err = db.Exec("UPDATE users SET last_name = '' WHERE last_name IS NULL")
		if err != nil {
			return err
		}

		log.Println("Столбец last_name успешно добавлен в таблицу users")
	} else {
		log.Println("Столбец last_name уже существует в таблице users")
	}

	return nil
}

// EnsureTestUser создает или обновляет тестового пользователя
func EnsureTestUser(db *sqlx.DB) error {
	// Проверяем наличие тестового пользователя
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = 'test@example.com')")
	if err != nil {
		return err
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if exists {
		// Обновляем пользователя
		_, err = db.Exec(`
			UPDATE users 
			SET name = $1, last_name = $2, password_hash = $3, is_verified = true
			WHERE email = $4
		`, "Тестовый", "Пользователь", string(hashedPassword), "test@example.com")
	} else {
		// Создаем пользователя
		_, err = db.Exec(`
			INSERT INTO users (email, password_hash, name, last_name, city, is_verified, created_at)
			VALUES ($1, $2, $3, $4, $5, true, NOW())
		`, "test@example.com", string(hashedPassword), "Тестовый", "Пользователь", "Москва")
	}

	if err != nil {
		return err
	}

	log.Println("Тестовый пользователь успешно создан/обновлен")
	return nil
}
