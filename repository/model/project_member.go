package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

// CreateProjectMemberInput is the input for creating a project member and deleting an invitation
type CreateProjectMemberInput struct {
	UserID      uuid.UUID `json:"p_user_id"`
	ProjectID   uuid.UUID `json:"p_project_id"`
	Email       string    `json:"p_email"`
	ProjectRole string    `json:"p_project_role"`
	CreatedAt   time.Time `json:"p_created_at"`
	UpdatedAt   time.Time `json:"p_updated_at"`
}

type UpdateProjectMemberInput struct {
	ProjectRole string    `json:"project_role"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectMember struct {
	User        User      `json:"users"` // Supabase returns joined table data in json array named after joined table, "users" in this case
	ProjectID   uuid.UUID `json:"project_id"`
	ProjectRole string    `json:"project_role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToCreateProjectMemberInput(m svcmodel.ProjectMember) CreateProjectMemberInput {
	return CreateProjectMemberInput{
		UserID:      m.User.ID,
		ProjectID:   m.ProjectID,
		Email:       m.User.Email,
		ProjectRole: string(m.ProjectRole),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func ToSvcProjectMember(p ProjectMember) svcmodel.ProjectMember {
	return svcmodel.ProjectMember{
		User:        ToSvcUser(p.User),
		ProjectID:   p.ProjectID,
		ProjectRole: svcmodel.ProjectRole(p.ProjectRole),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToSvcProjectMembers(members []ProjectMember) []svcmodel.ProjectMember {
	m := make([]svcmodel.ProjectMember, 0, len(members))
	for _, member := range members {
		m = append(m, ToSvcProjectMember(member))
	}
	return m
}

func ToUpdateProjectMemberInput(p svcmodel.ProjectMember) UpdateProjectMemberInput {
	return UpdateProjectMemberInput{
		ProjectRole: string(p.ProjectRole),
		UpdatedAt:   p.UpdatedAt,
	}
}
