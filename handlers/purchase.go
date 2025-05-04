package handlers

import (
	"FurniSwap/models"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// BuyListingHandler обрабатывает запрос на покупку объявления
func BuyListingHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID пользователя из контекста (установлен middleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
			return
		}

		// Получаем ID объявления из пути
		listingID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID объявления"})
			return
		}

		// Начинаем транзакцию
		tx, err := db.Beginx()
		if err != nil {
			log.Printf("Ошибка начала транзакции: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сервера"})
			return
		}

		// Отложенная функция для завершения транзакции
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				log.Printf("Паника при обработке покупки: %v", r)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка сервера"})
			}
		}()

		// Проверяем, существует ли объявление и доступно ли оно для покупки
		var listing models.Listing
		err = tx.Get(&listing, `
			SELECT l.id, l.user_id, l.title, l.price, l.status, u.name as user_name
			FROM listings l
			JOIN users u ON l.user_id = u.id
			WHERE l.id = $1
			FOR UPDATE`, listingID)

		if err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "объявление не найдено"})
			} else {
				log.Printf("Ошибка получения объявления: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сервера"})
			}
			return
		}

		// Проверяем, что объявление не принадлежит покупателю
		if listing.UserID == userID.(int) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "вы не можете купить своё собственное объявление"})
			return
		}

		// Проверяем статус объявления
		if listing.Status != "available" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "объявление уже продано"})
			return
		}

		// Создаем запись о покупке
		_, err = tx.Exec(`
			INSERT INTO purchases (listing_id, buyer_id, seller_id, price, purchased_at)
			VALUES ($1, $2, $3, $4, $5)`,
			listingID, userID, listing.UserID, listing.Price, time.Now())

		if err != nil {
			tx.Rollback()
			log.Printf("Ошибка создания записи о покупке: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сервера"})
			return
		}

		// Обновляем статус объявления
		_, err = tx.Exec("UPDATE listings SET status = 'sold' WHERE id = $1", listingID)
		if err != nil {
			tx.Rollback()
			log.Printf("Ошибка обновления статуса объявления: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сервера"})
			return
		}

		// Фиксируем транзакцию
		if err := tx.Commit(); err != nil {
			log.Printf("Ошибка фиксации транзакции: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сервера"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "объявление успешно куплено",
			"listing": gin.H{
				"id":        listing.ID,
				"title":     listing.Title,
				"seller":    listing.UserName,
				"price":     listing.Price,
				"purchased": time.Now(),
			},
		})
	}
}

// GetUserPurchasesHandler возвращает историю покупок пользователя
func GetUserPurchasesHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID пользователя из контекста (установлен middleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
			return
		}

		// Получаем покупки пользователя
		var purchases []struct {
			models.Purchase
			ListingTitle string `db:"listing_title"`
			SellerName   string `db:"seller_name"`
		}

		err := db.Select(&purchases, `
			SELECT p.id, p.listing_id, p.buyer_id, p.seller_id, p.price, p.purchased_at,
				   l.title as listing_title, u.name as seller_name
			FROM purchases p
			JOIN listings l ON p.listing_id = l.id
			JOIN users u ON p.seller_id = u.id
			WHERE p.buyer_id = $1
			ORDER BY p.purchased_at DESC`, userID)

		if err != nil {
			log.Printf("Ошибка получения покупок: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сервера"})
			return
		}

		c.JSON(http.StatusOK, purchases)
	}
}

// GetUserSalesHandler возвращает историю продаж пользователя
func GetUserSalesHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем ID пользователя из контекста (установлен middleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
			return
		}

		// Получаем продажи пользователя
		var sales []struct {
			models.Purchase
			ListingTitle string `db:"listing_title"`
			BuyerName    string `db:"buyer_name"`
		}

		err := db.Select(&sales, `
			SELECT p.id, p.listing_id, p.buyer_id, p.seller_id, p.price, p.purchased_at,
				   l.title as listing_title, u.name as buyer_name
			FROM purchases p
			JOIN listings l ON p.listing_id = l.id
			JOIN users u ON p.buyer_id = u.id
			WHERE p.seller_id = $1
			ORDER BY p.purchased_at DESC`, userID)

		if err != nil {
			log.Printf("Ошибка получения продаж: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сервера"})
			return
		}

		c.JSON(http.StatusOK, sales)
	}
}
