package service

import (
	"context"
	"errors"
	"testing"

	reperr "release-manager/repository/errors"
	"release-manager/repository/mocks"
	svcerr "release-manager/service/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"release-manager/service/model"
)

func TestProjectMemberService_Create(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()
	role := model.ProjectRoleAdmin()
	invitedByUserID := uuid.New()

	testCases := []struct {
		name           string
		mockRepoReturn model.ProjectMember
		mockRepoErr    error
		expectedErr    error
	}{
		{
			name:           "Member successfully created",
			mockRepoReturn: model.ProjectMember{},
			mockRepoErr:    reperr.ErrResourceNotFound,
			expectedErr:    nil,
		},
		{
			name:           "User is already a member",
			mockRepoReturn: model.ProjectMember{},
			mockRepoErr:    nil,
			expectedErr:    svcerr.ErrUserIsAlreadyMember,
		},
		{
			name:           "Error reading from repository",
			mockRepoReturn: model.ProjectMember{},
			mockRepoErr:    errors.New("unexpected error"),
			expectedErr:    errors.New("unexpected error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockProjectMemberRepository)
			mockRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockRepoReturn, tc.mockRepoErr)

			if tc.expectedErr == nil {
				mockRepo.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectMember{}, nil)
			}

			service := ProjectMemberService{
				repository: mockRepo,
			}

			_, err := service.Create(context.Background(), projectID, userID, role, invitedByUserID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectMemberService_UpdateRole(t *testing.T) {
	testCases := []struct {
		name           string
		member         model.ProjectMember
		updatedBy      model.ProjectMember
		newRole        model.ProjectRole
		mockRepoReturn model.ProjectMember
		mockRepoErr    error
		expectedErr    error
	}{
		{
			name:           "Role successfully updated",
			member:         model.ProjectMember{Role: model.ProjectRoleEditor()},
			updatedBy:      model.ProjectMember{Role: model.ProjectRoleAdmin()},
			newRole:        model.ProjectRoleViewer(),
			mockRepoReturn: model.ProjectMember{},
			mockRepoErr:    nil,
			expectedErr:    nil,
		},
		{
			name:           "Updater does not have permission to update the member",
			member:         model.ProjectMember{Role: model.ProjectRoleAdmin()},
			updatedBy:      model.ProjectMember{Role: model.ProjectRoleEditor()},
			newRole:        model.ProjectRoleViewer(),
			mockRepoReturn: model.ProjectMember{},
			mockRepoErr:    nil,
			expectedErr:    svcerr.ErrProjectMemberUpdateNotAllowed,
		},
		{
			name:           "Updater does not have permission to grant the new role",
			member:         model.ProjectMember{Role: model.ProjectRoleEditor()},
			updatedBy:      model.ProjectMember{Role: model.ProjectRoleAdmin()},
			newRole:        model.ProjectRoleAdmin(),
			mockRepoReturn: model.ProjectMember{},
			mockRepoErr:    nil,
			expectedErr:    svcerr.ErrProjectMemberRoleCannotBeGranted,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockProjectMemberRepository)

			if tc.expectedErr == nil {
				mockRepo.On("Update", mock.Anything, mock.Anything).Return(tc.mockRepoReturn, nil)
			}

			service := ProjectMemberService{
				repository: mockRepo,
			}

			_, err := service.UpdateRole(context.Background(), tc.member, tc.updatedBy, tc.newRole)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
