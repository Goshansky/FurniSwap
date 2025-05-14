package model

import (
	"time"
)

// Chat represents a chat between users
type Chat struct {
	ID            int       `db:"id" json:"id"`
	User1ID       int       `db:"user1_id" json:"user1_id"`
	User2ID       int       `db:"user2_id" json:"user2_id"`
	ListingID     *int      `db:"listing_id" json:"listing_id,omitempty"`
	LastMessageAt time.Time `db:"last_message_at" json:"last_message_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	User1Name     string    `db:"user1_name" json:"user1_name,omitempty"`
	User2Name     string    `db:"user2_name" json:"user2_name,omitempty"`
	ListingTitle  string    `db:"listing_title" json:"listing_title,omitempty"`
	LastMessage   string    `db:"last_message" json:"last_message,omitempty"`
}

// Message represents a chat message
type Message struct {
	ID         int       `db:"id" json:"id"`
	ChatID     int       `db:"chat_id" json:"chat_id"`
	SenderID   int       `db:"sender_id" json:"sender_id"`
	Content    string    `db:"content" json:"content"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	SenderName string    `db:"sender_name" json:"sender_name,omitempty"`
}

// InitiateChatRequest represents the data needed to start a chat
type InitiateChatRequest struct {
	RecipientID int    `json:"recipient_id" binding:"required"`
	ListingID   *int   `json:"listing_id"`
	Message     string `json:"message" binding:"required"`
}

// SendMessageRequest represents the data needed to send a message
type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// ChatResponse represents a list of chats with pagination
type ChatResponse struct {
	Chats       []Chat `json:"chats"`
	TotalCount  int    `json:"total_count"`
	CurrentPage int    `json:"current_page"`
	TotalPages  int    `json:"total_pages"`
}

// MessageResponse represents a list of messages with pagination
type MessageResponse struct {
	Messages    []Message `json:"messages"`
	TotalCount  int       `json:"total_count"`
	CurrentPage int       `json:"current_page"`
	TotalPages  int       `json:"total_pages"`
}
