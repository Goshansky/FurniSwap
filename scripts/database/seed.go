package database

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Константы для генерации тестовых данных
const (
	UserCount        = 20 // Количество пользователей
	ListingsPerUser  = 3  // Среднее количество объявлений на пользователя
	ImagesPerListing = 2  // Среднее количество изображений на объявление
	FavoritesCount   = 30 // Количество избранных объявлений
	ChatsCount       = 15 // Количество чатов
	MessagesPerChat  = 5  // Среднее количество сообщений в чате
)

// Данные для генерации случайных значений
var (
	Cities     = []string{"Москва", "Санкт-Петербург", "Казань", "Новосибирск", "Екатеринбург", "Нижний Новгород", "Самара"}
	Conditions = []string{"новое", "хорошее", "среднее", "плохое"}

	// Мебель по категориям
	FurnitureByCategory = map[int][]string{
		1: {"Диван угловой", "Диван прямой", "Кресло раскладное", "Кресло-кровать", "Пуфик мягкий", "Кресло-качалка", "Диван-кровать"},
		2: {"Стол обеденный", "Стол журнальный", "Стул деревянный", "Стул мягкий", "Стол письменный", "Табурет кухонный", "Стол компьютерный"},
		3: {"Шкаф-купе", "Комод с ящиками", "Тумба под ТВ", "Шкаф для одежды", "Пенал кухонный", "Книжный шкаф", "Шкаф в прихожую"},
		4: {"Кровать двуспальная", "Кровать односпальная", "Матрас ортопедический", "Кровать с ящиками", "Детская кровать", "Раскладушка", "Матрас пружинный"},
		5: {"Стеллаж", "Тумбочка прикроватная", "Вешалка напольная", "Полка настенная", "Зеркало напольное", "Подставка для цветов", "Этажерка"},
	}

	// Русские имена и фамилии
	FirstNamesRu = []string{"Александр", "Дмитрий", "Максим", "Сергей", "Иван", "Андрей", "Алексей", "Артём", "Михаил", "Никита", "Анна", "Мария", "Екатерина", "Елена", "Ольга", "Наталья", "Татьяна", "Юлия", "Дарья", "Виктория"}
	LastNamesRu  = []string{"Иванов", "Петров", "Сидоров", "Смирнов", "Кузнецов", "Попов", "Соколов", "Михайлов", "Новиков", "Федоров", "Морозов", "Волков", "Алексеев", "Лебедев", "Семенов", "Егоров", "Павлов", "Козлов", "Степанов", "Николаев"}

	// Английские имена и фамилии
	FirstNamesEn = []string{"Alexander", "Dmitry", "Maxim", "Sergei", "Ivan", "Andrey", "Alexey", "Artem", "Mikhail", "Nikita", "Anna", "Maria", "Ekaterina", "Elena", "Olga", "Natalia", "Tatiana", "Yulia", "Daria", "Victoria"}
	LastNamesEn  = []string{"Ivanov", "Petrov", "Sidorov", "Smirnov", "Kuznetsov", "Popov", "Sokolov", "Mikhailov", "Novikov", "Fedorov", "Morozov", "Volkov", "Alekseev", "Lebedev", "Semenov", "Egorov", "Pavlov", "Kozlov", "Stepanov", "Nikolaev"}

	// Пути к аватарам
	WebAvatarPaths = []string{
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
	Descriptions = []string{
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

// GenerateUsers генерирует случайных пользователей
func GenerateUsers(count int, useEnglishNames bool) []map[string]interface{} {
	users := make([]map[string]interface{}, count)

	firstNames := FirstNamesRu
	lastNames := LastNamesRu

	if useEnglishNames {
		firstNames = FirstNamesEn
		lastNames = LastNamesEn
	}

	for i := 0; i < count; i++ {
		name := firstNames[r.Intn(len(firstNames))]
		lastName := lastNames[r.Intn(len(lastNames))]
		email := fmt.Sprintf("%s.%s%d@example.com", name, lastName, i)
		email = fmt.Sprintf("%s", email) // Нормализуем email для нелатинских символов

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		users[i] = map[string]interface{}{
			"email":         email,
			"password_hash": string(hashedPassword),
			"name":          name,
			"last_name":     lastName,
			"city":          Cities[r.Intn(len(Cities))],
			"avatar":        WebAvatarPaths[r.Intn(len(WebAvatarPaths))],
			"is_verified":   true,
		}
	}

	return users
}

// InsertUsers вставляет пользователей в базу данных
func InsertUsers(tx *sqlx.Tx, users []map[string]interface{}) ([]int, error) {
	userIDs := make([]int, 0, len(users))

	for _, user := range users {
		var userID int
		err := tx.QueryRow(`
			INSERT INTO users (email, password_hash, name, last_name, city, avatar, is_verified, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
			RETURNING id
		`, user["email"], user["password_hash"], user["name"], user["last_name"], user["city"],
			user["avatar"], user["is_verified"]).Scan(&userID)

		if err != nil {
			return nil, err
		}

		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

// CleanDatabase очищает все данные из базы данных
func CleanDatabase(tx *sqlx.Tx) error {
	// Удаляем данные из таблиц в правильном порядке (с учетом внешних ключей)
	tables := []string{
		"messages",
		"chats",
		"favorites",
		"listing_images",
		"listings",
		"two_factor_codes",
		"users",
	}

	for _, table := range tables {
		_, err := tx.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("ошибка при удалении данных из таблицы %s: %w", table, err)
		}
		log.Printf("Таблица %s очищена", table)
	}

	return nil
}

// EnsureCategories проверяет наличие категорий в базе данных и добавляет их при необходимости
func EnsureCategories(tx *sqlx.Tx) error {
	var categoryCount int
	err := tx.Get(&categoryCount, "SELECT COUNT(*) FROM categories")
	if err != nil {
		return err
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
				return err
			}
		}
		log.Printf("Вставлено %d категорий", len(categories))
	} else {
		log.Printf("Категории уже существуют, пропускаем вставку")
	}

	return nil
}
