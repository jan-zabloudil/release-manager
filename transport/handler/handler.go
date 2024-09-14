package handler

import (
	"context"

	cryptox "release-manager/pkg/crypto"
	svcmodel "release-manager/service/model"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type projectService interface {
	CreateProject(ctx context.Context, c svcmodel.CreateProjectInput, authUserID uuid.UUID) (svcmodel.Project, error)
	GetProject(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (svcmodel.Project, error)
	ListProjects(ctx context.Context, authUserID uuid.UUID) ([]svcmodel.Project, error)
	UpdateProject(ctx context.Context, u svcmodel.UpdateProjectInput, projectID, authUserID uuid.UUID) (svcmodel.Project, error)
	DeleteProject(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) error

	CreateEnvironment(ctx context.Context, c svcmodel.CreateEnvironmentInput, authUserID uuid.UUID) (svcmodel.Environment, error)
	GetEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) (svcmodel.Environment, error)
	ListEnvironments(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.Environment, error)
	DeleteEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) error
	UpdateEnvironment(ctx context.Context, u svcmodel.UpdateEnvironmentInput, projectID, envID, authUserID uuid.UUID) (svcmodel.Environment, error)

	SetGithubRepoForProject(ctx context.Context, rawRepoURL string, projectID uuid.UUID, authUserID uuid.UUID) error
	GetGithubRepoForProject(ctx context.Context, projectID, authUserID uuid.UUID) (svcmodel.GithubRepo, error)
	ListGithubRepoTags(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.GitTag, error)

	Invite(ctx context.Context, c svcmodel.CreateProjectInvitationInput, authUserID uuid.UUID) (svcmodel.ProjectInvitation, error)
	ListInvitations(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.ProjectInvitation, error)
	CancelInvitation(ctx context.Context, projectID, invitationID, authUserID uuid.UUID) error
	AcceptInvitation(ctx context.Context, tkn cryptox.Token) error
	RejectInvitation(ctx context.Context, tkn cryptox.Token) error

	ListMembersForProject(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.ProjectMember, error)
	ListMembersForUser(ctx context.Context, authUserID uuid.UUID) ([]svcmodel.ProjectMember, error)
	DeleteMember(ctx context.Context, projectID, userID, authUserID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, newRole svcmodel.ProjectRole, projectID, userID, authUserID uuid.UUID) (svcmodel.ProjectMember, error)
}

type userService interface {
	Get(ctx context.Context, userID uuid.UUID) (svcmodel.User, error)
	GetForAdmin(ctx context.Context, userID, authUserID uuid.UUID) (svcmodel.User, error)
	ListAllForAdmin(ctx context.Context, authUserID uuid.UUID) ([]svcmodel.User, error)
	DeleteForAdmin(ctx context.Context, userID, authUserID uuid.UUID) error
}

type settingsService interface {
	Update(ctx context.Context, u svcmodel.UpdateSettingsInput, authUserID uuid.UUID) (svcmodel.Settings, error)
	Get(ctx context.Context, authUserID uuid.UUID) (svcmodel.Settings, error)
	GetGithubWebhookSecret(ctx context.Context) (string, error)
}

type releaseService interface {
	CreateRelease(ctx context.Context, input svcmodel.CreateReleaseInput, projectID, authUserID uuid.UUID) (svcmodel.Release, error)
	GetRelease(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) (svcmodel.Release, error)
	DeleteRelease(ctx context.Context, input svcmodel.DeleteReleaseInput, projectID, releaseID, authUserID uuid.UUID) error
	DeleteReleaseByGitTag(ctx context.Context, githubOwnerSlug, githubRepoSlug, gitTag string) error
	UpdateRelease(ctx context.Context, input svcmodel.UpdateReleaseInput, projectID, releaseID, authUserID uuid.UUID) (svcmodel.Release, error)
	ListReleasesForProject(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.Release, error)
	SendReleaseNotification(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) error
	UpsertGithubRelease(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) error
	GenerateGithubReleaseNotes(ctx context.Context, input svcmodel.GithubGeneratedReleaseNotesInput, projectID, authUserID uuid.UUID) (svcmodel.GithubGeneratedReleaseNotes, error)

	CreateDeployment(ctx context.Context, input svcmodel.CreateDeploymentInput, projectID, authUserID uuid.UUID) (svcmodel.Deployment, error)
	ListDeploymentsForProject(ctx context.Context, params svcmodel.DeploymentFilterParams, projectID, authUserID uuid.UUID) ([]svcmodel.Deployment, error)

	CreateReleaseAttachment(ctx context.Context, input svcmodel.CreateReleaseAttachmentInput, projectID, releaseID, authUserID uuid.UUID) (svcmodel.ReleaseAttachment, error)
}

type authClient interface {
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}

type Handler struct {
	Mux         *chi.Mux
	AuthClient  authClient
	UserSvc     userService
	ProjectSvc  projectService
	SettingsSvc settingsService
	ReleaseSvc  releaseService
}

func NewHandler(
	authClient authClient,
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
