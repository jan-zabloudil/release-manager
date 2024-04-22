package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ProjectRepository struct {
	mock.Mock
}

func (m *ProjectRepository) Create(ctx context.Context, p model.Project) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *ProjectRepository) Read(ctx context.Context, id uuid.UUID) (model.Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Project), args.Error(1)
}

func (m *ProjectRepository) ReadAll(ctx context.Context) ([]model.Project, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Project), args.Error(1)
}

func (m *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProjectRepository) Update(ctx context.Context, p model.Project) error {
	args := m.Called(ctx, p)
	return args.Error(1)
}
