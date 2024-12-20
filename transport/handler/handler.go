package handler

import (
	"context"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"
	"release-manager/transport/middleware"

	"github.com/go-chi/chi/v5"
)

type projectService interface {
	CreateProject(ctx context.Context, c svcmodel.CreateProjectInput, authUserID id.AuthUser) (svcmodel.Project, error)
	GetProject(ctx context.Context, projectID id.Project, authUserID id.AuthUser) (svcmodel.Project, error)
	ListProjects(ctx context.Context, authUserID id.AuthUser) ([]svcmodel.Project, error)
	UpdateProject(ctx context.Context, u svcmodel.UpdateProjectInput, projectID id.Project, authUserID id.AuthUser) error
	DeleteProject(ctx context.Context, projectID id.Project, authUserID id.AuthUser) error

	CreateEnvironment(ctx context.Context, c svcmodel.CreateEnvironmentInput, authUserID id.AuthUser) (svcmodel.Environment, error)
	GetEnvironment(ctx context.Context, projectID id.Project, envID id.Environment, authUserID id.AuthUser) (svcmodel.Environment, error)
	ListEnvironments(ctx context.Context, projectID id.Project, authUserID id.AuthUser) ([]svcmodel.Environment, error)
	DeleteEnvironment(ctx context.Context, projectID id.Project, envID id.Environment, authUserID id.AuthUser) error
	UpdateEnvironment(ctx context.Context, u svcmodel.UpdateEnvironmentInput, projectID id.Project, envID id.Environment, authUserID id.AuthUser) error

	SetGithubRepoForProject(ctx context.Context, rawRepoURL string, projectID id.Project, authUserID id.AuthUser) error
	GetGithubRepoForProject(ctx context.Context, projectID id.Project, authUserID id.AuthUser) (svcmodel.GithubRepo, error)
	ListGithubRepoTags(ctx context.Context, projectID id.Project, authUserID id.AuthUser) ([]svcmodel.GitTag, error)

	Invite(ctx context.Context, c svcmodel.CreateProjectInvitationInput, authUserID id.AuthUser) (svcmodel.ProjectInvitation, error)
	ListInvitations(ctx context.Context, projectID id.Project, authUserID id.AuthUser) ([]svcmodel.ProjectInvitation, error)
	CancelInvitation(ctx context.Context, projectID id.Project, invitationID id.ProjectInvitation, authUserID id.AuthUser) error
	AcceptInvitation(ctx context.Context, tkn svcmodel.ProjectInvitationToken) error
	RejectInvitation(ctx context.Context, tkn svcmodel.ProjectInvitationToken) error

	ListMembersForProject(ctx context.Context, projectID id.Project, authUserID id.AuthUser) ([]svcmodel.ProjectMember, error)
	ListMembersForUser(ctx context.Context, authUserID id.AuthUser) ([]svcmodel.ProjectMember, error)
	DeleteMember(ctx context.Context, projectID id.Project, userID id.User, authUserID id.AuthUser) error
	UpdateMemberRole(ctx context.Context, newRole svcmodel.ProjectRole, projectID id.Project, userID id.User, authUserID id.AuthUser) error
}

type userService interface {
	GetAuthenticated(ctx context.Context, authUserID id.AuthUser) (svcmodel.User, error)
	GetForAdmin(ctx context.Context, userID id.User, authUserID id.AuthUser) (svcmodel.User, error)
	ListAllForAdmin(ctx context.Context, authUserID id.AuthUser) ([]svcmodel.User, error)
	DeleteForAdmin(ctx context.Context, userID id.User, authUserID id.AuthUser) error
}

type settingsService interface {
	Update(ctx context.Context, u svcmodel.UpdateSettingsInput, authUserID id.AuthUser) error
	Get(ctx context.Context, authUserID id.AuthUser) (svcmodel.Settings, error)
}

type releaseService interface {
	CreateRelease(ctx context.Context, input svcmodel.CreateReleaseInput, projectID id.Project, authUserID id.AuthUser) (svcmodel.Release, error)
	GetRelease(ctx context.Context, releaseID id.Release, authUserID id.AuthUser) (svcmodel.Release, error)
	DeleteRelease(ctx context.Context, input svcmodel.DeleteReleaseInput, releaseID id.Release, authUserID id.AuthUser) error
	DeleteReleaseOnGitTagRemoval(ctx context.Context, input svcmodel.GithubTagDeletionWebhookInput) error
	UpdateRelease(ctx context.Context, input svcmodel.UpdateReleaseInput, releaseID id.Release, authUserID id.AuthUser) error
	ListReleasesForProject(ctx context.Context, projectID id.Project, authUserID id.AuthUser) ([]svcmodel.Release, error)
	SendReleaseNotification(ctx context.Context, releaseID id.Release, authUserID id.AuthUser) error
	UpsertGithubRelease(ctx context.Context, releaseID id.Release, authUserID id.AuthUser) error
	GenerateGithubReleaseNotes(ctx context.Context, input svcmodel.GithubReleaseNotesInput, projectID id.Project, authUserID id.AuthUser) (svcmodel.GithubReleaseNotes, error)

	CreateDeployment(ctx context.Context, input svcmodel.CreateDeploymentInput, projectID id.Project, authUserID id.AuthUser) (svcmodel.Deployment, error)
	ListDeploymentsForProject(ctx context.Context, params svcmodel.ListDeploymentsFilterParams, projectID id.Project, authUserID id.AuthUser) ([]svcmodel.Deployment, error)
}

type Handler struct {
	Mux         *chi.Mux
	AuthClient  middleware.AuthClient
	UserSvc     userService
	ProjectSvc  projectService
	SettingsSvc settingsService
	ReleaseSvc  releaseService
}

func NewHandler(
	authClient middleware.AuthClient,
	userSvc userService,
	projectSvc projectService,
	settingsSvc settingsService,
	releaseSvc releaseService,
) *Handler {
	h := &Handler{
		Mux:         chi.NewRouter(),
		AuthClient:  authClient,
		UserSvc:     userSvc,
		ProjectSvc:  projectSvc,
		SettingsSvc: settingsSvc,
		ReleaseSvc:  releaseSvc,
	}

	h.setupRoutes()

	return h
}
