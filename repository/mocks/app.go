package mocks

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAppRepository struct {
	mock.Mock
}

func (m *MockAppRepository) Insert(ctx context.Context, app svcmodel.App) (svcmodel.App, error) {
	args := m.Called(ctx, app)
	return args.Get(0).(svcmodel.App), args.Error(1)
}

func (m *MockAppRepository) Read(ctx context.Context, id uuid.UUID) (svcmodel.App, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(svcmodel.App), args.Error(1)
}

func (m *MockAppRepository) ReadAllForProject(ctx context.Context, projectID uuid.UUID) ([]svcmodel.App, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.App), args.Error(1)
}

func (m *MockAppRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAppRepository) Update(ctx context.Context, app svcmodel.App) (svcmodel.App, error) {
	args := m.Called(ctx, app)
	return args.Get(0).(svcmodel.App), args.Error(1)
}

func (m *MockAppRepository) InsertRepo(ctx context.Context, repo svcmodel.SCMRepo) (svcmodel.SCMRepo, error) {
	args := m.Called(ctx, repo)
	return args.Get(0).(svcmodel.SCMRepo), args.Error(1)
}

func (m *MockAppRepository) ReadRepo(ctx context.Context, appID uuid.UUID) (svcmodel.SCMRepo, error) {
	args := m.Called(ctx, appID)
	return args.Get(0).(svcmodel.SCMRepo), args.Error(1)
}

func (m *MockAppRepository) DeleteRepo(ctx context.Context, appID uuid.UUID) error {
	args := m.Called(ctx, appID)
	return args.Error(0)
}
