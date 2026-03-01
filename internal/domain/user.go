package domain

import (
	"time"
)

type User struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserName     string
	Email        string
	PasswordHash string
	ID           string
}

type BaseUser struct {
	ID       string
	Username string
	Email    string
}
