package repository

import (
	"FurniSwap/internal/modules/chat/model"
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for the chat module
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new chat repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetChatByUsers retrieves a chat between two users
func (r *Repository) GetChatByUsers(user1ID, user2ID int, listingID *int) (*model.Chat, error) {
	var chat model.Chat
	var err error

	if listingID != nil {
		err = r.db.Get(&chat, `
			SELECT id, buyer_id as user1_id, seller_id as user2_id, listing_id,
                   created_at, created_at as last_message_at
			FROM chats
			WHERE (buyer_id = $1 AND seller_id = $2 AND listing_id = $3)
			   OR (buyer_id = $2 AND seller_id = $1 AND listing_id = $3)
		`, user1ID, user2ID, listingID)
	} else {
		err = r.db.Get(&chat, `
			SELECT id, buyer_id as user1_id, seller_id as user2_id, listing_id,
                   created_at, created_at as last_message_at
			FROM chats
			WHERE ((buyer_id = $1 AND seller_id = $2) OR (buyer_id = $2 AND seller_id = $1))
			AND listing_id IS NULL
		`, user1ID, user2ID)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Chat doesn't exist
		}
		log.Printf("Error getting chat by users: %v", err)
		return nil, fmt.Errorf("error getting chat: %w", err)
	}

	return &chat, nil
}

// CreateChat creates a new chat between users
func (r *Repository) CreateChat(buyerID, sellerID int, listingID *int) (int, error) {
	var chatID int
	err := r.db.QueryRow(`
		INSERT INTO chats (buyer_id, seller_id, listing_id, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, buyerID, sellerID, listingID, time.Now()).Scan(&chatID)

	if err != nil {
		log.Printf("Error creating chat: %v", err)
		return 0, fmt.Errorf("error creating chat: %w", err)
	}

	return chatID, nil
}

// AddMessage adds a message to a chat
func (r *Repository) AddMessage(chatID, senderID int, content string) (int, error) {
	var messageID int
	err := r.db.QueryRow(`
		INSERT INTO messages (chat_id, user_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, chatID, senderID, content, time.Now()).Scan(&messageID)

	if err != nil {
		log.Printf("Error adding message: %v", err)
		return 0, fmt.Errorf("error adding message: %w", err)
	}

	return messageID, nil
}

// GetUserChats gets a user's chats with pagination
func (r *Repository) GetUserChats(userID, page, limit int) (*model.ChatResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get total count
	var totalCount int
	err := r.db.Get(&totalCount, `
		SELECT COUNT(*) FROM chats
		WHERE buyer_id = $1 OR seller_id = $1
	`, userID)
	if err != nil {
		log.Printf("Error getting chat count: %v", err)
		return nil, fmt.Errorf("error getting chat count: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Get chats
	var chats []model.Chat
	err = r.db.Select(&chats, `
		SELECT c.id, 
			   c.buyer_id as user1_id, 
			   c.seller_id as user2_id, 
			   c.listing_id,
			   c.created_at,
			   c.created_at as last_message_at,
			   u1.name || ' ' || COALESCE(u1.last_name, '') as user1_name,
			   u2.name || ' ' || COALESCE(u2.last_name, '') as user2_name,
			   l.title as listing_title,
			   COALESCE((SELECT content FROM messages 
				WHERE chat_id = c.id 
				ORDER BY created_at DESC LIMIT 1), '') as last_message
		FROM chats c
		JOIN users u1 ON c.buyer_id = u1.id
		JOIN users u2 ON c.seller_id = u2.id
		LEFT JOIN listings l ON c.listing_id = l.id
		WHERE c.buyer_id = $1 OR c.seller_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)

	if err != nil {
		log.Printf("Error getting chats: %v", err)
		return nil, fmt.Errorf("error getting chats: %w", err)
	}

	return &model.ChatResponse{
		Chats:       chats,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}

// GetChatMessages gets messages for a chat with pagination
func (r *Repository) GetChatMessages(chatID, page, limit int) (*model.MessageResponse, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get total count
	var totalCount int
	err := r.db.Get(&totalCount, "SELECT COUNT(*) FROM messages WHERE chat_id = $1", chatID)
	if err != nil {
		log.Printf("Error getting message count: %v", err)
		return nil, fmt.Errorf("error getting message count: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	// Get messages
	var messages []model.Message
	err = r.db.Select(&messages, `
		SELECT m.id, m.chat_id, m.user_id as sender_id, m.content, m.created_at,
			   u.name || ' ' || COALESCE(u.last_name, '') as sender_name
		FROM messages m
		JOIN users u ON m.user_id = u.id
		WHERE m.chat_id = $1
		ORDER BY m.created_at ASC
		LIMIT $2 OFFSET $3
	`, chatID, limit, offset)

	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return nil, fmt.Errorf("error getting messages: %w", err)
	}

	return &model.MessageResponse{
		Messages:    messages,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}

// GetChatByID retrieves a chat by ID
func (r *Repository) GetChatByID(chatID int) (*model.Chat, error) {
	var chat model.Chat
	err := r.db.Get(&chat, `
		SELECT c.id, c.buyer_id as user1_id, c.seller_id as user2_id, c.listing_id,
			   c.created_at, c.created_at as last_message_at,
			   u1.name || ' ' || COALESCE(u1.last_name, '') as user1_name,
			   u2.name || ' ' || COALESCE(u2.last_name, '') as user2_name,
			   l.title as listing_title
		FROM chats c
		JOIN users u1 ON c.buyer_id = u1.id
		JOIN users u2 ON c.seller_id = u2.id
		LEFT JOIN listings l ON c.listing_id = l.id
		WHERE c.id = $1
	`, chatID)

	if err != nil {
		log.Printf("Error getting chat by ID: %v", err)
		return nil, fmt.Errorf("error getting chat: %w", err)
	}

	return &chat, nil
}

// CheckChatAccess checks if a user has access to a chat
func (r *Repository) CheckChatAccess(chatID, userID int) (bool, error) {
	var count int
	err := r.db.Get(&count, `
		SELECT COUNT(*) FROM chats
		WHERE id = $1 AND (buyer_id = $2 OR seller_id = $2)
	`, chatID, userID)

	if err != nil {
		log.Printf("Error checking chat access: %v", err)
		return false, fmt.Errorf("error checking chat access: %w", err)
	}

	return count > 0, nil
}
