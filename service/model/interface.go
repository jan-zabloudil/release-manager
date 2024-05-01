package model

import (
	"context"

	cryptox "release-manager/pkg/crypto"

	"github.com/google/uuid"
)

type AuthRepository interface {
	ReadUserIDForToken(ctx context.Context, token string) (uuid.UUID, error)
}

type EnvironmentRepository interface {
	Create(ctx context.Context, env Environment) error
	Read(ctx context.Context, envID uuid.UUID) (Environment, error)
	ReadByNameForProject(ctx context.Context, projectID uuid.UUID, name string) (Environment, error)
	ReadAllForProject(ctx context.Context, projectID uuid.UUID) ([]Environment, error)
	Delete(ctx context.Context, envID uuid.UUID) error
	Update(ctx context.Context, env Environment) error
}

type ProjectRepository interface {
	Create(ctx context.Context, p Project) error
	Read(ctx context.Context, id uuid.UUID) (Project, error)
	ReadAll(ctx context.Context) ([]Project, error)
	Update(ctx context.Context, p Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	Read(ctx context.Context, id uuid.UUID) (User, error)
	ReadAll(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type AuthService interface {
	AuthorizeAdminRole(ctx context.Context, userID uuid.UUID) error
	AuthorizeRole(ctx context.Context, userID uuid.UUID, role UserRole) error
}

type SettingsRepository interface {
	Update(ctx context.Context, c Settings) error
	Read(ctx context.Context) (Settings, error)
}

type ProjectInvitationRepository interface {
	Create(ctx context.Context, i ProjectInvitation) error
	Read(ctx context.Context, id uuid.UUID) (ProjectInvitation, error)
	ReadByEmailForProject(ctx context.Context, email string, projectID uuid.UUID) (ProjectInvitation, error)
	ReadByTokenHashAndStatus(ctx context.Context, hash cryptox.Hash, status ProjectInvitationStatus) (ProjectInvitation, error)
	ReadAllForProject(ctx context.Context, projectID uuid.UUID) ([]ProjectInvitation, error)
	Update(ctx context.Context, i ProjectInvitation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectService interface {
	Get(ctx context.Context, projectID, authUserID uuid.UUID) (Project, error)
}

type GithubClient interface {
	ListTagsForRepository(ctx context.Context, repo GithubRepository) ([]GitTag, error)
	SetToken(token string)
}

type SettingsService interface {
	GetGithubSettings(ctx context.Context) (GithubSettings, error)
}