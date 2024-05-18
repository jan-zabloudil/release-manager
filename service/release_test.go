package service

import (
	"context"
	"testing"

	"release-manager/pkg/apierrors"
	"release-manager/pkg/dberrors"
	repo "release-manager/repository/mock"
	svc "release-manager/service/mock"
	"release-manager/service/model"
	slack "release-manager/slack/mock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReleaseService_Create(t *testing.T) {
	testCases := []struct {
		name                    string
		release                 model.CreateReleaseInput
		sendReleaseNotification bool
		mockSetup               func(*svc.ProjectService, *svc.SettingsService, *slack.Client, *repo.ReleaseRepository)
		wantErr                 bool
	}{
		{
			name: "Create release without sending notification",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
			},
			sendReleaseNotification: false,
			mockSetup: func(projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				releaseRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Create release with notification",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
			},
			sendReleaseNotification: true,
			mockSetup: func(projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "channel",
				}, nil)
				releaseRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("token", nil)
				slackClient.On("SendReleaseNotificationAsync", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			wantErr: false,
		},
		{
			name: "Create release (slack integration not enabled)",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
			},
			sendReleaseNotification: true,
			mockSetup: func(projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "channel",
				}, nil)
				releaseRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("token", apierrors.NewSlackIntegrationNotEnabledError())
			},
			wantErr: false,
		},
		{
			name: "Create release (project has no slack channel)",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
			},
			sendReleaseNotification: true,
			mockSetup: func(projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "",
				}, nil)
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
			mockSetup: func(projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
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
			mockSetup: func(projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			projectSvc := new(svc.ProjectService)
			settingsSvc := new(svc.SettingsService)
			releaseRepo := new(repo.ReleaseRepository)
			slackClient := new(slack.Client)
			service := NewReleaseService(projectSvc, settingsSvc, slackClient, releaseRepo)

			tc.mockSetup(projectSvc, settingsSvc, slackClient, releaseRepo)

			_, err := service.Create(context.TODO(), tc.release, tc.sendReleaseNotification, uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			projectSvc.AssertExpectations(t)
			settingsSvc.AssertExpectations(t)
			slackClient.AssertExpectations(t)
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
			settingsSvc := new(svc.SettingsService)
			releaseRepo := new(repo.ReleaseRepository)
			slackClient := new(slack.Client)
			service := NewReleaseService(projectSvc, settingsSvc, slackClient, releaseRepo)

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
			settingsSvc := new(svc.SettingsService)
			releaseRepo := new(repo.ReleaseRepository)
			slackClient := new(slack.Client)
			service := NewReleaseService(projectSvc, settingsSvc, slackClient, releaseRepo)

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
