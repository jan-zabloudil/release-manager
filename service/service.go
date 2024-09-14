package service

import (
	"context"
	"net/url"

	cryptox "release-manager/pkg/crypto"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type projectRepository interface {
	CreateProjectWithOwner(ctx context.Context, p model.Project, owner model.ProjectMember) error
	ReadProject(ctx context.Context, id uuid.UUID) (model.Project, error)
	ListProjects(ctx context.Context) ([]model.Project, error)
	ListProjectsForUser(ctx context.Context, userID uuid.UUID) ([]model.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
	UpdateProject(ctx context.Context, projectID uuid.UUID, fn model.UpdateProjectFunc) (model.Project, error)

	CreateEnvironment(ctx context.Context, env model.Environment) error
	ReadEnvironment(ctx context.Context, projectID, envID uuid.UUID) (model.Environment, error)
	UpdateEnvironment(ctx context.Context, projectID, envID uuid.UUID, fn model.UpdateEnvironmentFunc) (model.Environment, error)
	DeleteEnvironment(ctx context.Context, projectID, envID uuid.UUID) error
	ListEnvironmentsForProject(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error)

	CreateInvitation(ctx context.Context, i model.ProjectInvitation) error
	ListInvitationsForProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectInvitation, error)
	ReadPendingInvitationByHash(ctx context.Context, hash cryptox.Hash) (model.ProjectInvitation, error)
	DeleteInvitation(ctx context.Context, projectID, invitationID uuid.UUID) error
	DeleteInvitationByTokenHashAndStatus(ctx context.Context, hash cryptox.Hash, status model.ProjectInvitationStatus) error
	AcceptPendingInvitation(ctx context.Context, invitationID uuid.UUID, fn model.AcceptProjectInvitationFunc) error

	CreateMember(ctx context.Context, member model.ProjectMember) error
	ListMembersForProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectMember, error)
	ListMembersForUser(ctx context.Context, userID uuid.UUID) ([]model.ProjectMember, error)
	ReadMember(ctx context.Context, projectID, userID uuid.UUID) (model.ProjectMember, error)
	ReadMemberByEmail(ctx context.Context, projectID uuid.UUID, email string) (model.ProjectMember, error)
	DeleteMember(ctx context.Context, projectID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, projectID, userID uuid.UUID, fn model.UpdateProjectMemberFunc) (model.ProjectMember, error)
}

type userRepository interface {
	Read(ctx context.Context, id uuid.UUID) (model.User, error)
	ReadByEmail(ctx context.Context, email string) (model.User, error)
	ListAll(ctx context.Context) ([]model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type settingsRepository interface {
	Update(ctx context.Context, fn model.UpdateSettingsFunc) (model.Settings, error)
	Read(ctx context.Context) (model.Settings, error)
}

type releaseRepository interface {
	CreateRelease(ctx context.Context, r model.Release) error
	ReadRelease(ctx context.Context, projectID, releaseID uuid.UUID) (model.Release, error)
	DeleteRelease(ctx context.Context, projectID, releaseID uuid.UUID) error
	DeleteReleaseByGitTag(ctx context.Context, githubOwnerSlug, githubRepoSlug, gitTag string) error
	ListReleasesForProject(ctx context.Context, projectID uuid.UUID) ([]model.Release, error)
	UpdateRelease(ctx context.Context, projectID, releaseID uuid.UUID, fn model.UpdateReleaseFunc) (model.Release, error)

	CreateDeployment(ctx context.Context, d model.Deployment) error
	ListDeploymentsForProject(ctx context.Context, params model.DeploymentFilterParams, projectID uuid.UUID) ([]model.Deployment, error)
	ReadLastDeploymentForRelease(ctx context.Context, projectID, releaseID uuid.UUID) (model.Deployment, error)

	CreateReleaseAttachment(ctx context.Context, a model.ReleaseAttachment, releaseID uuid.UUID) error
}

type authGuard interface {
	AuthorizeUserRoleAdmin(ctx context.Context, userID uuid.UUID) error
	AuthorizeUserRoleUser(ctx context.Context, userID uuid.UUID) error
	AuthorizeProjectRoleEditor(ctx context.Context, projectID, userID uuid.UUID) error
	AuthorizeProjectRoleViewer(ctx context.Context, projectID, userID uuid.UUID) error
}

type settingsGetter interface {
	GetGithubToken(ctx context.Context) (string, error)
	GetSlackToken(ctx context.Context) (string, error)
	GetDefaultReleaseMessage(ctx context.Context) (string, error)
}

type userGetter interface {
	Get(ctx context.Context, userID uuid.UUID) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type projectGetter interface {
	GetProject(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (model.Project, error)
	ProjectExists(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (bool, error)
}

type environmentGetter interface {
	GetEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) (model.Environment, error)
	EnvironmentExists(ctx context.Context, projectID, envID, authUserID uuid.UUID) (bool, error)
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

type fileManager interface {
	FileExists(path string) (bool, error)
	GenerateFileURL(path string) (url.URL, error)
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
	fileManager fileManager,
) *Service {
	authSvc := NewAuthorizationService(userRepo, projectRepo)
	userSvc := NewUserService(authSvc, userRepo)
	settingsSvc := NewSettingsService(authSvc, settingsRepo)
	projectSvc := NewProjectService(authSvc, settingsSvc, userSvc, emailSender, githubManager, projectRepo)
	releaseSvc := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackNotifier, githubManager, fileManager, releaseRepo)

	return &Service{
		Authorization: authSvc,
		User:          userSvc,
		Project:       projectSvc,
		Settings:      settingsSvc,
		Release:       releaseSvc,
	}
}
