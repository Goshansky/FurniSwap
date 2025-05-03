package models

import "time"

type User struct {
	ID           int       `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Name         string    `db:"name" json:"name"`
	LastName     string    `db:"last_name" json:"last_name"`
	City         string    `db:"city" json:"city"`
	Avatar       string    `db:"avatar" json:"avatar"`
	IsVerified   bool      `db:"is_verified" json:"is_verified"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type TwoFactorCode struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Code      string    `db:"code"`
	ExpiresAt time.Time `db:"expires_at"`
}

// UserProfile представляет профиль пользователя для публичного доступа
type UserProfile struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	City      string    `json:"city"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}

// UserResponse возвращается при успешной аутентификации
type UserResponse struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Token    string `json:"token"`
}
