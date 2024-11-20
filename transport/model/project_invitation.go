package model

import (
	"time"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"
)

type CreateProjectInvitationInput struct {
	Email       string `json:"email" validate:"required,email"`
	ProjectRole string `json:"project_role" validate:"required"`
}

type ProjectInvitation struct {
	ID          id.ProjectInvitation `json:"id"`
	Email       string               `json:"email"`
	ProjectRole string               `json:"project_role"`
	Status      string               `json:"status"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type ProjectInvitationURLParams struct {
	ProjectID    id.Project           `param:"path=project_id"`
	InvitationID id.ProjectInvitation `param:"path=invitation_id"`
}

func ToSvcCreateProjectInvitationInput(i CreateProjectInvitationInput, projectID id.Project) svcmodel.CreateProjectInvitationInput {
	return svcmodel.CreateProjectInvitationInput{
		ProjectID:   projectID,
		Email:       i.Email,
		ProjectRole: i.ProjectRole,
	}
}

func ToProjectInvitation(i svcmodel.ProjectInvitation) ProjectInvitation {
	return ProjectInvitation{
		ID:          i.ID,
		Email:       i.Email,
		ProjectRole: string(i.ProjectRole),
		Status:      string(i.Status),
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func ToProjectInvitations(invitations []svcmodel.ProjectInvitation) []ProjectInvitation {
	i := make([]ProjectInvitation, 0, len(invitations))
	for _, v := range invitations {
		i = append(i, ToProjectInvitation(v))
	}
	return i
}
