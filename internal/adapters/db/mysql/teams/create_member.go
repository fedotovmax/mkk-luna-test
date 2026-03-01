package teams

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/google/uuid"
)

func (t *team) CreateMember(ctx context.Context, teamID, userID string, role domain.Role) (string, error) {

	const op = "adapters.db.mysql.teams.create_member"

	tx := t.txExtractor.ExtractTx(ctx)

	id := uuid.New().String()

	_, err := tx.ExecContext(ctx, createMember, id, teamID, userID, role)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return id, nil
}

const createMember = "insert into team_members (id, team_id, user_id, role) values (?, ?, ?, ?);"
