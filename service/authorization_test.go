package service

import (
	"context"
	"testing"

	repo "release-manager/repository/mock"
	svcerrors "release-manager/service/errors"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_AuthorizeUserRoleAdmin(t *testing.T) {
	adminUser := model.User{Role: model.UserRoleAdmin}
	user := model.User{Role: model.UserRoleUser}

	testCases := []struct {
		name      string
		mockSetup func(*repo.UserRepository)
		wantErr   bool
	}{
		{
			name: "User role admin",
			mockSetup: func(userRepo *repo.UserRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(adminUser, nil)
			},
			wantErr: false,
		},
		{
			name: "User role user",
			mockSetup: func(userRepo *repo.UserRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
			},
			wantErr: true,
		},
		{
			name: "User not found",
			mockSetup: func(userRepo *repo.UserRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(model.User{}, svcerrors.NewUserNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := new(repo.UserRepository)
			projectRepo := new(repo.ProjectRepository)
			service := NewAuthorizationService(userRepo, projectRepo)

			tc.mockSetup(userRepo)

			err := service.AuthorizeUserRoleAdmin(context.Background(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_AuthorizeUserRoleUser(t *testing.T) {
	adminUser := model.User{Role: model.UserRoleAdmin}
	user := model.User{Role: model.UserRoleUser}

	testCases := []struct {
		name      string
		mockSetup func(*repo.UserRepository)
		wantErr   bool
	}{
		{
			name: "User role admin",
			mockSetup: func(userRepo *repo.UserRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(adminUser, nil)
			},
			wantErr: false,
		},
		{
			name: "User role user",
			mockSetup: func(userRepo *repo.UserRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
			},
			wantErr: false,
		},
		{
			name: "User not found",
			mockSetup: func(userRepo *repo.UserRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(model.User{}, svcerrors.NewUserNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := new(repo.UserRepository)
			projectRepo := new(repo.ProjectRepository)
			service := NewAuthorizationService(userRepo, projectRepo)

			tc.mockSetup(userRepo)

			err := service.AuthorizeUserRoleUser(context.Background(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
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
		mockSetup func(*repo.UserRepository, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Project member editor",
			mockSetup: func(userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(editorProjectMember, nil)
			},
			wantErr: false,
		},
		{
			name: "Project member viewer",
			mockSetup: func(userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(viewerProjectMember, nil)
			},
			wantErr: true,
		},
		{
			name: "User not project member",
			mockSetup: func(userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectMember{}, svcerrors.NewProjectMemberNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "User admin",
			mockSetup: func(userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(adminUser, nil)
			},
			wantErr: false,
		},
		{
			name: "User not found",
			mockSetup: func(userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(model.User{}, svcerrors.NewUserNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := new(repo.UserRepository)
			projectRepo := new(repo.ProjectRepository)
			service := NewAuthorizationService(userRepo, projectRepo)

			tc.mockSetup(userRepo, projectRepo)

			err := service.AuthorizeProjectRoleEditor(context.Background(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
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
		mockSetup func(*repo.UserRepository, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Project member editor",
			mockSetup: func(userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(editorProjectMember, nil)
			},
			wantErr: false,
		},
		{
			name: "Project member viewer",
			mockSetup: func(userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(viewerProjectMember, nil)
			},
			wantErr: false,
		},
		{
			name: "User not project member",
			mockSetup: func(userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(user, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectMember{}, svcerrors.NewProjectMemberNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "User admin",
			mockSetup: func(userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(adminUser, nil)
			},
			wantErr: false,
		},
		{
			name: "User not found",
			mockSetup: func(userRepo *repo.UserRepository, projectRepo *repo.ProjectRepository) {
				userRepo.On("Read", mock.Anything, mock.Anything).Return(model.User{}, svcerrors.NewUserNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := new(repo.UserRepository)
			projectRepo := new(repo.ProjectRepository)
			service := NewAuthorizationService(userRepo, projectRepo)

			tc.mockSetup(userRepo, projectRepo)

			err := service.AuthorizeProjectRoleViewer(context.Background(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}
