package model

import svcerr "release-manager/service/errors"

const (
	basicUserRole = "user"
	adminUserRole = "admin"
)

type userRole struct {
	role string
}

type UserRole interface {
	Role() string
}

func (r *userRole) Role() string {
	return r.role
}

func NewUserRole(role string) (UserRole, error) {
	switch role {
	case basicUserRole, adminUserRole:
		return &userRole{role}, nil
	default:
		return nil, svcerr.ErrInvalidUserRole
	}
}

func NewBasicUserRole() UserRole {
	return &userRole{role: basicUserRole}
}
