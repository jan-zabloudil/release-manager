package mock

import (
	"context"

	"release-manager/pkg/id"
	svcmodel "release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Read(ctx context.Context, userID id.User) (svcmodel.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(svcmodel.User), args.Error(1)
}

func (m *UserRepository) ReadByEmail(ctx context.Context, email string) (svcmodel.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(svcmodel.User), args.Error(1)
}

func (m *UserRepository) ListAll(ctx context.Context) ([]svcmodel.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]svcmodel.User), args.Error(1)
}

func (m *UserRepository) Delete(ctx context.Context, id id.User) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
