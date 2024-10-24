package model

import (
	"errors"
	"time"

	"release-manager/pkg/id"
)

const (
	ProjectRoleOwner  ProjectRole = "owner"
	ProjectRoleEditor ProjectRole = "editor"
	ProjectRoleViewer ProjectRole = "viewer"

	projectRoleOwnerPriority  int = 1
	projectRoleEditorPriority int = 2
	projectRoleViewerPriority int = 3
)

var (
	validProjectRoles = map[ProjectRole]bool{
		ProjectRoleOwner:  true,
		ProjectRoleEditor: true,
		ProjectRoleViewer: true,
	}

	projectRolePriority = map[ProjectRole]int{
		ProjectRoleOwner:  projectRoleOwnerPriority,
		ProjectRoleEditor: projectRoleEditorPriority,
		ProjectRoleViewer: projectRoleViewerPriority,
	}

	errProjectRoleInvalid            = errors.New("invalid project role")
	errProjectMemberCannotGrantOwner = errors.New("cannot grant owner role")
)

type ProjectRole string

func (r ProjectRole) Validate() error {
	if _, exists := validProjectRoles[r]; exists {
		return nil
	}

	return errProjectRoleInvalid
}

func (r ProjectRole) IsRoleAtLeast(role ProjectRole) bool {
	return projectRolePriority[r] <= projectRolePriority[role]
}

type ProjectMember struct {
	User        User
	ProjectID   id.Project
	ProjectRole ProjectRole
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProjectOwner(u User, projectID id.Project) (ProjectMember, error) {
	return NewProjectMember(u, projectID, ProjectRoleOwner)
}

func NewProjectEditor(u User, projectID id.Project) (ProjectMember, error) {
	return NewProjectMember(u, projectID, ProjectRoleEditor)
}

func NewProjectViewer(u User, projectID id.Project) (ProjectMember, error) {
	return NewProjectMember(u, projectID, ProjectRoleViewer)
}

func NewProjectMember(u User, projectID id.Project, role ProjectRole) (ProjectMember, error) {
	now := time.Now()

	m := ProjectMember{
		User:        u,
		ProjectID:   projectID,
		ProjectRole: role,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := m.Validate(); err != nil {
		return ProjectMember{}, err
	}

	return m, nil
}

func (m *ProjectMember) Validate() error {
	if err := m.ProjectRole.Validate(); err != nil {
		return err
	}

	return nil
}

func (m *ProjectMember) UpdateProjectRole(role ProjectRole) error {
	// Cannot grant owner role, it can only be granted when creating a new project
	if role == ProjectRoleOwner {
		return errProjectMemberCannotGrantOwner
	}

	m.ProjectRole = role
	m.UpdatedAt = time.Now()

	return m.Validate()
}

// SatisfiesRequiredRole checks if the project member has a role that meets or exceeds the required role
func (m *ProjectMember) SatisfiesRequiredRole(role ProjectRole) bool {
	// Admin user role exceeds all project roles
	if m.User.IsAdmin() {
		return true
	}

	return m.ProjectRole.IsRoleAtLeast(role)
}
