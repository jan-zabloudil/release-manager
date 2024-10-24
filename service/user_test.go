package service

import (
	"context"
	"errors"
	"testing"

	"release-manager/pkg/id"
	repomock "release-manager/repository/mock"
	svcerrors "release-manager/service/errors"
	svcmock "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetAuthenticated(t *testing.T) {
	testCases := []struct {
		name        string
		setupMocks  func() *repomock.UserRepository
		expectedErr error
	}{
		{
			name: "Success",
			setupMocks: func() *repomock.UserRepository {
				repo := new(repomock.UserRepository)
				repo.On("Read", mock.Anything, mock.Anything).Return(model.User{}, nil)
				return repo
			},
			expectedErr: nil,
		},
		{
			name: "Unauthenticated user",
			setupMocks: func() *repomock.UserRepository {
				repo := new(repomock.UserRepository)
				repo.On("Read", mock.Anything, mock.Anything).Return(model.User{}, svcerrors.NewUserNotFoundError())
				return repo
			},
			expectedErr: svcerrors.NewUnauthenticatedUserError().Wrap(svcerrors.NewUserNotFoundError()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.setupMocks()
			userService := NewUserService(
				new(svcmock.AuthorizationService),
				repo,
			)

			_, err := userService.GetAuthenticated(context.Background(), id.AuthUser{})
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetForAdmin(t *testing.T) {
	testCases := []struct {
		name       string
		setupMocks func() (*svcmock.AuthorizationService, *repomock.UserRepository)
		expectErr  bool
	}{
		{
			name: "Success",
			setupMocks: func() (*svcmock.AuthorizationService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthorizationService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				mockUserRepo.On("Read", mock.Anything, mock.Anything).Return(model.User{}, nil)
				return mockAuthSvc, mockUserRepo
			},
			expectErr: false,
		},
		{
			name: "AuthorizationFailure",
			setupMocks: func() (*svcmock.AuthorizationService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthorizationService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(errors.New("test error"))
				return mockAuthSvc, mockUserRepo
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthSvc, mockUserRepo := tc.setupMocks()
			userService := NewUserService(mockAuthSvc, mockUserRepo)

			_, err := userService.GetForAdmin(context.Background(), id.User{}, id.AuthUser{})
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockAuthSvc.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_ListAllForAdmin(t *testing.T) {
	testCases := []struct {
		name       string
		setupMocks func() (*svcmock.AuthorizationService, *repomock.UserRepository)
		expectErr  bool
	}{
		{
			name: "Success",
			setupMocks: func() (*svcmock.AuthorizationService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthorizationService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				mockUserRepo.On("ListAll", mock.Anything).Return([]model.User{}, nil)
				return mockAuthSvc, mockUserRepo
			},
			expectErr: false,
		},
		{
			name: "AuthorizationFailure",
			setupMocks: func() (*svcmock.AuthorizationService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthorizationService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(errors.New("test error"))
				return mockAuthSvc, mockUserRepo
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthSvc, mockUserRepo := tc.setupMocks()
			userService := NewUserService(mockAuthSvc, mockUserRepo)

			_, err := userService.ListAllForAdmin(context.Background(), id.AuthUser{})
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockAuthSvc.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_DeleteForAdmin(t *testing.T) {
	testCases := []struct {
		name       string
		setupMocks func() (*svcmock.AuthorizationService, *repomock.UserRepository)
		expectErr  bool
	}{
		{
			name: "Success",
			setupMocks: func() (*svcmock.AuthorizationService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthorizationService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				mockUserRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.User{Role: model.UserRoleUser}, nil)
				mockUserRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)
				return mockAuthSvc, mockUserRepo
			},
			expectErr: false,
		},
		{
			name: "Cannot delete admin user",
			setupMocks: func() (*svcmock.AuthorizationService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthorizationService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				mockUserRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.User{Role: model.UserRoleAdmin}, nil)
				return mockAuthSvc, mockUserRepo
			},
			expectErr: true,
		},
		{
			name: "AuthorizationFailure",
			setupMocks: func() (*svcmock.AuthorizationService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthorizationService)
				mockUserRepo := new(repomock.UserRepository) // Even if not used, included for consistency
				mockAuthSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(errors.New("test error"))
				return mockAuthSvc, mockUserRepo
			},
			expectErr: true,
		},
		{
			name: "UserNotFound",
			setupMocks: func() (*svcmock.AuthorizationService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthorizationService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				mockUserRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.User{}, svcerrors.NewUserNotFoundError())
				return mockAuthSvc, mockUserRepo
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthSvc, mockUserRepo := tc.setupMocks()
			userService := NewUserService(mockAuthSvc, mockUserRepo)

			err := userService.DeleteForAdmin(context.Background(), id.User{}, id.AuthUser{})
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockAuthSvc.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
