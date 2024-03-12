package model

import (
	svcerr "release-manager/service/errors"
)

const (
	adminProjectRole  = "admin"
	editorProjectRole = "editor"
	viewerProjectRole = "viewer"
)

var rolePriority = map[string]int{
	adminProjectRole:  1,
	editorProjectRole: 2,
	viewerProjectRole: 3,
}

type projectRole struct {
	role string
}

type ProjectRole interface {
	String() string
	IsEqualOrSuperiorTo(ProjectRole) bool
	IsSuperiorTo(ProjectRole) bool
}

func NewProjectRole(role string) (ProjectRole, error) {
	switch role {
	case adminProjectRole, editorProjectRole, viewerProjectRole:
		return &projectRole{role}, nil
	default:
		return nil, svcerr.ErrInvalidProjectRole
	}
}

func ProjectRoleAdmin() ProjectRole {
	return &projectRole{role: adminProjectRole}
}

func ProjectRoleEditor() ProjectRole {
	return &projectRole{role: editorProjectRole}
}

func ProjectRoleViewer() ProjectRole {
	return &projectRole{role: viewerProjectRole}
}

func (r *projectRole) String() string {
	return r.role
}

func (r *projectRole) IsEqualOrSuperiorTo(other ProjectRole) bool {
	return rolePriority[r.role] <= rolePriority[other.String()]
}

func (r *projectRole) IsSuperiorTo(other ProjectRole) bool {
	return rolePriority[r.role] < rolePriority[other.String()]
}
