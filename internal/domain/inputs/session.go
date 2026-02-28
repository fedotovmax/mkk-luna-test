package inputs

type CreateSession struct {
	ID          string
	UserID      string
	RefreshHash string
	ExpiresAt   string
}
