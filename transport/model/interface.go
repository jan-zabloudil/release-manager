package model

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
)

type ProjectService interface {
	Create(ctx context.Context, c svcmodel.ProjectCreation, authUserID uuid.UUID) (svcmodel.Project, error)
	Get(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) (svcmodel.Project, error)
	GetAll(ctx context.Context, authUserID uuid.UUID) ([]svcmodel.Project, error)
	Update(ctx context.Context, u svcmodel.ProjectUpdate, projectID, authUserID uuid.UUID) (svcmodel.Project, error)
	Delete(ctx context.Context, projectID uuid.UUID, authUserID uuid.UUID) error

	CreateEnvironment(ctx context.Context, c svcmodel.EnvironmentCreation, authUserID uuid.UUID) (svcmodel.Environment, error)
	GetEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) (svcmodel.Environment, error)
	GetEnvironments(ctx context.Context, projectID, authUserID uuid.UUID) ([]svcmodel.Environment, error)
	DeleteEnvironment(ctx context.Context, projectID, envID, authUserID uuid.UUID) error
	UpdateEnvironment(ctx context.Context, u svcmodel.EnvironmentUpdate, projectID, envID, authUserID uuid.UUID) (svcmodel.Environment, error)
}

type UserService interface {
	Get(ctx context.Context, id, authUserID uuid.UUID) (svcmodel.User, error)
	GetAll(ctx context.Context, authUserID uuid.UUID) ([]svcmodel.User, error)
	Delete(ctx context.Context, id, authUserID uuid.UUID) error
}

type AuthService interface {
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}

type SettingsService interface {
	Update(ctx context.Context, u svcmodel.UpdateSettingsInput, authUserID uuid.UUID) (svcmodel.Settings, error)
	Get(ctx context.Context, authUserID uuid.UUID) (svcmodel.Settings, error)
}
