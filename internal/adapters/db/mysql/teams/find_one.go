package teams

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/fields"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (t *team) FindOne(ctx context.Context, field fields.TeamField, value string) (*domain.Team, error) {

	const op = "adapters.db.mysql.teams.find_one"

	err := fields.IsTeamEntityField(field)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tx := t.txExtractor.ExtractTx(ctx)

	rows, err := tx.QueryContext(ctx, findByQuery(field), value)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	defer rows.Close()

	var team *domain.Team
	memberMap := make(map[string]struct{})

	for rows.Next() {

		if team == nil {
			team = &domain.Team{}
		}

		var (
			memberID       sql.NullString
			memberUsername sql.NullString
			memberEmail    sql.NullString

			memberRowID sql.NullString
			role        sql.NullString
			joinedAt    sql.NullTime
		)

		err = rows.Scan(
			&team.ID,
			&team.Name,
			&team.CreatedAt,
			&team.UpdatedAt,

			&team.Owner.ID,
			&team.Owner.Username,
			&team.Owner.Email,

			&memberRowID,
			&role,
			&joinedAt,

			&memberID,
			&memberUsername,
			&memberEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, adapters.ErrInternal)
		}

		if memberID.Valid {
			if _, exists := memberMap[memberRowID.String]; !exists {

				member := domain.Member{
					ID:       memberRowID.String,
					Role:     domain.Role(role.String),
					JoinedAt: joinedAt.Time,
					User: domain.BaseUser{
						ID:       memberID.String,
						Username: memberUsername.String,
						Email:    memberEmail.String,
					},
				}

				team.Members = append(team.Members, member)
				memberMap[memberRowID.String] = struct{}{}
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	if team == nil {
		return nil, fmt.Errorf("%s: %s: %w", op, field, adapters.ErrNotFound)
	}

	return team, nil
}

func findByQuery(field fields.TeamField) string {
	return fmt.Sprintf(`
		select
			t.id,
			t.name,
			t.created_at,
			t.updated_at,

			owner.id,
			owner.username,
			owner.email,

			tm.id,
			tm.role,
			tm.joined_at,

			member.id,
			member.username,
			member.email

		from teams t
		join users owner on owner.id = t.created_by
		left join team_members tm on tm.team_id = t.id
		left join users member on member.id = tm.user_id
		where t.%s = ?
	`, field)
}
