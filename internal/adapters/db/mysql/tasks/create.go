package tasks

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/google/uuid"
)

func (t *task) Create(ctx context.Context, ownerID string, in *inputs.CreateTask) (string, error) {

	const op = "adapters.db.mysql.tasks.create"

	tx := t.txExtractor.ExtractTx(ctx)

	id := uuid.New().String()

	_, err := tx.ExecContext(ctx, create, id, in.TeamID, in.Title, in.Description, in.AssigneeID, ownerID)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return id, nil

}

const create = `
insert into tasks 
(id, team_id, title, description, assignee_id, created_by)
values (?, ?, ?, ?, ?, ?);
`
