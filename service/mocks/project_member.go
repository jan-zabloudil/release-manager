package mocks

import (
	"context"

	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProjectMemberService struct {
	mock.Mock
}

func (m *MockProjectMemberService) Get(ctx context.Context, projectID, userID uuid.UUID) (model.ProjectMember, error) {
	args := m.Called(ctx, projectID, userID)
	return args.Get(0).(model.ProjectMember), args.Error(1)
}

func (m *MockProjectMemberService) Create(ctx context.Context, projectID, userID uuid.UUID, role model.ProjectRole, invitedByUserID uuid.UUID) (model.ProjectMember, error) {
	args := m.Called(ctx, projectID, userID, role, invitedByUserID)
	return args.Get(0).(model.ProjectMember), args.Error(1)
}
