package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectMember struct {
	UserID        uuid.UUID `db:"user_id"`
	UserEmail     string    `db:"user_email"`
	UserName      string    `db:"user_name"`
	UserAvatarURL string    `db:"user_avatar_url"`
	UserRole      string    `db:"user_role"`
	UserCreatedAt time.Time `db:"user_created_at"`
	UserUpdatedAt time.Time `db:"user_updated_at"`
	ProjectID     uuid.UUID `db:"project_id"`
	ProjectRole   string    `db:"project_role"`
	CreatedAt     time.Time `db:"member_created_at"`
	UpdatedAt     time.Time `db:"member_updated_at"`
}

func ToSvcProjectMember(m ProjectMember) svcmodel.ProjectMember {
	return svcmodel.ProjectMember{
		User: svcmodel.User{
			ID:        m.UserID,
			Email:     m.UserEmail,
			Name:      m.UserName,
			AvatarURL: m.UserAvatarURL,
			Role:      svcmodel.UserRole(m.UserRole),
			CreatedAt: m.UserCreatedAt,
			UpdatedAt: m.UserUpdatedAt,
		},
		ProjectID:   m.ProjectID,
		ProjectRole: svcmodel.ProjectRole(m.ProjectRole),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func ToSvcProjectMembers(members []ProjectMember) []svcmodel.ProjectMember {
	m := make([]svcmodel.ProjectMember, 0, len(members))

	for _, member := range members {
		m = append(m, ToSvcProjectMember(member))
	}

	return m
}
