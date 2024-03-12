package model

import (
	"fmt"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
	"github.com/relvacode/iso8601"
)

type ProjectMemberInput struct {
	ProjectID       uuid.UUID `json:"project_id"`
	UserID          uuid.UUID `json:"user_id"`
	Role            string    `json:"role"`
	InvitedByUserID uuid.UUID `json:"invited_by_user_id"`
}

type ProjectMember struct {
	User            supabase.User `json:"user_data"`
	ProjectID       uuid.UUID     `json:"project_id"`
	Role            string        `json:"role"`
	InvitedByUserId uuid.UUID     `json:"invited_by_user_id"`
	CreatedAt       iso8601.Time  `json:"created_at"`
	UpdatedAt       iso8601.Time  `json:"updated_at"`
}

func ToProjectMemberInput(projectID, userID uuid.UUID, role svcmodel.ProjectRole, invitedByUserID uuid.UUID) ProjectMemberInput {
	return ProjectMemberInput{
		ProjectID:       projectID,
		UserID:          userID,
		Role:            role.String(),
		InvitedByUserID: invitedByUserID,
	}
}

func ToSvcProjectMember(m ProjectMember) (svcmodel.ProjectMember, error) {
	u, err := ToSvcUser(
		m.User.ID,
		m.User.Email,
		m.User.AppMetadata["is_admin"],
		m.User.UserMetadata["name"],
		m.User.UserMetadata["picture"],
		m.User.CreatedAt,
		m.User.UpdatedAt,
	)
	if err != nil {
		fmt.Println("to model part 1", err)
		return svcmodel.ProjectMember{}, err
	}

	role, err := svcmodel.NewProjectRole(m.Role)
	if err != nil {
		fmt.Println("to model part 2", err)
		return svcmodel.ProjectMember{}, err
	}

	return svcmodel.ProjectMember{
		User:            u,
		ProjectID:       m.ProjectID,
		Role:            role,
		InvitedByUserId: m.InvitedByUserId,
		CreatedAt:       m.CreatedAt.Time,
		UpdatedAt:       m.UpdatedAt.Time,
	}, nil
}

func ToSvcProjectMembers(members []ProjectMember) ([]svcmodel.ProjectMember, error) {
	m := make([]svcmodel.ProjectMember, 0, len(members))
	for _, member := range members {
		svcMember, err := ToSvcProjectMember(member)
		if err != nil {
			return nil, err
		}

		m = append(m, svcMember)
	}

	return m, nil
}
