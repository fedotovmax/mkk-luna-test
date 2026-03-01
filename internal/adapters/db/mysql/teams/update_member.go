package teams

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (s *team) UpdateMember(ctx context.Context, id string, newRole domain.Role) error {

	const op = "adapters.db.mysql.teams.update_member"

	tx := s.txExtractor.ExtractTx(ctx)

	res, err := tx.ExecContext(ctx, updateMember, id, newRole)

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

const updateMember = "update team_members set role = ? where id = ?;"
