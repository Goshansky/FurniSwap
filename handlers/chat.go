package handlers

import (
	"FurniSwap/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// InitiateChatHandler создает новый чат или возвращает существующий
func InitiateChatHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		buyerID := c.GetInt("userID")

		var req models.InitiateChatRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные"})
			return
		}

		// Получаем данные объявления и продавца
		var listing struct {
			ID       int    `db:"id"`
			SellerID int    `db:"user_id"`
			Title    string `db:"title"`
		}

		err := db.Get(&listing, "SELECT id, user_id, title FROM listings WHERE id = $1", req.ListingID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "объявление не найдено"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			}
			return
		}

		sellerID := listing.SellerID

		// Проверяем, что покупатель не является продавцом
		if buyerID == sellerID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "вы не можете начать чат с самим собой"})
			return
		}

		// Начинаем транзакцию
		tx, err := db.Beginx()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}
		defer tx.Rollback()

		// Проверяем, существует ли уже чат
		var chatID int
		err = tx.Get(&chatID, `
			SELECT id FROM chats 
			WHERE listing_id = $1 AND buyer_id = $2 AND seller_id = $3
		`, req.ListingID, buyerID, sellerID)

		// Если чат не существует, создаем новый
		if err == sql.ErrNoRows {
			err = tx.QueryRow(`
				INSERT INTO chats (listing_id, buyer_id, seller_id)
				VALUES ($1, $2, $3)
				RETURNING id
			`, req.ListingID, buyerID, sellerID).Scan(&chatID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания чата"})
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка проверки чата"})
			return
		}

		// Создаем сообщение
		var messageID int
		err = tx.QueryRow(`
			INSERT INTO messages (chat_id, user_id, content)
			VALUES ($1, $2, $3)
			RETURNING id
		`, chatID, buyerID, req.Message).Scan(&messageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания сообщения"})
			return
		}

		// Фиксируем транзакцию
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "чат создан",
			"chat_id": chatID,
		})
	}
}

// GetChatsHandler возвращает список чатов пользователя
func GetChatsHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")

		// Получаем все чаты пользователя (как покупателя, так и продавца)
		query := `
			SELECT c.id, c.listing_id, c.buyer_id, c.seller_id, c.created_at,
				   l.title as listing_title,
				   CASE
					   WHEN c.buyer_id = $1 THEN c.seller_id
					   ELSE c.buyer_id
				   END as other_user_id,
				   CASE
					   WHEN c.buyer_id = $1 THEN s.name
					   ELSE b.name
				   END as other_user_name,
				   m.content as last_message,
				   m.created_at as last_message_time,
				   (SELECT COUNT(*) FROM messages 
					WHERE chat_id = c.id AND user_id != $1 AND is_read = false) as unread_count
			FROM chats c
			JOIN listings l ON c.listing_id = l.id
			JOIN users b ON c.buyer_id = b.id
			JOIN users s ON c.seller_id = s.id
			LEFT JOIN (
				SELECT chat_id, content, created_at, user_id,
				       ROW_NUMBER() OVER (PARTITION BY chat_id ORDER BY created_at DESC) as rn
				FROM messages
			) m ON m.chat_id = c.id AND m.rn = 1
			WHERE c.buyer_id = $1 OR c.seller_id = $1
			ORDER BY m.created_at DESC
		`

		var chats []models.ChatResponse
		err := db.Select(&chats, query, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения чатов"})
			return
		}

		// Получаем основное изображение для каждого объявления
		for i := range chats {
			var imagePath string
			err := db.Get(&imagePath, `
				SELECT image_path
				FROM listing_images
				WHERE listing_id = $1 AND is_main = true
				LIMIT 1
			`, chats[i].ListingID)
			if err == nil {
				chats[i].ImageURL = imagePath
			}
		}

		c.JSON(http.StatusOK, gin.H{"chats": chats})
	}
}

// GetChatMessagesHandler возвращает сообщения в чате
func GetChatMessagesHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		chatID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID чата"})
			return
		}

		// Проверяем, принадлежит ли чат пользователю
		var chat models.Chat
		err = db.Get(&chat, `
			SELECT id, listing_id, buyer_id, seller_id, created_at
			FROM chats 
			WHERE id = $1 AND (buyer_id = $2 OR seller_id = $2)
		`, chatID, userID)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "чат не найден"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			}
			return
		}

		// Получаем информацию об объявлении
		var listing models.Listing
		err = db.Get(&listing, `
			SELECT id, user_id, title, description, price, condition, city, category_id, created_at, updated_at
			FROM listings
			WHERE id = $1
		`, chat.ListingID)
		if err == nil {
			chat.Listing = &listing
		}

		// Получаем имена участников чата
		var buyerName, sellerName string
		_ = db.Get(&buyerName, "SELECT name FROM users WHERE id = $1", chat.BuyerID)
		_ = db.Get(&sellerName, "SELECT name FROM users WHERE id = $1", chat.SellerID)
		chat.BuyerName = buyerName
		chat.SellerName = sellerName

		// Получаем все сообщения в чате
		query := `
			SELECT m.id, m.chat_id, m.user_id, m.content, m.created_at, m.is_read,
			       u.name as user_name
			FROM messages m
			JOIN users u ON m.user_id = u.id
			WHERE m.chat_id = $1
			ORDER BY m.created_at ASC
		`

		var messages []models.Message
		err = db.Select(&messages, query, chatID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения сообщений"})
			return
		}

		// Помечаем все непрочитанные сообщения как прочитанные
		_, err = db.Exec(`
			UPDATE messages 
			SET is_read = true 
			WHERE chat_id = $1 AND user_id != $2 AND is_read = false
		`, chatID, userID)
		if err != nil {
			// Логируем ошибку, но продолжаем выполнение
			//log.Printf("Ошибка обновления статуса сообщений: %v", err)
		}

		// Возвращаем информацию о чате и сообщения
		c.JSON(http.StatusOK, gin.H{
			"chat":     chat,
			"messages": messages,
		})
	}
}

// SendMessageHandler отправляет новое сообщение в чат
func SendMessageHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		chatID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID чата"})
			return
		}

		var req models.CreateMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверные данные"})
			return
		}

		// Проверяем, принадлежит ли чат пользователю
		var exists bool
		err = db.Get(&exists, `
			SELECT EXISTS(
				SELECT 1 FROM chats 
				WHERE id = $1 AND (buyer_id = $2 OR seller_id = $2)
			)
		`, chatID, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
			return
		}

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "чат не найден"})
			return
		}

		// Сохраняем сообщение
		var messageID int
		err = db.QueryRow(`
			INSERT INTO messages (chat_id, user_id, content)
			VALUES ($1, $2, $3)
			RETURNING id
		`, chatID, userID, req.Content).Scan(&messageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка отправки сообщения"})
			return
		}

		// Получаем созданное сообщение
		var message models.Message
		err = db.Get(&message, `
			SELECT m.id, m.chat_id, m.user_id, m.content, m.created_at, m.is_read,
			       u.name as user_name
			FROM messages m
			JOIN users u ON m.user_id = u.id
			WHERE m.id = $1
		`, messageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения сообщения"})
			return
		}

		c.JSON(http.StatusOK, message)
	}
}
