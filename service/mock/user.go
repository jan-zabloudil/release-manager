package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserService struct {
	mock.Mock
}

func (m *UserService) Get(ctx context.Context, id uuid.UUID, authUserID uuid.UUID) (model.User, error) {
	args := m.Called(ctx, id, authUserID)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *UserService) GetByEmail(ctx context.Context, email string) (model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}
