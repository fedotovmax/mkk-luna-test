package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
)

type Invite struct {
	log     *slog.Logger
	storage ports.TeamStorage
	query   queries.Teams
}

func NewInvite(log *slog.Logger, storage ports.TeamStorage, query queries.Teams) *Invite {
	return &Invite{log: log, storage: storage, query: query}
}

func (u *Invite) Execute(ctx context.Context, inviterID string, teamID string, userID string) (string, error) {

	const op = "usecases.invite"

	member, err := u.query.FindMember(ctx, inviterID, teamID)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !member.CanInvite() {
		return "", fmt.Errorf("%s: %w", op, errs.ErrNoRightsToInviteMember)
	}

	_, err = u.query.FindMember(ctx, userID, teamID)

	if err != nil && !errors.Is(err, errs.ErrTeamMemberNotFound) {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err == nil {
		return "", fmt.Errorf("%s: %w", op, errs.ErrUserAlreadyInTeam)
	}

	memberID, err := u.storage.CreateMember(ctx, teamID, userID, domain.RoleMember)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return memberID, nil

}
