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

const findMember = `
select
    tm.id,
    tm.role,
    tm.joined_at,
    u.id,
    u.username,
    u.email
from team_members tm
join users u on u.id = tm.user_id
where tm.user_id = ? and tm.team_id = ?;`
