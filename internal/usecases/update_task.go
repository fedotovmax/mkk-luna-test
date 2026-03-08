package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/transaction/mysql"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
)

type UpdateTask struct {
	log     *slog.Logger
	tx      mysql.Manager
	storage ports.TaskStorage
	tasks   queries.Tasks
	teams   queries.Teams
}

func NewUpdateTask(
	log *slog.Logger,
	tx mysql.Manager,
	storage ports.TaskStorage,
	tasks queries.Tasks,
	teams queries.Teams,
) *UpdateTask {
	return &UpdateTask{
		log:     log,
		tx:      tx,
		storage: storage,
		tasks:   tasks,
		teams:   teams,
	}
}

func (u *UpdateTask) Execute(ctx context.Context, userID string, taskID string, in *inputs.UpdateTask) error {

	const op = "usecases.update_task"

	return u.tx.Wrap(ctx, func(txctx context.Context) error {

		task, err := u.tasks.FindByID(txctx, taskID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		member, err := u.teams.FindMember(txctx, userID, task.TeamID)

		if err != nil {
			if errors.Is(err, errs.ErrTeamMemberNotFound) {
				return fmt.Errorf("%s: %w", op, errs.ErrUserNotInTaskTeam)
			}
			return fmt.Errorf("%s: %w", op, err)
		}

		if !member.CanUpdateTask(task) {
			return fmt.Errorf("%s: %w", op, errs.ErrNoRightsToUpdateTask)
		}

		historyInput := &inputs.CreateHistory{
			OldTask:     task,
			ChangedByID: userID,
		}

		_, err = u.storage.CreateHistory(txctx, historyInput)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = u.storage.Update(txctx, taskID, in)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})
}
