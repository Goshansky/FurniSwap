package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение к базе данных
	connectionString := fmt.Sprintf("host=localhost port=5431 user=postgres password=password dbname=furni_swap sslmode=disable")
	log.Println("Connecting to database:", connectionString)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Check if we can access the database
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM users")
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	log.Printf("Current user count: %d", count)

	// Get all users
	type User struct {
		ID       int    `db:"id"`
		Email    string `db:"email"`
		Name     string `db:"name"`
		LastName string `db:"last_name"`
	}

	var users []User
	err = db.Select(&users, "SELECT id, email, name, last_name FROM users")
	if err != nil {
		log.Fatalf("Failed to fetch users: %v", err)
	}

	// Update emails for non-English addresses
	updateCount := 0
	for _, user := range users {
		// Skip already English emails
		if strings.Contains(user.Email, "@example.com") && !strings.Contains(user.Email, "ё") && !strings.Contains(user.Email, "ж") {
			continue
		}

		// Create anglicized email
		name := convertToEnglish(user.Name)
		lastName := convertToEnglish(user.LastName)
		newEmail := fmt.Sprintf("%s.%s%d@example.com", name, lastName, user.ID)
		newEmail = strings.ToLower(newEmail)

		// Update the user
		_, err := db.Exec("UPDATE users SET email = $1 WHERE id = $2", newEmail, user.ID)
		if err != nil {
			log.Printf("Failed to update user %d: %v", user.ID, err)
			continue
		}
		updateCount++
	}
	log.Printf("Updated %d users with English emails", updateCount)

	// Add missing db tags to models.ChatResponse
	// This doesn't need a database update, as you've already fixed the model
	// in models/chat.go by adding the db:"user_name" tag.

	log.Println("Chat issue should be fixed now. Try accessing the chats again.")
}

// Simple transliteration for Russian to English
func convertToEnglish(input string) string {
	if input == "" {
		return "User"
	}

	// Simple mapping for common Russian names
	mapping := map[string]string{
		"Александр": "Alexander",
		"Дмитрий":   "Dmitry",
		"Максим":    "Maxim",
		"Сергей":    "Sergei",
		"Иван":      "Ivan",
		"Андрей":    "Andrey",
		"Алексей":   "Alexey",
		"Артём":     "Artem",
		"Михаил":    "Mikhail",
		"Никита":    "Nikita",
		"Анна":      "Anna",
		"Мария":     "Maria",
		"Екатерина": "Ekaterina",
		"Елена":     "Elena",
		"Ольга":     "Olga",
		"Наталья":   "Natalia",
		"Татьяна":   "Tatiana",
		"Юлия":      "Yulia",
		"Дарья":     "Daria",
		"Виктория":  "Victoria",
		"Иванов":    "Ivanov",
		"Петров":    "Petrov",
		"Сидоров":   "Sidorov",
		"Смирнов":   "Smirnov",
		"Кузнецов":  "Kuznetsov",
		"Попов":     "Popov",
		"Соколов":   "Sokolov",
		"Михайлов":  "Mikhailov",
		"Новиков":   "Novikov",
		"Федоров":   "Fedorov",
		"Морозов":   "Morozov",
		"Волков":    "Volkov",
		"Алексеев":  "Alekseev",
		"Лебедев":   "Lebedev",
		"Семенов":   "Semenov",
		"Егоров":    "Egorov",
		"Павлов":    "Pavlov",
		"Козлов":    "Kozlov",
		"Степанов":  "Stepanov",
		"Николаев":  "Nikolaev",
	}

	if translated, ok := mapping[input]; ok {
		return translated
	}

	// If not in mapping, use input as-is
	return input
}
