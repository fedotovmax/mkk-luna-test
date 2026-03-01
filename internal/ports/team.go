package ports

import (
	"context"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/fields"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

type TeamStorage interface {
	FindMany(ctx context.Context, limit, offset int, userID string) ([]*domain.Team, error)
	FindMember(ctx context.Context, userID, teamID string) (*domain.Member, error)
	FindOne(ctx context.Context, field fields.TeamField, value string) (*domain.Team, error)
	Stats(ctx context.Context) ([]domain.TeamStats, error)
	TopUsers(ctx context.Context) ([]domain.TopUserInTeam, error)

	Create(ctx context.Context, ownerID, name string) (string, error)
	CreateMember(ctx context.Context, teamID, userID string, role domain.Role) (string, error)

	UpdateMember(ctx context.Context, id string, newRole domain.Role) error

	DeleteMember(ctx context.Context, id string) error
	Delete(ctx context.Context, ownerID string, teamID string) error
}
