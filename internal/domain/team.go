package domain

import "time"

type Team struct {
	Members   []Member
	Owner     BaseUser
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        string
	Name      string
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleOwner  Role = "owner"
)

type Member struct {
	User     BaseUser
	JoinedAt time.Time
	ID       string
	Role     Role
}

func (m *Member) CanInvite() bool {
	return m.Role == RoleAdmin || m.Role == RoleOwner
}

func (m *Member) CanDelete() bool {
	return m.Role == RoleOwner
}

type TeamStats struct {
	ID                     string
	Name                   string
	MembersCount           int
	DoneTasksLastSevenDays int
}

type TopUserInTeam struct {
	User         BaseUser
	TeamID       string
	TeamName     string
	CreatedTasks int
}
