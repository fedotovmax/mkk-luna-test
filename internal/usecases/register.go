package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/auth/password"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
)

type Register struct {
	log     *slog.Logger
	storage ports.UserStorage
	query   queries.Users
}

func NewRegister(log *slog.Logger, storage ports.UserStorage, query queries.Users) *Register {
	return &Register{
		log:     log,
		storage: storage,
		query:   query,
	}
}

func (u *Register) Execute(ctx context.Context, in *inputs.CreateUser) (string, error) {

	const op = "usecases.register"

	_, err := u.query.FindByEmail(ctx, in.Email)

	if err != nil && !errors.Is(err, errs.ErrUserNotFound) {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err == nil {
		return "", fmt.Errorf("%s: %w", op, errs.ErrUserAlreadyExists)
	}

	hashedPassword, err := password.HashPassword(in.Password)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	in.SetPasswordHash(hashedPassword)

	id, err := u.storage.Create(ctx, in)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}
