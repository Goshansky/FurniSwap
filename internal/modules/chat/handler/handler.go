package handler

import (
	"FurniSwap/internal/modules/chat/model"
	"FurniSwap/internal/modules/chat/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler provides chat handlers
type Handler struct {
	service *service.Service
}

// NewHandler creates a new chat handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers chat routes to router
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/chats", h.InitiateChat)
	router.GET("/chats", h.GetChats)
	router.GET("/chats/:id", h.GetChatMessages)
	router.POST("/chats/:id/messages", h.SendMessage)
}

// InitiateChat handles starting a new chat
func (h *Handler) InitiateChat(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var req model.InitiateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error parsing chat request: %v, body: %v", err, c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Логируем данные запроса для отладки
	log.Printf("Initiating chat: user %d with recipient %d, listing %v, message: %s",
		userID.(int), req.RecipientID, req.ListingID, req.Message)

	// Initiate chat
	chatID, err := h.service.InitiateChat(userID.(int), req)
	if err != nil {
		if err.Error() == "cannot initiate chat with yourself" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot initiate chat with yourself"})
			return
		}
		log.Printf("Error initiating chat: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error initiating chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": chatID, "message": "Chat initiated successfully"})
}

// GetChats handles getting the user's chats
func (h *Handler) GetChats(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// Get chats
	chats, err := h.service.GetUserChats(userID.(int), page, limit)
	if err != nil {
		log.Printf("Error getting chats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// GetChatMessages handles getting messages in a chat
func (h *Handler) GetChatMessages(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse chat ID
	chatID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// Get chat info first
	chat, err := h.service.GetChatByID(chatID, userID.(int))
	if err != nil {
		if err.Error() == "you don't have access to this chat" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this chat"})
			return
		}
		log.Printf("Error getting chat: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chat"})
		return
	}

	// Get messages
	messages, err := h.service.GetChatMessages(chatID, userID.(int), page, limit)
	if err != nil {
		if err.Error() == "you don't have access to this chat" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this chat"})
			return
		}
		log.Printf("Error getting chat messages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chat messages"})
		return
	}

	// Combine chat info with messages
	response := gin.H{
		"chat":     chat,
		"messages": messages.Messages,
		"pagination": gin.H{
			"total_count":  messages.TotalCount,
			"current_page": messages.CurrentPage,
			"total_pages":  messages.TotalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

// SendMessage handles sending a message in a chat
func (h *Handler) SendMessage(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse chat ID
	chatID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	// Parse request body
	var req model.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Send message
	messageID, err := h.service.SendMessage(chatID, userID.(int), req)
	if err != nil {
		if err.Error() == "you don't have access to this chat" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this chat"})
			return
		}
		log.Printf("Error sending message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": messageID, "message": "Message sent successfully"})
}
