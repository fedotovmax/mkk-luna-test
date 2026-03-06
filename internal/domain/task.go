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
	Assignee    *BaseUser
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Owner       BaseUser
	ID          string
	TeamID      string
	Title       string
	Status      Status
}

type History struct {
	ChangedBy BaseUser
	ChangedAt time.Time
	Shapshot  Task
	ID        string
	TaskID    string
}

type Comment struct {
	User      BaseUser
	CreatedAt time.Time
	ID        string
	TaskID    string
	Comment   string
}

type FindTasksResponse struct {
	Tasks []*Task
	Total int
}
