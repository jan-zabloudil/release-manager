package mock

import (
	"context"

	"release-manager/pkg/id"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type AuthorizationService struct {
	mock.Mock
}

func (m *AuthorizationService) AuthorizeUserRoleAdmin(ctx context.Context, userID id.AuthUser) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *AuthorizationService) AuthorizeUserRoleUser(ctx context.Context, userID id.AuthUser) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *AuthorizationService) AuthorizeProjectRoleViewer(ctx context.Context, projectID uuid.UUID, userID id.AuthUser) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *AuthorizationService) AuthorizeProjectRoleEditor(ctx context.Context, projectID uuid.UUID, userID id.AuthUser) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *AuthorizationService) AuthorizeReleaseViewer(ctx context.Context, releaseID uuid.UUID, userID id.AuthUser) error {
	args := m.Called(ctx, releaseID, userID)
	return args.Error(0)
}

func (m *AuthorizationService) AuthorizeReleaseEditor(ctx context.Context, releaseID uuid.UUID, userID id.AuthUser) error {
	args := m.Called(ctx, releaseID, userID)
	return args.Error(0)
}
