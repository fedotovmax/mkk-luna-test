package teams

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
)

func (t *team) Delete(ctx context.Context, ownerID string, teamID string) error {
	const op = "adapters.db.mysql.teams.delete"

	tx := t.txExtractor.ExtractTx(ctx)

	res, err := tx.ExecContext(ctx, delete, teamID, ownerID)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	if affected == 0 {
		return fmt.Errorf("%s: %w", op, adapters.ErrNotFound)
	}

	return nil
}

const delete = "delete from teams where team_id = ? and created_by = ?;"
