package query

import (
	_ "embed"
)

var (
	//go:embed scripts/create_release.sql
	CreateRelease string
	//go:embed scripts/read_release_for_project.sql
	ReadReleaseForProject string

	//go:embed scripts/read_user.sql
	ReadUser string
	//go:embed scripts/read_user_by_email.sql
	ReadUserByEmail string
	//go:embed scripts/list_users.sql
	ListUsers string

	//go:embed scripts/read_settings.sql
	ReadSettings string
	//go:embed scripts/upsert_settings.sql
	UpsertSettings string

	//go:embed scripts/read_project.sql
	ReadProject string
	//go:embed scripts/delete_project.sql
	DeleteProject string
	//go:embed scripts/create_project.sql
	CreateProject string

	//go:embed scripts/create_project_member.sql
	CreateProjectMember string

	//go:embed scripts/create_environment.sql
	CreateEnvironment string
)
