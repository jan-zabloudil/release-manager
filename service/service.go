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

type environmentRepository interface {
	Create(ctx context.Context, env model.Environment) error
	Read(ctx context.Context, envID uuid.UUID) (model.Environment, error)
	ReadByNameForProject(ctx context.Context, projectID uuid.UUID, name string) (model.Environment, error)
	ReadAllForProject(ctx context.Context, projectID uuid.UUID) ([]model.Environment, error)
	Delete(ctx context.Context, envID uuid.UUID) error
	Update(ctx context.Context, env model.Environment) error
}

type projectRepository interface {
	Create(ctx context.Context, p model.Project) error
	Read(ctx context.Context, id uuid.UUID) (model.Project, error)
	ReadAll(ctx context.Context) ([]model.Project, error)
	Update(ctx context.Context, p model.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository interface {
	Read(ctx context.Context, id uuid.UUID) (model.User, error)
	ReadAll(ctx context.Context) ([]model.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type settingsRepository interface {
	Update(ctx context.Context, c model.Settings) error
	Read(ctx context.Context) (model.Settings, error)
}

type projectInvitationRepository interface {
	Create(ctx context.Context, i model.ProjectInvitation) error
	Read(ctx context.Context, id uuid.UUID) (model.ProjectInvitation, error)
	ReadByEmailForProject(ctx context.Context, email string, projectID uuid.UUID) (model.ProjectInvitation, error)
	ReadByTokenHashAndStatus(ctx context.Context, hash cryptox.Hash, status model.ProjectInvitationStatus) (model.ProjectInvitation, error)
	ReadAllForProject(ctx context.Context, projectID uuid.UUID) ([]model.ProjectInvitation, error)
	Update(ctx context.Context, i model.ProjectInvitation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type authGuard interface {
	AuthorizeAdminRole(ctx context.Context, userID uuid.UUID) error
	AuthorizeRole(ctx context.Context, userID uuid.UUID, role model.UserRole) error
}

type projectGetter interface {
	Get(ctx context.Context, projectID, authUserID uuid.UUID) (model.Project, error)
}

type githubClient interface {
	ListTagsForRepository(ctx context.Context, repo model.GithubRepository) ([]model.GitTag, error)
	SetToken(token string)
}

type settingsGetter interface {
	GetGithubSettings(ctx context.Context) (model.GithubSettings, error)
}

type Service struct {
	Auth              *AuthService
	User              *UserService
	Project           *ProjectService
	Settings          *SettingsService
	ProjectMembership *ProjectMembershipService
}

func NewService(
	ar authRepository,
	ur userRepository,
	pr projectRepository,
	env environmentRepository,
	sr settingsRepository,
	pi projectInvitationRepository,
	gc githubClient,
) *Service {
	authSvc := NewAuthService(ar, ur)
	settingsSvc := NewSettingsService(authSvc, sr)
	projectSvc := NewProjectService(authSvc, settingsSvc, pr, env, pi, gc)

	return &Service{
		Auth:              authSvc,
		User:              NewUserService(authSvc, ur),
		Project:           projectSvc,
		Settings:          settingsSvc,
		ProjectMembership: NewProjectMembershipService(authSvc, projectSvc, pi),
	}
}
