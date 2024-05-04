package mock

import (
	"context"

	svcmodel "release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Read(ctx context.Context, userID uuid.UUID) (svcmodel.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(svcmodel.User), args.Error(1)
}

func (m *UserRepository) ReadByEmail(ctx context.Context, email string) (svcmodel.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(svcmodel.User), args.Error(1)
}

func (m *UserRepository) ReadAll(ctx context.Context) ([]svcmodel.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]svcmodel.User), args.Error(1)
}

func (m *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
