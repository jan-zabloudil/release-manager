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

func TestProjectInvitationService_Create(t *testing.T) {
	projectID := uuid.New()
	email := "test@example.com"
	role := model.ProjectRoleAdmin()
	invitedByUserID := uuid.New()

	testCases := []struct {
		name           string
		mockRepoReturn model.ProjectInvitation
		mockRepoErr    error
		expectedErr    error
	}{
		{
			name:           "Invitation successfully created",
			mockRepoReturn: model.ProjectInvitation{},
			mockRepoErr:    reperr.ErrResourceNotFound,
			expectedErr:    nil,
		},
		{
			name:           "Invitation already exists",
			mockRepoReturn: model.ProjectInvitation{},
			mockRepoErr:    nil,
			expectedErr:    svcerr.ErrInvitationAlreadyExists,
		},
		{
			name:           "Error reading from repository",
			mockRepoReturn: model.ProjectInvitation{},
			mockRepoErr:    errors.New("unexpected error"),
			expectedErr:    errors.New("unexpected error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockProjectInvitationRepository)
			mockRepo.On("ReadByEmail", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockRepoReturn, tc.mockRepoErr)

			if tc.expectedErr == nil {
				mockRepo.On("Insert", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, nil)
			}

			service := ProjectInvitationService{
				repository: mockRepo,
			}

			_, err := service.Create(context.Background(), projectID, email, role, invitedByUserID)

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
