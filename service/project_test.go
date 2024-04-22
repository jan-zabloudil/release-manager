package service

import (
	"context"
	"errors"
	"testing"

	"release-manager/pkg/dberrors"
	repo "release-manager/repository/mock"
	svc "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProjectService_Create(t *testing.T) {
	testCases := []struct {
		name      string
		project   model.ProjectCreation
		mockSetup func(*svc.AuthService, *repo.ProjectRepository, *repo.EnvironmentRepository)
		wantErr   bool
	}{
		{
			name: "Valid project",
			project: model.ProjectCreation{
				Name:                      "Test Project",
				SlackChannelID:            "c1234",
				ReleaseNotificationConfig: model.ReleaseNotificationConfig{},
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid project - missing name",
			project: model.ProjectCreation{
				Name:                      "",
				SlackChannelID:            "",
				ReleaseNotificationConfig: model.ReleaseNotificationConfig{},
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			envRepo := new(repo.EnvironmentRepository)
			auth := new(svc.AuthService)
			service := NewProjectService(auth, projectRepo, envRepo)

			tc.mockSetup(auth, projectRepo, envRepo)

			_, err := service.Create(context.Background(), tc.project, uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			envRepo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_Get(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository, *repo.EnvironmentRepository)
		wantErr   bool
	}{
		{
			name:      "Existing project",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "Non-existing project",
			projectID: uuid.Nil,
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			envRepo := new(repo.EnvironmentRepository)
			auth := new(svc.AuthService)
			service := NewProjectService(auth, projectRepo, envRepo)

			tc.mockSetup(auth, projectRepo, envRepo)

			_, err := service.Get(context.Background(), tc.projectID, uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			envRepo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		mockSetup func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository)
		wantErr   bool
	}{
		{
			name:      "Existing project",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Read", mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Non-existing project",
			projectID: uuid.Nil,
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Read", mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			envRepo := new(repo.EnvironmentRepository)
			auth := new(svc.AuthService)
			service := NewProjectService(auth, projectRepo, envRepo)

			tc.mockSetup(auth, projectRepo, envRepo)

			err := service.Delete(context.Background(), tc.projectID, uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			envRepo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_Update(t *testing.T) {
	validProjectName := "Project name"
	invalidProjectName := ""
	slackChannelID := "channelID"

	testCases := []struct {
		name      string
		update    model.ProjectUpdate
		mockSetup func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository)
		wantErr   bool
	}{
		{
			name: "Valid project update",
			update: model.ProjectUpdate{
				Name:           &validProjectName,
				SlackChannelID: &slackChannelID,
				ReleaseNotificationConfigUpdate: model.ReleaseNotificationConfigUpdate{
					Message:         new(string),
					ShowProjectName: new(bool),
					ShowReleaseName: new(bool),
					ShowChangelog:   new(bool),
					ShowDeployments: new(bool),
					ShowSourceCode:  new(bool),
				},
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything).Return(model.Project{}, nil)
				projectRepo.On("Update", mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid project update - missing name",
			update: model.ProjectUpdate{
				Name: &invalidProjectName,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
		{
			name:   "Non-existing project",
			update: model.ProjectUpdate{},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			envRepo := new(repo.EnvironmentRepository)
			auth := new(svc.AuthService)
			service := NewProjectService(auth, projectRepo, envRepo)

			tc.mockSetup(auth, projectRepo, envRepo)

			_, err := service.Update(context.Background(), tc.update, uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			envRepo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_CreateEnvironment(t *testing.T) {
	testCases := []struct {
		name      string
		envCreate model.EnvironmentCreation
		mockSetup func(*svc.AuthService, *repo.ProjectRepository, *repo.EnvironmentRepository)
		wantErr   bool
	}{
		{
			name: "Valid environment creation",
			envCreate: model.EnvironmentCreation{
				ProjectID:     uuid.New(),
				Name:          "Test Environment",
				ServiceRawURL: "http://example.com",
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("ReadByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, dberrors.NewNotFoundError())
				envRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid environment creation - duplicate name",
			envCreate: model.EnvironmentCreation{
				ProjectID:     uuid.New(),
				Name:          "Test Environment",
				ServiceRawURL: "http://example.com",
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("ReadByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid environment creation - not absolute service url",
			envCreate: model.EnvironmentCreation{
				ProjectID:     uuid.New(),
				Name:          "Test Environment",
				ServiceRawURL: "example",
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid environment creation - empty name",
			envCreate: model.EnvironmentCreation{
				ProjectID:     uuid.New(),
				Name:          "",
				ServiceRawURL: "example",
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			envRepo := new(repo.EnvironmentRepository)
			auth := new(svc.AuthService)
			service := NewProjectService(auth, projectRepo, envRepo)

			tc.mockSetup(auth, projectRepo, envRepo)

			_, err := service.CreateEnvironment(context.Background(), tc.envCreate, uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			envRepo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetEnvironment(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		envID     uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository, *repo.EnvironmentRepository)
		wantErr   bool
	}{
		{
			name:      "Existing environment",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "Project not found",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
		{
			name:      "Environment not found",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, errors.New("env not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			envRepo := new(repo.EnvironmentRepository)
			auth := new(svc.AuthService)
			service := NewProjectService(auth, projectRepo, envRepo)

			tc.mockSetup(auth, projectRepo, envRepo)

			_, err := service.GetEnvironment(context.Background(), tc.projectID, tc.envID, uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			envRepo.AssertExpectations(t)
			auth.AssertExpectations(t)
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
		envUpdate model.EnvironmentUpdate
		mockSetup func(*svc.AuthService, *repo.ProjectRepository, *repo.EnvironmentRepository)
		wantErr   bool
	}{
		{
			name: "Valid Environment Update",
			envUpdate: model.EnvironmentUpdate{
				Name:          &validName,
				ServiceRawURL: &validURL,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
				envRepo.On("ReadByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, dberrors.NewNotFoundError())
				envRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid Environment Update - duplicate name",
			envUpdate: model.EnvironmentUpdate{
				Name:          &validName,
				ServiceRawURL: &validURL,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
				envRepo.On("ReadByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid Environment Update - not absolute service url",
			envUpdate: model.EnvironmentUpdate{
				Name:          &validName,
				ServiceRawURL: &invalidURL,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Invalid Environment Update - missing name",
			envUpdate: model.EnvironmentUpdate{
				Name:          &invalidName,
				ServiceRawURL: &validURL,
			},
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			envRepo := new(repo.EnvironmentRepository)
			auth := new(svc.AuthService)
			service := NewProjectService(auth, projectRepo, envRepo)

			tc.mockSetup(auth, projectRepo, envRepo)

			_, err := service.UpdateEnvironment(context.Background(), tc.envUpdate, uuid.New(), uuid.New(), uuid.UUID{})

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			envRepo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetEnvironments(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository, *repo.EnvironmentRepository)
		wantErr   bool
	}{
		{
			name:      "Success",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("ReadAllForProject", mock.Anything, mock.Anything).Return([]model.Environment{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "Project not found",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			envRepo := new(repo.EnvironmentRepository)
			auth := new(svc.AuthService)
			service := NewProjectService(auth, projectRepo, envRepo)

			tc.mockSetup(auth, projectRepo, envRepo)

			_, err := service.GetEnvironments(context.Background(), tc.projectID, uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			envRepo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_DeleteEnvironment(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		envID     uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository, *repo.EnvironmentRepository)
		wantErr   bool
	}{
		{
			name:      "Success",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
				envRepo.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Project not found",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, errors.New("project not found"))
			},
			wantErr: true,
		},
		{
			name:      "Env not found",
			projectID: uuid.New(),
			envID:     uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				auth.On("AuthorizeAdminRole", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				envRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, errors.New("env not found"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			envRepo := new(repo.EnvironmentRepository)
			auth := new(svc.AuthService)
			service := NewProjectService(auth, projectRepo, envRepo)

			tc.mockSetup(auth, projectRepo, envRepo)

			err := service.DeleteEnvironment(context.Background(), tc.projectID, tc.envID, uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			envRepo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}

func TestProjectService_validateEnvironmentNameUnique(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		mockSetup func(*svc.AuthService, *repo.ProjectRepository, *repo.EnvironmentRepository)
		wantErr   bool
		isUnique  bool
	}{
		{
			name:      "Name is unique",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				envRepo.On("ReadByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, dberrors.NewNotFoundError())
			},
			isUnique: true,
			wantErr:  false,
		},
		{
			name:      "Name is duplicate",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				envRepo.On("ReadByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
			},
			isUnique: false,
			wantErr:  false,
		},
		{
			name:      "Unexpected error",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthService, projectRepo *repo.ProjectRepository, envRepo *repo.EnvironmentRepository) {
				envRepo.On("ReadByNameForProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, dberrors.NewUnknownError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectRepo := new(repo.ProjectRepository)
			envRepo := new(repo.EnvironmentRepository)
			auth := new(svc.AuthService)
			service := NewProjectService(auth, projectRepo, envRepo)

			tc.mockSetup(auth, projectRepo, envRepo)

			isUnique, err := service.isEnvironmentNameUnique(context.Background(), tc.projectID, "env")

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.isUnique, isUnique)
				assert.NoError(t, err)
			}

			projectRepo.AssertExpectations(t)
			envRepo.AssertExpectations(t)
			auth.AssertExpectations(t)
		})
	}
}
