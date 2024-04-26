package model

import (
	"context"

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
