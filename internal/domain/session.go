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

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
