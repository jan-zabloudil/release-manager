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

func TestReleaseService_CreateRelease(t *testing.T) {
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
				releaseRepo.On("CreateRelease", mock.Anything, mock.Anything).Return(nil)
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
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackClient, githubClient, releaseRepo)

			tc.mockSetup(authSvc, settingsSvc, projectSvc, githubClient, releaseRepo)

			_, err := service.CreateRelease(context.TODO(), tc.release, uuid.New(), uuid.New())

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

func TestReleaseService_GetRelease(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizeService, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Existing release",
			mockSetup: func(auth *svc.AuthorizeService, repo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
			},
			wantErr: false,
		},
		{
			name: "Non-existing release",
			mockSetup: func(auth *svc.AuthorizeService, repo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseNotFoundError())
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
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackClient, githubClient, releaseRepo)

			tc.mockSetup(authSvc, releaseRepo)

			_, err := service.GetRelease(context.TODO(), uuid.New(), uuid.New(), uuid.New())

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

func TestReleaseService_DeleteRelease(t *testing.T) {
	testCases := []struct {
		name               string
		mockSetup          func(*svc.AuthorizeService, *svc.SettingsService, *svc.ProjectService, *github.Client, *repo.ReleaseRepository)
		deleteReleaseInput model.DeleteReleaseInput
		wantErr            bool
	}{
		{
			name: "Success without deleting github release",
			deleteReleaseInput: model.DeleteReleaseInput{
				DeleteGithubRelease: false,
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("DeleteRelease", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Success with deleting github release",
			deleteReleaseInput: model.DeleteReleaseInput{
				DeleteGithubRelease: true,
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
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				github.On("DeleteReleaseByTag", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("DeleteRelease", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Github integration not enabled",
			deleteReleaseInput: model.DeleteReleaseInput{
				DeleteGithubRelease: true,
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("", svcerrors.NewGithubIntegrationNotEnabledError())
			},
			wantErr: true,
		},
		{
			name: "Project not found",
			deleteReleaseInput: model.DeleteReleaseInput{
				DeleteGithubRelease: true,
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Release not found",
			deleteReleaseInput: model.DeleteReleaseInput{
				DeleteGithubRelease: true,
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Repo not set for project",
			deleteReleaseInput: model.DeleteReleaseInput{
				DeleteGithubRelease: true,
			},
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Success (deleting non-existing github release)",
			deleteReleaseInput: model.DeleteReleaseInput{
				DeleteGithubRelease: true,
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
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				github.On("DeleteReleaseByTag", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(svcerrors.NewGithubReleaseNotFoundError())
				releaseRepo.On("DeleteRelease", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
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
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackClient, githubClient, releaseRepo)

			tc.mockSetup(authSvc, settingsSvc, projectSvc, githubClient, releaseRepo)

			err := service.DeleteRelease(context.TODO(), tc.deleteReleaseInput, uuid.New(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			settingsSvc.AssertExpectations(t)
			projectSvc.AssertExpectations(t)
			githubClient.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}

func TestReleaseService_ListReleasesForProject(t *testing.T) {
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
				releaseRepo.On("ListReleasesForProject", mock.Anything, mock.Anything).Return([]model.Release{
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
				releaseRepo.On("ListReleasesForProject", mock.Anything, mock.Anything).Return([]model.Release{}, nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
			},
			wantErr: false,
		},
		{
			name:      "project not found",
			projectID: uuid.New(),
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("ListReleasesForProject", mock.Anything, mock.Anything).Return([]model.Release{}, nil)
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
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackClient, githubClient, releaseRepo)

			tc.mockSetup(authSvc, projectSvc, releaseRepo)

			_, err := service.ListReleasesForProject(context.Background(), tc.projectID, uuid.New())

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

func TestReleaseService_UpdateRelease(t *testing.T) {
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
				repo.On("UpdateRelease", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
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
				repo.On("UpdateRelease", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseUnprocessableError())
			},
			wantErr: true,
		},
		{
			name:   "Non existing release",
			update: model.UpdateReleaseInput{},
			mockSetup: func(auth *svc.AuthorizeService, repo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("UpdateRelease", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseNotFoundError())
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
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackClient, githubClient, releaseRepo)

			tc.mockSetup(authSvc, releaseRepo)

			_, err := service.UpdateRelease(context.Background(), tc.update, uuid.New(), uuid.New(), uuid.New())

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
			name: "Send release notification with deployment",
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "channel",
				}, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				releaseRepo.On("ReadLastDeploymentForRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Deployment{}, nil)
				slackClient.On("SendReleaseNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Send release notification without deployment",
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "channel",
				}, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				releaseRepo.On("ReadLastDeploymentForRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Deployment{}, svcerrors.NewDeploymentNotFoundError())
				slackClient.On("SendReleaseNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Project not found",
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Release not found",
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "channel",
				}, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Slack integration not enabled",
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("", svcerrors.NewSlackIntegrationNotEnabledError())
			},
			wantErr: true,
		},
		{
			name: "Project has no slack channel",
			mockSetup: func(auth *svc.AuthorizeService, projectSvc *svc.ProjectService, settingsSvc *svc.SettingsService, slackClient *slack.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetSlackToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					SlackChannelID: "",
				}, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				releaseRepo.On("ReadLastDeploymentForRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Deployment{}, nil)
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
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackClient, githubClient, releaseRepo)

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

func TestReleaseService_UpsertGithubRelease(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizeService, *svc.SettingsService, *svc.ProjectService, *github.Client, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Success",
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					GithubRepo: &model.GithubRepo{
						OwnerSlug: "owner",
						RepoSlug:  "repo",
					},
				}, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				github.On("UpsertRelease", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Github integration not enabled",
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("", svcerrors.NewGithubIntegrationNotEnabledError())
			},
			wantErr: true,
		},
		{
			name: "Project not found",
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, svcerrors.NewProjectNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Release not found",
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{
					GithubRepo: &model.GithubRepo{
						OwnerSlug: "owner",
						RepoSlug:  "repo",
					},
				}, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "Github repo not set for project",
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
			},
			wantErr: true,
		},
		{
			name: "Github repo not set for project",
			mockSetup: func(auth *svc.AuthorizeService, settingsSvc *svc.SettingsService, projectSvc *svc.ProjectService, github *github.Client, releaseRepo *repo.ReleaseRepository) {
				auth.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				settingsSvc.On("GetGithubToken", mock.Anything).Return("token", nil)
				projectSvc.On("GetProject", mock.Anything, mock.Anything, mock.Anything).Return(model.Project{}, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
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
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackClient, githubClient, releaseRepo)

			tc.mockSetup(authSvc, settingsSvc, projectSvc, githubClient, releaseRepo)

			err := service.UpsertGithubRelease(context.TODO(), uuid.New(), uuid.New(), uuid.New())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authSvc.AssertExpectations(t)
			settingsSvc.AssertExpectations(t)
			projectSvc.AssertExpectations(t)
			githubClient.AssertExpectations(t)
			releaseRepo.AssertExpectations(t)
		})
	}
}

func TestReleaseService_CreateDeployment(t *testing.T) {
	testCases := []struct {
		name      string
		input     model.CreateDeploymentInput
		mockSetup func(*svc.AuthorizeService, *svc.ProjectService, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "success",
			input: model.CreateDeploymentInput{
				ReleaseID:     uuid.New(),
				EnvironmentID: uuid.New(),
			},
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				authSvc.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				projectSvc.On("GetEnvironment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, nil)
				releaseRepo.On("CreateDeployment", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid input",
			input: model.CreateDeploymentInput{
				ReleaseID:     uuid.Nil,
				EnvironmentID: uuid.Nil,
			},
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
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
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
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
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				authSvc.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, svcerrors.NewReleaseNotFoundError())
			},
			wantErr: true,
		},
		{
			name: "env not found",
			input: model.CreateDeploymentInput{
				ReleaseID:     uuid.New(),
				EnvironmentID: uuid.New(),
			},
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				authSvc.On("AuthorizeProjectRoleEditor", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
				releaseRepo.On("ReadRelease", mock.Anything, mock.Anything, mock.Anything).Return(model.Release{}, nil)
				projectSvc.On("GetEnvironment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.Environment{}, svcerrors.NewEnvironmentNotFoundError())
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
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackClient, githubClient, releaseRepo)

			tc.mockSetup(authSvc, projectSvc, releaseRepo)

			_, err := service.CreateDeployment(context.TODO(), tc.input, uuid.Nil, uuid.Nil)
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

func TestReleaseService_ListDeploymentsForProject(t *testing.T) {
	testCases := []struct {
		name      string
		mockSetup func(*svc.AuthorizeService, *svc.ProjectService, *repo.ReleaseRepository)
		wantErr   bool
	}{
		{
			name: "Deployments fetched successfully",
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				authSvc.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("ListDeploymentsForProject", mock.Anything, mock.Anything).Return([]model.Deployment{
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
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				authSvc.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("ListDeploymentsForProject", mock.Anything, mock.Anything).Return([]model.Deployment{}, nil)
				projectSvc.On("ProjectExists", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
			},
			wantErr: false,
		},
		{
			name: "Project not found",
			mockSetup: func(authSvc *svc.AuthorizeService, projectSvc *svc.ProjectService, releaseRepo *repo.ReleaseRepository) {
				authSvc.On("AuthorizeProjectRoleViewer", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				releaseRepo.On("ListDeploymentsForProject", mock.Anything, mock.Anything).Return([]model.Deployment{}, nil)
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
			service := NewReleaseService(authSvc, projectSvc, settingsSvc, projectSvc, slackClient, githubClient, releaseRepo)

			tc.mockSetup(authSvc, projectSvc, releaseRepo)

			_, err := service.ListDeploymentsForProject(context.TODO(), model.DeploymentFilterParams{}, uuid.Nil, uuid.Nil)
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
