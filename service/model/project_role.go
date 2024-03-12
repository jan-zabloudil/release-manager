package model

import (
	svcerr "release-manager/service/errors"
)

const (
	AdminProjectRole  = "admin"
	EditorProjectRole = "editor"
	ViewerProjectRole = "viewer"
)

type projectRole struct {
	role string
}

type ProjectRole interface {
	Role() string
}

func NewProjectRole(role string) (ProjectRole, error) {
	switch role {
	case AdminProjectRole, EditorProjectRole, ViewerProjectRole:
		return &projectRole{role}, nil
	default:
		return nil, svcerr.ErrInvalidProjectRole
	}
}

func (r *projectRole) Role() string {
	return r.role
}
