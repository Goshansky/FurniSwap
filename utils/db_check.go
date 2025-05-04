package utils

import (
	"log"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// CheckDatabase выполняет проверку базы данных и выводит информацию о таблицах
func CheckDatabase(db *sqlx.DB) {
	log.Println("Проверка базы данных...")

	// Проверка подключения
	if err := db.Ping(); err != nil {
		log.Fatalf("Невозможно подключиться к базе данных: %v", err)
	}
	log.Println("Подключение к базе данных успешно")

	// Проверка таблицы users
	var userCount int
	err := db.Get(&userCount, "SELECT COUNT(*) FROM users")
	if err != nil {
		log.Printf("Ошибка при проверке таблицы users: %v", err)
	} else {
		log.Printf("Количество пользователей в базе: %d", userCount)
	}

	// Проверка таблицы two_factor_codes
	var codesCount int
	err = db.Get(&codesCount, "SELECT COUNT(*) FROM two_factor_codes")
	if err != nil {
		log.Printf("Ошибка при проверке таблицы two_factor_codes: %v", err)
	} else {
		log.Printf("Количество кодов в базе: %d", codesCount)
	}

	// Проверка таблицы categories
	var categoriesCount int
	err = db.Get(&categoriesCount, "SELECT COUNT(*) FROM categories")
	if err != nil {
		log.Printf("Ошибка при проверке таблицы categories: %v", err)
	} else {
		log.Printf("Количество категорий в базе: %d", categoriesCount)
	}

	// Проверка таблицы purchases
	var purchasesCount int
	err = db.Get(&purchasesCount, "SELECT COUNT(*) FROM purchases")
	if err != nil {
		log.Printf("Таблица purchases не существует или другая ошибка: %v", err)
		// Проверяем, существует ли столбец status в таблице listings
		var hasStatusColumn bool
		err = db.Get(&hasStatusColumn, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'listings' AND column_name = 'status'
			)
		`)
		if err != nil {
			log.Printf("Ошибка при проверке столбца status в listings: %v", err)
		} else if !hasStatusColumn {
			log.Println("Необходимо выполнить миграцию add_purchases.sql")
		}
	} else {
		log.Printf("Количество записей в таблице purchases: %d", purchasesCount)
	}

	// Проверка структуры таблицы users
	var userColumns []string
	err = db.Select(&userColumns, `
		SELECT column_name FROM information_schema.columns 
		WHERE table_name = 'users'
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Printf("Ошибка при проверке структуры таблицы users: %v", err)
	} else {
		log.Printf("Структура таблицы users: %v", userColumns)
	}

	// Если пользователей нет, создаем тестового пользователя
	if userCount == 0 {
		CreateTestUser(db)
	}
}

// CreateTestUser создает тестового пользователя в базе данных
func CreateTestUser(db *sqlx.DB) {
	log.Println("Создание тестового пользователя...")

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		return
	}

	// Проверяем, существует ли пользователь
	var exists bool
	err = db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", "test@example.com")
	if err != nil {
		log.Printf("Ошибка проверки существования пользователя: %v", err)
		return
	}

	if exists {
		log.Println("Тестовый пользователь уже существует")
		return
	}

	// Создаем пользователя
	_, err = db.Exec(`
		INSERT INTO users (email, password_hash, name, last_name, city, is_verified, created_at)
		VALUES ($1, $2, $3, $4, $5, true, NOW())
	`, "test@example.com", string(hashedPassword), "Тестовый", "Пользователь", "Москва")

	if err != nil {
		log.Printf("Ошибка создания тестового пользователя: %v", err)
		return
	}

	log.Println("Тестовый пользователь успешно создан")

	// Выводим информацию о созданном пользователе
	var userID int
	err = db.Get(&userID, "SELECT id FROM users WHERE email = $1", "test@example.com")
	if err != nil {
		log.Printf("Ошибка получения ID пользователя: %v", err)
		return
	}

	log.Printf("Тестовый пользователь ID: %d", userID)
}
