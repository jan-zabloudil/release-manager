package model

import (
	"errors"
	"time"

	cryptox "release-manager/pkg/crypto"
	validator "release-manager/pkg/validator"

	"github.com/google/uuid"
)

const (
	ProjectRoleOwner  ProjectRole = "owner"
	ProjectRoleEditor ProjectRole = "editor"
	ProjectRoleViewer ProjectRole = "viewer"

	projectRoleOwnerPriority  int = 1
	projectRoleEditorPriority int = 2
	projectRoleViewerPriority int = 3

	InvitationStatusPending  ProjectInvitationStatus = "pending"
	InvitationStatusAccepted ProjectInvitationStatus = "accepted_awaiting_registration"
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

	validInvitationStatuses = map[ProjectInvitationStatus]bool{
		InvitationStatusPending:  true,
		InvitationStatusAccepted: true,
	}

	errProjectRoleInvalid                    = errors.New("invalid project role")
	errProjectInvitationStatusInvalid        = errors.New("invalid invitation status")
	errProjectInvitationEmailRequired        = errors.New("email is required")
	errProjectInvitationInvalidEmail         = errors.New("invalid email")
	errProjectInvitationCannotGrantOwnerRole = errors.New("cannot grant owner role to project member")
)

type ProjectInvitation struct {
	ID            uuid.UUID
	ProjectID     uuid.UUID
	Email         string
	ProjectRole   ProjectRole
	Status        ProjectInvitationStatus
	TokenHash     cryptox.Hash
	InviterUserID uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateProjectInvitationInput struct {
	ProjectID   uuid.UUID
	Email       string
	ProjectRole string
}

type ProjectRole string
type ProjectInvitationStatus string

func NewProjectInvitation(c CreateProjectInvitationInput, tkn cryptox.Token, inviterUserID uuid.UUID) (ProjectInvitation, error) {
	now := time.Now()

	i := ProjectInvitation{
		ID:            uuid.New(),
		ProjectID:     c.ProjectID,
		Email:         c.Email,
		ProjectRole:   ProjectRole(c.ProjectRole),
		Status:        InvitationStatusPending,
		TokenHash:     tkn.ToHash(),
		InviterUserID: inviterUserID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := i.Validate(); err != nil {
		return ProjectInvitation{}, err
	}

	return i, nil
}

func (i *ProjectInvitation) Accept() {
	i.Status = InvitationStatusAccepted
	i.UpdatedAt = time.Now()
}

func (i *ProjectInvitation) Validate() error {
	if i.Email == "" {
		return errProjectInvitationEmailRequired
	}
	if !validator.IsValidEmail(i.Email) {
		return errProjectInvitationInvalidEmail
	}
	if err := i.ProjectRole.Validate(); err != nil {
		return err
	}
	if i.ProjectRole == ProjectRoleOwner {
		return errProjectInvitationCannotGrantOwnerRole
	}
	if err := i.Status.Validate(); err != nil {
		return err
	}

	return nil
}

func (r ProjectRole) Validate() error {
	if _, exists := validProjectRoles[r]; exists {
		return nil
	}

	return errProjectRoleInvalid
}

func (i ProjectInvitationStatus) Validate() error {
	if _, exists := validInvitationStatuses[i]; exists {
		return nil
	}

	return errProjectInvitationStatusInvalid
}

func (r ProjectRole) IsRoleAtLeast(role ProjectRole) bool {
	return projectRolePriority[r] <= projectRolePriority[role]
}

type ProjectMember struct {
	User        User
	ProjectID   uuid.UUID
	ProjectRole ProjectRole
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProjectOwner(u User, projectID uuid.UUID) (ProjectMember, error) {
	return NewProjectMember(u, projectID, ProjectRoleOwner)
}

func NewProjectEditor(u User, projectID uuid.UUID) (ProjectMember, error) {
	return NewProjectMember(u, projectID, ProjectRoleEditor)
}

func NewProjectViewer(u User, projectID uuid.UUID) (ProjectMember, error) {
	return NewProjectMember(u, projectID, ProjectRoleViewer)
}

func NewProjectMember(u User, projectID uuid.UUID, role ProjectRole) (ProjectMember, error) {
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
	m.ProjectRole = role
	m.UpdatedAt = time.Now()

	return m.Validate()
}

func (m *ProjectMember) HasAtLeastProjectRole(role ProjectRole) bool {
	return m.ProjectRole.IsRoleAtLeast(role)
}
