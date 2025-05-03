package main

import (
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

	// Обновляем фамилию для тестового пользователя
	_, err = db.Exec("UPDATE users SET last_name = 'Пользователь' WHERE email = 'test@example.com'")
	if err != nil {
		log.Fatalf("Ошибка при обновлении фамилии: %v", err)
	}

	// Проверяем результат
	var name, lastName string
	err = db.QueryRow("SELECT name, last_name FROM users WHERE email = 'test@example.com'").Scan(&name, &lastName)
	if err != nil {
		log.Fatalf("Ошибка при получении данных пользователя: %v", err)
	}

	log.Printf("Пользователь обновлен: %s %s", name, lastName)
}
