package query

import (
	_ "embed"
)

var (
	//go:embed scripts/create_release.sql
	CreateRelease string
	//go:embed scripts/read_release.sql
	ReadRelease string
	//go:embed scripts/delete_release.sql
	DeleteRelease string

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
	//go:embed scripts/read_project_for_update.sql
	ReadProjectForUpdate string
	//go:embed scripts/delete_project.sql
	DeleteProject string
	//go:embed scripts/create_project.sql
	CreateProject string
	//go:embed scripts/update_project.sql
	UpdateProject string
	//go:embed scripts/list_projects.sql
	ListProjects string
	//go:embed scripts/list_projects_for_user.sql
	ListProjectsForUser string

	//go:embed scripts/read_invitation_by_email.sql
	ReadInvitationByEmail string
	//go:embed scripts/list_invitations_for_project.sql
	ListInvitationsForProject string
	//go:embed scripts/delete_invitation.sql
	DeleteInvitation string
	//go:embed scripts/delete_invitation_by_email_and_project_id.sql
	DeleteInvitationByEmailAndProjectID string
	//go:embed scripts/delete_invitation_by_hash_and_status.sql
	DeleteInvitationByHashAndStatus string

	//go:embed scripts/create_member.sql
	CreateMember string
	//go:embed scripts/delete_member.sql
	DeleteMember string
	//go:embed scripts/list_members_for_project.sql
	ListMembersForProject string
	//go:embed scripts/read_member.sql
	ReadMember string
	//go:embed scripts/read_member_by_email.sql
	ReadMemberByEmail string

	//go:embed scripts/create_environment.sql
	CreateEnvironment string
	//go:embed scripts/list_environments_for_project.sql
	ListEnvironmentsForProject string
	//go:embed scripts/delete_environment.sql
	DeleteEnvironment string
	//go:embed scripts/read_environment.sql
	ReadEnvironment string
	//go:embed scripts/read_environment_by_name.sql
	ReadEnvironmentByName string
)
