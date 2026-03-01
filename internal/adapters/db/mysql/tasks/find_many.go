package tasks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
)

func (t *task) FindMany(
	ctx context.Context,
	in *inputs.FindManyTasks,
) (*domain.FindTasksResponse, error) {

	const op = "adapters.db.mysql.tasks.find_many"

	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 20
	}

	tx := t.txExtractor.ExtractTx(ctx)

	offset := (in.Page - 1) * in.PageSize

	countQuery, countArgs := buildCountTasksQuery(in.TeamID, in.Status, in.AssigneeID)
	var total int

	err := tx.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	query, args := buildFindTasksQuery(in.TeamID, in.Status, in.AssigneeID, in.PageSize, offset)
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}
	defer rows.Close()

	tasks := make([]*domain.Task, 0, in.PageSize)

	for rows.Next() {

		task := &domain.Task{}
		var assigneeID, assigneeUsername, assigneeEmail sql.NullString

		err := rows.Scan(
			&task.ID,
			&task.TeamID,
			&task.Title,
			&task.Description,
			&task.Status,

			&task.Owner.ID,
			&task.Owner.Username,
			&task.Owner.Email,

			&assigneeID,
			&assigneeUsername,
			&assigneeEmail,

			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
		}

		if assigneeID.Valid {
			task.Assignee = &domain.BaseUser{
				ID:       assigneeID.String,
				Username: assigneeUsername.String,
				Email:    assigneeEmail.String,
			}
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, adapters.ErrInternal)
	}

	return &domain.FindTasksResponse{
		Total: total,
		Tasks: tasks,
	}, nil
}

func buildFindTasksQuery(teamID string, status domain.Status, assigneeID string, limit, offset int) (string, []any) {
	whereParts := make([]string, 0)
	args := make([]any, 0)

	if teamID != "" {
		whereParts = append(whereParts, "t.team_id = ?")
		args = append(args, teamID)
	}

	if status != "" {
		whereParts = append(whereParts, "t.status = ?")
		args = append(args, status)
	}

	if assigneeID != "" {
		whereParts = append(whereParts, "t.assignee_id = ?")
		args = append(args, assigneeID)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "where " + strings.Join(whereParts, " and ")
	}

	query := fmt.Sprintf(`
select
    t.id,
    t.team_id,
    t.title,
    t.description,
    t.status,

    owner.id,
    owner.username,
    owner.email,

    assignee.id,
    assignee.username,
    assignee.email,

    t.created_at,
    t.updated_at

from tasks t
join users owner on owner.id = t.created_by
left join users assignee on assignee.id = t.assignee_id

%s
order by t.created_at desc
limit ? offset ?;
`, whereClause)

	args = append(args, limit, offset)
	return query, args
}

func buildCountTasksQuery(teamID string, status domain.Status, assigneeID string) (string, []any) {
	whereParts := make([]string, 0)
	args := make([]any, 0)

	if teamID != "" {
		whereParts = append(whereParts, "team_id = ?")
		args = append(args, teamID)
	}

	if status != "" {
		whereParts = append(whereParts, "status = ?")
		args = append(args, status)
	}

	if assigneeID != "" {
		whereParts = append(whereParts, "assignee_id = ?")
		args = append(args, assigneeID)
	}

	whereClause := ""
	if len(whereParts) > 0 {
		whereClause = "where " + strings.Join(whereParts, " and ")
	}

	query := fmt.Sprintf(`select count(*) from tasks %s;`, whereClause)
	return query, args
}
