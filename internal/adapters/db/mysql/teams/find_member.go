package teams

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (t *team) FindMember(ctx context.Context, userID, teamID string) (*domain.Member, error) {

	const op = "adapters.db.mysql.teams.find_member"

	tx := t.txExtractor.ExtractTx(ctx)

	row := tx.QueryRowContext(ctx, findMember, userID, teamID)

	member := &domain.Member{}

	err := row.Scan(
		&member.ID,
		&member.Role,
		&member.JoinedAt,
		&member.User.ID,
		&member.User.Username,
		&member.User.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, adapters.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	return member, nil
}
