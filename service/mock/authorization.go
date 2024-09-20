package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type AuthorizationService struct {
	mock.Mock
}

func (m *AuthorizationService) AuthorizeUserRoleAdmin(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *AuthorizationService) AuthorizeUserRoleUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *AuthorizationService) AuthorizeProjectRoleViewer(ctx context.Context, projectID, userID uuid.UUID) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *AuthorizationService) AuthorizeProjectRoleEditor(ctx context.Context, projectID, userID uuid.UUID) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *AuthorizationService) GetAuthorizedUser(ctx context.Context, userID uuid.UUID) (model.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(model.User), args.Error(1)
}
