package model

import (
	"errors"
	"time"

	cryptox "release-manager/pkg/crypto"
	"release-manager/pkg/id"
	validator "release-manager/pkg/validator"

	"github.com/google/uuid"
)

const (
	InvitationStatusPending  ProjectInvitationStatus = "pending"
	InvitationStatusAccepted ProjectInvitationStatus = "accepted_awaiting_registration"
)

var (
	validInvitationStatuses = map[ProjectInvitationStatus]bool{
		InvitationStatusPending:  true,
		InvitationStatusAccepted: true,
	}

	errProjectInvitationStatusInvalid        = errors.New("invalid invitation status")
	errProjectInvitationEmailRequired        = errors.New("email is required")
	errProjectInvitationInvalidEmail         = errors.New("invalid email")
	errProjectInvitationCannotGrantOwnerRole = errors.New("cannot grant owner role to project member")
	ErrProjectInvitationAlreadyAccepted      = errors.New("invitation already accepted")
)

type ProjectInvitation struct {
	ID            id.ProjectInvitation
	ProjectID     uuid.UUID
	Email         string
	ProjectRole   ProjectRole
	Status        ProjectInvitationStatus
	TokenHash     ProjectInvitationTokenHash
	InviterUserID id.AuthUser
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateProjectInvitationInput struct {
	ProjectID   uuid.UUID
	Email       string
	ProjectRole string
}

type ProjectInvitationStatus string
type ProjectInvitationToken cryptox.Token
type ProjectInvitationTokenHash cryptox.Hash

func NewProjectInvitation(c CreateProjectInvitationInput, tkn ProjectInvitationToken, inviterUserID id.AuthUser) (ProjectInvitation, error) {
	now := time.Now()

	i := ProjectInvitation{
		ID:            id.NewProjectInvitation(),
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

func (i *ProjectInvitation) Accept() error {
	if i.Status == InvitationStatusAccepted {
		return ErrProjectInvitationAlreadyAccepted
	}

	i.Status = InvitationStatusAccepted
	i.UpdatedAt = time.Now()

	return nil
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

func (i ProjectInvitationStatus) Validate() error {
	if _, exists := validInvitationStatuses[i]; exists {
		return nil
	}

	return errProjectInvitationStatusInvalid
}

func NewProjectInvitationToken() (ProjectInvitationToken, error) {
	tkn, err := cryptox.NewToken()
	if err != nil {
		return "", err
	}

	return ProjectInvitationToken(tkn), nil
}

func (i ProjectInvitationToken) ToHash() ProjectInvitationTokenHash {
	return ProjectInvitationTokenHash(cryptox.Token(i).ToHash())
}

func (i ProjectInvitationTokenHash) ToBase64() string {
	return cryptox.Hash(i).ToBase64()
}

type ProjectInvitationEmailData struct {
	ProjectName string
	Token       ProjectInvitationToken
}

func NewProjectInvitationEmailData(projectName string, tkn ProjectInvitationToken) ProjectInvitationEmailData {
	return ProjectInvitationEmailData{
		ProjectName: projectName,
		Token:       tkn,
	}
}
