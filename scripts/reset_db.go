// Package main provides utilities for resetting and seeding the database
package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Проверяет наличие директории, если отсутствует - создает ее
func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

// Создает пустой файл если он не существует
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

// ResetDatabase performs a full database reset and seeding
func ResetDatabase() {
	// Убедимся, что директория uploads существует
	err := ensureDir("uploads")
	if err != nil {
		log.Fatalf("Ошибка при создании директории uploads: %v", err)
	}

	// Список аватаров, убедимся что все файлы существуют
	avatarPaths := []string{
		"uploads/avatar1.jpg",
		"uploads/avatar2.jpg",
		"uploads/avatar3.jpg",
		"uploads/avatar4.jpg",
		"uploads/avatar5.jpg",
		"uploads/avatar6.jpg",
		"uploads/avatar7.jpg",
		"uploads/avatar8.jpg",
	}

	// Убедимся, что все файлы аватаров существуют
	for _, path := range avatarPaths {
		// Исправляем путь для файловой системы (без начального слеша)
		osPath := filepath.FromSlash(path)
		err := ensureFile(osPath)
		if err != nil {
			log.Printf("Предупреждение: не удалось создать файл %s: %v", osPath, err)
		}
	}

	// Подключение к базе данных
	connectionString := fmt.Sprintf("host=localhost port=5431 user=postgres password=password dbname=furni_swap sslmode=disable")
	log.Println("Подключение к базе данных:", connectionString)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Начинаем транзакцию
	tx, err := db.Beginx()
	if err != nil {
		log.Fatalf("Ошибка при создании транзакции: %v", err)
	}
	defer tx.Rollback()

	// Удаляем все данные из таблиц
	log.Println("Удаление существующих данных...")

	// Сначала удаляем данные из таблиц с внешними ключами
	_, err = tx.Exec("DELETE FROM messages")
	if err != nil {
		log.Printf("Ошибка при удалении сообщений: %v", err)
	}

	_, err = tx.Exec("DELETE FROM chats")
	if err != nil {
		log.Printf("Ошибка при удалении чатов: %v", err)
	}

	_, err = tx.Exec("DELETE FROM favorites")
	if err != nil {
		log.Printf("Ошибка при удалении избранного: %v", err)
	}

	_, err = tx.Exec("DELETE FROM listing_images")
	if err != nil {
		log.Printf("Ошибка при удалении изображений: %v", err)
	}

	_, err = tx.Exec("DELETE FROM listings")
	if err != nil {
		log.Printf("Ошибка при удалении объявлений: %v", err)
	}

	_, err = tx.Exec("DELETE FROM two_factor_codes")
	if err != nil {
		log.Printf("Ошибка при удалении кодов двухфакторной аутентификации: %v", err)
	}

	_, err = tx.Exec("DELETE FROM users")
	if err != nil {
		log.Printf("Ошибка при удалении пользователей: %v", err)
	}

	log.Println("Все существующие данные удалены")

	// Настройки для генерации данных
	userCount := 20
	favoritesCount := 30
	chatsCount := 15

	// Данные для генерации случайных значений
	cities := []string{"Москва", "Санкт-Петербург", "Казань", "Новосибирск", "Екатеринбург", "Нижний Новгород", "Самара"}
	conditions := []string{"новое", "хорошее", "среднее", "плохое"}

	// Мебель по категориям
	furnitureByCategory := map[int][]string{
		1: {"Диван угловой", "Диван прямой", "Кресло раскладное", "Кресло-кровать", "Пуфик мягкий", "Кресло-качалка", "Диван-кровать"},
		2: {"Стол обеденный", "Стол журнальный", "Стул деревянный", "Стул мягкий", "Стол письменный", "Табурет кухонный", "Стол компьютерный"},
		3: {"Шкаф-купе", "Комод с ящиками", "Тумба под ТВ", "Шкаф для одежды", "Пенал кухонный", "Книжный шкаф", "Шкаф в прихожую"},
		4: {"Кровать двуспальная", "Кровать односпальная", "Матрас ортопедический", "Кровать с ящиками", "Детская кровать", "Раскладушка", "Матрас пружинный"},
		5: {"Стеллаж", "Тумбочка прикроватная", "Вешалка напольная", "Полка настенная", "Зеркало напольное", "Подставка для цветов", "Этажерка"},
	}

	// Имена и фамилии для генерации пользователей
	firstNames := []string{"Alexander", "Dmitry", "Maxim", "Sergei", "Ivan", "Andrey", "Alexey", "Artem", "Mikhail", "Nikita", "Anna", "Maria", "Ekaterina", "Elena", "Olga", "Natalia", "Tatiana", "Yulia", "Daria", "Victoria"}
	lastNames := []string{"Ivanov", "Petrov", "Sidorov", "Smirnov", "Kuznetsov", "Popov", "Sokolov", "Mikhailov", "Novikov", "Fedorov", "Morozov", "Volkov", "Alekseev", "Lebedev", "Semenov", "Egorov", "Pavlov", "Kozlov", "Stepanov", "Nikolaev"}

	// Пути к аватарам для пользователей с веб-путями (с начальным слешем)
	webAvatarPaths := []string{
		"/uploads/avatar1.jpg",
		"/uploads/avatar2.jpg",
		"/uploads/avatar3.jpg",
		"/uploads/avatar4.jpg",
		"/uploads/avatar5.jpg",
		"/uploads/avatar6.jpg",
		"/uploads/avatar7.jpg",
		"/uploads/avatar8.jpg",
	}

	// Описания для объявлений
	descriptions := []string{
		"В отличном состоянии, использовался мало. Продаю в связи с переездом.",
		"Состояние почти новое, без дефектов. Самовывоз из района %s.",
		"Удобная модель, практичная. Есть небольшие потертости. Торг уместен.",
		"Качественная мебель от известного производителя. Гарантия качества.",
		"Продаю срочно, в связи с ремонтом. Цена договорная, звоните.",
		"В использовании был меньше года. Никаких дефектов нет.",
		"Стильный дизайн, подойдет под любой интерьер. Доставка возможна.",
		"Функциональная модель с множеством полезных деталей.",
		"Продается в связи с обновлением интерьера. Качество отличное.",
		"Приобретался в фирменном магазине. Есть все документы и чеки.",
	}

	// Генератор случайных данных
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Вставляем пользователей
	users := make([]map[string]interface{}, userCount)
	for i := 0; i < userCount; i++ {
		name := firstNames[r.Intn(len(firstNames))]
		lastName := lastNames[r.Intn(len(lastNames))]
		email := fmt.Sprintf("%s.%s%d@example.com", name, lastName, i)
		email = strings.ToLower(email) // Make email lowercase
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		users[i] = map[string]interface{}{
			"email":         email,
			"password_hash": string(hashedPassword),
			"name":          name,
			"last_name":     lastName,
			"city":          cities[r.Intn(len(cities))],
			"avatar":        webAvatarPaths[r.Intn(len(webAvatarPaths))],
			"is_verified":   true,
		}
	}

	// Вставляем пользователей и получаем их ID
	userIDs := make([]int, 0, len(users))
	for _, user := range users {
		var id int
		err := tx.QueryRow(`
			INSERT INTO users (email, password_hash, name, last_name, city, avatar, is_verified, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW() - INTERVAL '1 day' * $8)
			RETURNING id
		`, user["email"], user["password_hash"], user["name"], user["last_name"],
			user["city"], user["avatar"], user["is_verified"], r.Intn(60)).Scan(&id)

		if err != nil {
			log.Fatalf("Ошибка при вставке пользователя: %v", err)
		}

		userIDs = append(userIDs, id)
	}
	log.Printf("Вставлено %d пользователей", len(userIDs))

	// Вставляем тестового пользователя с известными данными
	testPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	var testUserID int
	err = tx.QueryRow(`
		INSERT INTO users (email, password_hash, name, last_name, city, avatar, is_verified, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, true, NOW())
		RETURNING id
	`, "test@example.com", string(testPassword), "Test", "User", "Moscow",
		webAvatarPaths[r.Intn(len(webAvatarPaths))]).Scan(&testUserID)

	if err != nil {
		log.Printf("Ошибка при создании тестового пользователя: %v", err)
	} else {
		log.Printf("Создан тестовый пользователь с ID: %d", testUserID)
		userIDs = append(userIDs, testUserID)
	}

	// Проверяем наличие категорий
	var categoryCount int
	err = tx.Get(&categoryCount, "SELECT COUNT(*) FROM categories")
	if err != nil {
		log.Fatalf("Ошибка при проверке категорий: %v", err)
	}

	if categoryCount == 0 {
		// Категории отсутствуют, вставляем их
		categories := []string{
			"Диваны и кресла",
			"Столы и стулья",
			"Шкафы и комоды",
			"Кровати и матрасы",
			"Другое",
		}
		for i, cat := range categories {
			_, err = tx.Exec("INSERT INTO categories (id, name) VALUES ($1, $2)", i+1, cat)
			if err != nil {
				log.Fatalf("Ошибка при вставке категории %s: %v", cat, err)
			}
		}
		log.Printf("Вставлено %d категорий", len(categories))
	} else {
		log.Printf("Категории уже существуют, пропускаем вставку")
	}

	// Вставляем объявления
	var totalListings int
	for _, userID := range userIDs {
		listingCount := r.Intn(5) + 1 // От 1 до 5 объявлений на пользователя
		for i := 0; i < listingCount; i++ {
			categoryID := r.Intn(5) + 1 // От 1 до 5 категорий

			// Получаем случайное название из категории
			furnitureItems := furnitureByCategory[categoryID]
			title := furnitureItems[r.Intn(len(furnitureItems))]

			// Получаем случайное описание
			description := descriptions[r.Intn(len(descriptions))]

			// Если в описании есть плейсхолдер для города, заменяем его
			if fmt.Sprintf(description, "test") != description {
				var city string
				err := tx.Get(&city, "SELECT city FROM users WHERE id = $1", userID)
				if err != nil {
					city = cities[r.Intn(len(cities))]
				}
				description = fmt.Sprintf(description, city)
			}

			price := float64(r.Intn(15000) + 1000) // Цена от 1000 до 16000
			condition := conditions[r.Intn(len(conditions))]
			city := cities[r.Intn(len(cities))]

			var listingID int
			err := tx.QueryRow(`
				INSERT INTO listings (user_id, title, description, price, condition, city, 
								     category_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, 
						NOW() - INTERVAL '1 day' * $8, 
						NOW() - INTERVAL '1 day' * $8)
				RETURNING id
			`, userID, title, description, price, condition, city, categoryID, r.Intn(30)).Scan(&listingID)

			if err != nil {
				log.Fatalf("Ошибка при вставке объявления: %v", err)
			}

			// Вставляем изображения для объявления
			imageCount := r.Intn(3) + 1 // От 1 до 3 изображений на объявление
			for j := 0; j < imageCount; j++ {
				isMain := j == 0 // Первое изображение делаем главным

				// Генерируем случайное имя файла для изображения
				imageType := r.Intn(3)
				var imageDir string
				if imageType == 0 {
					imageDir = "sofa"
				} else if imageType == 1 {
					imageDir = "table"
				} else {
					imageDir = "chair"
				}

				imagePath := fmt.Sprintf("/uploads/%s%d.jpg", imageDir, r.Intn(5)+1)

				_, err := tx.Exec(`
					INSERT INTO listing_images (listing_id, image_path, is_main, created_at)
					VALUES ($1, $2, $3, NOW() - INTERVAL '1 hour' * $4)
				`, listingID, imagePath, isMain, r.Intn(24))

				if err != nil {
					log.Fatalf("Ошибка при вставке изображения: %v", err)
				}
			}

			totalListings++
		}
	}
	log.Printf("Вставлено %d объявлений с изображениями", totalListings)

	// Вставляем избранное
	for i := 0; i < favoritesCount; i++ {
		userID := userIDs[r.Intn(len(userIDs))]

		// Получаем случайное объявление, не принадлежащее этому пользователю
		var listingID int
		err = tx.Get(&listingID, `
			SELECT id FROM listings 
			WHERE user_id != $1 
			ORDER BY RANDOM() 
			LIMIT 1
		`, userID)
		if err != nil {
			log.Printf("Не удалось найти объявление для избранного: %v", err)
			continue
		}

		// Проверяем, не добавлено ли уже это объявление в избранное данным пользователем
		var exists bool
		err = tx.Get(&exists, `
			SELECT EXISTS(
				SELECT 1 FROM favorites 
				WHERE user_id = $1 AND listing_id = $2
			)
		`, userID, listingID)
		if err != nil {
			log.Printf("Ошибка при проверке избранного: %v", err)
			continue
		}

		if !exists {
			_, err = tx.Exec(`
				INSERT INTO favorites (user_id, listing_id, created_at)
				VALUES ($1, $2, NOW())
			`, userID, listingID)
			if err != nil {
				log.Printf("Ошибка при вставке избранного: %v", err)
			}
		}
	}

	// Считаем сколько вставили избранных
	var favoritesInserted int
	err = tx.Get(&favoritesInserted, "SELECT COUNT(*) FROM favorites")
	if err != nil {
		log.Printf("Ошибка при подсчете избранных: %v", err)
	} else {
		log.Printf("Вставлено %d записей избранного", favoritesInserted)
	}

	// Вставляем чаты и сообщения
	for i := 0; i < chatsCount; i++ {
		// Выбираем случайное объявление
		var listing struct {
			ID     int `db:"id"`
			UserID int `db:"user_id"` // Продавец
		}
		err = tx.Get(&listing, "SELECT id, user_id FROM listings ORDER BY RANDOM() LIMIT 1")
		if err != nil {
			log.Printf("Ошибка при выборе объявления для чата: %v", err)
			continue
		}

		// Выбираем случайного покупателя (не владельца объявления)
		var buyerID int
		for {
			candidateID := userIDs[r.Intn(len(userIDs))]
			if candidateID != listing.UserID {
				buyerID = candidateID
				break
			}
		}

		// Проверяем, не существует ли уже чат между этими пользователями по этому объявлению
		var exists bool
		err = tx.Get(&exists, `
			SELECT EXISTS(
				SELECT 1 FROM chats 
				WHERE listing_id = $1 AND buyer_id = $2 AND seller_id = $3
			)
		`, listing.ID, buyerID, listing.UserID)
		if err != nil {
			log.Printf("Ошибка при проверке чата: %v", err)
			continue
		}

		if exists {
			continue
		}

		// Вставляем чат
		var chatID int
		err = tx.QueryRow(`
			INSERT INTO chats (listing_id, buyer_id, seller_id, created_at)
			VALUES ($1, $2, $3, NOW())
			RETURNING id
		`, listing.ID, buyerID, listing.UserID).Scan(&chatID)
		if err != nil {
			log.Printf("Ошибка при создании чата: %v", err)
			continue
		}

		// Вставляем сообщения в чат
		messageCount := r.Intn(10) + 1 // От 1 до 10 сообщений
		for j := 0; j < messageCount; j++ {
			// Определяем отправителя (чередуем покупателя и продавца)
			senderID := buyerID
			if j%2 == 1 {
				senderID = listing.UserID
			}

			// Генерируем текст сообщения
			var content string
			if j == 0 && senderID == buyerID {
				// Первое сообщение от покупателя
				firstMessages := []string{
					"Здравствуйте, это объявление еще актуально?",
					"Добрый день! Мебель еще продается?",
					"Приветствую! Можно узнать подробнее об этом предложении?",
					"Здравствуйте! Интересуюсь вашим объявлением. Все еще доступно?",
					"Добрый день, мебель еще доступна для покупки?",
				}
				content = firstMessages[r.Intn(len(firstMessages))]
			} else if j == 1 && senderID == listing.UserID {
				// Первый ответ продавца
				sellerResponses := []string{
					"Добрый день! Да, объявление актуально.",
					"Здравствуйте! Да, всё ещё продаётся.",
					"Приветствую! Да, мебель в наличии.",
					"Добрый день! Конечно, что именно вас интересует?",
					"Здравствуйте! Да, всё в наличии. Есть вопросы?",
				}
				content = sellerResponses[r.Intn(len(sellerResponses))]
			} else if senderID == buyerID {
				// Последующие сообщения покупателя
				buyerMessages := []string{
					"Можно договориться о скидке?",
					"Какие у вас есть варианты доставки?",
					"Можно посмотреть вживую перед покупкой?",
					"В какое время можно подъехать посмотреть?",
					"Есть какие-нибудь дефекты?",
					"На выходных смогу забрать, если договоримся.",
					"А торг уместен?",
					"Можно еще фото с других ракурсов?",
					"Давно у вас эта мебель?",
					"Как насчет 10% скидки?",
				}
				content = buyerMessages[r.Intn(len(buyerMessages))]
			} else {
				// Последующие сообщения продавца
				sellerMessages := []string{
					"Да, можно, но не больше 5%.",
					"Могу организовать доставку за дополнительную плату.",
					"Конечно, можете приехать и посмотреть. Когда вам удобно?",
					"Могу встретиться в любой день после 18:00.",
					"Нет, никаких дефектов, всё в отличном состоянии.",
					"Хорошо, буду ждать вас на выходных.",
					"Небольшой торг возможен при осмотре.",
					"Хорошо, вечером пришлю дополнительные фото.",
					"Мебель у меня около двух лет, состояние как новое.",
					"Давайте встретимся посередине, скидка 7%?",
				}
				content = sellerMessages[r.Intn(len(sellerMessages))]
			}

			// Вставляем сообщение
			isRead := r.Float32() < 0.7 // 70% сообщений прочитаны
			_, err = tx.Exec(`
				INSERT INTO messages (chat_id, user_id, content, created_at, is_read)
				VALUES ($1, $2, $3, NOW() - INTERVAL '1 minute' * $4, $5)
			`, chatID, senderID, content, messageCount-j, isRead)
			if err != nil {
				log.Printf("Ошибка при вставке сообщения: %v", err)
			}
		}
	}

	// Подсчет вставленных чатов и сообщений
	var chatsInserted, messagesCount int
	tx.Get(&chatsInserted, "SELECT COUNT(*) FROM chats")
	tx.Get(&messagesCount, "SELECT COUNT(*) FROM messages")
	log.Printf("Вставлено %d чатов и %d сообщений", chatsInserted, messagesCount)

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Ошибка при коммите транзакции: %v", err)
	}

	log.Println("База данных успешно обновлена")
	log.Println("ВНИМАНИЕ: У всех пользователей установлен пароль 'password123'")

	// Создадим папки и файлы для изображений мебели при необходимости
	furnitureDirs := []string{"sofa", "table", "chair"}
	for _, dir := range furnitureDirs {
		dirPath := filepath.Join("uploads", dir)
		err := ensureDir(dirPath)
		if err != nil {
			log.Printf("Предупреждение: не удалось создать директорию %s: %v", dirPath, err)
			continue
		}

		// Создадим 5 файлов изображений в каждой папке
		for i := 1; i <= 5; i++ {
			imgPath := filepath.Join(dirPath, fmt.Sprintf("%d.jpg", i))
			err := ensureFile(imgPath)
			if err != nil {
				log.Printf("Предупреждение: не удалось создать файл %s: %v", imgPath, err)
			}
		}
	}

	log.Println("Все необходимые файлы и папки для изображений созданы")
}

// Run executes the database reset
func main() {
	ResetDatabase()
}
