package main

import (
	"FurniSwap/pkg/config"
	"FurniSwap/pkg/database"
	"FurniSwap/pkg/middleware"

	// Auth module
	authHandler "FurniSwap/internal/modules/auth/handler"
	authRepo "FurniSwap/internal/modules/auth/repository"
	authService "FurniSwap/internal/modules/auth/service"

	// Profile module
	profileHandler "FurniSwap/internal/modules/profile/handler"
	profileRepo "FurniSwap/internal/modules/profile/repository"
	profileService "FurniSwap/internal/modules/profile/service"

	// Listing module
	listingHandler "FurniSwap/internal/modules/listing/handler"
	listingRepo "FurniSwap/internal/modules/listing/repository"
	listingService "FurniSwap/internal/modules/listing/service"

	// Favorite module
	favoriteHandler "FurniSwap/internal/modules/favorite/handler"
	favoriteRepo "FurniSwap/internal/modules/favorite/repository"
	favoriteService "FurniSwap/internal/modules/favorite/service"

	// Purchase module
	purchaseHandler "FurniSwap/internal/modules/purchase/handler"
	purchaseRepo "FurniSwap/internal/modules/purchase/repository"
	purchaseService "FurniSwap/internal/modules/purchase/service"

	// Chat module
	chatHandler "FurniSwap/internal/modules/chat/handler"
	chatRepo "FurniSwap/internal/modules/chat/repository"
	chatService "FurniSwap/internal/modules/chat/service"

	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.Load()

	// Initialize database
	db, err := database.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Check and update database structure
	database.UpdateDatabaseStructure(db)
	database.CheckDatabase(db)

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatalf("Error creating uploads directory: %v", err)
	}

	// Initialize router with default middleware
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.Config.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Serve static files
	r.Static("/uploads", "./uploads")

	// Initialize module repositories
	authRepository := authRepo.NewRepository(db)
	profileRepository := profileRepo.NewRepository(db)
	listingRepository := listingRepo.NewRepository(db)
	favoriteRepository := favoriteRepo.NewRepository(db)
	purchaseRepository := purchaseRepo.NewRepository(db)
	chatRepository := chatRepo.NewRepository(db)

	// Initialize module services
	authSvc := authService.NewService(authRepository)
	profileSvc := profileService.NewService(profileRepository)
	listingSvc := listingService.NewService(listingRepository)
	favoriteSvc := favoriteService.NewService(favoriteRepository)
	purchaseSvc := purchaseService.NewService(purchaseRepository, listingRepository)
	chatSvc := chatService.NewService(chatRepository)

	// Initialize module handlers
	authHandler := authHandler.NewHandler(authSvc)
	profileHandler := profileHandler.NewHandler(profileSvc)
	listingHandler := listingHandler.NewHandler(listingSvc)
	favoriteHandler := favoriteHandler.NewHandler(favoriteSvc)
	purchaseHandler := purchaseHandler.NewHandler(purchaseSvc)
	chatHandler := chatHandler.NewHandler(chatSvc)

	// Register public routes (no auth required)
	authHandler.RegisterRoutes(r.Group(""))

	// User public endpoints
	r.GET("/users/:id", func(c *gin.Context) {
		// Получаем ID пользователя из URL-параметра
		userIDStr := c.Param("id")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Printf("Error converting user ID: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Получаем профиль пользователя из БД
		var userProfile struct {
			ID        int       `db:"id" json:"id"`
			Name      string    `db:"name" json:"name"`
			LastName  string    `db:"last_name" json:"last_name"`
			City      string    `db:"city" json:"city"`
			Avatar    string    `db:"avatar" json:"avatar"`
			CreatedAt time.Time `db:"created_at" json:"created_at"`
		}
		query := `
			SELECT id, name, last_name, city, avatar, created_at
			FROM users
			WHERE id = $1
		`
		err = db.Get(&userProfile, query, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("User with ID %d not found\n", userID)
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			log.Printf("Error getting user with ID %d: %v\n", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user data"})
			return
		}

		c.JSON(http.StatusOK, userProfile)
	})

	// Categories endpoint (public)
	r.GET("/categories", func(c *gin.Context) {
		var categories []struct {
			ID   int    `db:"id" json:"id"`
			Name string `db:"name" json:"name"`
		}
		err := db.Select(&categories, "SELECT id, name FROM categories ORDER BY name")
		if err != nil {
			log.Printf("Error getting categories: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting categories"})
			return
		}
		c.JSON(http.StatusOK, categories)
	})

	// Public listing routes
	publicListings := r.Group("/listings")
	listingHandler.RegisterPublicRoutes(publicListings)

	// Protected API routes (auth required)
	api := r.Group("/api")
	api.Use(middleware.AuthRequired(db))
	{
		// Register module routes to protected API group
		profileHandler.RegisterRoutes(api)
		listingHandler.RegisterProtectedRoutes(api)
		favoriteHandler.RegisterRoutes(api)
		purchaseHandler.RegisterRoutes(api)
		chatHandler.RegisterRoutes(api)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + config.Config.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server started on http://localhost:%s", config.Config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
