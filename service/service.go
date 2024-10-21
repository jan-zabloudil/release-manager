package service

import (
	"context"
	"net/url"

	cryptox "release-manager/pkg/crypto"
	"release-manager/pkg/id"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type projectRepository interface {
	CreateProjectWithOwner(ctx context.Context, p model.Project, owner model.ProjectMember) error
	ReadProject(ctx context.Context, id uuid.UUID) (model.Project, error)
	ListProjects(ctx context.Context) ([]model.Project, error)
	ListProjectsForUser(ctx context.Context, userID id.AuthUser) ([]model.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
	UpdateProject(
		ctx context.Context,
		projectID uuid.UUID,
		updateFn func(p model.Project) (model.Project, error),
	) error

	CreateEnvironment(ctx context.Context, env model.Environment) error
	ReadEnvironment(ctx context.Context, projectID, envID uuid.UUID) (model.Environment, error)
	UpdateEnvironment(
		ctx context.Context,
		projectID,
		envID uuid.UUID,
		updateFn func(e model.Environment) (model.Environment, error),
	) error
	DeleteEnvironment(ctx context.Context, projectID, envID uuid.UUID) error
	ListEnvironmentsForProject(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error)

	CreateInvitation(ctx context.Context, i model.ProjectInvitation) error
	ListInvitationsForProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectInvitation, error)
	DeleteInvitation(ctx context.Context, projectID, invitationID uuid.UUID) error
	DeleteInvitationByTokenHashAndStatus(ctx context.Context, hash cryptox.Hash, status model.ProjectInvitationStatus) error
	UpdateInvitation(
		ctx context.Context,
		invitationHash cryptox.Hash,
		updateFn func(i model.ProjectInvitation) (model.ProjectInvitation, error),
	) error

	CreateMember(
		ctx context.Context,
		invitationHash cryptox.Hash,
		createMemberFn func(i model.ProjectInvitation) (model.ProjectMember, error),
	) error
	ListMembersForProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectMember, error)
	ListMembersForUser(ctx context.Context, userID id.AuthUser) ([]model.ProjectMember, error)
	ReadMember(ctx context.Context, projectID uuid.UUID, userID id.User) (model.ProjectMember, error)
	ReadMemberByEmail(ctx context.Context, projectID uuid.UUID, email string) (model.ProjectMember, error)
	DeleteMember(ctx context.Context, projectID uuid.UUID, userID id.User) error
	UpdateMember(
		ctx context.Context,
		projectID uuid.UUID,
		userID id.User,
		updateFn func(m model.ProjectMember) (model.ProjectMember, error),
	) error
}

type userRepository interface {
	Read(ctx context.Context, id id.User) (model.User, error)
	ReadByEmail(ctx context.Context, email string) (model.User, error)
	ListAll(ctx context.Context) ([]model.User, error)
	Delete(ctx context.Context, id id.User) error
}

type settingsRepository interface {
	Upsert(
		ctx context.Context,
		upsertFn func(s model.Settings) (model.Settings, error),
	) error
	Read(ctx context.Context) (model.Settings, error)
}

type releaseRepository interface {
	CreateRelease(ctx context.Context, r model.Release) error
	ReadRelease(ctx context.Context, releaseID uuid.UUID) (model.Release, error)
	ReadReleaseForProject(ctx context.Context, projectID, releaseID uuid.UUID) (model.Release, error)
	DeleteRelease(ctx context.Context, releaseID uuid.UUID) error
	DeleteReleaseByGitTag(ctx context.Context, githubOwnerSlug, githubRepoSlug, gitTag string) error
	ListReleasesForProject(ctx context.Context, projectID uuid.UUID) ([]model.Release, error)
	UpdateRelease(
		ctx context.Context,
		releaseID uuid.UUID,
		updateFn func(r model.Release) (model.Release, error),
	) error

	CreateDeployment(ctx context.Context, d model.Deployment) error
	ListDeploymentsForProject(ctx context.Context, params model.DeploymentFilterParams, projectID uuid.UUID) ([]model.Deployment, error)
	ReadLastDeploymentForRelease(ctx context.Context, releaseID uuid.UUID) (model.Deployment, error)
}

type authGuard interface {
	AuthorizeUserRoleAdmin(ctx context.Context, userID id.AuthUser) error
	AuthorizeUserRoleUser(ctx context.Context, userID id.AuthUser) error
	AuthorizeProjectRoleEditor(ctx context.Context, projectID uuid.UUID, userID id.AuthUser) error
	AuthorizeProjectRoleViewer(ctx context.Context, projectID uuid.UUID, userID id.AuthUser) error
	AuthorizeReleaseEditor(ctx context.Context, releaseID uuid.UUID, userID id.AuthUser) error
	AuthorizeReleaseViewer(ctx context.Context, releaseID uuid.UUID, userID id.AuthUser) error
}

type settingsGetter interface {
	GetGithubToken(ctx context.Context) (string, error)
	GetSlackToken(ctx context.Context) (string, error)
	GetDefaultReleaseMessage(ctx context.Context) (string, error)
}

type userGetter interface {
	GetAuthenticated(ctx context.Context, userID id.AuthUser) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type projectGetter interface {
	GetProject(ctx context.Context, projectID uuid.UUID, authUserID id.AuthUser) (model.Project, error)
}

type environmentGetter interface {
	GetEnvironment(ctx context.Context, projectID, envID uuid.UUID, authUserID id.AuthUser) (model.Environment, error)
}

type githubManager interface {
	ReadRepo(ctx context.Context, token, rawRepoURL string) (model.GithubRepo, error)
	ReadTagsForRepo(ctx context.Context, token string, repo model.GithubRepo) ([]model.GitTag, error)
	DeleteReleaseByTag(ctx context.Context, token string, repo model.GithubRepo, tagName string) error
	GenerateGitTagURL(ownerSlug, repoSlug, tagName string) (url.URL, error)
	TagExists(ctx context.Context, token string, repo model.GithubRepo, tagName string) (bool, error)
	UpsertRelease(ctx context.Context, token string, repo model.GithubRepo, rls model.Release) error
	GenerateReleaseNotes(ctx context.Context, token string, repo model.GithubRepo, input model.GithubGeneratedReleaseNotesInput) (model.GithubGeneratedReleaseNotes, error)
}

type emailSender interface {
	SendProjectInvitationEmailAsync(ctx context.Context, data model.ProjectInvitationEmailData, recipient string)
}

type slackNotifier interface {
	SendReleaseNotification(ctx context.Context, token, channel string, notification model.ReleaseNotification) error
}

type Service struct {
	Authorization *AuthorizationService
	User          *UserService
	Project       *ProjectService
	Settings      *SettingsService
	Release       *ReleaseService
}

func NewService(
	userRepo userRepository,
	projectRepo projectRepository,
	settingsRepo settingsRepository,
	releaseRepo releaseRepository,
	githubManager githubManager,
	emailSender emailSender,
	slackNotifier slackNotifier,
) *Service {
	authSvc := NewAuthorizationService(userRepo, projectRepo, releaseRepo)
	userSvc := NewUserService(authSvc, userRepo)
	settingsSvc := NewSettingsService(authSvc, settingsRepo)
	projectSvc := NewProjectService(authSvc, settingsSvc, userSvc, emailSender, githubManager, projectRepo)
	releaseSvc := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackNotifier, githubManager, releaseRepo)

	return &Service{
		Authorization: authSvc,
		User:          userSvc,
		Project:       projectSvc,
		Settings:      settingsSvc,
		Release:       releaseSvc,
	}
}
