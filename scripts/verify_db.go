package main

import (
	"FurniSwap/utils"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение к базе данных
	log.Println("Подключение к базе данных...")

	// Используем ту же строку подключения, что и в main.go
	connectionString := fmt.Sprintf("host=localhost port=5431 user=postgres password=password dbname=furni_swap sslmode=disable")
	log.Println("Подключение к базе данных:", connectionString)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		log.Fatalf("Проверка подключения к базе данных не удалась: %v", err)
	}
	log.Println("Успешное подключение к базе данных")

	// Проверяем структуру базы данных
	utils.CheckDatabase(db)

	// Создаем тестового пользователя, если он не существует
	utils.CreateTestUser(db)

	log.Println("Проверка базы данных завершена")
}
