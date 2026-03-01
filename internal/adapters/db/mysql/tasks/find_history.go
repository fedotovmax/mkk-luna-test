package tasks

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (t *task) FindTaskHistory(ctx context.Context, taskID string) ([]*domain.History, error) {
	const op = "adapters.db.mysql.tasks.find_history"

	tx := t.txExtractor.ExtractTx(ctx)

	rows, err := tx.QueryContext(ctx, findTaskHistory, taskID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	defer rows.Close()

	historyList := make([]*domain.History, 0)

	for rows.Next() {

		h := &domain.History{}

		err := rows.Scan(
			&h.ID,
			&h.TaskID,
			&h.Shapshot,
			&h.ChangedAt,

			&h.ChangedBy.ID,
			&h.ChangedBy.Username,
			&h.ChangedBy.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
		}

		historyList = append(historyList, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	return historyList, nil
}

const findTaskHistory = `
select
    h.id,
    h.task_id,
    h.snapshot,
    h.changed_at,

    u.id,
    u.username,
    u.email

from task_history h
join users u on u.id = h.changed_by
where h.task_id = ?
order by h.changed_at asc;
`
