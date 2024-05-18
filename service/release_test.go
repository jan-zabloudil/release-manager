package service

import (
	"context"
	"testing"

	"release-manager/pkg/apierrors"
	"release-manager/pkg/dberrors"
	repo "release-manager/repository/mock"
	svc "release-manager/service/mock"
	"release-manager/service/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReleaseService_Create(t *testing.T) {
	testCases := []struct {
		name      string
		release   model.CreateReleaseInput
		mockSetup func(*svc.ProjectService, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Valid release",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
			},
			mockSetup: func(projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				releaseRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Unknown project",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Release",
				ReleaseNotes: "Test release notes",
			},
			mockSetup: func(projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, dberrors.NewNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Missing release title",
			release: model.CreateReleaseInput{
				ReleaseTitle: "",
				ReleaseNotes: "Test release notes",
			},
			mockSetup: func(projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectSvc := new(svc.ProjectService)
			releaseRepo := new(repo.ReleaseRepository)
			service := NewReleaseService(projectSvc, releaseRepo)

			tc.mockSetup(projectSvc, releaseRepo)

			_, err := service.Create(context.TODO(), tc.release, uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectSvc.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}

func TestReleaseService_Get(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Existing release",
			mockSetup: func(repo *repo.ReleaseRepository) {
				repo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Non-existing release",
			mockSetup: func(repo *repo.ReleaseRepository) {
				repo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, apierrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectSvc := new(svc.ProjectService)
			releaseRepo := new(repo.ReleaseRepository)
			service := NewReleaseService(projectSvc, releaseRepo)

			tc.mockSetup(releaseRepo)

			_, err := service.Get(context.TODO(), uuid.New(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectSvc.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}

func TestReleaseService_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func(repo *repo.ReleaseRepository) {
				repo.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Non-existing release",
			mockSetup: func(repo *repo.ReleaseRepository) {
				repo.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(apierrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectSvc := new(svc.ProjectService)
			releaseRepo := new(repo.ReleaseRepository)
			service := NewReleaseService(projectSvc, releaseRepo)

			tc.mockSetup(releaseRepo)

			err := service.Delete(context.TODO(), uuid.New(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectSvc.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}
