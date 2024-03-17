package mocks

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProjectInvitationService struct {
	mock.Mock
}

func (m *MockProjectInvitationService) GetByEmail(ctx context.Context, projectID uuid.UUID, email string) (model.ProjectInvitation, error) {
	args := m.Called(ctx, projectID, email)
	return args.Get(0).(model.ProjectInvitation), args.Error(1)
}

func (m *MockProjectInvitationService) Create(ctx context.Context, projectID uuid.UUID, email string, role model.ProjectRole, invitedByUserID uuid.UUID) (model.ProjectInvitation, error) {
	args := m.Called(ctx, projectID, email, role, invitedByUserID)
	return args.Get(0).(model.ProjectInvitation), args.Error(1)
}
