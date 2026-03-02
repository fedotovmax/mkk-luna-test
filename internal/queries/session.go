package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
)

type Sessions interface {
	FindOne(ctx context.Context, hash string) (*domain.Session, error)
}

type sessionDB struct {
	sessionStorage ports.SessionStorage
}

func NewSessions(sessionStorage ports.SessionStorage) Sessions {
	return &sessionDB{sessionStorage: sessionStorage}
}

func (q *sessionDB) FindOne(ctx context.Context, hash string) (*domain.Session, error) {

	s, err := q.sessionStorage.FindOne(ctx, hash)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrSessionNotFound, err)
		}
		return nil, err
	}

	return s, nil

}
