package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/google/uuid"
)

func (t *task) CreateHistory(ctx context.Context, in *inputs.CreateHistory) (string, error) {

	const op = "adapters.db.mysql.tasks.create_history"

	tx := t.txExtractor.ExtractTx(ctx)

	snapshot, err := json.Marshal(in.OldTask)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	id := uuid.New().String()

	_, err = tx.ExecContext(ctx, createHistory, id, in.OldTask.ID, in.ChangedByID, snapshot)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return id, nil

}

const createHistory = "insert into task_history (id, task_id, changed_by, snapshot) values (?, ?, ?, ?);"
