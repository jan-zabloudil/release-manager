package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	svcmodel "release-manager/service/model"
)

type MockProjectMemberRepository struct {
	mock.Mock
}

func (m *MockProjectMemberRepository) Insert(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, role svcmodel.ProjectRole, invitedByUserID uuid.UUID) (svcmodel.ProjectMember, error) {
	args := m.Called(ctx, projectID, userID, role, invitedByUserID)
	return args.Get(0).(svcmodel.ProjectMember), args.Error(1)
}

func (m *MockProjectMemberRepository) Read(ctx context.Context, projectID, userID uuid.UUID) (svcmodel.ProjectMember, error) {
	args := m.Called(ctx, projectID, userID)
	return args.Get(0).(svcmodel.ProjectMember), args.Error(1)
}

func (m *MockProjectMemberRepository) Delete(ctx context.Context, projectID, userID uuid.UUID) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func (m *MockProjectMemberRepository) ReadAll(ctx context.Context, projectID uuid.UUID) ([]svcmodel.ProjectMember, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).([]svcmodel.ProjectMember), args.Error(1)
}

func (m *MockProjectMemberRepository) Update(ctx context.Context, member svcmodel.ProjectMember) (svcmodel.ProjectMember, error) {
	args := m.Called(ctx, member)
	return args.Get(0).(svcmodel.ProjectMember), args.Error(1)
}
