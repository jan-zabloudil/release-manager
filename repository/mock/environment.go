package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type EnvironmentRepository struct {
	mock.Mock
}

func (m *EnvironmentRepository) Create(ctx context.Context, e svcmodel.Environment) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *EnvironmentRepository) Read(ctx context.Context, envID uuid.UUID) (svcmodel.Environment, error) {
	args := m.Called(ctx, envID)
	return args.Get(0).(svcmodel.Environment), args.Error(1)
}

func (m *EnvironmentRepository) ReadByNameForProject(ctx context.Context, projectID uuid.UUID, name string) (svcmodel.Environment, error) {
	args := m.Called(ctx, projectID, name)
	return args.Get(0).(svcmodel.Environment), args.Error(1)
}

func (m *EnvironmentRepository) ReadAllForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.Environment, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.Environment), args.Error(1)
}

func (m *EnvironmentRepository) Delete(ctx context.Context, envID uuid.UUID) error {
	args := m.Called(ctx, envID)
	return args.Error(0)
}

func (m *EnvironmentRepository) Update(ctx context.Context, e svcmodel.Environment) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}
