package teams

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (t *team) FindMany(
	ctx context.Context,
	limit, offset int,
	userID string,
) (*domain.FindTeamsResponse, error) {

	const op = "adapters.db.mysql.teams.find_many"

	tx := t.txExtractor.ExtractTx(ctx)

	countQuery, countArgs := buildCount(userID)

	var total int
	err := tx.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	query, args := buildFindMany(limit, offset, userID)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if _, exists := teamMap[teamID]; !exists {
			teamMap[teamID] = &domain.Team{
				ID:        teamID,
				Name:      teamName,
				CreatedAt: created,
				UpdatedAt: updated,
				Owner: domain.BaseUser{
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
					User: domain.BaseUser{
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result := make([]*domain.Team, 0, len(teamOrder))
	for _, id := range teamOrder {
		result = append(result, teamMap[id])
	}

	return &domain.FindTeamsResponse{
		Total: total,
		Teams: result,
	}, nil
}

func buildFindMany(limit, offset int, userID string) (string, []any) {
	whereParts := []string{}
	args := []any{}

	if userID != "" {
		whereParts = append(whereParts, `
			exists (
				select 1
				from team_members tm_filter
				where tm_filter.team_id = t.id
				and tm_filter.user_id = ?
			)
		`)
		args = append(args, userID)
	}

	//Other filters (mb later)

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "where " + strings.Join(whereParts, " and ")
	}

	limitClause := ""
	offsetClause := ""
	if limit > 0 {
		limitClause = " limit ?"
		args = append(args, limit)
		if offset > 0 {
			offsetClause = " offset ?"
			args = append(args, offset)
		}
	}

	query := fmt.Sprintf(`
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

		from (
			select t.id
			from teams t
			%s
			order by t.created_at desc
			%s
			%s
		) paged
		join teams t on t.id = paged.id
		join users owner on owner.id = t.created_by
		left join team_members tm on tm.team_id = t.id
		left join users member on member.id = tm.user_id
		order by t.created_at desc;
	`, whereClause, limitClause, offsetClause)

	return query, args
}

func buildCount(userID string) (string, []any) {

	whereParts := []string{}
	args := []any{}

	if userID != "" {
		whereParts = append(whereParts, `
			exists (
				select 1
				from team_members tm_filter
				where tm_filter.team_id = t.id
				and tm_filter.user_id = ?
			)
		`)
		args = append(args, userID)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "where " + strings.Join(whereParts, " and ")
	}

	query := fmt.Sprintf(`
		select count(*)
		from teams t
		%s;
	`, whereClause)

	return query, args
}
