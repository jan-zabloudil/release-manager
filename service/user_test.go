package service

import (
	"context"
	"errors"
	"testing"

	repomock "release-manager/repository/mock"
	svcerrors "release-manager/service/errors"
	svcmock "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Get(t *testing.T) {
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

			_, err := userService.GetForAdmin(context.Background(), uuid.New(), uuid.New())
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

func TestUserService_GetAll(t *testing.T) {
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

			_, err := userService.ListAllForAdmin(context.Background(), uuid.New())
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

func TestUserService_Delete(t *testing.T) {
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

			err := userService.DeleteForAdmin(context.Background(), uuid.New(), uuid.New())
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
