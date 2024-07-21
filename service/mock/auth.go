package mock

import (
	"context"

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

func (m *AuthorizeService) AuthorizeUserRoleUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
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
