package mock

import (
	"context"

	"release-manager/pkg/id"
	"release-manager/service/model"

	"github.com/stretchr/testify/mock"
)

type UserService struct {
	mock.Mock
}

func (m *UserService) GetAuthenticated(ctx context.Context, userID id.AuthUser) (model.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserService) GetByEmail(ctx context.Context, email string) (model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}
