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
	Create(ctx context.Context, r model.Release) error
	Read(ctx context.Context, projectID, releaseID uuid.UUID) (model.Release, error)
	Delete(ctx context.Context, projectID, releaseID uuid.UUID) error
	ListForProject(ctx context.Context, projectID uuid.UUID) ([]model.Release, error)
	Update(ctx context.Context, projectID, releaseID uuid.UUID, fn model.UpdateReleaseFunc) (model.Release, error)
}

type authGuard interface {
	AuthorizeUserRoleAdmin(ctx context.Context, userID uuid.UUID) error
	AuthorizeUserRole(ctx context.Context, userID uuid.UUID, model model.UserRole) error
	AuthorizeProjectRoleEditor(ctx context.Context, projectID, userID uuid.UUID) error
	AuthorizeProjectRoleViewer(ctx context.Context, projectID, userID uuid.UUID) error
}

type settingsGetter interface {
	GetGithubToken(ctx context.Context) (string, error)
	GetSlackToken(ctx context.Context) (string, error)
	GetDefaultReleaseMessage(ctx context.Context) (string, error)
}

type userGetter interface {
	Get(ctx context.Context, userID, authUserID uuid.UUID) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type projectGetter interface {
	GetProject(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (model.Project, error)
	ProjectExists(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (bool, error)
}

type githubManager interface {
	ReadRepo(ctx context.Context, token, rawRepoURL string) (model.GithubRepo, error)
	ReadTagsForRepo(ctx context.Context, token string, repo model.GithubRepo) ([]model.GitTag, error)
	DeleteReleaseByTag(ctx context.Context, token string, repo model.GithubRepo, tagName string) error
	GenerateGitTagURL(repo model.GithubRepo, tagName string) (url.URL, error)
	TagExists(ctx context.Context, token string, repo model.GithubRepo, tagName string) (bool, error)
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
	authSvc := NewAuthorizationService(userRepo, projectRepo)
	userSvc := NewUserService(authSvc, userRepo)
	settingsSvc := NewSettingsService(authSvc, settingsRepo)
	projectSvc := NewProjectService(authSvc, settingsSvc, userSvc, emailSender, githubManager, projectRepo)
	releaseSvc := NewReleaseService(authSvc, projectSvc, settingsSvc, slackNotifier, githubManager, releaseRepo)

	return &Service{
		Authorization: authSvc,
		User:          userSvc,
		Project:       projectSvc,
		Settings:      settingsSvc,
		Release:       releaseSvc,
	}
}
