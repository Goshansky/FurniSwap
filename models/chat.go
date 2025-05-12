package models

import (
	"time"
)

// Chat представляет чат между покупателем и продавцом
type Chat struct {
	ID          int       `db:"id" json:"id"`
	ListingID   int       `db:"listing_id" json:"listing_id"`
	BuyerID     int       `db:"buyer_id" json:"buyer_id"`
	SellerID    int       `db:"seller_id" json:"seller_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Listing     *Listing  `json:"listing,omitempty"` // Для отображения в списке чатов
	BuyerName   string    `json:"buyer_name,omitempty"`
	SellerName  string    `json:"seller_name,omitempty"`
	LastMessage *Message  `json:"last_message,omitempty"`
}

// Message представляет сообщение в чате
type Message struct {
	ID        int       `db:"id" json:"id"`
	ChatID    int       `db:"chat_id" json:"chat_id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	IsRead    bool      `db:"is_read" json:"is_read"`
	UserName  string    `db:"user_name" json:"user_name,omitempty"`
}

// CreateMessageRequest структура для создания нового сообщения
type CreateMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// ChatResponse структура для ответа со списком чатов пользователя
type ChatResponse struct {
	ID              int       `db:"id" json:"id"`
	ListingID       int       `db:"listing_id" json:"listing_id"`
	BuyerID         int       `db:"buyer_id" json:"buyer_id,omitempty"`
	SellerID        int       `db:"seller_id" json:"seller_id,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at,omitempty"`
	ListingTitle    string    `db:"listing_title" json:"listing_title"`
	ImageURL        string    `json:"image_url,omitempty"`
	OtherUserID     int       `db:"other_user_id" json:"other_user_id"`
	OtherUserName   string    `db:"other_user_name" json:"other_user_name"`
	LastMessage     string    `db:"last_message" json:"last_message,omitempty"`
	LastMessageTime time.Time `db:"last_message_time" json:"last_message_time,omitempty"`
	UnreadCount     int       `db:"unread_count" json:"unread_count"`
}

// InitiateChatRequest структура для создания нового чата
type InitiateChatRequest struct {
	ListingID int    `json:"listing_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
}
