package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение к базе данных
	connectionString := fmt.Sprintf("host=localhost port=5431 user=postgres password=password dbname=furni_swap sslmode=disable")
	log.Println("Подключение к базе данных:", connectionString)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Проверяем структуру таблицы users
	var columns []string
	err = db.Select(&columns, `
		SELECT column_name FROM information_schema.columns 
		WHERE table_name = 'users'
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Fatalf("Ошибка при проверке структуры таблицы: %v", err)
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
			log.Fatalf("Ошибка при добавлении столбца last_name: %v", err)
		}

		// Обновляем существующие записи
		_, err = db.Exec("UPDATE users SET last_name = '' WHERE last_name IS NULL")
		if err != nil {
			log.Fatalf("Ошибка при обновлении существующих записей: %v", err)
		}

		log.Println("Столбец last_name успешно добавлен в таблицу users")
	} else {
		log.Println("Столбец last_name уже существует в таблице users")
	}

	// Проверяем наличие тестового пользователя
	var userExists bool
	err = db.Get(&userExists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = 'test@example.com')")
	if err != nil {
		log.Fatalf("Ошибка при проверке наличия тестового пользователя: %v", err)
	}

	if userExists {
		// Если пользователь существует, проверяем наличие фамилии
		var lastName sql.NullString
		err = db.Get(&lastName, "SELECT last_name FROM users WHERE email = 'test@example.com'")
		if err != nil {
			log.Fatalf("Ошибка при получении last_name: %v", err)
		}

		if !lastName.Valid || lastName.String == "" {
			// Обновляем фамилию
			_, err = db.Exec("UPDATE users SET last_name = 'Пользователь' WHERE email = 'test@example.com'")
			if err != nil {
				log.Fatalf("Ошибка при обновлении фамилии: %v", err)
			}
			log.Println("Фамилия тестового пользователя обновлена")
		} else {
			log.Printf("Тестовый пользователь уже имеет фамилию: %s", lastName.String)
		}
	}

	log.Println("Проверка и исправление БД завершены")
}
