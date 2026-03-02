package domain

import (
	"time"
)

type Local struct {
	UserID    string
	SessionID string
}

type Session struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   time.Time
	RefreshHash string
	ID          string
	UserID      string
}

type LoginResponse struct {
	AccessToken    string
	RefreshToken   string
	AccessExpTime  time.Time
	RefreshExpTime time.Time
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
