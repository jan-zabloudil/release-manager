package service

import (
	"context"

	cryptox "release-manager/pkg/crypto"
	"release-manager/service/model"

	"github.com/google/uuid"
)

type authRepository interface {
	ReadUserIDForToken(ctx context.Context, token string) (uuid.UUID, error)
}

type projectRepository interface {
	CreateProject(ctx context.Context, p model.Project) error
	ReadProject(ctx context.Context, id uuid.UUID) (model.Project, error)
	ReadAllProjects(ctx context.Context) ([]model.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
	UpdateProject(ctx context.Context, p model.Project) error

	CreateEnvironment(ctx context.Context, env model.Environment) error
	ReadEnvironment(ctx context.Context, envID uuid.UUID) (model.Environment, error)
	ReadEnvironmentByNameForProject(ctx context.Context, projectID uuid.UUID, name string) (model.Environment, error)
	UpdateEnvironment(ctx context.Context, env model.Environment) error
	DeleteEnvironment(ctx context.Context, envID uuid.UUID) error
	ReadAllEnvironmentsForProject(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error)

	CreateInvitation(ctx context.Context, i model.ProjectInvitation) error
	ReadInvitation(ctx context.Context, invitationID uuid.UUID) (model.ProjectInvitation, error)
	ReadAllInvitationsForProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectInvitation, error)
	ReadInvitationByTokenHashAndStatus(ctx context.Context, hash cryptox.Hash, status model.ProjectInvitationStatus) (model.ProjectInvitation, error)
	ReadInvitationByEmailForProject(ctx context.Context, email string, projectID uuid.UUID) (model.ProjectInvitation, error)
	DeleteInvitation(ctx context.Context, invitationID uuid.UUID) error
	UpdateInvitation(ctx context.Context, i model.ProjectInvitation) error

	CreateMember(ctx context.Context, member model.ProjectMember) error
	ReadMembersForProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectMember, error)
	ReadMember(ctx context.Context, projectID, userID uuid.UUID) (model.ProjectMember, error)
	DeleteMember(ctx context.Context, projectID, userID uuid.UUID) error
}

type userRepository interface {
	Read(ctx context.Context, id uuid.UUID) (model.User, error)
	ReadByEmail(ctx context.Context, email string) (model.User, error)
	ReadAll(ctx context.Context) ([]model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type settingsRepository interface {
	Update(ctx context.Context, c model.Settings) error
	Read(ctx context.Context) (model.Settings, error)
}

type authGuard interface {
	AuthorizeAdminRole(ctx context.Context, userID uuid.UUID) error
	AuthorizeRole(ctx context.Context, userID uuid.UUID, role model.UserRole) error
}

type settingsGetter interface {
	GetGithubSettings(ctx context.Context) (model.GithubSettings, error)
}

type userGetter interface {
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type projectInvitationSender interface {
	SendProjectInvitation(ctx context.Context, input model.ProjectInvitationInput)
}

type githubRepositoryManager interface {
	ListTagsForRepository(ctx context.Context, repo model.GithubRepository) ([]model.GitTag, error)
	SetToken(token string)
}

type emailSender interface {
	SendEmailAsync(ctx context.Context, subject, text, html string, recipients ...string)
}

type Service struct {
	Auth     *AuthService
	User     *UserService
	Project  *ProjectService
	Settings *SettingsService
}

func NewService(
	authRepo authRepository,
	userRepo userRepository,
	projectRepo projectRepository,
	settingsRepo settingsRepository,
	githubRepoManager githubRepositoryManager,
	emailSender emailSender,
) *Service {
	authSvc := NewAuthService(authRepo, userRepo)
	userSvc := NewUserService(authSvc, userRepo)
	settingsSvc := NewSettingsService(authSvc, settingsRepo)
	emailSvc := NewEmailService(emailSender)
	projectSvc := NewProjectService(authSvc, settingsSvc, userSvc, githubRepoManager, emailSvc, projectRepo)

	return &Service{
		Auth:     authSvc,
		User:     userSvc,
		Project:  projectSvc,
		Settings: settingsSvc,
	}
}
