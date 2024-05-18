package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ReleaseRepository struct {
	mock.Mock
}

func (m *ReleaseRepository) Create(ctx context.Context, rls svcmodel.Release) error {
	args := m.Called(ctx, rls)
	return args.Error(0)
}

func (m *ReleaseRepository) Read(ctx context.Context, projectID, releaseID uuid.UUID) (svcmodel.Release, error) {
	args := m.Called(ctx, projectID, releaseID)
	return args.Get(0).(svcmodel.Release), args.Error(1)
}

func (m *ReleaseRepository) Delete(ctx context.Context, projectID, releaseID uuid.UUID) error {
	args := m.Called(ctx, projectID, releaseID)
	return args.Error(0)
}
