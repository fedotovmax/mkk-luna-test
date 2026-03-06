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

type CreateTask struct {
	log     *slog.Logger
	tx      mysql.Manager
	storage ports.TaskStorage
	teams   queries.Teams
}

func NewCreateTask(
	log *slog.Logger,
	tx mysql.Manager,
	storage ports.TaskStorage,
	teams queries.Teams,
) *CreateTask {
	return &CreateTask{
		log:     log,
		tx:      tx,
		storage: storage,
		teams:   teams,
	}
}

func (u *CreateTask) Execute(ctx context.Context, ownerID string, in *inputs.CreateTask) (string, error) {

	const op = "usecases.create_task"

	member, err := u.teams.FindMember(ctx, ownerID, in.TeamID)

	if err != nil {
		if errors.Is(err, errs.ErrTeamMemberNotFound) {
			return "", fmt.Errorf("%s: %w", op, errs.ErrUserNotInTaskTeam)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !member.CanCreateTask() {
		return "", fmt.Errorf("%s: %w", op, errs.ErrNoRightsToCreateTask)
	}

	taskID, err := u.storage.Create(ctx, ownerID, in)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return taskID, nil
}
