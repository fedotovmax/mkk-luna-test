package ports

import (
	"context"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
)

type SessionStorage interface {
	FindOne(ctx context.Context, hash string) (*domain.Session, error)
	Create(ctx context.Context, in *inputs.CreateSession) error
	Update(ctx context.Context, id string, newHash string, newExpires time.Time) error
}
