package tasks

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (t *task) FindTaskComments(ctx context.Context, taskID string) ([]*domain.Comment, error) {

	const op = "adapters.db.mysql.tasks.find_comments"

	tx := t.txExtractor.ExtractTx(ctx)

	rows, err := tx.QueryContext(ctx, findTaskComments, taskID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}
	defer rows.Close()

	comments := make([]*domain.Comment, 0)

	for rows.Next() {
		c := &domain.Comment{}

		err := rows.Scan(
			&c.ID,
			&c.TaskID,
			&c.Comment,
			&c.CreatedAt,

			&c.User.ID,
			&c.User.Username,
			&c.User.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
		}

		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	return comments, nil
}

const findTaskComments = `
select
    c.id,
    c.task_id,
    c.comment,
    c.created_at,

    u.id,
    u.username,
    u.email

from task_comments c
join users u on u.id = c.user_id
where c.task_id = ?
order by c.created_at asc
`
