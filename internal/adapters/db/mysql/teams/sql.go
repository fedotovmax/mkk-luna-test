package teams

import (
	"fmt"
	"strings"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/fields"
)

const create = "insert into teams (id, name, created_by) values (?, ?, ?);"

const createMember = "insert into team_members (id, team_id, user_id, role) values (?, ?, ?, ?);"

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
	}
	if offset > 0 {
		offsetClause = " offset ?"
		args = append(args, offset)
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

const findMember = `
select
    tm.id,
    tm.role,
    tm.joined_at,
    u.id,
    u.username,
    u.email
from team_members tm
join users u on u.id = tm.user_id
where tm.user_id = ? and tm.team_id = ?;`
