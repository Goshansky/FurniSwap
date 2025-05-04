package handlers

import (
	"FurniSwap/models"
	"FurniSwap/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GetListingsHandler возвращает список объявлений с фильтрацией и пагинацией
func GetListingsHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter models.ListingFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверные параметры фильтрации"})
			return
		}

		// Устанавливаем значения по умолчанию, если они не указаны
		if filter.Page == 0 {
			filter.Page = 1
		}
		if filter.Limit == 0 {
			filter.Limit = 10
		}

		// Построение запроса с фильтрами
		query := `
			SELECT l.id, l.user_id, l.title, l.description, l.price, l.condition, l.city, 
			       l.category_id, l.status, l.created_at, l.updated_at, u.name as user_name
			FROM listings l
			JOIN users u ON l.user_id = u.id
			WHERE l.status = 'available'
		`
		countQuery := `SELECT COUNT(*) FROM listings l WHERE l.status = 'available'`

		var args []interface{}
		argIndex := 1

		// Добавляем условия фильтрации
		if filter.CategoryID != nil {
			query += fmt.Sprintf(" AND l.category_id = $%d", argIndex)
			countQuery += fmt.Sprintf(" AND l.category_id = $%d", argIndex)
			args = append(args, *filter.CategoryID)
			argIndex++
		}

		if filter.City != "" {
			query += fmt.Sprintf(" AND l.city = $%d", argIndex)
			countQuery += fmt.Sprintf(" AND l.city = $%d", argIndex)
			args = append(args, filter.City)
			argIndex++
		}

		if filter.Condition != "" {
			query += fmt.Sprintf(" AND l.condition = $%d", argIndex)
			countQuery += fmt.Sprintf(" AND l.condition = $%d", argIndex)
			args = append(args, filter.Condition)
			argIndex++
		}

		if filter.MinPrice != nil {
			query += fmt.Sprintf(" AND l.price >= $%d", argIndex)
			countQuery += fmt.Sprintf(" AND l.price >= $%d", argIndex)
			args = append(args, *filter.MinPrice)
			argIndex++
		}

		if filter.MaxPrice != nil {
			query += fmt.Sprintf(" AND l.price <= $%d", argIndex)
			countQuery += fmt.Sprintf(" AND l.price <= $%d", argIndex)
			args = append(args, *filter.MaxPrice)
			argIndex++
		}

		// Сортировка
		if filter.SortBy == "" {
			filter.SortBy = "-date" // По умолчанию сортируем по дате (новые сначала)
		}

		switch filter.SortBy {
		case "date":
			query += " ORDER BY l.created_at ASC"
		case "-date":
			query += " ORDER BY l.created_at DESC"
		case "price":
			query += " ORDER BY l.price ASC"
		case "-price":
			query += " ORDER BY l.price DESC"
		default:
			query += " ORDER BY l.created_at DESC" // Если указан неизвестный способ сортировки
		}

		// Копируем список аргументов для подсчета общего количества
		countArgs := make([]interface{}, len(args))
		copy(countArgs, args)

		// Пагинация
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, filter.Limit, (filter.Page-1)*filter.Limit)

		// Логируем SQL запрос и параметры для отладки
		fmt.Printf("SQL Query: %s\n", query)
		fmt.Printf("SQL Args: %v\n", args)

		// Выполняем запрос
		var listings []models.Listing
		err := db.Select(&listings, query, args...)
		if err != nil {
			fmt.Printf("Error executing query: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения объявлений"})
			return
		}

		// Получаем общее количество объявлений для пагинации
		var totalCount int
		err = db.Get(&totalCount, countQuery, countArgs...)
		if err != nil {
			fmt.Printf("Error counting listings: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка подсчета объявлений"})
			return
		}

		// Получаем основное изображение для каждого объявления
		for i := range listings {
			var image models.Image
			err := db.Get(&image, `
				SELECT id, listing_id, image_path, is_main, created_at
				FROM listing_images
				WHERE listing_id = $1 AND is_main = true
				LIMIT 1
			`, listings[i].ID)
			if err == nil {
				listings[i].Images = []models.Image{image}
			}
		}

		// Вычисляем количество страниц
		totalPages := 1
		if totalCount > 0 {
			totalPages = (totalCount + filter.Limit - 1) / filter.Limit
		}

		c.JSON(http.StatusOK, gin.H{
			"listings":    listings,
			"total_count": totalCount,
			"page":        filter.Page,
			"limit":       filter.Limit,
			"total_pages": totalPages,
		})
	}
}

// GetListingHandler возвращает детальную информацию об объявлении
func GetListingHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		listingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID объявления"})
			return
		}

		// Получаем объявление
		var listing models.Listing
		err = db.Get(&listing, `
			SELECT l.id, l.user_id, l.title, l.description, l.price, 
			       l.condition, l.city, l.category_id, l.status, l.created_at, l.updated_at,
			       u.name as user_name
			FROM listings l
			JOIN users u ON l.user_id = u.id
			WHERE l.id = $1
		`, listingID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "объявление не найдено"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения объявления"})
			}
			return
		}

		// Получаем все изображения объявления
		var images []models.Image
		err = db.Select(&images, `
			SELECT id, listing_id, image_path, is_main, created_at
			FROM listing_images
			WHERE listing_id = $1
			ORDER BY is_main DESC, created_at ASC
		`, listingID)
		if err == nil {
			listing.Images = images
		}

		// Получаем категорию
		var category models.Category
		err = db.Get(&category, "SELECT id, name FROM categories WHERE id = $1", listing.CategoryID)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"listing":  listing,
				"category": category.Name,
			})
		} else {
			c.JSON(http.StatusOK, listing)
		}
	}
}

// CreateListingHandler создает новое объявление
func CreateListingHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")

		var req models.CreateListingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные"})
			return
		}

		// Начинаем транзакцию
		tx, err := db.Beginx()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}
		defer tx.Rollback()

		// Создаем объявление
		var listingID int
		err = tx.QueryRow(`
			INSERT INTO listings (user_id, title, description, price, condition, city, category_id, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'available')
			RETURNING id
		`, userID, req.Title, req.Description, req.Price, req.Condition, req.City, req.CategoryID).Scan(&listingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания объявления"})
			return
		}

		// Фиксируем транзакцию
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":    "объявление успешно создано",
			"listing_id": listingID,
		})
	}
}

// UpdateListingHandler обновляет существующее объявление
func UpdateListingHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		listingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID объявления"})
			return
		}

		// Проверяем, принадлежит ли объявление пользователю
		var ownerID int
		err = db.Get(&ownerID, "SELECT user_id FROM listings WHERE id = $1", listingID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "объявление не найдено"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			}
			return
		}

		if ownerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "нет прав на редактирование"})
			return
		}

		var req models.UpdateListingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные"})
			return
		}

		// Начинаем транзакцию
		tx, err := db.Beginx()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}
		defer tx.Rollback()

		// Обновляем объявление с учетом возможных пустых полей
		_, err = tx.Exec(`
			UPDATE listings
			SET title = COALESCE(NULLIF($1, ''), title),
				description = COALESCE(NULLIF($2, ''), description),
				price = CASE WHEN $3 > 0 THEN $3 ELSE price END,
				condition = COALESCE(NULLIF($4, ''), condition),
				city = COALESCE(NULLIF($5, ''), city),
				category_id = CASE WHEN $6 > 0 THEN $6 ELSE category_id END,
				updated_at = NOW()
			WHERE id = $7
		`, req.Title, req.Description, req.Price, req.Condition, req.City, req.CategoryID, listingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления объявления"})
			return
		}

		// Фиксируем транзакцию
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "объявление успешно обновлено"})
	}
}

// DeleteListingHandler удаляет объявление
func DeleteListingHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		listingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID объявления"})
			return
		}

		// Проверяем, принадлежит ли объявление пользователю
		var ownerID int
		err = db.Get(&ownerID, "SELECT user_id FROM listings WHERE id = $1", listingID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "объявление не найдено"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			}
			return
		}

		if ownerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "нет прав на удаление"})
			return
		}

		// Начинаем транзакцию
		tx, err := db.Beginx()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}
		defer tx.Rollback()

		// Получаем пути всех изображений для последующего удаления файлов
		var imagePaths []string
		err = tx.Select(&imagePaths, "SELECT image_path FROM listing_images WHERE listing_id = $1", listingID)
		if err != nil && err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения изображений"})
			return
		}

		// Удаляем объявление (каскадно удаляются изображения и избранное)
		_, err = tx.Exec("DELETE FROM listings WHERE id = $1", listingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка удаления объявления"})
			return
		}

		// Фиксируем транзакцию
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		// Удаляем файлы изображений
		for _, path := range imagePaths {
			utils.DeleteImage(path)
		}

		c.JSON(http.StatusOK, gin.H{"message": "объявление успешно удалено"})
	}
}

// UploadListingImageHandler загружает изображение для объявления
func UploadListingImageHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		listingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID объявления"})
			return
		}

		// Проверяем, принадлежит ли объявление пользователю
		var ownerID int
		err = db.Get(&ownerID, "SELECT user_id FROM listings WHERE id = $1", listingID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "объявление не найдено"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			}
			return
		}

		if ownerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "нет прав на добавление изображений"})
			return
		}

		// Получаем загружаемый файл
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "файл не найден"})
			return
		}

		// Сохраняем файл
		filePath, err := utils.UploadImage(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверяем, есть ли уже главное изображение
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM listing_images WHERE listing_id = $1", listingID)
		if err != nil {
			utils.DeleteImage(filePath)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		// Если это первое изображение, делаем его главным
		isMain := count == 0

		// Добавляем изображение в БД
		var imageID int
		err = db.QueryRow(`
			INSERT INTO listing_images (listing_id, image_path, is_main)
			VALUES ($1, $2, $3)
			RETURNING id
		`, listingID, filePath, isMain).Scan(&imageID)
		if err != nil {
			utils.DeleteImage(filePath)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения изображения"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "изображение успешно загружено",
			"image_id":   imageID,
			"image_path": filePath,
			"is_main":    isMain,
		})
	}
}

// DeleteListingImageHandler удаляет изображение объявления
func DeleteListingImageHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		imageID, err := strconv.Atoi(c.Param("imageId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID изображения"})
			return
		}

		// Проверяем, принадлежит ли изображение пользователю
		var image struct {
			ListingID int    `db:"listing_id"`
			OwnerID   int    `db:"owner_id"`
			ImagePath string `db:"image_path"`
			IsMain    bool   `db:"is_main"`
		}

		err = db.Get(&image, `
			SELECT i.listing_id, l.user_id as owner_id, i.image_path, i.is_main
			FROM listing_images i
			JOIN listings l ON i.listing_id = l.id
			WHERE i.id = $1
		`, imageID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "изображение не найдено"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			}
			return
		}

		if image.OwnerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "нет прав на удаление изображения"})
			return
		}

		tx, err := db.Beginx()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}
		defer tx.Rollback()

		// Удаляем изображение из БД
		_, err = tx.Exec("DELETE FROM listing_images WHERE id = $1", imageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка удаления изображения"})
			return
		}

		// Если удаляемое изображение было главным, делаем главным первое найденное
		if image.IsMain {
			var newMainImageID int
			err = tx.Get(&newMainImageID, `
				SELECT id FROM listing_images 
				WHERE listing_id = $1 
				ORDER BY created_at ASC 
				LIMIT 1
			`, image.ListingID)

			if err == nil {
				_, err = tx.Exec("UPDATE listing_images SET is_main = true WHERE id = $1", newMainImageID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления главного изображения"})
					return
				}
			}
		}

		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		// Удаляем файл изображения
		if err := utils.DeleteImage(image.ImagePath); err != nil {
			// Логируем ошибку, но продолжаем выполнение
			fmt.Printf("Ошибка удаления файла %s: %v\n", image.ImagePath, err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "изображение успешно удалено"})
	}
}

// SetMainImageHandler устанавливает главное изображение для объявления
func SetMainImageHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		imageID, err := strconv.Atoi(c.Param("imageId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID изображения"})
			return
		}

		// Проверяем, принадлежит ли изображение пользователю
		var image struct {
			ListingID int `db:"listing_id"`
			OwnerID   int `db:"owner_id"`
		}

		err = db.Get(&image, `
			SELECT i.listing_id, l.user_id as owner_id
			FROM listing_images i
			JOIN listings l ON i.listing_id = l.id
			WHERE i.id = $1
		`, imageID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "изображение не найдено"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			}
			return
		}

		if image.OwnerID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "нет прав на изменение изображения"})
			return
		}

		// Начинаем транзакцию
		tx, err := db.Beginx()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}
		defer tx.Rollback()

		// Сбрасываем флаг главного изображения для всех изображений объявления
		_, err = tx.Exec("UPDATE listing_images SET is_main = false WHERE listing_id = $1", image.ListingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обновления изображений"})
			return
		}

		// Устанавливаем новое главное изображение
		_, err = tx.Exec("UPDATE listing_images SET is_main = true WHERE id = $1", imageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка установки главного изображения"})
			return
		}

		// Фиксируем транзакцию
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "главное изображение успешно установлено"})
	}
}
