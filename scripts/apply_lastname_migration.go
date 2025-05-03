package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

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

	// Читаем SQL файл
	sqlContent, err := ioutil.ReadFile("migrations/add_lastname.sql")
	if err != nil {
		log.Fatalf("Ошибка чтения файла миграции: %v", err)
	}

	// Разделяем файл на отдельные SQL-команды
	commands := strings.Split(string(sqlContent), ";")

	for _, cmd := range commands {
		// Пропускаем пустые команды
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}

		// Выполняем каждую команду
		log.Printf("Выполнение команды: %s", cmd)
		_, err := db.Exec(cmd)
		if err != nil {
			log.Printf("Ошибка выполнения команды: %v", err)
		} else {
			log.Println("Команда выполнена успешно")
		}
	}

	log.Println("Миграция успешно применена")
}
