package domain

import (
	"time"
)

type Session struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	RefreshHash string
	ID          string
	UserID      string
}
