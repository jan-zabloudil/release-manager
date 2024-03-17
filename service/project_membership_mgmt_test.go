package service

import (
	"context"
	"testing"

	reperr "release-manager/repository/errors"
	svcerr "release-manager/service/errors"
	"release-manager/service/mocks"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testCase struct {
	name                  string
	membershipRequestRole model.ProjectRole
	requestedBy           model.ProjectMember
	mockUserReturn        model.User
	mockUserErr           error
	expectedErr           error
	expectedResponse      model.ProjectMembershipResponse
}

func TestProjectMembershipManagementService_Create(t *testing.T) {
	projectID := uuid.New()
	email := "test@example.com"
	requestedByUserID := uuid.New()

	testCases := []testCase{
		{
			name:                  "User already exists, member successfully created",
			membershipRequestRole: model.ProjectRoleEditor(),
			requestedBy:           model.ProjectMember{Role: model.ProjectRoleAdmin()},
			mockUserReturn:        model.User{ID: uuid.New()},
			mockUserErr:           nil,
			expectedResponse: model.ProjectMembershipResponse{
				Status:   model.MemberCreatedStatus,
				Resource: model.ProjectMember{},
			},
		},
		{
			name:                  "User does not exist, invitation successfully sent",
			membershipRequestRole: model.ProjectRoleEditor(),
			requestedBy:           model.ProjectMember{Role: model.ProjectRoleAdmin()},
			mockUserReturn:        model.User{},
			mockUserErr:           reperr.ErrResourceNotFound,
			expectedResponse: model.ProjectMembershipResponse{
				Status:   model.InvitationSentStatus,
				Resource: model.ProjectInvitation{},
			},
		},
		{
			name:                  "Requested role cannot be granted",
			membershipRequestRole: model.ProjectRoleViewer(),
			requestedBy:           model.ProjectMember{Role: model.ProjectRoleViewer()},
			expectedErr:           svcerr.ErrProjectMemberRoleCannotBeGranted,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserSvc, mockMemberSvc, mockInvitationSvc := setupMocks(tc, projectID, requestedByUserID, email)

			service := ProjectMembershipManagementService{
				userSvc:       mockUserSvc,
				memberSvc:     mockMemberSvc,
				invitationSvc: mockInvitationSvc,
			}

			response, err := service.Create(
				context.Background(),
				model.ProjectMembershipRequest{
					ProjectID:         projectID,
					Email:             email,
					Role:              tc.membershipRequestRole,
					RequestedByUserID: requestedByUserID,
				},
				tc.requestedBy,
			)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)
			}

			mockUserSvc.AssertExpectations(t)
			mockMemberSvc.AssertExpectations(t)
			mockInvitationSvc.AssertExpectations(t)
		})
	}
}

func setupMocks(tc testCase, projectID, requestedByUserID uuid.UUID, email string) (*mocks.MockUserService, *mocks.MockProjectMemberService, *mocks.MockProjectInvitationService) {
	mockUserSvc := new(mocks.MockUserService)
	mockMemberSvc := new(mocks.MockProjectMemberService)
	mockInvitationSvc := new(mocks.MockProjectInvitationService)

	switch tc.name {
	case "User already exists, member successfully created":
		mockUserSvc.On("GetByEmail", mock.Anything, email).Return(tc.mockUserReturn, tc.mockUserErr)
		mockMemberSvc.On("Create", mock.Anything, projectID, tc.mockUserReturn.ID, tc.membershipRequestRole, requestedByUserID).Return(model.ProjectMember{}, tc.expectedErr)
	case "User does not exist, invitation successfully sent":
		mockUserSvc.On("GetByEmail", mock.Anything, email).Return(tc.mockUserReturn, tc.mockUserErr)
		mockInvitationSvc.On("Create", mock.Anything, projectID, email, tc.membershipRequestRole, requestedByUserID).Return(model.ProjectInvitation{}, tc.expectedErr)
	case "Requested role cannot be granted":
		// No mocks to set up for this case
	}

	return mockUserSvc, mockMemberSvc, mockInvitationSvc
}
