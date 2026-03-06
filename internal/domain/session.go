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
	AccessToken    string    `json:"access_token" validate:"required" example:"sfjosdofj0293j4029jehsohsdofisdhoifhsdo"`
	RefreshToken   string    `json:"refresh_token" validate:"required" example:"sokdf11gpodjio23h9hsdoifhnso1321"`
	AccessExpTime  time.Time `json:"access_exp_time" validate:"required" example:"2026-03-05T09:50:42.108996Z"`
	RefreshExpTime time.Time `json:"refresh_exp_time" validate:"required" example:"2026-03-05T09:50:42.108996Z"`
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
