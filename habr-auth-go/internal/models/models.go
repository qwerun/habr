package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Nickname     string    `json:"nickname"`
}

func NewUser(email, passwordHash, nickname string) *User {
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		Nickname:     nickname,
	}
}
