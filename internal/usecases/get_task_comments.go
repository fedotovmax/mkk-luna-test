package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
)

type GetTaskComments struct {
	log   *slog.Logger
	tasks queries.Tasks
	teams queries.Teams
}

func NewGetTaskComments(
	log *slog.Logger,
	tasks queries.Tasks,
	teams queries.Teams,
) *GetTaskComments {
	return &GetTaskComments{
		log:   log,
		tasks: tasks,
		teams: teams,
	}
}

func (u *GetTaskComments) Execute(
	ctx context.Context,
	userID string,
	taskID string,
) ([]*domain.Comment, error) {

	const op = "usecases.get_task_comments"

	task, err := u.tasks.FindByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrTaskNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = u.teams.FindMember(ctx, userID, task.TeamID)
	if err != nil {
		if errors.Is(err, errs.ErrTeamMemberNotFound) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrUserNotInTaskTeam)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	comments, err := u.tasks.FindComments(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}
