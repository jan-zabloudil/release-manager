package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type CreateProjectInvitationInput struct {
	Email       string `json:"email"`
	ProjectRole string `json:"project_role"`
}

type ProjectInvitation struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	ProjectRole string    `json:"project_role"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToSvcCreateProjectInvitationInput(i CreateProjectInvitationInput, projectID uuid.UUID) svcmodel.CreateProjectInvitationInput {
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
		CreatedAt:   i.CreatedAt.Local(),
		UpdatedAt:   i.UpdatedAt.Local(),
	}
}

func ToProjectInvitations(invitations []svcmodel.ProjectInvitation) []ProjectInvitation {
	i := make([]ProjectInvitation, 0, len(invitations))
	for _, v := range invitations {
		i = append(i, ToProjectInvitation(v))
	}
	return i
}
