package query

import (
	_ "embed"
	"strconv"
)

var (
	//go:embed scripts/create_release.sql
	CreateRelease string
	//go:embed scripts/read_release.sql
	ReadRelease string
	//go:embed scripts/read_release_for_project.sql
	ReadReleaseForProject string
	//go:embed scripts/delete_release.sql
	DeleteRelease string
	//go:embed scripts/delete_release_by_git_tag.sql
	DeleteReleaseByGitTag string
	//go:embed scripts/list_releases_for_project.sql
	ListReleasesForProject string
	//go:embed scripts/update_release.sql
	UpdateRelease string

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
	//go:embed scripts/update_project.sql
	UpdateProject string
	//go:embed scripts/list_projects.sql
	ListProjects string
	//go:embed scripts/list_projects_for_user.sql
	ListProjectsForUser string

	//go:embed scripts/read_invitation_by_hash.sql
	ReadInvitationByHash string
	//go:embed scripts/list_invitations_for_project.sql
	ListInvitationsForProject string
	//go:embed scripts/delete_invitation.sql
	DeleteInvitation string
	//go:embed scripts/delete_invitation_by_hash_and_status.sql
	DeleteInvitationByHashAndStatus string
	//go:embed scripts/create_invitation.sql
	CreateInvitation string
	//go:embed scripts/update_invitation.sql
	UpdateInvitation string

	//go:embed scripts/create_member.sql
	CreateMember string
	//go:embed scripts/delete_member.sql
	DeleteMember string
	//go:embed scripts/list_members_for_project.sql
	ListMembersForProject string
	//go:embed scripts/list_members_for_user.sql
	ListMembersForUser string
	//go:embed scripts/read_member.sql
	ReadMember string
	//go:embed scripts/read_member_by_email.sql
	ReadMemberByEmail string
	//go:embed scripts/update_member.sql
	UpdateMember string

	//go:embed scripts/create_environment.sql
	CreateEnvironment string
	//go:embed scripts/list_environments_for_project.sql
	ListEnvironmentsForProject string
	//go:embed scripts/delete_environment.sql
	DeleteEnvironment string
	//go:embed scripts/read_environment.sql
	ReadEnvironment string
	//go:embed scripts/update_environment.sql
	UpdateEnvironment string

	//go:embed scripts/create_deployment.sql
	CreateDeployment string
	//go:embed scripts/list_deployments_for_project.sql
	ListDeploymentsForProject string
	//go:embed scripts/read_last_deployment_for_release.sql
	ReadLastDeploymentForRelease string
)

func AppendForUpdate(query string) string {
	return query + "\nFOR UPDATE"
}

func AppendLimit(query string, limit int) string {
	return query + "\nLIMIT " + strconv.Itoa(limit)
}
