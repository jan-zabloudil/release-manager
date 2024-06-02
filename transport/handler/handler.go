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
	SetGithubRepoForProject(ctx context.Context, rawRepoURL string, projectID uuid.UUID, authUserID uuid.UUID) error

	CreateEnvironment(ctx context.Context, c svcmodel.CreateEnvironmentInput, authUserID uuid.UUID) (svcmodel.Environment, error)
	GetEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) (svcmodel.Environment, error)
	ListEnvironments(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.Environment, error)
	DeleteEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) error
	UpdateEnvironment(ctx context.Context, u svcmodel.UpdateEnvironmentInput, projectID, envID, authUserID uuid.UUID) (svcmodel.Environment, error)

	ListGithubRepoTags(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.GitTag, error)

	Invite(ctx context.Context, c svcmodel.CreateProjectInvitationInput, authUserID uuid.UUID) (svcmodel.ProjectInvitation, error)
	ListInvitations(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.ProjectInvitation, error)
	CancelInvitation(ctx context.Context, projectID, invitationID, authUserID uuid.UUID) error
	AcceptInvitation(ctx context.Context, tkn cryptox.Token) error
	RejectInvitation(ctx context.Context, tkn cryptox.Token) error

	ListMembers(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.ProjectMember, error)
	DeleteMember(ctx context.Context, projectID, userID, authUserID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, newRole svcmodel.ProjectRole, projectID, userID, authUserID uuid.UUID) (svcmodel.ProjectMember, error)
}

type userService interface {
	Get(ctx context.Context, id, authUserID uuid.UUID) (svcmodel.User, error)
	ListAll(ctx context.Context, authUserID uuid.UUID) ([]svcmodel.User, error)
	Delete(ctx context.Context, id, authUserID uuid.UUID) error
}

type settingsService interface {
	Update(ctx context.Context, u svcmodel.UpdateSettingsInput, authUserID uuid.UUID) (svcmodel.Settings, error)
	Get(ctx context.Context, authUserID uuid.UUID) (svcmodel.Settings, error)
}

type releaseService interface {
	CreateRelease(ctx context.Context, input svcmodel.CreateReleaseInput, projectID, authUserID uuid.UUID) (svcmodel.Release, error)
	GetRelease(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) (svcmodel.Release, error)
	DeleteRelease(ctx context.Context, input svcmodel.DeleteReleaseInput, projectID, releaseID, authUserID uuid.UUID) error
	UpdateRelease(ctx context.Context, input svcmodel.UpdateReleaseInput, projectID, releaseID, authUserID uuid.UUID) (svcmodel.Release, error)
	ListReleasesForProject(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.Release, error)
	SendReleaseNotification(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) error
	UpsertGithubRelease(ctx context.Context, projectID, releaseID, authUserID uuid.UUID) error

	CreateDeployment(ctx context.Context, input svcmodel.CreateDeploymentInput, projectID, authUserID uuid.UUID) (svcmodel.Deployment, error)
	ListDeploymentsForProject(ctx context.Context, input svcmodel.DeploymentFilterParams, projectID, authUserID uuid.UUID) ([]svcmodel.Deployment, error)
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
