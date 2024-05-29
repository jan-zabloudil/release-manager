package service

import (
	"context"
	"net/url"
	"testing"

	github "release-manager/github/mock"
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
		name      string
		release   model.CreateReleaseInput
		mockSetup func(*svc.AuthorizeService, *svc.SettingsService, *svc.ProjectService, *github.Client, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Create release",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
				GitTagName:   "v1.0.0",
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					GithubRepo: &model.GithubRepo{
						OwnerSlug: "owner",
						RepoSlug:  "repo",
					},
				}, nil)
				github.On("GenerateGitTagURL", mock.Anything, mock.Anything).Return(url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/owner/repo/releases/tag/v1.0.0",
				}, nil)
				github.On("TagExists", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
				releaseRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Missing git tag",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					GithubRepo: &model.GithubRepo{
						OwnerSlug: "owner",
						RepoSlug:  "repo",
					},
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "Github integration not enabled",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
				GitTagName:   "v1.0.0",
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("", svcerrors.NewGithubIntegrationNotEnabledError())
			},
			wantErr: true,
		},
		{
			name: "Github integration not enabled",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("", svcerrors.NewGithubIntegrationNotEnabledError())
			},
			wantErr: true,
		},
		{
			name: "Github repo not set",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
				GitTagName:   "v1.0.0",
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Git tag not found",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
				GitTagName:   "v1.0.0",
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					GithubRepo: &model.GithubRepo{
						OwnerSlug: "owner",
						RepoSlug:  "repo",
					},
				}, nil)
				github.On("GenerateGitTagURL", mock.Anything, mock.Anything).Return(url.URL{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/owner/repo/releases/tag/v1.0.0",
				}, nil)
				github.On("TagExists", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
			},
			wantErr: true,
		},
		{
			name: "Project not found",
			release: model.CreateReleaseInput{
				ReleaseTitle: "Test release",
				ReleaseNotes: "Test release notes",
				GitTagName:   "v1.0.0",
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
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
			githubClient := new(github.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, githubClient, releaseRepo)

			tc.mockSetup(authSvc, settingsSvc, projectSvc, githubClient, releaseRepo)

			_, err := service.Create(context.TODO(), tc.release, uuid.New(), uuid.New())

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
			githubClient := new(github.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, githubClient, releaseRepo)

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
			githubClient := new(github.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, githubClient, releaseRepo)

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
			githubClient := new(github.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, githubClient, releaseRepo)

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
			githubClient := new(github.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, githubClient, releaseRepo)

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
				slackClient.On("SendReleaseNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
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
			githubClient := new(github.Client)
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, slackClient, githubClient, releaseRepo)

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
