package service

import (
	"context"
	"testing"

	repo "release-manager/repository/mock"
	svcerrors "release-manager/service/errors"
	svc "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeploymentService_Create(t *testing.T) {
	testCases := []struct {
		name      string
		input     model.CreateDeploymentInput
		mockSetup func(*svc.AuthorizeService, *svc.ProjectService, *svc.ReleaseService, *repo.DeploymentRepository)
		wantErr   bool
	}{
		{
			name: "success",
			input: model.CreateDeploymentInput{
				ReleaseID:     uuid.New(),
				EnvironmentID: uuid.New(),
			},
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseSvc *svc.ReleaseService, releaseRepo *repo.DeploymentRepository) {
				authSvc.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
				releaseSvc.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				projectSvc.On("GetEnvironment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
				releaseRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid input",
			input: model.CreateDeploymentInput{
				ReleaseID:     uuid.Nil,
				EnvironmentID: uuid.Nil,
			},
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseSvc *svc.ReleaseService, releaseRepo *repo.DeploymentRepository) {
				authSvc.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "project not found",
			input: model.CreateDeploymentInput{
				ReleaseID:     uuid.New(),
				EnvironmentID: uuid.New(),
			},
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseSvc *svc.ReleaseService, releaseRepo *repo.DeploymentRepository) {
				authSvc.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
			},
			wantErr: true,
		},
		{
			name: "release not found",
			input: model.CreateDeploymentInput{
				ReleaseID:     uuid.New(),
				EnvironmentID: uuid.New(),
			},
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseSvc *svc.ReleaseService, releaseRepo *repo.DeploymentRepository) {
				authSvc.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
				releaseSvc.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "env not found",
			input: model.CreateDeploymentInput{
				ReleaseID:     uuid.New(),
				EnvironmentID: uuid.New(),
			},
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseSvc *svc.ReleaseService, releaseRepo *repo.DeploymentRepository) {
				authSvc.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
				releaseSvc.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				projectSvc.On("GetEnvironment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, svcerrors.NewEnvironmentNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := new(svc.AuthorizeService)
			projectSvc := new(svc.ProjectService)
			releaseSvc := new(svc.ReleaseService)
			releaseRepo := new(repo.DeploymentRepository)
			service := NewDeploymentService(authSvc, projectSvc, releaseSvc, projectSvc, releaseRepo)

			tc.mockSetup(authSvc, projectSvc, releaseSvc, releaseRepo)

			_, err := service.Create(context.TODO(), tc.input, uuid.Nil, uuid.Nil)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			projectSvc.AssertExpectations(t)
			releaseSvc.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}

func TestDeploymentService_ListForProject(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizeService, *svc.ProjectService, *repo.DeploymentRepository)
		wantErr   bool
	}{
		{
			name: "Deployments fetched successfully",
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.DeploymentRepository) {
				authSvc.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("ListForProject", mock.Anything, mock.Anything).Return([]model.Deployment{
					{
						ID: uuid.New(),
					},
					{
						ID: uuid.New(),
					},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "Deployments fetched successfully but no deployments found",
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.DeploymentRepository) {
				authSvc.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("ListForProject", mock.Anything, mock.Anything).Return([]model.Deployment{}, nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
			},
			wantErr: false,
		},
		{
			name: "Project not found",
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.DeploymentRepository) {
				authSvc.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("ListForProject", mock.Anything, mock.Anything).Return([]model.Deployment{}, nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := new(svc.AuthorizeService)
			projectSvc := new(svc.ProjectService)
			releaseSvc := new(svc.ReleaseService)
			releaseRepo := new(repo.DeploymentRepository)
			service := NewDeploymentService(authSvc, projectSvc, releaseSvc, projectSvc, releaseRepo)

			tc.mockSetup(authSvc, projectSvc, releaseRepo)

			_, err := service.ListForProject(context.TODO(), uuid.Nil, uuid.Nil)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			projectSvc.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}
