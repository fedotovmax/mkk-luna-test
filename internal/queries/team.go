package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/fields"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
)

type Teams interface {
	FindByID(ctx context.Context, id string) (*domain.Team, error)
	FindByName(ctx context.Context, name string) (*domain.Team, error)
	FindMany(ctx context.Context, page, pageSize int, uid string) (*domain.FindTeamsResponse, error)
	FindMember(ctx context.Context, userID, teamID string) (*domain.Member, error)
	Stats(ctx context.Context) ([]domain.TeamStats, error)
	TopUsers(ctx context.Context) ([]domain.TopUserInTeam, error)
}

type teams struct {
	storage ports.TeamStorage
}

func NewTeams(storage ports.TeamStorage) Teams {
	return &teams{storage: storage}
}

func (q *teams) FindByID(ctx context.Context, id string) (*domain.Team, error) {
	t, err := q.storage.FindOne(ctx, fields.TeamFieldID, id)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrTeamNotFound, err)
		}
		return nil, err
	}

	return t, nil
}

func (q *teams) FindByName(ctx context.Context, name string) (*domain.Team, error) {
	t, err := q.storage.FindOne(ctx, fields.TeamFieldName, name)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrTeamNotFound, err)
		}
		return nil, err
	}

	return t, nil
}

func (q *teams) FindMany(ctx context.Context, offset, limit int, uid string) (
	*domain.FindTeamsResponse, error) {

	res, err := q.storage.FindMany(ctx, limit, offset, uid)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (q *teams) FindMember(ctx context.Context, userID, teamID string) (*domain.Member, error) {

	m, err := q.storage.FindMember(ctx, userID, teamID)

	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			return nil, fmt.Errorf("%w: %v", errs.ErrTeamMemberNotFound, err)
		}
		return nil, err
	}

	return m, nil
}

func (q *teams) Stats(ctx context.Context) ([]domain.TeamStats, error) {
	return q.storage.Stats(ctx)
}

func (q *teams) TopUsers(ctx context.Context) ([]domain.TopUserInTeam, error) {
	return q.storage.TopUsers(ctx)
}
