package model

import (
	"errors"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

const (
	userRoleAdminPriority int = 1
	userRoleUserPriority  int = 2
)

var (
	errUserInvalidRole = errors.New("invalid user role")
)

var userRolePriority = map[UserRole]int{
	UserRoleAdmin: userRoleAdminPriority,
	UserRoleUser:  userRoleUserPriority,
}

func (r UserRole) IsRoleAtLeast(role UserRole) bool {
	return userRolePriority[r] <= userRolePriority[role]
}

func NewUserRole(role string) (UserRole, error) {
	switch role {
	case string(UserRoleAdmin):
		return UserRoleAdmin, nil
	case string(UserRoleUser):
		return UserRoleUser, nil
	default:
		return "", errUserInvalidRole
	}
}
