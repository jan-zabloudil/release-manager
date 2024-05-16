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
	ReadAllProjects(ctx context.Context) ([]model.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
	UpdateProject(ctx context.Context, projectID uuid.UUID, fn model.UpdateProjectFunc) (model.Project, error)

	CreateEnvironment(ctx context.Context, env model.Environment) error
	ReadEnvironment(ctx context.Context, envID uuid.UUID) (model.Environment, error)
	ReadEnvironmentByNameForProject(ctx context.Context, projectID uuid.UUID, name string) (model.Environment, error)
	UpdateEnvironment(ctx context.Context, env model.Environment) error
	DeleteEnvironment(ctx context.Context, projectID, envID uuid.UUID) error
	ListEnvironmentsForProject(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error)

	CreateInvitation(ctx context.Context, i model.ProjectInvitation) error
	ReadAllInvitationsForProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectInvitation, error)
	ReadInvitationByTokenHashAndStatus(ctx context.Context, hash cryptox.Hash, status model.ProjectInvitationStatus) (model.ProjectInvitation, error)
	ReadInvitationByEmailForProject(ctx context.Context, email string, projectID uuid.UUID) (model.ProjectInvitation, error)
	DeleteInvitation(ctx context.Context, projectID, invitationID uuid.UUID) error
	DeleteInvitationByTokenHashAndStatus(ctx context.Context, hash cryptox.Hash, status model.ProjectInvitationStatus) error
	UpdateInvitation(ctx context.Context, i model.ProjectInvitation) error

	CreateMember(ctx context.Context, member model.ProjectMember) error
	ListMembersForProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectMember, error)
	ReadMember(ctx context.Context, projectID, userID uuid.UUID) (model.ProjectMember, error)
	ReadMemberByEmail(ctx context.Context, projectID uuid.UUID, email string) (model.ProjectMember, error)
	DeleteMember(ctx context.Context, projectID, userID uuid.UUID) error
	UpdateMember(ctx context.Context, m model.ProjectMember) error
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
	ReadForProject(ctx context.Context, projectID, releaseID uuid.UUID) (model.Release, error)
}

type authGuard interface {
	AuthorizeUserRoleAdmin(ctx context.Context, userID uuid.UUID) error
	AuthorizeUserRole(ctx context.Context, userID uuid.UUID, model model.UserRole) error
	AuthorizeProjectRoleEditor(ctx context.Context, projectID, userID uuid.UUID) error
	AuthorizeProjectRoleViewer(ctx context.Context, projectID, userID uuid.UUID) error
}

type settingsGetter interface {
	GetGithubToken(ctx context.Context) (string, error)
}

type userGetter interface {
	Get(ctx context.Context, userID, authUserID uuid.UUID) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type projectGetter interface {
	GetProject(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (model.Project, error)
}

type githubManager interface {
	ReadTagsForRepository(ctx context.Context, token string, repoURL url.URL) ([]model.GitTag, error)
}

type emailSender interface {
	SendEmailAsync(ctx context.Context, email model.Email)
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
) *Service {
	authSvc := NewAuthorizationService(userRepo, projectRepo)
	userSvc := NewUserService(authSvc, userRepo)
	settingsSvc := NewSettingsService(authSvc, settingsRepo)
	projectSvc := NewProjectService(authSvc, settingsSvc, userSvc, emailSender, githubManager, projectRepo)
	releaseSvc := NewReleaseService(projectSvc, releaseRepo)

	return &Service{
		Authorization: authSvc,
		User:          userSvc,
		Project:       projectSvc,
		Settings:      settingsSvc,
		Release:       releaseSvc,
	}
}
