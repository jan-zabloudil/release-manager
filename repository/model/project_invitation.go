package model

import (
	"time"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"
)

type ProjectInvitation struct {
	ID            id.ProjectInvitation `db:"id"`
	ProjectID     id.Project           `db:"project_id"`
	Email         string               `db:"email"`
	ProjectRole   string               `db:"project_role"`
	Status        string               `db:"status"`
	TokenHash     []byte               `db:"token_hash"`
	InviterUserID id.AuthUser          `db:"invited_by"`
	CreatedAt     time.Time            `db:"created_at"`
	UpdatedAt     time.Time            `db:"updated_at"`
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
