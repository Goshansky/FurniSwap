package main

import (
	"FurniSwap/scripts/database"
	"log"
	"os"
	"path/filepath"
)

// ensureDir проверяет наличие директории, если отсутствует - создает ее
func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// ensureFile создает пустой файл если он не существует
func ensureFile(fileName string) error {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		log.Printf("Создан файл: %s", fileName)
	}
	return nil
}

func main() {
	log.Println("Запуск сброса и заполнения базы данных...")

	// Убедимся, что директория uploads существует
	err := ensureDir("uploads")
	if err != nil {
		log.Fatalf("Ошибка при создании директории uploads: %v", err)
	}

	// Убедимся, что все файлы аватаров существуют
	for _, path := range database.WebAvatarPaths {
		// Убираем начальный слеш для файловой системы
		localPath := filepath.Join(".", path[1:])
		err := ensureFile(localPath)
		if err != nil {
			log.Printf("Предупреждение: не удалось создать файл %s: %v", localPath, err)
		}
	}

	// Подключение к базе данных через общую утилиту
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Начинаем транзакцию
	tx, err := db.Beginx()
	if err != nil {
		log.Fatalf("Ошибка при создании транзакции: %v", err)
	}
	defer tx.Rollback()

	// Проверка и добавление колонки last_name
	err = database.CheckLastNameColumn(db)
	if err != nil {
		log.Fatalf("Ошибка при работе с колонкой last_name: %v", err)
	}

	// Очищаем базу данных
	log.Println("Очистка базы данных...")
	err = database.CleanDatabase(tx)
	if err != nil {
		log.Fatalf("Ошибка при очистке базы данных: %v", err)
	}

	// Проверка наличия категорий
	err = database.EnsureCategories(tx)
	if err != nil {
		log.Fatalf("Ошибка при проверке категорий: %v", err)
	}

	// Генерация и вставка пользователей
	log.Println("Генерация пользователей...")
	users := database.GenerateUsers(database.UserCount, true) // true - использовать английские имена
	userIDs, err := database.InsertUsers(tx, users)
	if err != nil {
		log.Fatalf("Ошибка при вставке пользователей: %v", err)
	}
	log.Printf("Вставлено %d пользователей", len(userIDs))

	// Создание/обновление тестового пользователя
	log.Println("Создание тестового пользователя...")
	var testUserID int
	err = tx.QueryRow(`
		INSERT INTO users (email, password_hash, name, last_name, city, is_verified, created_at)
		VALUES ($1, $2, $3, $4, $5, true, NOW())
		RETURNING id
	`, "test@example.com", "$2a$10$ZKqp8jEG9S.Di4VaKTobNeQbrK1GrPYWiAPZZYJEOK7jEuQgDGIG2", "Тестовый", "Пользователь", "Москва").Scan(&testUserID)
	if err != nil {
		log.Printf("Ошибка при создании тестового пользователя: %v", err)
	} else {
		log.Printf("Создан тестовый пользователь с ID: %d", testUserID)
		userIDs = append(userIDs, testUserID)
	}

	// Сохраняем транзакцию
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Ошибка при сохранении транзакции: %v", err)
	}

	log.Println("База данных успешно сброшена и заполнена")
}
