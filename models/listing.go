package models

import (
	"time"
)

// Listing представляет объявление о продаже мебели
type Listing struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price" json:"price"`
	Condition   string    `db:"condition" json:"condition"`
	City        string    `db:"city" json:"city"`
	CategoryID  int       `db:"category_id" json:"category_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Images      []Image   `json:"images,omitempty"`                   // Не хранится в БД, заполняется отдельным запросом
	UserName    string    `db:"user_name" json:"user_name,omitempty"` // Имя продавца
}

// Image представляет изображение для объявления
type Image struct {
	ID        int       `db:"id" json:"id"`
	ListingID int       `db:"listing_id" json:"listing_id"`
	ImagePath string    `db:"image_path" json:"image_path"`
	IsMain    bool      `db:"is_main" json:"is_main"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Category представляет категорию мебели
type Category struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// CreateListingRequest структура для создания нового объявления
type CreateListingRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Condition   string  `json:"condition" binding:"required"`
	City        string  `json:"city" binding:"required"`
	CategoryID  int     `json:"category_id" binding:"required"`
}

// UpdateListingRequest структура для обновления объявления
type UpdateListingRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"min=0"`
	Condition   string  `json:"condition"`
	City        string  `json:"city"`
	CategoryID  int     `json:"category_id"`
}

// ListingFilter структура для фильтрации объявлений
type ListingFilter struct {
	CategoryID *int     `form:"category_id"`
	City       string   `form:"city"`
	Condition  string   `form:"condition"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
	SortBy     string   `form:"sort_by" binding:"omitempty,oneof=date price -date -price"`
	Page       int      `form:"page,default=1" binding:"min=1"`
	Limit      int      `form:"limit,default=10" binding:"min=1,max=50"`
}

// Favorite представляет избранное объявление пользователя
type Favorite struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	ListingID int       `db:"listing_id" json:"listing_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
