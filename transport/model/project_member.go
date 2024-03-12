package model

import (
	"context"
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectMember struct {
	User            User      `json:"user"`
	Role            string    `json:"role"`
	InvitedByUserID uuid.UUID `json:"invited_by_user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PatchProjectMember struct {
	Role *string `json:"role" validate:"required"`
}

type ProjectMemberService interface {
	ListAll(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectMember, error)
	Get(ctx context.Context, projectID, userID uuid.UUID) (svcmodel.ProjectMember, error)
	Delete(ctx context.Context, projectID, userID uuid.UUID) error
	Update(ctx context.Context, member svcmodel.ProjectMember) (svcmodel.ProjectMember, error)
}

func ToNetProjectMember(m svcmodel.ProjectMember) ProjectMember {
	return ProjectMember{
		User:            ToNetUser(m.User.ID, m.User.IsAdmin, m.User.Email, m.User.Name, m.User.AvatarUrl, m.User.CreatedAt, m.User.UpdatedAt),
		Role:            m.Role.Role(),
		InvitedByUserID: m.InvitedByUserId,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func ToNetProjectMembers(members []svcmodel.ProjectMember) []ProjectMember {
	m := make([]ProjectMember, 0, len(members))
	for _, member := range members {
		m = append(m, ToNetProjectMember(member))
	}

	return m
}

func PatchToSvcProjectMember(patch PatchProjectMember, m svcmodel.ProjectMember) (svcmodel.ProjectMember, error) {
	if patch.Role != nil {
		newRole, err := svcmodel.NewProjectRole(*patch.Role)
		if err != nil {
			return svcmodel.ProjectMember{}, err
		}

		m.Role = newRole
	}

	return m, nil
}
