package models

import "time"

type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	IsVerified   bool      `db:"is_verified"`
	CreatedAt    time.Time `db:"created_at"`
}

type TwoFactorCode struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Code      string    `db:"code"`
	ExpiresAt time.Time `db:"expires_at"`
}
