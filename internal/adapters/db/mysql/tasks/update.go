package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
)

func (t *task) Update(ctx context.Context, id string, in *inputs.UpdateTask) error {

	const op = "adapters.db.mysql.tasks.update"

	tx := t.txExtractor.ExtractTx(ctx)

	q, args := buildUpdateTask(id, in)

	if q == "" {
		return fmt.Errorf("%s: %w", op, adapters.ErrNoFieldsToUpdate)
	}

	res, err := tx.ExecContext(ctx, q, args...)

	if err != nil {
		return fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	if affected == 0 {
		return fmt.Errorf("%s: %w", op, adapters.ErrNotFound)
	}

	return nil
}

func buildUpdateTask(id string, in *inputs.UpdateTask) (string, []any) {
	setParts := make([]string, 0)
	args := make([]any, 0)

	if in.Title != nil {
		setParts = append(setParts, "title = ?")
		args = append(args, *in.Title)
	}

	if in.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *in.Description)
	}

	if in.Role != nil {
		setParts = append(setParts, "status = ?")
		args = append(args, *in.Role)
	}

	if len(setParts) == 0 {
		return "", nil
	}

	query := fmt.Sprintf(`
update tasks
set %s
where id = ?;
`, strings.Join(setParts, ", "))

	args = append(args, id)

	return query, args
}
