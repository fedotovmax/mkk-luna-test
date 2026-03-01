package teams

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (t *team) Stats(ctx context.Context) ([]domain.TeamStats, error) {

	const op = "adapters.db.mysql.teams.stats"

	tx := t.txExtractor.ExtractTx(ctx)

	rows, err := tx.QueryContext(ctx, stats)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}
	defer rows.Close()

	result := make([]domain.TeamStats, 0)

	for rows.Next() {

		var s domain.TeamStats

		if err := rows.Scan(&s.ID, &s.Name, &s.MembersCount, &s.DoneTasksLastSevenDays); err != nil {
			return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
		}

		result = append(result, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	return result, nil
}

const stats = `
select
    t.id,
    t.name,
    count(distinct tm.user_id) as members_count,
    count(distinct case 
        when ts.status = 'done' and ts.updated_at >= (now() - interval 7 day)
        then ts.id 
    end
    ) as done_tasks_last7d
from teams t
left join team_members tm on tm.team_id = t.id
left join tasks ts on ts.team_id = t.id
group by t.id, t.name
order by t.name;
`
