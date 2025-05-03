package main

import (
	"FurniSwap/handlers"
	"FurniSwap/middlewares"
	"FurniSwap/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Функция для проверки и обновления структуры БД
func updateDatabaseStructure(db *sqlx.DB) {
	// Проверяем структуру таблицы users
	var columns []string
	err := db.Select(&columns, `
		SELECT column_name FROM information_schema.columns 
		WHERE table_name = 'users'
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Printf("Ошибка при проверке структуры таблицы: %v", err)
		return
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
			log.Printf("Ошибка при добавлении столбца last_name: %v", err)
			return
		}

		// Обновляем существующие записи
		_, err = db.Exec("UPDATE users SET last_name = '' WHERE last_name IS NULL")
		if err != nil {
			log.Printf("Ошибка при обновлении существующих записей: %v", err)
			return
		}

		log.Println("Столбец last_name успешно добавлен в таблицу users")
	} else {
		log.Println("Столбец last_name уже существует в таблице users")
	}

	// Обновляем данные тестового пользователя
	result, err := db.Exec("UPDATE users SET name = $1, last_name = $2 WHERE email = $3",
		"Тестовый", "Пользователь", "test@example.com")
	if err != nil {
		log.Printf("Ошибка при обновлении данных пользователя: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Обновлено пользователей: %d", rowsAffected)
}

func main() {
	// Загружаем .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Предупреждение: Ошибка загрузки .env файла. Используются значения по умолчанию.")
	}

	// Создаем директорию для загрузок, если ее нет
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal("Ошибка создания директории uploads:", err)
	}

	connectionString := fmt.Sprintf("host=localhost port=5431 user=postgres password=password dbname=furni_swap sslmode=disable")
	log.Println("Подключение к базе данных:", connectionString)

	// Подключение к БД
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// Проверка структуры базы данных
	utils.CheckDatabase(db)

	// Обновляем структуру БД при необходимости
	updateDatabaseStructure(db)

	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Замени на адрес фронтенда
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Статические файлы для загруженных изображений
	r.Static("/uploads", "./uploads")

	// Группа маршрутов для аутентификации
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.RegisterHandler(db))
		auth.POST("/verify", handlers.VerifyHandler(db))
		auth.POST("/login", handlers.LoginHandler(db))
		auth.POST("/verify-2fa", handlers.Verify2FAHandler(db))
	}

	// Группа маршрутов для категорий (без аутентификации)
	categories := r.Group("/categories")
	{
		categories.GET("", func(c *gin.Context) {
			var categories []struct {
				ID   int    `db:"id" json:"id"`
				Name string `db:"name" json:"name"`
			}
			err := db.Select(&categories, "SELECT id, name FROM categories ORDER BY name")
			if err != nil {
				log.Printf("Ошибка получения категорий: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения категорий"})
				return
			}
			c.JSON(http.StatusOK, categories)
		})
	}

	// Группа маршрутов для объявлений (без аутентификации)
	listings := r.Group("/listings")
	{
		listings.GET("", handlers.GetListingsHandler(db))
		listings.GET("/:id", handlers.GetListingHandler(db))
	}

	// API маршруты, требующие аутентификации
	api := r.Group("/api")
	api.Use(middlewares.AuthRequired(db))
	{
		// Пользовательский профиль
		api.GET("/profile", handlers.GetProfileHandler(db))
		api.PUT("/profile", handlers.UpdateProfileHandler(db))
		api.POST("/profile/avatar", handlers.UploadAvatarHandler(db))

		// Управление объявлениями пользователя
		api.POST("/listings", handlers.CreateListingHandler(db))
		api.PUT("/listings/:id", handlers.UpdateListingHandler(db))
		api.DELETE("/listings/:id", handlers.DeleteListingHandler(db))
		api.POST("/listings/:id/images", handlers.UploadListingImageHandler(db))
		api.DELETE("/listings/:id/images/:imageId", handlers.DeleteListingImageHandler(db))
		api.PUT("/listings/:id/images/:imageId/main", handlers.SetMainImageHandler(db))

		// Избранное
		api.POST("/listings/:id/favorite", handlers.AddFavoriteHandler(db))
		api.DELETE("/listings/:id/favorite", handlers.RemoveFavoriteHandler(db))
		api.GET("/listings/:id/favorite", handlers.IsFavoriteHandler(db))
		api.GET("/favorites", handlers.GetFavoritesHandler(db))

		// Чаты и сообщения
		api.POST("/chats", handlers.InitiateChatHandler(db))
		api.GET("/chats", handlers.GetChatsHandler(db))
		api.GET("/chats/:id", handlers.GetChatMessagesHandler(db))
		api.POST("/chats/:id/messages", handlers.SendMessageHandler(db))
	}

	log.Println("Сервер запущен на http://localhost:8080")
	r.Run(":8080")
}
