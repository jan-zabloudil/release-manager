package service

import (
	"context"
	"testing"

	"release-manager/pkg/dberrors"
	repo "release-manager/repository/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuth_AuthorizeRole(t *testing.T) {
	mockAuthRepo := new(repo.AuthRepository)
	mockUserRepo := new(repo.UserRepository)
	mockProjectRepo := new(repo.ProjectRepository)
	authService := NewAuthService(mockAuthRepo, mockUserRepo, mockProjectRepo)

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

			err := authService.AuthorizeUserRole(context.Background(), tc.userID, tc.requireRole)

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
	mockAuthRepo := new(repo.AuthRepository)
	mockUserRepo := new(repo.UserRepository)
	mockProjectRepo := new(repo.ProjectRepository)
	authService := NewAuthService(mockAuthRepo, mockUserRepo, mockProjectRepo)

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

			err := authService.AuthorizeUserRoleAdmin(context.Background(), tc.userID)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuth_AuthorizeProjectRoleEditor(t *testing.T) {
	adminUser := model.User{Role: model.UserRoleAdmin}
	user := model.User{Role: model.UserRoleUser}
	editorProjectMember := model.ProjectMember{ProjectRole: model.ProjectRoleEditor}
	viewerProjectMember := model.ProjectMember{ProjectRole: model.ProjectRoleViewer}

	testCases := []struct {
		name      string
		mockSetup func(*repo.AuthRepository, *repo.UserRepository, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Project member editor",
			mockSetup: func(authRepo *repo.AuthRepository, userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(editorProjectMember, nil)
			},
			wantErr: false,
		},
		{
			name: "Project member viewer",
			mockSetup: func(authRepo *repo.AuthRepository, userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(viewerProjectMember, nil)
			},
			wantErr: true,
		},
		{
			name: "User not project member",
			mockSetup: func(authRepo *repo.AuthRepository, userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectMember{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "User admin",
			mockSetup: func(authRepo *repo.AuthRepository, userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(adminUser, nil)
			},
			wantErr: false,
		},
		{
			name: "User not found",
			mockSetup: func(authRepo *repo.AuthRepository, userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(model.User{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authRepo := new(repo.AuthRepository)
			userRepo := new(repo.UserRepository)
			projectRepo := new(repo.ProjectRepository)
			service := NewAuthService(authRepo, userRepo, projectRepo)

			tc.mockSetup(authRepo, userRepo, projectRepo)

			err := service.AuthorizeProjectRoleEditor(context.Background(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			authRepo.AssertExpectations(t)
		})
	}
}

func TestAuth_AuthorizeProjectRoleViewer(t *testing.T) {
	adminUser := model.User{Role: model.UserRoleAdmin}
	user := model.User{Role: model.UserRoleUser}
	editorProjectMember := model.ProjectMember{ProjectRole: model.ProjectRoleEditor}
	viewerProjectMember := model.ProjectMember{ProjectRole: model.ProjectRoleViewer}

	testCases := []struct {
		name      string
		mockSetup func(*repo.AuthRepository, *repo.UserRepository, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Project member editor",
			mockSetup: func(authRepo *repo.AuthRepository, userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(editorProjectMember, nil)
			},
			wantErr: false,
		},
		{
			name: "Project member viewer",
			mockSetup: func(authRepo *repo.AuthRepository, userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(viewerProjectMember, nil)
			},
			wantErr: false,
		},
		{
			name: "User not project member",
			mockSetup: func(authRepo *repo.AuthRepository, userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectMember{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "User admin",
			mockSetup: func(authRepo *repo.AuthRepository, userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(adminUser, nil)
			},
			wantErr: false,
		},
		{
			name: "User not found",
			mockSetup: func(authRepo *repo.AuthRepository, userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(model.User{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authRepo := new(repo.AuthRepository)
			userRepo := new(repo.UserRepository)
			projectRepo := new(repo.ProjectRepository)
			service := NewAuthService(authRepo, userRepo, projectRepo)

			tc.mockSetup(authRepo, userRepo, projectRepo)

			err := service.AuthorizeProjectRoleViewer(context.Background(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			authRepo.AssertExpectations(t)
		})
	}
}
