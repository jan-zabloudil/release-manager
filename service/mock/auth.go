package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type AuthService struct {
	mock.Mock
}

func (m *AuthService) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *AuthService) AuthorizeUserRoleAdmin(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *AuthService) AuthorizeUserRole(ctx context.Context, userID uuid.UUID, role model.UserRole) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *AuthService) AuthorizeProjectRoleViewer(ctx context.Context, projectID, userID uuid.UUID) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *AuthService) AuthorizeProjectRoleEditor(ctx context.Context, projectID, userID uuid.UUID) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}
