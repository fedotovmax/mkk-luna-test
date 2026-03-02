package inputs

import (
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

type CreateTask struct {
	AssigneeID  *string `json:"assignee_id" validate:"optional"`
	Description *string `json:"description" validate:"optional"`
	TeamID      string  `json:"team_id" validate:"required"`
	Title       string  `json:"title" validate:"required"`
}

type UpdateTask struct {
	Role        *domain.Status `json:"role" validate:"optional"`
	Title       *string        `json:"title" validate:"optional"`
	Description *string        `json:"description" validate:"optional"`
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
