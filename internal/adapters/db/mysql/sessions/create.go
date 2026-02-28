package sessions

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
)

func (s *session) Create(ctx context.Context, in *inputs.CreateSession) error {

	const op = "adapters.db.mysql.sessions.create"

	tx := s.txExtractor.ExtractTx(ctx)

	_, err := tx.ExecContext(ctx, create, in.ID, in.UserID, in.RefreshHash, in.ExpiresAt)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
