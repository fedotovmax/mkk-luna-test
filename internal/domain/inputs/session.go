package inputs

import "time"

type CreateSession struct {
	ID          string
	UserID      string
	RefreshHash string
	ExpiresAt   time.Time
}
