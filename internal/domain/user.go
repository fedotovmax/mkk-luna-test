package domain

import (
	"time"
)

type User struct {
	CreatedAt    time.Time `json:"created_at" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at" validate:"required"`
	UserName     string    `json:"username" validate:"required"`
	Email        string    `json:"email" validate:"required"`
	PasswordHash string    `json:"password_hash" validate:"required"`
	ID           string    `json:"id" validate:"required"`
}

type BaseUser struct {
	ID       string `json:"id" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
}
