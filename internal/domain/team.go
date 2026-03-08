package domain

import "time"

type Team struct {
	Members   []Member  `json:"members" validate:"required"`
	Owner     BaseUser  `json:"owner" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
	ID        string    `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleOwner  Role = "owner"
)

type Member struct {
	User     BaseUser  `json:"user" validate:"required"`
	JoinedAt time.Time `json:"joined_at" validate:"required"`
	ID       string    `json:"id" validate:"required"`
	Role     Role      `json:"role" validate:"required"`
}

func (m *Member) CanInvite() bool {
	return m.Role == RoleAdmin || m.Role == RoleOwner
}

func (m *Member) CanDelete() bool {
	return m.Role == RoleOwner
}

func (m *Member) CanCreateTask() bool {
	return m.Role == RoleAdmin || m.Role == RoleOwner
}

func (m *Member) CanUpdateTask(task *Task) bool {

	if m.Role == RoleOwner || m.Role == RoleAdmin {
		return true
	}

	if task.Assignee != nil {
		if task.Assignee.ID == m.User.ID {
			return true
		}
		return false
	}

	return true
}

type TeamStats struct {
	ID                     string `json:"id" validate:"required"`
	Name                   string `json:"name" validate:"required"`
	MembersCount           int    `json:"members_count" validate:"required"`
	DoneTasksLastSevenDays int    `json:"done_tasks_last_seven_days" validate:"required"`
}

type TopUserInTeam struct {
	User         BaseUser `json:"user" validate:"required"`
	TeamID       string   `json:"team_id" validate:"required"`
	TeamName     string   `json:"team_name" validate:"required"`
	CreatedTasks int      `json:"created_tasks" validate:"required"`
}

type FindTeamsResponse struct {
	Teams []*Team `json:"teams" validate:"required"`
	Total int     `json:"total" validate:"required"`
}
