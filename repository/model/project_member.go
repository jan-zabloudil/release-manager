package model

import (
	svcmodel "release-manager/service/model"

	"github.com/jackc/pgx/v5"
)

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
