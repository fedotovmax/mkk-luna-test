package domain

import (
	"time"
)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	switch s {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

type Task struct {
	Assignee    *BaseUser `json:"assignee" validate:"optional"`
	Description *string   `json:"description" validate:"optional"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
	Owner       BaseUser  `json:"owner" validate:"required"`
	ID          string    `json:"id" validate:"required"`
	TeamID      string    `json:"team_id" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Status      Status    `json:"status" validate:"required"`
}

type History struct {
	ChangedBy BaseUser  `json:"changed_by" validate:"required"`
	ChangedAt time.Time `json:"changed_at" validate:"required"`
	Shapshot  Task      `json:"snapshot" validate:"required"`
	ID        string    `json:"id" validate:"required"`
	TaskID    string    `json:"task_id" validate:"required"`
}

type Comment struct {
	User      BaseUser  `json:"user" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	ID        string    `json:"id" validate:"required"`
	TaskID    string    `json:"task_id" validate:"required"`
	Comment   string    `json:"comment" validate:"required"`
}

type FindTasksResponse struct {
	Tasks []*Task `json:"tasks" validate:"required"`
	Total int     `json:"total" validate:"required"`
}
