package model

import (
	"time"

	"github.com/google/uuid"
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
	ID        uuid.UUID
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

func (r UserRole) IsRoleAtLeast(role UserRole) bool {
	return userRolePriority[r] <= userRolePriority[role]
}
