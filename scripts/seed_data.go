package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Настройки для генерации данных
const (
	UserCount        = 20 // Количество пользователей
	ListingsPerUser  = 3  // Среднее количество объявлений на пользователя
	ImagesPerListing = 2  // Среднее количество изображений на объявление
	FavoritesCount   = 30 // Количество избранных объявлений
	ChatsCount       = 15 // Количество чатов
	MessagesPerChat  = 5  // Среднее количество сообщений в чате
)

var (
	cities     = []string{"Москва", "Санкт-Петербург", "Казань", "Новосибирск", "Екатеринбург", "Нижний Новгород", "Самара"}
	conditions = []string{"новое", "хорошее", "среднее", "плохое"}

	// Мебель по категориям
	furnitureByCategory = map[int][]string{
		1: {"Диван угловой", "Диван прямой", "Кресло раскладное", "Кресло-кровать", "Пуфик мягкий", "Кресло-качалка", "Диван-кровать"},
		2: {"Стол обеденный", "Стол журнальный", "Стул деревянный", "Стул мягкий", "Стол письменный", "Табурет кухонный", "Стол компьютерный"},
		3: {"Шкаф-купе", "Комод с ящиками", "Тумба под ТВ", "Шкаф для одежды", "Пенал кухонный", "Книжный шкаф", "Шкаф в прихожую"},
		4: {"Кровать двуспальная", "Кровать односпальная", "Матрас ортопедический", "Кровать с ящиками", "Детская кровать", "Раскладушка", "Матрас пружинный"},
		5: {"Стеллаж", "Тумбочка прикроватная", "Вешалка напольная", "Полка настенная", "Зеркало напольное", "Подставка для цветов", "Этажерка"},
	}

	// Имена и фамилии для генерации пользователей
	firstNames = []string{"Александр", "Дмитрий", "Максим", "Сергей", "Иван", "Андрей", "Алексей", "Артём", "Михаил", "Никита", "Анна", "Мария", "Екатерина", "Елена", "Ольга", "Наталья", "Татьяна", "Юлия", "Дарья", "Виктория"}
	lastNames  = []string{"Иванов", "Петров", "Сидоров", "Смирнов", "Кузнецов", "Попов", "Соколов", "Михайлов", "Новиков", "Федоров", "Морозов", "Волков", "Алексеев", "Лебедев", "Семенов", "Егоров", "Павлов", "Козлов", "Степанов", "Николаев"}

	// Описания для объявлений
	descriptions = []string{
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
)

// Генератор случайных данных
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
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

	// Вставляем пользователей
	users := generateUsers(UserCount)
	userIDs, err := insertUsers(tx, users)
	if err != nil {
		log.Fatalf("Ошибка при вставке пользователей: %v", err)
	}
	log.Printf("Вставлено %d пользователей", len(userIDs))

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
			listingID, err := insertListing(tx, userID, categoryID)
			if err != nil {
				log.Fatalf("Ошибка при вставке объявления: %v", err)
			}

			// Вставляем изображения для объявления
			imageCount := r.Intn(3) + 1 // От 1 до 3 изображений на объявление
			for j := 0; j < imageCount; j++ {
				isMain := j == 0 // Первое изображение делаем главным
				err = insertImage(tx, listingID, isMain)
				if err != nil {
					log.Fatalf("Ошибка при вставке изображения: %v", err)
				}
			}
			totalListings++
		}
	}
	log.Printf("Вставлено %d объявлений с изображениями", totalListings)

	// Вставляем избранное
	for i := 0; i < FavoritesCount; i++ {
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
	var favoritesCount int
	err = tx.Get(&favoritesCount, "SELECT COUNT(*) FROM favorites")
	if err != nil {
		log.Printf("Ошибка при подсчете избранных: %v", err)
	} else {
		log.Printf("Вставлено %d записей избранного", favoritesCount)
	}

	// Вставляем чаты и сообщения
	for i := 0; i < ChatsCount; i++ {
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
			content := getRandomMessage(j, senderID == buyerID)

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
	var chatsCount, messagesCount int
	tx.Get(&chatsCount, "SELECT COUNT(*) FROM chats")
	tx.Get(&messagesCount, "SELECT COUNT(*) FROM messages")
	log.Printf("Вставлено %d чатов и %d сообщений", chatsCount, messagesCount)

	// Фиксируем транзакцию
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Ошибка при коммите транзакции: %v", err)
	}

	log.Println("Данные успешно загружены в базу данных")
}

// Функция для генерации пользователей
func generateUsers(count int) []map[string]interface{} {
	users := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		name := firstNames[r.Intn(len(firstNames))]
		lastName := lastNames[r.Intn(len(lastNames))]
		email := fmt.Sprintf("%s.%s%d@example.com", name, lastName, i)
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		users[i] = map[string]interface{}{
			"email":         email,
			"password_hash": string(hashedPassword),
			"name":          name,
			"last_name":     lastName,
			"city":          cities[r.Intn(len(cities))],
			"is_verified":   true,
		}
	}
	return users
}

// Функция для вставки пользователей и возврата их ID
func insertUsers(tx *sqlx.Tx, users []map[string]interface{}) ([]int, error) {
	userIDs := make([]int, 0, len(users))

	for _, user := range users {
		var id int
		err := tx.QueryRow(`
			INSERT INTO users (email, password_hash, name, last_name, city, is_verified, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW() - INTERVAL '1 day' * $7)
			RETURNING id
		`, user["email"], user["password_hash"], user["name"], user["last_name"],
			user["city"], user["is_verified"], r.Intn(60)).Scan(&id)

		if err != nil {
			// Проверяем, существует ли уже пользователь с таким email
			var exists bool
			tx.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user["email"])
			if exists {
				// Если пользователь уже существует, получаем его ID
				tx.Get(&id, "SELECT id FROM users WHERE email = $1", user["email"])
				userIDs = append(userIDs, id)
				continue
			}
			return nil, err
		}

		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}

// Функция для вставки объявления
func insertListing(tx *sqlx.Tx, userID, categoryID int) (int, error) {
	furnitureItems := furnitureByCategory[categoryID]
	title := furnitureItems[r.Intn(len(furnitureItems))]
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

	var id int
	err := tx.QueryRow(`
		INSERT INTO listings (user_id, title, description, price, condition, city, 
							 category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 
				NOW() - INTERVAL '1 day' * $8, 
				NOW() - INTERVAL '1 day' * $8)
		RETURNING id
	`, userID, title, description, price, condition, city, categoryID, r.Intn(30)).Scan(&id)

	return id, err
}

// Функция для вставки изображения к объявлению
func insertImage(tx *sqlx.Tx, listingID int, isMain bool) error {
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

	return err
}

// Функция для генерации случайного сообщения в чате
func getRandomMessage(position int, isBuyer bool) string {
	if position == 0 && isBuyer {
		// Первое сообщение от покупателя
		firstMessages := []string{
			"Здравствуйте, это объявление еще актуально?",
			"Добрый день! Мебель еще продается?",
			"Приветствую! Можно узнать подробнее об этом предложении?",
			"Здравствуйте! Интересуюсь вашим объявлением. Все еще доступно?",
			"Добрый день, мебель еще доступна для покупки?",
		}
		return firstMessages[r.Intn(len(firstMessages))]
	} else if position == 1 && !isBuyer {
		// Первый ответ продавца
		sellerResponses := []string{
			"Добрый день! Да, объявление актуально.",
			"Здравствуйте! Да, всё ещё продаётся.",
			"Приветствую! Да, мебель в наличии.",
			"Добрый день! Конечно, что именно вас интересует?",
			"Здравствуйте! Да, всё в наличии. Есть вопросы?",
		}
		return sellerResponses[r.Intn(len(sellerResponses))]
	} else if isBuyer {
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
		return buyerMessages[r.Intn(len(buyerMessages))]
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
		return sellerMessages[r.Intn(len(sellerMessages))]
	}
}
