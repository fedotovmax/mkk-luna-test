package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/fields"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
)

type Users interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type usersDB struct {
	usersStorage ports.UserStorage
}

func NewUsers(usersStorage ports.UserStorage) Users {
	return &usersDB{
		usersStorage: usersStorage,
	}
}

func (q *usersDB) FindByEmail(ctx context.Context, email string) (*domain.User, error) {

	user, err := q.usersStorage.FindOne(ctx, fields.UserFieldEmail, email)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, err
	}

	return user, nil
}

func (q *usersDB) FindByID(ctx context.Context, id string) (*domain.User, error) {

	user, err := q.usersStorage.FindOne(ctx, fields.UserFieldID, id)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, err
	}

	return user, nil
}
