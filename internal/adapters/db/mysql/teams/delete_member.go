package teams

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
)

func (t *team) DeleteMember(ctx context.Context, id string) error {

	const op = "adapters.db.mysql.teams.delete_member"

	tx := t.txExtractor.ExtractTx(ctx)

	res, err := tx.ExecContext(ctx, deleteMember, id)

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

const deleteMember = "delete from team_members where id = ?;"
