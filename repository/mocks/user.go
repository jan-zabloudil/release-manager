package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	svcmodel "release-manager/service/model"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) ReadForToken(ctx context.Context, token string) (svcmodel.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(svcmodel.User), args.Error(1)
}

func (m *MockUserRepository) Read(ctx context.Context, id uuid.UUID) (svcmodel.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(svcmodel.User), args.Error(1)
}

func (m *MockUserRepository) ReadByEmail(ctx context.Context, email string) (svcmodel.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(svcmodel.User), args.Error(1)
}
