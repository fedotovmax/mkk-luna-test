package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/cache"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
)

type GetTasks struct {
	log   *slog.Logger
	tasks queries.Tasks
	teams queries.Teams
	cache ports.TaskCache
}

func NewGetTasks(
	log *slog.Logger,
	tasks queries.Tasks,
	teams queries.Teams,
	cache ports.TaskCache,

) *GetTasks {
	return &GetTasks{
		log:   log,
		tasks: tasks,
		cache: cache,
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

	l := u.log.With(slog.String("op", op))

	_, err := u.teams.FindMember(ctx, userID, in.TeamID)

	if err != nil {
		if errors.Is(err, errs.ErrTeamMemberNotFound) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrUserNotInTaskTeam)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	key := cache.TasksListKey(limit, offset, in)
	res, err := u.cache.Get(ctx, key)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			res, err := u.tasks.FindMany(ctx, limit, offset, in)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			err = u.cache.Set(ctx, key, res)
			if err != nil {
				l.Error("failed to set tasks to cache", "error", err)
			}
			l.Info("TASKS FROM DB")
			return res, nil
		}
	}

	l.Info("TASKS FROM CACHE")
	return res, nil
}
