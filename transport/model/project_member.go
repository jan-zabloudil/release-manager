package model

import (
	"time"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"
)

type ProjectMember struct {
	UserID      id.User   `json:"user_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	ProjectRole string    `json:"project_role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateProjectMemberRoleInput struct {
	ProjectRole string `json:"project_role"`
}

type ProjectMemberURLParams struct {
	ProjectID id.Project `param:"path=project_id"`
	UserID    id.User    `param:"path=user_id"`
}

func ToProjectMember(p svcmodel.ProjectMember) ProjectMember {
	return ProjectMember{
		UserID:      p.User.ID,
		Name:        p.User.Name,
		Email:       p.User.Email,
		ProjectRole: string(p.ProjectRole),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToProjectMembers(members []svcmodel.ProjectMember) []ProjectMember {
	m := make([]ProjectMember, 0, len(members))
	for _, member := range members {
		m = append(m, ToProjectMember(member))
	}
	return m
}
