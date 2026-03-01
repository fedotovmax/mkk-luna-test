package teams

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (t *team) TopUsers(ctx context.Context) ([]domain.TopUserInTeam, error) {
	const op = "adapters.db.mysql.teams.top_users"

	tx := t.txExtractor.ExtractTx(ctx)

	rows, err := tx.QueryContext(ctx, topUsers)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}
	defer rows.Close()

	result := make([]domain.TopUserInTeam, 0)

	for rows.Next() {
		var t domain.TopUserInTeam
		if err := rows.Scan(&t.TeamID, &t.TeamName, &t.User.ID, &t.User.Username, &t.User.Email, &t.CreatedTasks); err != nil {
			return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
		}
		result = append(result, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	return result, nil
}

const topUsers = `
with user_task_counts as (
    select
        t.team_id,
        t.created_by as user_id,
        count(*) as tasks_created
    from tasks t
    where t.created_at >= now() - interval 1 month
    group by t.team_id, t.created_by
),
ranked as (
    select
        utc.team_id,
        utc.user_id,
        utc.tasks_created,
        row_number() over (
            partition by utc.team_id
            order by utc.tasks_created desc
        ) as rn
    from user_task_counts utc
)
select
    r.team_id,
    tm.name as team_name,
    u.id,
    u.username,
    u.email,
    r.tasks_created
from ranked r
join teams tm on tm.id = r.team_id
join users u on u.id = r.user_id
where r.rn <= 3
order by r.team_id, r.tasks_created desc;
`
