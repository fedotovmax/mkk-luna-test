package tasks

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (t *task) FindOne(ctx context.Context, id string) (*domain.Task, error) {

	const op = "adapters.db.mysql.tasks.find_one"

	tx := t.txExtractor.ExtractTx(ctx)

	row := tx.QueryRowContext(ctx, findOne, id)

	task := &domain.Task{}

	var (
		assigneeID       *string
		assigneeUsername *string
		assigneeEmail    *string
	)

	err := row.Scan(
		&task.ID,
		&task.TeamID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,

		&task.Owner.ID,
		&task.Owner.Username,
		&task.Owner.Email,

		&assigneeID,
		&assigneeUsername,
		&assigneeEmail,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, adapters.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	if assigneeID != nil {
		task.Assignee = &domain.BaseUser{
			ID:       *assigneeID,
			Username: *assigneeUsername,
			Email:    *assigneeEmail,
		}
	}

	return task, nil
}

const findOne = `
select
    t.id,
    t.team_id,
    t.title,
    t.description,
    t.status,
    t.created_at,
    t.updated_at,

    owner.id,
    owner.username,
    owner.email,

    assignee.id,
    assignee.username,
    assignee.email

from tasks t
join users owner on owner.id = t.created_by
left join users assignee on assignee.id = t.assignee_id

where t.id = ?;
`
