package service

import (
	"context"
	"testing"

	repomock "release-manager/repository/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuth_AuthorizeRole(t *testing.T) {
	mockAuthRepo := new(repomock.AuthRepository)
	mockUserRepo := new(repomock.UserRepository)
	authService := NewAuthService(mockAuthRepo, mockUserRepo)

	adminRole := model.UserRoleAdmin
	userRole := model.UserRoleUser

	testCases := []struct {
		name        string
		userID      uuid.UUID
		role        model.UserRole
		requireRole model.UserRole
		expectErr   bool
	}{
		{
			name:        "Admin - user role required",
			userID:      uuid.New(),
			role:        adminRole,
			requireRole: userRole,
			expectErr:   false,
		},
		{
			name:        "Admin - admin role required",
			userID:      uuid.New(),
			role:        adminRole,
			requireRole: adminRole,
			expectErr:   false,
		},
		{
			name:        "User - admin role required",
			userID:      uuid.New(),
			role:        userRole,
			requireRole: adminRole,
			expectErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUser := model.User{Role: tc.role}
			mockUserRepo.On("Read", mock.Anything, tc.userID).Return(mockUser, nil)

			err := authService.AuthorizeRole(context.Background(), tc.userID, tc.requireRole)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuth_AuthorizeAdminRole(t *testing.T) {
	mockAuthRepo := new(repomock.AuthRepository)
	mockUserRepo := new(repomock.UserRepository)
	authService := NewAuthService(mockAuthRepo, mockUserRepo)

	testCases := []struct {
		name      string
		userID    uuid.UUID
		userRole  model.UserRole
		expectErr bool
	}{
		{
			name:      "Admin role success",
			userID:    uuid.New(),
			userRole:  model.UserRoleAdmin,
			expectErr: false,
		},
		{
			name:      "User role denied",
			userID:    uuid.New(),
			userRole:  model.UserRoleUser,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUser := model.User{Role: tc.userRole}
			mockUserRepo.On("Read", mock.Anything, tc.userID).Return(mockUser, nil)

			err := authService.AuthorizeAdminRole(context.Background(), tc.userID)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}
