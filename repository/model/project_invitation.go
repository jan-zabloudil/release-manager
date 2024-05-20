package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

// TODO remove json tags once all functions use db pool
type ProjectInvitation struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ProjectID     uuid.UUID `json:"project_id" db:"project_id"`
	Email         string    `json:"email" db:"email"`
	ProjectRole   string    `json:"project_role" db:"project_role"`
	Status        string    `json:"status" db:"status"`
	TokenHash     []byte    `json:"token_hash" db:"token_hash"`
	InviterUserID uuid.UUID `json:"invited_by" db:"invited_by"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
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
