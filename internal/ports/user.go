package ports

import (
	"context"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/fields"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
)

type UserStorage interface {
	FindOne(ctx context.Context, field fields.UserField, value string) (*domain.User, error)
	Create(ctx context.Context, in *inputs.CreateUser) (string, error)
}
