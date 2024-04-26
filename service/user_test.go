package service

import (
	"context"
	"errors"
	"testing"

	"release-manager/pkg/dberrors"
	repomock "release-manager/repository/mock"
	svcmock "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Get(t *testing.T) {
	authUserID := uuid.New()
	testUserID := uuid.New()
	testUser := model.User{ID: testUserID}

	testCases := []struct {
		name       string
		setupMocks func() (*svcmock.AuthService, *repomock.UserRepository)
		expectErr  bool
	}{
		{
			name: "Success",
			setupMocks: func() (*svcmock.AuthService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeAdminRole", mock.Anything, authUserID).Return(nil)
				mockUserRepo.On("Read", mock.Anything, testUserID).Return(testUser, nil)
				return mockAuthSvc, mockUserRepo
			},
			expectErr: false,
		},
		{
			name: "AuthorizationFailure",
			setupMocks: func() (*svcmock.AuthService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeAdminRole", mock.Anything, authUserID).Return(errors.New("test error"))
				return mockAuthSvc, mockUserRepo
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthSvc, mockUserRepo := tc.setupMocks()
			userService := NewUserService(mockAuthSvc, mockUserRepo)

			_, err := userService.Get(context.Background(), testUserID, authUserID)
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
	authUserID := uuid.New()
	users := []model.User{{ID: uuid.New()}, {ID: uuid.New()}}

	testCases := []struct {
		name       string
		setupMocks func() (*svcmock.AuthService, *repomock.UserRepository)
		expectErr  bool
	}{
		{
			name: "Success",
			setupMocks: func() (*svcmock.AuthService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeAdminRole", mock.Anything, authUserID).Return(nil)
				mockUserRepo.On("ReadAll", mock.Anything).Return(users, nil)
				return mockAuthSvc, mockUserRepo
			},
			expectErr: false,
		},
		{
			name: "AuthorizationFailure",
			setupMocks: func() (*svcmock.AuthService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeAdminRole", mock.Anything, authUserID).Return(errors.New("test error"))
				return mockAuthSvc, mockUserRepo
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthSvc, mockUserRepo := tc.setupMocks()
			userService := NewUserService(mockAuthSvc, mockUserRepo)

			_, err := userService.ListAll(context.Background(), authUserID)
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
	testUserID := uuid.New()
	authUserID := uuid.New()

	testCases := []struct {
		name       string
		setupMocks func() (*svcmock.AuthService, *repomock.UserRepository)
		expectErr  bool
	}{
		{
			name: "Success",
			setupMocks: func() (*svcmock.AuthService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeAdminRole", mock.Anything, authUserID).Return(nil)
				mockUserRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.User{}, nil)
				mockUserRepo.On("Delete", mock.Anything, testUserID).Return(nil)
				return mockAuthSvc, mockUserRepo
			},
			expectErr: false,
		},
		{
			name: "AuthorizationFailure",
			setupMocks: func() (*svcmock.AuthService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthService)
				mockUserRepo := new(repomock.UserRepository) // Even if not used, included for consistency
				mockAuthSvc.On("AuthorizeAdminRole", mock.Anything, authUserID).Return(errors.New("test error"))
				return mockAuthSvc, mockUserRepo
			},
			expectErr: true,
		},
		{
			name: "UserNotFound",
			setupMocks: func() (*svcmock.AuthService, *repomock.UserRepository) {
				mockAuthSvc := new(svcmock.AuthService)
				mockUserRepo := new(repomock.UserRepository)
				mockAuthSvc.On("AuthorizeAdminRole", mock.Anything, authUserID).Return(nil)
				mockUserRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.User{}, dberrors.NewNotFoundError())
				return mockAuthSvc, mockUserRepo
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthSvc, mockUserRepo := tc.setupMocks()
			userService := NewUserService(mockAuthSvc, mockUserRepo)

			err := userService.Delete(context.Background(), testUserID, authUserID)
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
