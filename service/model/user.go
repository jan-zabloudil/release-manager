package model

import (
	"time"

	"release-manager/pkg/id"
)

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"

	userRoleAdminPriority int = 1
	userRoleUserPriority  int = 2
)

var userRolePriority = map[UserRole]int{
	UserRoleAdmin: userRoleAdminPriority,
	UserRoleUser:  userRoleUserPriority,
}

type UserRole string

type User struct {
	ID        id.User
	Email     string
	Name      string
	AvatarURL string
	Role      UserRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s User) HasAtLeastRole(role UserRole) bool {
	return s.Role.IsRoleAtLeast(role)
}

func (s User) IsAdmin() bool {
	return s.Role == UserRoleAdmin
}

func (r UserRole) IsRoleAtLeast(role UserRole) bool {
	return userRolePriority[r] <= userRolePriority[role]
}
