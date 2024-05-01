package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectInvitation struct {
	ID            uuid.UUID `json:"id"`
	ProjectID     uuid.UUID `json:"project_id"`
	Email         string    `json:"email"`
	ProjectRole   string    `json:"project_role"`
	Status        string    `json:"status"`
	TokenHash     []byte    `json:"token_hash"`
	InviterUserID uuid.UUID `json:"invited_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UpdateProjectInvitationInput struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToProjectInvitation(i svcmodel.ProjectInvitation) ProjectInvitation {
	return ProjectInvitation{
		ID:            i.ID,
		ProjectID:     i.ProjectID,
		Email:         i.Email,
		ProjectRole:   string(i.ProjectRole),
		Status:        string(i.Status),
		TokenHash:     i.TokenHash,
		InviterUserID: i.InviterUserID,
		CreatedAt:     i.CreatedAt,
		UpdatedAt:     i.UpdatedAt,
	}
}

func ToUpdateProjectInvitationInput(u svcmodel.ProjectInvitation) UpdateProjectInvitationInput {
	return UpdateProjectInvitationInput{
		Status:    string(u.Status),
		UpdatedAt: u.UpdatedAt,
	}
}

func ToSvcProjectInvitation(i ProjectInvitation) svcmodel.ProjectInvitation {
	return svcmodel.ProjectInvitation{
		ID:            i.ID,
		ProjectID:     i.ProjectID,
		Email:         i.Email,
		ProjectRole:   svcmodel.ProjectRole(i.ProjectRole),
		Status:        svcmodel.ProjectInvitationStatus(i.Status),
		TokenHash:     i.TokenHash,
		InviterUserID: i.InviterUserID,
		CreatedAt:     i.CreatedAt,
		UpdatedAt:     i.UpdatedAt,
	}
}

func ToSvcProjectInvitations(invitations []ProjectInvitation) []svcmodel.ProjectInvitation {
	i := make([]svcmodel.ProjectInvitation, 0, len(invitations))
	for _, invitation := range invitations {
		i = append(i, ToSvcProjectInvitation(invitation))
	}
	return i
}
