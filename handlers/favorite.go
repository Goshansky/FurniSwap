package handlers

import (
	"FurniSwap/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// AddFavoriteHandler добавляет объявление в избранное
func AddFavoriteHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		listingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID объявления"})
			return
		}

		// Проверяем, существует ли объявление
		var exists bool
		err = db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM listings WHERE id = $1)", listingID)
		if err != nil || !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "объявление не найдено"})
			return
		}

		// Проверяем, не принадлежит ли объявление пользователю
		var ownerID int
		err = db.Get(&ownerID, "SELECT user_id FROM listings WHERE id = $1", listingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		if ownerID == userID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "нельзя добавить своё объявление в избранное"})
			return
		}

		// Пытаемся добавить в избранное
		_, err = db.Exec(`
			INSERT INTO favorites (user_id, listing_id) 
			VALUES ($1, $2) 
			ON CONFLICT (user_id, listing_id) DO NOTHING
		`, userID, listingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка добавления в избранное"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "объявление добавлено в избранное"})
	}
}

// RemoveFavoriteHandler удаляет объявление из избранного
func RemoveFavoriteHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		listingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID объявления"})
			return
		}

		// Удаляем из избранного
		result, err := db.Exec("DELETE FROM favorites WHERE user_id = $1 AND listing_id = $2", userID, listingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка удаления из избранного"})
			return
		}

		// Проверяем, была ли запись удалена
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "объявление не найдено в избранном"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "объявление удалено из избранного"})
	}
}

// GetFavoritesHandler возвращает список избранных объявлений пользователя
func GetFavoritesHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")

		// Получаем все избранные объявления пользователя с информацией о них
		query := `
			SELECT l.id, l.user_id, l.title, l.description, l.price, 
			       l.condition, l.city, l.category_id, l.created_at, l.updated_at,
			       u.name as user_name, f.created_at as favorite_at
			FROM favorites f
			JOIN listings l ON f.listing_id = l.id
			JOIN users u ON l.user_id = u.id
			WHERE f.user_id = $1
			ORDER BY f.created_at DESC
		`

		var listings []struct {
			models.Listing
			FavoriteAt sql.NullTime `db:"favorite_at" json:"favorite_at"`
		}
		err := db.Select(&listings, query, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения избранных объявлений"})
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

		c.JSON(http.StatusOK, gin.H{"favorites": listings})
	}
}

// IsFavoriteHandler проверяет, добавлено ли объявление в избранное пользователя
func IsFavoriteHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		listingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID объявления"})
			return
		}

		var exists bool
		err = db.Get(&exists, `
			SELECT EXISTS(
				SELECT 1 FROM favorites 
				WHERE user_id = $1 AND listing_id = $2
			)
		`, userID, listingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"is_favorite": exists})
	}
}
