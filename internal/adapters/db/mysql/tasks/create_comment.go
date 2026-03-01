package tasks

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/google/uuid"
)

func (t *task) CreateComment(ctx context.Context, in *inputs.CreateComment) (string, error) {

	const op = "adapters.db.mysql.tasks.create_comment"

	tx := t.txExtractor.ExtractTx(ctx)

	id := uuid.NewString()

	_, err := tx.ExecContext(
		ctx,
		createComment,
		id,
		in.TaskID,
		in.UserID,
		in.Text,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	return id, nil
}

const createComment = `
insert into task_comments (id, task_id, user_id, comment) values (?, ?, ?, ?);
`
