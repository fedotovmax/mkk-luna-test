package inputs

import (
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

type CreateTask struct {
	AssigneeID  *string
	Description *string
	TeamID      string
	Title       string
}

type UpdateTask struct {
	Role        *domain.Status
	Title       *string
	Description *string
}

type FindManyTasks struct {
	Status     domain.Status
	TeamID     string
	AssigneeID string
	Page       int
	PageSize   int
}

type CreateHistory struct {
	OldTask     *domain.Task
	ChangedByID string
}

type CreateComment struct {
	TaskID string
	UserID string
	Text   string
}
