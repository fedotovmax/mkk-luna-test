package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
)

type GetTasks struct {
	log   *slog.Logger
	tasks queries.Tasks
	teams queries.Teams
}

func NewGetTasks(
	log *slog.Logger,
	tasks queries.Tasks,
	teams queries.Teams,
) *GetTasks {
	return &GetTasks{
		log:   log,
		tasks: tasks,
		teams: teams,
	}
}

func (u *GetTasks) Execute(
	ctx context.Context,
	userID string,
	limit int,
	offset int,
	in *inputs.FindManyTasks,
) (*domain.FindTasksResponse, error) {

	const op = "usecases.get_tasks"

	_, err := u.teams.FindMember(ctx, userID, in.TeamID)
	if err != nil {

		if errors.Is(err, errs.ErrTeamMemberNotFound) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrUserNotInTaskTeam)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := u.tasks.FindMany(ctx, limit, offset, in)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
