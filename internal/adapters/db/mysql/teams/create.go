package teams

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/google/uuid"
)

func (t *team) Create(ctx context.Context, ownerID, name string) (string, error) {

	const op = "adapters.db.mysql.teams.create"

	tx := t.txExtractor.ExtractTx(ctx)

	id := uuid.New().String()

	_, err := tx.ExecContext(ctx, create, id, name, ownerID)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return id, nil
}
