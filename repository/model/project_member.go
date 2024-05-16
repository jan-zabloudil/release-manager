package model

import (
	"time"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

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

func ScanToSvcProjectMember(row pgx.Row) (svcmodel.ProjectMember, error) {
	var m svcmodel.ProjectMember

	err := row.Scan(
		&m.User.ID,
		&m.User.Email,
		&m.User.Name,
		&m.User.AvatarURL,
		&m.User.Role,
		&m.User.CreatedAt,
		&m.User.UpdatedAt,
		&m.ProjectID,
		&m.ProjectRole,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return svcmodel.ProjectMember{}, err
	}

	return m, nil
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

func ToUpdateProjectMemberInput(p svcmodel.ProjectMember) UpdateProjectMemberInput {
	return UpdateProjectMemberInput{
		ProjectRole: string(p.ProjectRole),
		UpdatedAt:   p.UpdatedAt,
	}
}
