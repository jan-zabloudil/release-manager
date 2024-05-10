package mock

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type AuthorizeService struct {
	mock.Mock
}

func (m *AuthorizeService) AuthorizeUserRoleAdmin(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *AuthorizeService) AuthorizeUserRole(ctx context.Context, userID uuid.UUID, role model.UserRole) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *AuthorizeService) AuthorizeProjectRoleViewer(ctx context.Context, projectID, userID uuid.UUID) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *AuthorizeService) AuthorizeProjectRoleEditor(ctx context.Context, projectID, userID uuid.UUID) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}
