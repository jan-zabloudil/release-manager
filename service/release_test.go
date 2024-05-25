package service

import (
	"context"
	"testing"

	repo "release-manager/repository/mock"
	svcerrors "release-manager/service/errors"
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
		mockSetup               func(*svc.AuthorizeService, *svc.ProjectService, *svc.SettingsService, *slack.Client, *repo.ReleaseRepository)
		wantErr                 bool
	}{
		{
			name: "Create release without sending notification",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
				GitTagName:   "v1.0.0",
			},
			sendReleaseNotification: false,
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
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
				GitTagName:   "v1.0.0",
			},
			sendReleaseNotification: true,
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
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
				GitTagName:   "v1.0.0",
			},
			sendReleaseNotification: true,
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "channel",
				}, nil)
				releaseRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("token", svcerrors.NewSlackIntegrationNotEnabledError())
			},
			wantErr: false,
		},
		{
			name: "Create release (project has no slack channel)",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
				GitTagName:   "v1.0.0",
			},
			sendReleaseNotification: true,
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
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
				GitTagName:   "v1.0.0",
			},
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Missing release title",
			release: model.CreateReleaseInput{
				ReleaseTitle: "",
				ReleaseNotes: "Test release notes",
				GitTagName:   "v1.0.0",
			},
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := new(svc.AuthorizeService)
			projectSvc := new(svc.ProjectService)
			settingsSvc := new(svc.SettingsService)
			releaseRepo := new(repo.ReleaseRepository)
			slackClient := new(slack.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, releaseRepo)

			tc.mockSetup(authSvc, projectSvc, settingsSvc, slackClient, releaseRepo)

			_, err := service.Create(context.TODO(), tc.release, tc.sendReleaseNotification, uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
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
		mockSetup func(*svc.AuthorizeService, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Existing release",
			mockSetup: func(auth *svc.AuthorizeService, repo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Non-existing release",
			mockSetup: func(auth *svc.AuthorizeService, repo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := new(svc.AuthorizeService)
			projectSvc := new(svc.ProjectService)
			settingsSvc := new(svc.SettingsService)
			releaseRepo := new(repo.ReleaseRepository)
			slackClient := new(slack.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, releaseRepo)

			tc.mockSetup(authSvc, releaseRepo)

			_, err := service.Get(context.TODO(), uuid.New(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}

func TestReleaseService_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizeService, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizeService, repo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Non-existing release",
			mockSetup: func(auth *svc.AuthorizeService, repo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := new(svc.AuthorizeService)
			projectSvc := new(svc.ProjectService)
			settingsSvc := new(svc.SettingsService)
			releaseRepo := new(repo.ReleaseRepository)
			slackClient := new(slack.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, releaseRepo)

			tc.mockSetup(authSvc, releaseRepo)

			err := service.Delete(context.TODO(), uuid.New(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}

func TestReleaseService_ListForProject(t *testing.T) {
	testCases := []struct {
		name      string
		projectID uuid.UUID
		mockSetup func(*svc.AuthorizeService, *svc.ProjectService, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name:      "Success",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("ListForProject", mock.Anything, mock.Anything).Return([]model.Release{
					{ID: uuid.New()},
					{ID: uuid.New()},
				}, nil)
			},
			wantErr: false,
		},
		{
			name:      "no releases",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("ListForProject", mock.Anything, mock.Anything).Return([]model.Release{}, nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
			},
			wantErr: false,
		},
		{
			name:      "project not found",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("ListForProject", mock.Anything, mock.Anything).Return([]model.Release{}, nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := new(svc.AuthorizeService)
			projectSvc := new(svc.ProjectService)
			settingsSvc := new(svc.SettingsService)
			releaseRepo := new(repo.ReleaseRepository)
			slackClient := new(slack.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, releaseRepo)

			tc.mockSetup(authSvc, projectSvc, releaseRepo)

			_, err := service.ListForProject(context.Background(), tc.projectID, uuid.New())

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

func TestReleaseService_Update(t *testing.T) {
	validName := "Test release"
	validNotes := "Test release notes"
	invalidName := ""

	testCases := []struct {
		name      string
		update    model.UpdateReleaseInput
		mockSetup func(*svc.AuthorizeService, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Valid release update",
			update: model.UpdateReleaseInput{
				ReleaseTitle: &validName,
				ReleaseNotes: &validNotes,
			},
			mockSetup: func(auth *svc.AuthorizeService, repo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Empty release title",
			update: model.UpdateReleaseInput{
				ReleaseTitle: &invalidName,
				ReleaseNotes: &validNotes,
			},
			mockSetup: func(auth *svc.AuthorizeService, repo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseUnprocessableError())
			},
			wantErr: true,
		},
		{
			name:   "Non existing release",
			update: model.UpdateReleaseInput{},
			mockSetup: func(auth *svc.AuthorizeService, repo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := new(svc.AuthorizeService)
			projectSvc := new(svc.ProjectService)
			settingsSvc := new(svc.SettingsService)
			releaseRepo := new(repo.ReleaseRepository)
			slackClient := new(slack.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, releaseRepo)

			tc.mockSetup(authSvc, releaseRepo)

			_, err := service.Update(context.Background(), tc.update, uuid.New(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}

func TestReleaseService_SendReleaseNotification(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizeService, *svc.ProjectService, *svc.SettingsService, *slack.Client, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Send release notification",
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "channel",
				}, nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("token", nil)
				slackClient.On("SendReleaseNotificationAsync", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			},
			wantErr: false,
		},
		{
			name: "Release not found",
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Slack integration not enabled",
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "channel",
				}, nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("", svcerrors.NewSlackIntegrationNotEnabledError())
			},
			wantErr: true,
		},
		{
			name: "Project has no slack channel",
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "",
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authSvc := new(svc.AuthorizeService)
			projectSvc := new(svc.ProjectService)
			settingsSvc := new(svc.SettingsService)
			releaseRepo := new(repo.ReleaseRepository)
			slackClient := new(slack.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, releaseRepo)

			tc.mockSetup(authSvc, projectSvc, settingsSvc, slackClient, releaseRepo)

			err := service.SendReleaseNotification(context.TODO(), uuid.New(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			projectSvc.AssertExpectations(t)
			settingsSvc.AssertExpectations(t)
			slackClient.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}
