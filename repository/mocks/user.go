package mocks

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Read(ctx context.Context, userID uuid.UUID) (svcmodel.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(svcmodel.User), args.Error(1)
}

func (m *MockUserRepository) ReadAll(ctx context.Context) ([]svcmodel.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]svcmodel.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
