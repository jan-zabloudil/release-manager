package service

import (
	"context"
	"errors"
	"testing"

	github "release-manager/github/mock"
	"release-manager/pkg/apierrors"
	cryptox "release-manager/pkg/crypto"
	"release-manager/pkg/dberrors"
	repo "release-manager/repository/mock"
	svc "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProjectService_CreateProject(t *testing.T) {
	testCases := []struct {
		name      string
		project   model.CreateProjectInput
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Valid project",
			project: model.CreateProjectInput{
				Name:                      "Test projectGetter",
				SlackChannelID:            "c1234",
				ReleaseNotificationConfig: model.ReleaseNotificationConfig{},
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("CreateProject", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid project - missing name",
			project: model.CreateProjectInput{
				Name:                      "",
				SlackChannelID:            "",
				ReleaseNotificationConfig: model.ReleaseNotificationConfig{},
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid project - invalid github repository url",
			project: model.CreateProjectInput{
				Name:                      "",
				SlackChannelID:            "",
				ReleaseNotificationConfig: model.ReleaseNotificationConfig{},
				GithubRepositoryRawURL:    "https://github.com/owner",
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.CreateProject(context.Background(), tc.project, uuid.New())

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

func TestProjectService_GetProject(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name:      "Existing project",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "Non-existing project",
			projectID: uuid.Nil,
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.GetProject(context.Background(), tc.projectID, uuid.New())

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
		projectID uuid.UUID
		mockSetup func(auth *svc.AuthService, projectRepo *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name:      "Existing project",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("DeleteProject", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Non-existing project",
			projectID: uuid.Nil,
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.DeleteProject(context.Background(), tc.projectID, uuid.New())

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
	validProjectName := "projectGetter name"
	invalidProjectName := ""
	slackChannelID := "channelID"
	invalidGithubRepositoryURL := "https://github.com/owner"

	testCases := []struct {
		name      string
		update    model.UpdateProjectInput
		mockSetup func(auth *svc.AuthService, projectRepo *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Valid project update",
			update: model.UpdateProjectInput{
				Name:           &validProjectName,
				SlackChannelID: &slackChannelID,
				ReleaseNotificationConfigUpdate: model.UpdateReleaseNotificationConfigInput{
					Message:         new(string),
					ShowProjectName: new(bool),
					ShowReleaseName: new(bool),
					ShowChangelog:   new(bool),
					ShowDeployments: new(bool),
					ShowSourceCode:  new(bool),
				},
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("UpdateProject", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid project update - invalid github repository url",
			update: model.UpdateProjectInput{
				Name:           &validProjectName,
				SlackChannelID: &slackChannelID,
				ReleaseNotificationConfigUpdate: model.UpdateReleaseNotificationConfigInput{
					Message:         new(string),
					ShowProjectName: new(bool),
					ShowReleaseName: new(bool),
					ShowChangelog:   new(bool),
					ShowDeployments: new(bool),
					ShowSourceCode:  new(bool),
				},
				GithubRepositoryRawURL: &invalidGithubRepositoryURL,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid project update - missing name",
			update: model.UpdateProjectInput{
				Name: &invalidProjectName,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
		{
			name:   "Non-existing project",
			update: model.UpdateProjectInput{},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.UpdateProject(context.Background(), tc.update, uuid.New(), uuid.New())

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
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Valid environment creation",
			envCreate: model.CreateEnvironmentInput{
				ProjectID:     uuid.New(),
				Name:          "Test Environment",
				ServiceRawURL: "http://example.com",
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadEnvironmentByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, dberrors.NewNotFoundError())
				projectRepo.On("CreateEnvironment", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid environment creation - duplicate name",
			envCreate: model.CreateEnvironmentInput{
				ProjectID:     uuid.New(),
				Name:          "Test Environment",
				ServiceRawURL: "http://example.com",
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadEnvironmentByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid environment creation - not absolute service url",
			envCreate: model.CreateEnvironmentInput{
				ProjectID:     uuid.New(),
				Name:          "Test Environment",
				ServiceRawURL: "example",
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid environment creation - empty name",
			envCreate: model.CreateEnvironmentInput{
				ProjectID:     uuid.New(),
				Name:          "",
				ServiceRawURL: "example",
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.CreateEnvironment(context.Background(), tc.envCreate, uuid.New())

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
		projectID uuid.UUID
		envID     uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name:      "Success",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "project not found",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
		{
			name:      "environment not found",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, errors.New("env not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.GetEnvironment(context.Background(), tc.projectID, tc.envID, uuid.New())

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
	validURL := "http://example.com"
	validName := "Test Environment"
	invalidURL := "example"
	invalidName := ""

	testCases := []struct {
		name      string
		envUpdate model.UpdateEnvironmentInput
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			envUpdate: model.UpdateEnvironmentInput{
				Name:          &validName,
				ServiceRawURL: &validURL,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
				projectRepo.On("ReadEnvironmentByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, dberrors.NewNotFoundError())
				projectRepo.On("UpdateEnvironment", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid environment update - duplicate name",
			envUpdate: model.UpdateEnvironmentInput{
				Name:          &validName,
				ServiceRawURL: &validURL,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
				projectRepo.On("ReadEnvironmentByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid environment update - not absolute service url",
			envUpdate: model.UpdateEnvironmentInput{
				Name:          &validName,
				ServiceRawURL: &invalidURL,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid environment update - missing name",
			envUpdate: model.UpdateEnvironmentInput{
				Name:          &invalidName,
				ServiceRawURL: &validURL,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.UpdateEnvironment(context.Background(), tc.envUpdate, uuid.New(), uuid.New(), uuid.UUID{})

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
		projectID uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name:      "Success",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadAllEnvironmentsForProject", mock.Anything, mock.Anything).Return([]model.Environment{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "project not found",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.ListEnvironments(context.Background(), tc.projectID, uuid.New())

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
		projectID uuid.UUID
		envID     uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name:      "Success",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
				projectRepo.On("DeleteEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "project not found",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
		{
			name:      "env not found",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadEnvironment", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, errors.New("env not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.DeleteEnvironment(context.Background(), tc.projectID, tc.envID, uuid.New())

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

func TestProjectService_validateEnvironmentNameUnique(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
		isUnique  bool
	}{
		{
			name:      "Name is unique",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadEnvironmentByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, dberrors.NewNotFoundError())
			},
			isUnique: true,
			wantErr:  false,
		},
		{
			name:      "Name is duplicate",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadEnvironmentByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			isUnique: false,
			wantErr:  false,
		},
		{
			name:      "Unexpected error",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadEnvironmentByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, dberrors.NewUnknownError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			isUnique, err := service.isEnvironmentNameUnique(context.Background(), tc.projectID, "env")

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.isUnique, isUnique)
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
		mockSetup func(*svc.AuthService, *svc.EmailService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name:     "Unknown project",
			creation: model.CreateProjectInvitationInput{},
			mockSetup: func(auth *svc.AuthService, email *svc.EmailService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, apierrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Invalid project invitation - missing email",
			creation: model.CreateProjectInvitationInput{
				Email:       "",
				ProjectRole: "editor",
				ProjectID:   uuid.New(),
			},
			mockSetup: func(auth *svc.AuthService, email *svc.EmailService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Invalid project invitation - invalid role",
			creation: model.CreateProjectInvitationInput{
				Email:       "",
				ProjectRole: "new",
				ProjectID:   uuid.New(),
			},
			mockSetup: func(auth *svc.AuthService, email *svc.EmailService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invitation already exists",
			creation: model.CreateProjectInvitationInput{
				Email:       "test@test.tt",
				ProjectRole: "viewer",
				ProjectID:   uuid.New(),
			},
			mockSetup: func(auth *svc.AuthService, email *svc.EmailService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadInvitationByEmailForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Success",
			creation: model.CreateProjectInvitationInput{
				Email:       "test@test.tt",
				ProjectRole: "viewer",
				ProjectID:   uuid.New(),
			},
			mockSetup: func(auth *svc.AuthService, email *svc.EmailService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadInvitationByEmailForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, dberrors.NewNotFoundError())
				projectRepo.On("CreateInvitation", mock.Anything, mock.Anything).Return(nil)
				email.On("SendProjectInvitation", mock.Anything, mock.Anything)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, emailSvc, projectRepo)

			_, err := service.Invite(context.Background(), tc.creation, uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			emailSvc.AssertExpectations(t)
			authSvc.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetInvitations(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Unknown project",
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadAllInvitationsForProject", mock.Anything, mock.Anything).Return(
					[]model.ProjectInvitation{
						{Email: "test@test.tt", ProjectRole: model.ProjectRoleEditor, Status: model.InvitationStatusPending, ProjectID: uuid.New()},
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
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.ListInvitations(context.Background(), uuid.New(), uuid.New())

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
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Unknown project",
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Unknown invitation",
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadInvitation", mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadInvitation", mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, nil)
				projectRepo.On("DeleteInvitation", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.CancelInvitation(context.Background(), uuid.New(), uuid.New(), uuid.New())

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
		mockSetup func(repository *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name: "Unknown invitation",
			mockSetup: func(projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadInvitationByTokenHashAndStatus", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadInvitationByTokenHashAndStatus", mock.Anything, mock.Anything, mock.Anything).Return(
					model.ProjectInvitation{
						Email: "test@test.tt", ProjectRole: model.ProjectRoleEditor, Status: model.InvitationStatusPending, ProjectID: uuid.New(),
					},
					nil,
				)
				projectRepo.On("UpdateInvitation", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(projectRepo)

			tkn, err := cryptox.NewToken()
			if err != nil {
				t.Fatal(err)
			}

			err = service.AcceptInvitation(context.Background(), tkn)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

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
				projectRepo.On("ReadInvitationByTokenHashAndStatus", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectInvitation{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Success",
			mockSetup: func(projectRepo *repo.ProjectRepository) {
				projectRepo.On("ReadInvitationByTokenHashAndStatus", mock.Anything, mock.Anything, mock.Anything).Return(
					model.ProjectInvitation{
						Email: "test@test.tt", ProjectRole: model.ProjectRoleEditor, Status: model.InvitationStatusPending, ProjectID: uuid.New(),
					},
					nil,
				)
				projectRepo.On("DeleteInvitation", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(projectRepo)

			tkn, err := cryptox.NewToken()
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

func TestProjectService_ListMembers(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name:      "Non existing project",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name:      "Success",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadMembersForProject", mock.Anything, mock.Anything).Return([]model.ProjectMember{}, nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			_, err := service.ListMembers(context.Background(), tc.projectID, uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_DeleteMember(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository)
		wantErr   bool
	}{
		{
			name:      "Non existing project",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name:      "Non existing member",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectMember{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name:      "Success",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("ReadProject", mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("ReadMember", mock.Anything, mock.Anything, mock.Anything).Return(model.ProjectMember{}, nil)
				projectRepo.On("DeleteMember", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			githubClient := new(github.Client)
			emailSvc := new(svc.EmailService)
			settingsSvc := new(svc.SettingsService)
			authSvc := new(svc.AuthService)
			service := NewProjectService(authSvc, settingsSvc, githubClient, emailSvc, projectRepo)

			tc.mockSetup(authSvc, projectRepo)

			err := service.DeleteMember(context.Background(), tc.projectID, uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
		})
	}
}
