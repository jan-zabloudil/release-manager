package service

import (
	"context"
	"errors"
	"testing"

	githubmock "release-manager/github/mock"
	"release-manager/pkg/id"
	"release-manager/pkg/pointer"
	repo "release-manager/repository/mock"
	resendmock "release-manager/resend/mock"
	svcerrors "release-manager/service/errors"
	svc "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProjectService_CreateProject(t *testing.T) {
	testCases := []struct {
		name      string
		project   model.CreateProjectInput
		mockSetup func(*svc.AuthorizationService, *svc.SettingsService, *svc.UserService, *githubmock.Client, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Valid project with default config",
			project: model.CreateProjectInput{
				Name:                      "Test projectGetter",
				SlackChannelID:            "c1234",
				ReleaseNotificationConfig: model.ReleaseNotificationConfig{},
			},
			mockSetup: func(auth *svc.AuthorizationService, settings *svc.SettingsService, user *svc.UserService, github *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				settings.On("GetDefaultReleaseMessage", mock.Anything).Return("message", nil)
				user.On("GetAuthenticated", mock.Anything, mock.Anything).Return(model.User{}, nil)
				projectRepo.On("CreateProjectWithOwner", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Valid project with custom config",
			project: model.CreateProjectInput{
				Name:                      "Test projectGetter",
				SlackChannelID:            "c1234",
				ReleaseNotificationConfig: model.ReleaseNotificationConfig{Message: "test message"},
			},
			mockSetup: func(auth *svc.AuthorizationService, settings *svc.SettingsService, user *svc.UserService, github *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				user.On("GetAuthenticated", mock.Anything, mock.Anything).Return(model.User{}, nil)
				projectRepo.On("CreateProjectWithOwner", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid config",
			project: model.CreateProjectInput{
				Name:                      "Test projectGetter",
				SlackChannelID:            "c1234",
				ReleaseNotificationConfig: model.ReleaseNotificationConfig{Message: "", ShowProjectName: true},
			},
			mockSetup: func(auth *svc.AuthorizationService, settings *svc.SettingsService, user *svc.UserService, github *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid project - missing name",
			project: model.CreateProjectInput{
				Name:                      "",
				SlackChannelID:            "",
				ReleaseNotificationConfig: model.ReleaseNotificationConfig{Message: "test message"},
			},
			mockSetup: func(auth *svc.AuthorizationService, settings *svc.SettingsService, user *svc.UserService, github *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, settingsSvc, userSvc, github, projectRepo)

			_, err := service.CreateProject(context.Background(), tc.project, id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			settingsSvc.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
			userSvc.AssertExpectations(t)
			authSvc.AssertExpectations(t)
			github.AssertExpectations(t)
		})
	}
}

func TestProjectService_ListProjects(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.UserService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Non admin user",
			mockSetup: func(user *svc.UserService, projectRepo *repo.ProjectRepository) {
				user.On("GetAuthenticated", mock.Anything, mock.Anything).Return(model.User{
					Role: model.UserRoleUser,
				}, nil)
				projectRepo.On("ListProjectsForUser", mock.Anything, mock.Anything).Return([]model.Project{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Admin user",
			mockSetup: func(user *svc.UserService, projectRepo *repo.ProjectRepository) {
				user.On("GetAuthenticated", mock.Anything, mock.Anything).Return(model.User{
					Role: model.UserRoleAdmin,
				}, nil)
				projectRepo.On("ListProjects", mock.Anything).Return([]model.Project{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Unauthenticated user",
			mockSetup: func(user *svc.UserService, projectRepo *repo.ProjectRepository) {
				user.On("GetAuthenticated", mock.Anything, mock.Anything).Return(model.User{}, svcerrors.NewUnauthenticatedUserError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(userSvc, projectRepo)

			_, err := service.ListProjects(context.Background(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
			userSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetProject(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Existing project",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Non-existing project",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.GetProject(context.Background(), id.NewProject(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_DeleteProject(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Existing project",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("DeleteProject", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Non-existing project",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("DeleteProject", mock.Anything, mock.Anything).Return(svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.DeleteProject(context.Background(), id.NewProject(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_UpdateProject(t *testing.T) {
	testCases := []struct {
		name      string
		update    model.UpdateProjectInput
		mockSetup func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Valid project update",
			update: model.UpdateProjectInput{
				Name:           pointer.StringPtr("new name"),
				SlackChannelID: pointer.StringPtr("new channelID"),
				ReleaseNotificationConfigUpdate: model.UpdateReleaseNotificationConfigInput{
					Message:            new(string),
					ShowProjectName:    new(bool),
					ShowReleaseTitle:   new(bool),
					ShowReleaseNotes:   new(bool),
					ShowLastDeployment: new(bool),
					ShowSourceCode:     new(bool),
				},
			},
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("UpdateProject", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid project update",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("UpdateProject", mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewProjectInvalidError())
			},
			wantErr: true,
		},
		{
			name:   "Non-existing-project",
			update: model.UpdateProjectInput{},
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("UpdateProject", mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.UpdateProject(context.Background(), tc.update, id.NewProject(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_CreateEnvironment(t *testing.T) {
	testCases := []struct {
		name      string
		envCreate model.CreateEnvironmentInput
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Valid environment creation",
			envCreate: model.CreateEnvironmentInput{
				ProjectID:     id.NewProject(),
				Name:          "Test Environment",
				ServiceRawURL: "http://example.com",
			},
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("CreateEnvironment", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid environment creation - not absolute service url",
			envCreate: model.CreateEnvironmentInput{
				ProjectID:     id.NewProject(),
				Name:          "Test Environment",
				ServiceRawURL: "example",
			},
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid environment creation - empty name",
			envCreate: model.CreateEnvironmentInput{
				ProjectID:     id.NewProject(),
				Name:          "",
				ServiceRawURL: "example",
			},
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.CreateEnvironment(context.Background(), tc.envCreate, id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetEnvironment(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: false,
		},
		{
			name: "env not found",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, svcerrors.NewEnvironmentNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.GetEnvironment(context.Background(), id.NewProject(), id.NewEnvironment(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_UpdateEnvironment(t *testing.T) {
	testCases := []struct {
		name      string
		envUpdate model.UpdateEnvironmentInput
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			envUpdate: model.UpdateEnvironmentInput{
				Name:          pointer.StringPtr("New name"),
				ServiceRawURL: pointer.StringPtr("https://new.example.com"),
			},
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("UpdateEnvironment", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Unknown environment",
			envUpdate: model.UpdateEnvironmentInput{},
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("UpdateEnvironment", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewEnvironmentNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.UpdateEnvironment(context.Background(), tc.envUpdate, id.NewProject(), id.NewEnvironment(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetEnvironments(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ListEnvironmentsForProject", mock.Anything, mock.Anything).Return([]model.Environment{
					{ID: id.NewEnvironment()},
					{ID: id.NewEnvironment()},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "no environments",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ListEnvironmentsForProject", mock.Anything, mock.Anything).Return([]model.Environment{}, nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.ListEnvironments(context.Background(), id.NewProject(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_DeleteEnvironment(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("DeleteEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "env not found",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("DeleteEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewEnvironmentNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.DeleteEnvironment(context.Background(), id.NewProject(), id.NewEnvironment(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_Invite(t *testing.T) {
	testCases := []struct {
		name      string
		creation  model.CreateProjectInvitationInput
		mockSetup func(*svc.AuthorizationService, *svc.UserService, *resendmock.Client, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name:     "Unknown project",
			creation: model.CreateProjectInvitationInput{},
			mockSetup: func(auth *svc.AuthorizationService, user *svc.UserService, email *resendmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Invalid project invitation - missing email",
			creation: model.CreateProjectInvitationInput{
				Email:       "",
				ProjectRole: "editor",
				ProjectID:   id.NewProject(),
			},
			mockSetup: func(auth *svc.AuthorizationService, user *svc.UserService, email *resendmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Invalid project invitation - invalid role",
			creation: model.CreateProjectInvitationInput{
				Email:       "",
				ProjectRole: "new",
				ProjectID:   id.NewProject(),
			},
			mockSetup: func(auth *svc.AuthorizationService, user *svc.UserService, email *resendmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Member already exists",
			creation: model.CreateProjectInvitationInput{
				Email:       "test@test.tt",
				ProjectRole: "viewer",
				ProjectID:   id.NewProject(),
			},
			mockSetup: func(auth *svc.AuthorizationService, user *svc.UserService, email *resendmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadMemberByEmail", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectMember{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invitation already exists",
			creation: model.CreateProjectInvitationInput{
				Email:       "test@test.tt",
				ProjectRole: "viewer",
				ProjectID:   id.NewProject(),
			},
			mockSetup: func(auth *svc.AuthorizationService, user *svc.UserService, email *resendmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadMemberByEmail", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectMember{}, svcerrors.NewProjectMemberNotFoundError())
				projectRepo.On("CreateInvitation", mock.Anything, mock.Anything).Return(svcerrors.NewProjectInvitationAlreadyExistsError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			creation: model.CreateProjectInvitationInput{
				Email:       "test@test.tt",
				ProjectRole: "viewer",
				ProjectID:   id.NewProject(),
			},
			mockSetup: func(auth *svc.AuthorizationService, user *svc.UserService, email *resendmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadMemberByEmail", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectMember{}, svcerrors.NewProjectMemberNotFoundError())
				projectRepo.On("CreateInvitation", mock.Anything, mock.Anything).Return(nil)
				email.On("SendProjectInvitationEmailAsync", mock.Anything, mock.Anything, mock.Anything)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, userSvc, email, projectRepo)

			_, err := service.Invite(context.Background(), tc.creation, id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			email.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetInvitations(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Unknown project",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ListInvitationsForProject", mock.Anything, mock.Anything).Return([]model.ProjectInvitation{}, nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "No invitations",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ListInvitationsForProject", mock.Anything, mock.Anything).Return([]model.ProjectInvitation{}, nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ListInvitationsForProject", mock.Anything, mock.Anything).Return(
					[]model.ProjectInvitation{
						{Email: "test@test.tt", ProjectRole: model.ProjectRoleEditor, Status: model.InvitationStatusPending, ProjectID: id.NewProject()},
					},
					nil,
				)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.ListInvitations(context.Background(), id.NewProject(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_CancelInvitation(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Unknown invitation",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("DeleteInvitation", mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewProjectInvitationNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("DeleteInvitation", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.CancelInvitation(context.Background(), id.NewProject(), id.NewProjectInvitation(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_AcceptInvitation(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(user *svc.UserService, repository *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Unknown invitation",
			mockSetup: func(user *svc.UserService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("UpdateInvitation", mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewProjectInvitationNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Invitation accepted, user not exists yet",
			mockSetup: func(user *svc.UserService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("UpdateInvitation", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("CreateMember", mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewUserNotFoundError())
			},
			wantErr: false,
		},
		{
			name: "Invitation accepted, member created",
			mockSetup: func(user *svc.UserService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("UpdateInvitation", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("CreateMember", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(userSvc, projectRepo)

			tkn, err := model.NewProjectInvitationToken()
			if err != nil {
				t.Fatal(err)
			}

			err = service.AcceptInvitation(context.Background(), tkn)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			userSvc.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_RejectInvitation(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(repository *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Unknown invitation",
			mockSetup: func(projectRepo *repo.ProjectRepository) {
				projectRepo.On("DeleteInvitationByTokenHashAndStatus", mock.Anything, mock.Anything, mock.Anything).
					Return(svcerrors.NewProjectInvitationNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(projectRepo *repo.ProjectRepository) {
				projectRepo.On("DeleteInvitationByTokenHashAndStatus", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(projectRepo)

			tkn, err := model.NewProjectInvitationToken()
			if err != nil {
				t.Fatal(err)
			}

			err = service.RejectInvitation(context.Background(), tkn)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_ListMembersForProject(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Non existing project",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ListMembersForProject", mock.Anything, mock.Anything).Return([]model.ProjectMember{}, nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "No members in project",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ListMembersForProject", mock.Anything, mock.Anything).Return([]model.ProjectMember{}, nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ListMembersForProject", mock.Anything, mock.Anything).Return([]model.ProjectMember{
					{ProjectID: id.NewProject()},
					{ProjectID: id.NewProject()},
				}, nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.ListMembersForProject(context.Background(), id.NewProject(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_ListMembersForUser(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleUser", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ListMembersForUser", mock.Anything, mock.Anything).Return([]model.ProjectMember{}, nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.ListMembersForUser(context.Background(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_DeleteMember(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Non existing member",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("DeleteMember", mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewProjectMemberNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("DeleteMember", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.DeleteMember(context.Background(), id.NewProject(), id.User{}, id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_UpdateMemberRole(t *testing.T) {
	testCases := []struct {
		name      string
		newRole   model.ProjectRole
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Non existing member",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("UpdateMember", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewProjectMemberNotFoundError())
			},
			wantErr: true,
		},
		{
			name:    "Success",
			newRole: model.ProjectRoleEditor,
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeUserRoleAdmin", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("UpdateMember", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.UpdateMemberRole(context.Background(), tc.newRole, id.NewProject(), id.User{}, id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_SetGithubRepoForProject(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *svc.SettingsService, *githubmock.Client, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, settingsSvc *svc.SettingsService, githubClient *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return(model.GithubToken("token"), nil)
				githubClient.On("ReadRepo", mock.Anything, mock.Anything, mock.Anything).Return(model.GithubRepo{}, nil)
				projectRepo.On("UpdateProject", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Github integration not enabled",
			mockSetup: func(auth *svc.AuthorizationService, settingsSvc *svc.SettingsService, githubClient *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return(model.GithubToken(""), svcerrors.NewGithubIntegrationNotEnabledError())
			},
			wantErr: true,
		},
		{
			name: "Github repo not found",
			mockSetup: func(auth *svc.AuthorizationService, settingsSvc *svc.SettingsService, githubClient *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return(model.GithubToken("token"), nil)
				githubClient.On("ReadRepo", mock.Anything, mock.Anything, mock.Anything).Return(model.GithubRepo{}, svcerrors.NewGithubRepoNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, settingsSvc, github, projectRepo)

			err := service.SetGithubRepoForProject(context.Background(), "https://github.com/test/test", id.NewProject(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
			settingsSvc.AssertExpectations(t)
			github.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetGithubRepoForProject(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					GithubRepo: &model.GithubRepo{
						OwnerSlug: "test",
						RepoSlug:  "test",
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "Project not found",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Project has no github repo",
			mockSetup: func(auth *svc.AuthorizationService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.GetGithubRepoForProject(context.Background(), id.NewProject(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_ListGithubRepoTags(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizationService, *svc.SettingsService, *githubmock.Client, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizationService, settingsSvc *svc.SettingsService, githubClient *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return(model.GithubToken("token"), nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					GithubRepo: &model.GithubRepo{
						OwnerSlug: "test",
						RepoSlug:  "test",
					},
				}, nil)
				githubClient.On("ReadTagsForRepo", mock.Anything, mock.Anything, mock.Anything).Return([]model.GitTag{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Github integration not enabled",
			mockSetup: func(auth *svc.AuthorizationService, settingsSvc *svc.SettingsService, githubClient *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return(model.GithubToken(""), svcerrors.NewGithubIntegrationNotEnabledError())
			},
			wantErr: true,
		},
		{
			name: "Project not found",
			mockSetup: func(auth *svc.AuthorizationService, settingsSvc *svc.SettingsService, githubClient *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return(model.GithubToken("token"), nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Github repo not set",
			mockSetup: func(auth *svc.AuthorizationService, settingsSvc *svc.SettingsService, githubClient *githubmock.Client, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return(model.GithubToken("token"), nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			github := new(githubmock.Client)
			email := new(resendmock.Client)
			userSvc := new(svc.UserService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthorizationService)
			service := NewProjectService(authSvc, settingsSvc, userSvc, email, github, projectRepo)

			tc.mockSetup(authSvc, settingsSvc, github, projectRepo)

			_, err := service.ListGithubRepoTags(context.Background(), id.NewProject(), id.AuthUser{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
			settingsSvc.AssertExpectations(t)
			github.AssertExpectations(t)
		})
	}
}
