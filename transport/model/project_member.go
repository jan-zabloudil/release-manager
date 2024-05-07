package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectMember struct {
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	ProjectRole string    `json:"project_role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
