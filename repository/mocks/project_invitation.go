package mocks

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProjectInvitationRepository struct {
	mock.Mock
}

func (m *MockProjectInvitationRepository) Insert(ctx context.Context, projectID uuid.UUID, email string, role model.ProjectRole, invitedByUserID uuid.UUID) (model.ProjectInvitation, error) {
	args := m.Called(ctx, projectID, email, role, invitedByUserID)
	return args.Get(0).(model.ProjectInvitation), args.Error(1)
}

func (m *MockProjectInvitationRepository) ReadByEmail(ctx context.Context, projectID uuid.UUID, email string) (model.ProjectInvitation, error) {
	args := m.Called(ctx, projectID, email)
	return args.Get(0).(model.ProjectInvitation), args.Error(1)
}

func (m *MockProjectInvitationRepository) Read(ctx context.Context, projectID, invitationID uuid.UUID) (model.ProjectInvitation, error) {
	args := m.Called(ctx, projectID, invitationID)
	return args.Get(0).(model.ProjectInvitation), args.Error(1)
}

func (m *MockProjectInvitationRepository) ReadAll(ctx context.Context, projectID uuid.UUID) ([]model.ProjectInvitation, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]model.ProjectInvitation), args.Error(1)
}

func (m *MockProjectInvitationRepository) Delete(ctx context.Context, projectID, invitationID uuid.UUID) error {
	args := m.Called(ctx, projectID, invitationID)
	return args.Error(0)
}
