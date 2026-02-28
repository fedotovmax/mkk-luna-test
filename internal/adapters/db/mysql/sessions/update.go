package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
)

func (s *session) Update(
	ctx context.Context,
	id string,
	newHash string,
	newExpires time.Time,
) error {

	const op = "adapters.db.mysql.sessions.update"

	tx := s.txExtractor.ExtractTx(ctx)

	_, err := tx.ExecContext(ctx, update, newHash, newExpires, id)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
