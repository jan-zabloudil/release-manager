package model

import (
	"context"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectInvitationService interface {
	ListAll(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectInvitation, error)
	Get(ctx context.Context, invitationID uuid.UUID) (svcmodel.ProjectInvitation, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectInvitation struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	InvitedByUserID uuid.UUID `json:"invited_by_user_id"`
	CreatedAt       time.Time `json:"created_at"`
}

func ToNetProjectInvitation(i svcmodel.ProjectInvitation) ProjectInvitation {
	return ProjectInvitation{
		ID:              i.ID,
		Email:           i.Email,
		Role:            i.Role.Role(),
		InvitedByUserID: i.InvitedByUserID,
		CreatedAt:       i.CreatedAt,
	}
}

func ToNetProjectInvitations(invitations []svcmodel.ProjectInvitation) []ProjectInvitation {
	i := make([]ProjectInvitation, 0, len(invitations))
	for _, invitation := range invitations {
		i = append(i, ToNetProjectInvitation(invitation))
	}

	return i
}

/*
type ProjectInvitation struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email" validate:"required"`
	Role      string    `json:"role" validate:"required,eq=admin|eq=editor|eq=viewer"`
	InvitedByUserID uuid.UUID `json:"invited_by"`
	CreatedAt time.Time `json:"created_at"`
}

func ToSvcProjectInvitation(i ProjectInvitation, projectID uuid.UUID, userID uuid.UUID) (svcmodel.ProjectInvitation, error) {
	role, err := svcmodel.NewProjectRole(i.Role)
	if err != nil {
		return svcmodel.ProjectInvitation{}, err
	}

	return svcmodel.ProjectInvitation{
		ID:        i.ID,
		Email:     i.Email,
		Role:      role,
		ProjectID: projectID,
		InvitedByUserID: userID,
	}, nil
}

func ToNetProjectInvitation(i svcmodel.ProjectInvitation) ProjectInvitation {
	return ProjectInvitation{
		ID:        i.ID,
		Email:     i.Email,
		Role:      i.Role.Role(),
		InvitedByUserID: i.InvitedByUserID,
		CreatedAt: i.CreatedAt,
	}
}

func ToNetProjectInvitations(invitations []svcmodel.ProjectInvitation) []ProjectInvitation {
	i := make([]ProjectInvitation, 0, len(invitations))
	for _, invitation := range invitations {
		i = append(i, ToNetProjectInvitation(invitation))
	}

	return i
}

*/
