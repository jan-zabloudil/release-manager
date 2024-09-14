package model

import (
	"testing"

	"release-manager/pkg/pointer"

	"github.com/stretchr/testify/assert"
)

func TestSettings_Update(t *testing.T) {
	tests := []struct {
		name     string
		settings Settings
		update   UpdateSettingsInput
		wantErr  bool
	}{
		{
			name: "Valid Update",
			settings: Settings{
				OrganizationName:      "Old Organization",
				DefaultReleaseMessage: "Old Message",
				Slack:                 SlackSettings{Enabled: false, Token: ""},
				Github:                GithubSettings{Enabled: false, Token: ""},
			},
			update: UpdateSettingsInput{
				OrganizationName:  pointer.StringPtr("New Organization"),
				DefaultReleaseMsg: pointer.StringPtr("New Message"),
				Slack:             UpdateSlackSettingsInput{Enabled: pointer.BoolPtr(true), Token: pointer.StringPtr("newToken")},
				Github:            UpdateGithubSettingsInput{Enabled: pointer.BoolPtr(true), Token: pointer.StringPtr("newToken")},
			},
			wantErr: false,
		},
		{
			name: "Invalid Update - missing slack token",
			settings: Settings{
				OrganizationName:      "Old Organization",
				DefaultReleaseMessage: "Old Message",
				Slack:                 SlackSettings{Enabled: false, Token: ""},
				Github:                GithubSettings{Enabled: false, Token: ""},
			},
			update: UpdateSettingsInput{
				Slack: UpdateSlackSettingsInput{Enabled: pointer.BoolPtr(true), Token: nil},
			},
			wantErr: true,
		},
		{
			name: "Invalid Update - missing github token",
			settings: Settings{
				OrganizationName:      "Old Organization",
				DefaultReleaseMessage: "Old Message",
				Slack:                 SlackSettings{Enabled: false, Token: ""},
				Github:                GithubSettings{Enabled: false, Token: ""},
			},
			update: UpdateSettingsInput{
				Github: UpdateGithubSettingsInput{Enabled: pointer.BoolPtr(true), Token: nil},
			},
			wantErr: true,
		},
		{
			name: "Invalid Update - missing org name",
			settings: Settings{
				OrganizationName:      "Old Organization",
				DefaultReleaseMessage: "Old Message",
				Slack:                 SlackSettings{Enabled: false, Token: ""},
				Github:                GithubSettings{Enabled: false, Token: ""},
			},
			update: UpdateSettingsInput{
				OrganizationName: pointer.StringPtr(""),
			},
			wantErr: true,
		},
		{
			name: "Invalid Update - missing default release message",
			settings: Settings{
				OrganizationName:      "Old Organization",
				DefaultReleaseMessage: "Old Message",
				Slack:                 SlackSettings{Enabled: false, Token: ""},
				Github:                GithubSettings{Enabled: false, Token: ""},
			},
			update: UpdateSettingsInput{
				DefaultReleaseMsg: pointer.StringPtr(""),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.settings.Update(tt.update)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, *tt.update.OrganizationName, tt.settings.OrganizationName)
				assert.Equal(t, *tt.update.DefaultReleaseMsg, tt.settings.DefaultReleaseMessage)
				assert.Equal(t, *tt.update.Slack.Enabled, tt.settings.Slack.Enabled)
				assert.Equal(t, *tt.update.Slack.Token, tt.settings.Slack.Token)
				assert.Equal(t, *tt.update.Github.Enabled, tt.settings.Github.Enabled)
				assert.Equal(t, *tt.update.Github.Token, tt.settings.Github.Token)
			}
		})
	}
}

func TestSettings_Validate(t *testing.T) {
	tests := []struct {
		name     string
		settings Settings
		wantErr  bool
	}{
		{
			name: "Valid Settings",
			settings: Settings{
				OrganizationName:      "Test Organization",
				DefaultReleaseMessage: "Test Message",
				Slack:                 SlackSettings{Enabled: true, Token: "slackToken"},
				Github:                GithubSettings{Enabled: true, Token: "githubToken"},
			},
			wantErr: false,
		},
		{
			name: "Invalid Settings - missing organization name",
			settings: Settings{
				OrganizationName:      "",
				DefaultReleaseMessage: "Test Message",
				Slack:                 SlackSettings{Enabled: true, Token: "slackToken"},
				Github:                GithubSettings{Enabled: true, Token: "githubToken"},
			},
			wantErr: true,
		},
		{
			name: "Invalid Settings - missing default release message",
			settings: Settings{
				OrganizationName:      "Test Organization",
				DefaultReleaseMessage: "",
				Slack:                 SlackSettings{Enabled: true, Token: "slackToken"},
				Github:                GithubSettings{Enabled: true, Token: "githubToken"},
			},
			wantErr: true,
		},
		{
			name: "Invalid Settings - missing slack token",
			settings: Settings{
				OrganizationName:      "Test Organization",
				DefaultReleaseMessage: "",
				Slack:                 SlackSettings{Enabled: true, Token: ""},
				Github:                GithubSettings{Enabled: true, Token: "githubToken"},
			},
			wantErr: true,
		},
		{
			name: "Invalid Settings - missing github token",
			settings: Settings{
				OrganizationName:      "Test Organization",
				DefaultReleaseMessage: "",
				Slack:                 SlackSettings{Enabled: true, Token: "slackToken"},
				Github:                GithubSettings{Enabled: true, Token: ""},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.settings.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
