package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/transaction/mysql"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
)

type CreateTeam struct {
	log     *slog.Logger
	tx      mysql.Manager
	storage ports.TeamStorage
	query   queries.Teams
}

func NewCreateTeam(
	log *slog.Logger,
	tx mysql.Manager,
	storage ports.TeamStorage,
	query queries.Teams,
) *CreateTeam {
	return &CreateTeam{
		log:     log,
		tx:      tx,
		storage: storage,
		query:   query,
	}
}

func (u *CreateTeam) Execute(ctx context.Context, ownerID string, in *inputs.CreateTeam) (string, error) {

	const op = "usecases.create_team"

	var teamID string

	err := u.tx.Wrap(ctx, func(txctx context.Context) error {

		var err error

		_, err = u.query.FindByName(txctx, in.Name)

		if err != nil && !errors.Is(err, errs.ErrTeamNotFound) {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err == nil {
			return fmt.Errorf("%s: %w", op, errs.ErrTeamAlreadyExists)
		}

		teamID, err = u.storage.Create(txctx, ownerID, in.Name)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		_, err = u.storage.CreateMember(txctx, teamID, ownerID, domain.RoleOwner)

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil

	})

	if err != nil {
		return "", err
	}

	return teamID, err
}
