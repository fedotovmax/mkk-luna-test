package teams

import (
	"context"
	"database/sql"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (t *team) FindMany(ctx context.Context, limit, offset int, userID string) ([]*domain.Team, error) {
	const op = "adapters.db.mysql.teams.find_many"

	tx := t.txExtractor.ExtractTx(ctx)

	q, args := buildFindMany(limit, offset, userID)

	rows, err := tx.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teamMap := make(map[string]*domain.Team)
	memberMap := make(map[string]map[string]struct{})
	teamOrder := make([]string, 0, limit)

	for rows.Next() {

		var (
			teamID   string
			teamName string
			created  time.Time
			updated  time.Time

			ownerID       string
			ownerUsername string
			ownerEmail    string

			memberRowID sql.NullString
			role        sql.NullString
			joinedAt    sql.NullTime

			memberID       sql.NullString
			memberUsername sql.NullString
			memberEmail    sql.NullString
		)

		err := rows.Scan(
			&teamID,
			&teamName,
			&created,
			&updated,

			&ownerID,
			&ownerUsername,
			&ownerEmail,

			&memberRowID,
			&role,
			&joinedAt,

			&memberID,
			&memberUsername,
			&memberEmail,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := teamMap[teamID]; !exists {
			teamMap[teamID] = &domain.Team{
				ID:        teamID,
				Name:      teamName,
				CreatedAt: created,
				UpdatedAt: updated,
				Owner: domain.TeamUser{
					ID:       ownerID,
					Username: ownerUsername,
					Email:    ownerEmail,
				},
			}
			memberMap[teamID] = make(map[string]struct{})
			teamOrder = append(teamOrder, teamID)
		}

		if memberRowID.Valid {
			if _, exists := memberMap[teamID][memberRowID.String]; !exists {

				member := domain.Member{
					ID:       memberRowID.String,
					Role:     domain.Role(role.String),
					JoinedAt: joinedAt.Time,
					User: domain.TeamUser{
						ID:       memberID.String,
						Username: memberUsername.String,
						Email:    memberEmail.String,
					},
				}

				teamMap[teamID].Members = append(teamMap[teamID].Members, member)
				memberMap[teamID][memberRowID.String] = struct{}{}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]*domain.Team, 0, len(teamOrder))
	for _, id := range teamOrder {
		result = append(result, teamMap[id])
	}

	return result, nil
}
