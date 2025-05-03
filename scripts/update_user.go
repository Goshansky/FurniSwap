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

	// Обновляем данные для тестового пользователя (используем более безопасный метод)
	result, err := db.Exec("UPDATE users SET name = $1, last_name = $2 WHERE email = $3",
		"Тестовый", "Пользователь", "test@example.com")
	if err != nil {
		log.Fatalf("Ошибка при обновлении данных пользователя: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Обновлено строк: %d", rowsAffected)

	// Выводим все данные из таблицы users
	type User struct {
		ID         int    `db:"id"`
		Email      string `db:"email"`
		Name       string `db:"name"`
		LastName   string `db:"last_name"`
		IsVerified bool   `db:"is_verified"`
	}

	var users []User
	err = db.Select(&users, "SELECT id, email, name, last_name, is_verified FROM users")
	if err != nil {
		log.Fatalf("Ошибка при получении данных пользователей: %v", err)
	}

	log.Println("Пользователи в базе данных:")
	for _, user := range users {
		log.Printf("ID: %d, Email: %s, Имя: %s, Фамилия: %s, Подтвержден: %t",
			user.ID, user.Email, user.Name, user.LastName, user.IsVerified)
	}
}
