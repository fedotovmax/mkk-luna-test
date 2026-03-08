package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
)

type CreateComment struct {
	log     *slog.Logger
	storage ports.TaskStorage
	tasks   queries.Tasks
	teams   queries.Teams
}

func NewCreateComment(
	log *slog.Logger,
	storage ports.TaskStorage,
	tasks queries.Tasks,
	teams queries.Teams,
) *CreateComment {
	return &CreateComment{
		log:     log,
		storage: storage,
		tasks:   tasks,
		teams:   teams,
	}
}

func (u *CreateComment) Execute(ctx context.Context, userID, taskID, text string) (string, error) {
	const op = "usecases.create_comment"

	task, err := u.tasks.FindByID(ctx, taskID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = u.teams.FindMember(ctx, userID, task.TeamID)

	if err != nil {
		if errors.Is(err, errs.ErrTeamMemberNotFound) {
			return "", fmt.Errorf("%s: %w", op, errs.ErrUserNotInTaskTeam)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	commentID, err := u.storage.CreateComment(ctx, userID, taskID, text)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return commentID, nil
}
