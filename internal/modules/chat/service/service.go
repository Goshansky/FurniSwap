package service

import (
	"FurniSwap/internal/modules/chat/model"
	"FurniSwap/internal/modules/chat/repository"
	"errors"
	"fmt"
)

// Service provides chat operations
type Service struct {
	repo *repository.Repository
}

// NewService creates a new chat service
func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// InitiateChat starts a new chat or retrieves an existing one
func (s *Service) InitiateChat(userID int, req model.InitiateChatRequest) (int, error) {
	// Don't allow chat with self
	if userID == req.RecipientID {
		return 0, errors.New("cannot initiate chat with yourself")
	}

	// Validate recipient user exists (to be implemented)
	// This would check for the existence of the recipient in the user database

	// Check if chat already exists
	existingChat, err := s.repo.GetChatByUsers(userID, req.RecipientID, req.ListingID)
	if err != nil {
		return 0, fmt.Errorf("error checking existing chat: %w", err)
	}

	var chatID int
	if existingChat != nil {
		chatID = existingChat.ID
	} else {
		// Create new chat
		chatID, err = s.repo.CreateChat(userID, req.RecipientID, req.ListingID)
		if err != nil {
			return 0, fmt.Errorf("error creating chat: %w", err)
		}
	}

	// Add initial message
	_, err = s.repo.AddMessage(chatID, userID, req.Message)
	if err != nil {
		return 0, fmt.Errorf("error adding message: %w", err)
	}

	return chatID, nil
}

// SendMessage sends a message in a chat
func (s *Service) SendMessage(chatID, userID int, req model.SendMessageRequest) (int, error) {
	// Check if user has access to the chat
	hasAccess, err := s.repo.CheckChatAccess(chatID, userID)
	if err != nil {
		return 0, fmt.Errorf("error checking chat access: %w", err)
	}

	if !hasAccess {
		return 0, errors.New("you don't have access to this chat")
	}

	// Add message
	messageID, err := s.repo.AddMessage(chatID, userID, req.Content)
	if err != nil {
		return 0, fmt.Errorf("error adding message: %w", err)
	}

	return messageID, nil
}

// GetUserChats gets all chats for a user
func (s *Service) GetUserChats(userID, page, limit int) (*model.ChatResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	return s.repo.GetUserChats(userID, page, limit)
}

// GetChatMessages gets all messages in a chat
func (s *Service) GetChatMessages(chatID, userID, page, limit int) (*model.MessageResponse, error) {
	// Check if user has access to the chat
	hasAccess, err := s.repo.CheckChatAccess(chatID, userID)
	if err != nil {
		return nil, fmt.Errorf("error checking chat access: %w", err)
	}

	if !hasAccess {
		return nil, errors.New("you don't have access to this chat")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 50
	}
	return s.repo.GetChatMessages(chatID, page, limit)
}

// GetChatByID gets a chat by ID
func (s *Service) GetChatByID(chatID, userID int) (*model.Chat, error) {
	// Check if user has access to the chat
	hasAccess, err := s.repo.CheckChatAccess(chatID, userID)
	if err != nil {
		return nil, fmt.Errorf("error checking chat access: %w", err)
	}

	if !hasAccess {
		return nil, errors.New("you don't have access to this chat")
	}

	return s.repo.GetChatByID(chatID)
}
