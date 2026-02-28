package sessions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (s *session) FindOne(ctx context.Context, hash string) (*domain.Session, error) {

	const op = "adapters.db.mysql.sessions.find_one"

	tx := s.txExtractor.ExtractTx(ctx)

	sess := &domain.Session{}

	row := tx.QueryRowContext(ctx, findOne, hash)

	err := row.Scan(
		&sess.ID,
		&sess.UserID,
		&sess.RefreshHash,
		&sess.CreatedAt,
		&sess.UpdatedAt,
		&sess.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return sess, nil

}
