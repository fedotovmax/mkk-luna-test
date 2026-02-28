package domain

import "time"

type TeamUser struct {
	ID       string
	Username string
	Email    string
}

type Team struct {
	Members   []Member
	Owner     TeamUser
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
	User     TeamUser
	JoinedAt time.Time
	ID       string
	Role     Role
}
