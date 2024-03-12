package model

import (
	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/relvacode/iso8601"
)

type ProjectInvitationInput struct {
	ProjectID       uuid.UUID `json:"project_id"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	InvitedByUserID uuid.UUID `json:"invited_by_user_id"`
}

type ProjectInvitationResponse struct {
	ID              uuid.UUID    `json:"id"`
	ProjectID       uuid.UUID    `json:"project_id"`
	Email           string       `json:"email"`
	Role            string       `json:"role"`
	InvitedByUserID uuid.UUID    `json:"invited_by_user_id"`
	CreatedAt       iso8601.Time `json:"created_at"`
}

type ProjectInvitationPatch struct {
	Role string `json:"role"`
}

func ToProjectInvitationInput(projectID uuid.UUID, email string, role svcmodel.ProjectRole, invitedByUserID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"p_project_id":         projectID,
		"p_email":              email,
		"p_role":               role.Role(),
		"p_invited_by_user_id": invitedByUserID,
	}
}

func ToSvcProjectInvitation(r ProjectInvitationResponse) (svcmodel.ProjectInvitation, error) {
	role, err := svcmodel.NewProjectRole(r.Role)
	if err != nil {
		return svcmodel.ProjectInvitation{}, err
	}

	return svcmodel.ProjectInvitation{
		ID:              r.ID,
		ProjectID:       r.ProjectID,
		Email:           r.Email,
		Role:            role,
		InvitedByUserID: r.InvitedByUserID,
		CreatedAt:       r.CreatedAt.Time,
	}, nil
}

func ToSvcProjectInvitations(invitations []ProjectInvitationResponse) ([]svcmodel.ProjectInvitation, error) {
	i := make([]svcmodel.ProjectInvitation, 0, len(invitations))
	for _, invitation := range invitations {
		svcInvitation, err := ToSvcProjectInvitation(invitation)
		if err != nil {
			return nil, err
		}

		i = append(i, svcInvitation)
	}

	return i, nil
}
