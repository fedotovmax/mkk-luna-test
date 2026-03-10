package cache

import (
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
)

func TasksListKey(limit, offset int, in *inputs.FindManyTasks) string {

	return fmt.Sprintf("TASKS_LIST_KEY:limit=%d,offset=%d,status=%s,team_id=%s,assignee_id=%s", limit, offset, in.Status, in.TeamID, in.AssigneeID)

}
